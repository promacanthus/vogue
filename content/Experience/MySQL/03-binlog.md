---
title: 03 binlog
date: 2020-04-14T10:09:16.242627+08:00
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

## 全量导入优化

加快source的一些MySQL参数：

1. log_bin=OFF
2. innodb_flush_log_at_trx_commit=0
3. sync_binlog
4. max_allowed_packet=500M

> 注意，全局变量和会话级别变量的区别，使用global参数，即`set global max_allowed_packet=500M`。

### innodb_flush_log_at_trx_commit

提交事务的时候将 `redo` 日志写入磁盘中，(所谓的 redo 日志，就是记录下来你对数据做了什么修改)。如果要提交一个事务，此时就会根据一定的策略把 redo 日志从 redo log buffer 里刷入到磁盘文件里去。此时这个策略是通过 `innodb_flush_log_at_trx_commit` 来配置的，它的配置选项：

- 值为0 : 提交事务的时候，不立即把 redo log buffer 刷入磁盘文件，而是依靠 InnoDB 的主线程每秒执行一次刷新到磁盘。
- 值为1 : 提交事务的时候，就必须把 redo log buffer 刷入磁盘文件，只要事务提交成功，redo log 必然在磁盘。注意，因为操作系统的“延迟写”特性，此时的刷入只是写到了操作系统的缓冲区中，因此执行**同步**操作才能保证一定持久化到了硬盘中。
- 值为2 : 提交事务的时候，把 redo 日志写入磁盘文件对应的 os cache 缓存里，而不是直接进入磁盘文件，可能 1 秒后才会把 os cache 里的数据写入到磁盘文件里去。

### sync_binlog

该参数控制着二进制日志写入磁盘的过程，他的配置选项：

- 0：默认值。事务提交后，将二进制日志从缓冲写入磁盘，但是不进行刷新操作（`fsync()`），此时只是写入了操作系统缓冲，若操作系统宕机则会丢失部分二进制日志。
- 1：事务提交后，将二进制文件写入磁盘并立即执行刷新操作，相当于是同步写入磁盘，不经过操作系统的缓存。
- N：每写N次操作系统缓冲就执行一次刷新操作。

### max_allowed_packet

`max_allowed_packet` 参数（单位字节 ）限制Server接受的数据包大小。
