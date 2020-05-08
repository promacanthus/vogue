---
title: 03-binlog
date: 2020-04-14T10:09:14.242627+08:00
draft: false
---

从binlog恢复数据步骤：

1. 根据出现问题的事件定位到需要恢复的binlog位置
2. 清空数据库，将全量备份恢复
3. 根据binlog的位置恢复

## 定位binlog的位置

查看binlog是否开启：

```bash
mysql> show variables like 'log_%';
+----------------------------------------+----------------------------------------+
| Variable_name                          | Value                                  |
+----------------------------------------+----------------------------------------+
| log_bin                                | ON                                     |
| log_bin_basename                       | /var/lib/mysql/binlog                  |
| log_bin_index                          | /var/lib/mysql/binlog.index            |
| log_bin_trust_function_creators        | OFF                                    |
| log_bin_use_v1_row_events              | OFF                                    |
| log_error                              | stderr                                 |
| log_error_services                     | log_filter_internal; log_sink_internal |
| log_error_suppression_list             |                                        |
| log_error_verbosity                    | 2                                      |
| log_output                             | FILE                                   |
| log_queries_not_using_indexes          | OFF                                    |
| log_slave_updates                      | ON                                     |
| log_slow_admin_statements              | OFF                                    |
| log_slow_extra                         | OFF                                    |
| log_slow_slave_statements              | OFF                                    |
| log_statements_unsafe_for_binlog       | ON                                     |
| log_throttle_queries_not_using_indexes | 0                                      |
| log_timestamps                         | UTC                                    |
+----------------------------------------+----------------------------------------+
18 rows in set (0.00 sec)
```

查看binlog日志列表：

```bash
mysql> show logs;
ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'logs' at line 1

mysql> show master logs;
+---------------+-----------+-----------+
| Log_name      | File_size | Encrypted |
+---------------+-----------+-----------+
| binlog.000001 |   3091158 | No        |
| binlog.000002 | 141156437 | No        |
+---------------+-----------+-----------+
2 rows in set (0.17 sec)

```

查看master状态，也就是最新一个binlog日志编号名称和最后一个操作事件pos结束位置：

```bash
mysql> show master status;
+---------------+-----------+--------------+------------------+-------------------+
| File          | Position  | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------+-----------+--------------+------------------+-------------------+
| binlog.000002 | 141156437 |              |                  |                   |
+---------------+-----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

```

刷新log日志，将会产生一个新编号的binlog日志文件：

```bash
mysql> flush logs;
```

如果需要清空原来的binlog：

```bash
mysql> reset master;

```

查看binlog内容：

```bash
# shell方式
mysqlbinlog -v --base64-output=decode-rows /var/lib/mysql/master.000003

# 带条件查询
mysqlbinlog -v --base64-output=decode-rows /var/lib/mysql/master.000003 \
    --start-datetime="2019-03-01 00:00:00"  \
    --stop-datetime="2019-03-10 00:00:00"   \
    --start-position="5000"    \
    --stop-position="20000"

# binlog输入出下
# at 21019
#190308 10:10:09 server id 1  end_log_pos 21094 CRC32 0x7a405abc Query thread_id=113 exec_time=0 error_code=0
SET TIMESTAMP=1552011009/*!*/;
BEGIN
/*!*/;

```

每个字段的含义：

- position: 位于文件中的位置，即第一行的（# at 21019）,说明该事件记录从文件第21019个字节开始
- timestamp: 事件发生的时间戳，即第二行的（#190308 10:10:09）
- server id: 服务器标识（1）
- end_log_pos 表示下一个事件开始的位置（即当前事件的结束位置+1）
- thread_id: 执行该事件的线程id （thread_id=113）
- exec_time: 事件执行的花费时间
- error_code: 错误码，0意味着没有发生错误
- type:事件类型Query

mysql客户端查看binlog：

```bash
mysql> show binlog events in 'binlog.000002' from 968 limit 10;
+---------------+------+----------------+-----------+-------------+---------------------------------------------------------------------------------------------------------------------------------+
| Log_name      | Pos  | Event_type     | Server_id | End_log_pos | Info                                                                                                                            |
+---------------+------+----------------+-----------+-------------+---------------------------------------------------------------------------------------------------------------------------------+
| binlog.000002 |  968 | Anonymous_Gtid |         1 |        1047 | SET @@SESSION.GTID_NEXT= 'ANONYMOUS'                                                                                            |
| binlog.000002 | 1047 | Query          |         1 |        1256 | ALTER USER 'cuishifeng'@'%' IDENTIFIED WITH 'mysql_native_password' AS '*10320381F36BE49A18F09B06A4BC005223975101' /* xid=12 */ |
| binlog.000002 | 1256 | Anonymous_Gtid |         1 |        1333 | SET @@SESSION.GTID_NEXT= 'ANONYMOUS'                                                                                            |
| binlog.000002 | 1333 | Query          |         1 |        1423 | flush privileges                                                                                                                |
| binlog.000002 | 1423 | Anonymous_Gtid |         1 |        1500 | SET @@SESSION.GTID_NEXT= 'ANONYMOUS'                                                                                            |
| binlog.000002 | 1500 | Query          |         1 |        1646 | GRANT ALL PRIVILEGES ON *.* TO 'cuishifeng'@'%' /* xid=70 */                                                                    |
| binlog.000002 | 1646 | Anonymous_Gtid |         1 |        1723 | SET @@SESSION.GTID_NEXT= 'ANONYMOUS'                                                                                            |
| binlog.000002 | 1723 | Query          |         1 |        1813 | flush privileges                                                                                                                |
| binlog.000002 | 1813 | Anonymous_Gtid |         1 |        1890 | SET @@SESSION.GTID_NEXT= 'ANONYMOUS'                                                                                            |
| binlog.000002 | 1890 | Query          |         1 |        1968 | FLUSH TABLES                                                                                                                    |
+---------------+------+----------------+-----------+-------------+---------------------------------------------------------------------------------------------------------------------------------+
10 rows in set (0.00 sec)

# 从最早的binlog开始
mysql> show binlog events;
```

## 恢复binlog

```bash
# 全量导入
mysql> source  /var/lib/mysql/database-backup.sql

# 根据时间恢复
mysqlbinlog --start-datetime="2013-11-29 13:18:54" --stop-datetime="2013-11-29 13:21:53" --database=zyyshop binlog.000002 | mysql -uroot -p123456

# 根据位置恢复
mysqlbinlog  --start-position=293963814  --stop-position=346091760 --database=academy | mysql -uroot -pmysql
```
