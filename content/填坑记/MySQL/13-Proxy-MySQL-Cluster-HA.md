---
title: "13 Proxy MySQL Cluster HA"
date: 2020-05-27T14:44:19+08:00
draft: true
---

## 配置文件

`proxysql.cnf`

```cnf
datadir="/var/lib/proxysql"

admin_variables=
{
    admin_credentials="admin:admin;radmin:radmin"
    mysql_ifaces="0.0.0.0:6032"
}
```

## 启动容器

```bash
docker run --name proxysql -d -v /path/to/proxysql.cnf:/etc/proxysql.cnf proxysql/proxysql:2.0.12
```

使用下面的配置文件来运行容器，其中包括管理员账户和管理端口。

## MySQL配置

在数据库中创建用于代理监控的用户。

```sql
CREATE USER 'monitor'@'%' IDENTIFIED BY 'monitor';
GRANT SELECT on sys.* to 'monitor'@'%';
FLUSH PRIVILEGES;
```

Proxysql具有组复制功能，只要将用于运行状况监视的用户添加到一个MySQL节点上，便会在所有三个节点上对其进行完全配置。

## Proxysql配置

### 登录管理端

- admin/admin：本地用户
- radmin/radmin：远程用户
- 6032：proxysql管理端口

```bash
# 本地
mysql -uradmin -pradmin -h127.0.0.1 -P6032 --prompt='ProxySQLAdmin> '

# 远程
mysql -uradmin -pradmin -h172.17.0.2 -P6032 --prompt='ProxySQLAdmin> '

# --prompt是一个可选标志，用于更改默认提示，通常为mysql>。
# 在这里，将其更改为ProxySQLAdmin>，以明确表明我们已连接到ProxySQL管理界面。
# 这将有助于避免以后在连接复制数据库服务器上的MySQL接口时造成混乱。
```

### 修改密码

通过更新（`UPDATE`）数据库中的全局变量`admin-admin_credentials`来更改管理帐户密码，多个帐号用`;`分隔。

```sql
UPDATE global_variables SET variable_value='admin:password' WHERE variable_name='admin-admin_credentials';
```

由于ProxySQL的配置系统的工作原理，此更改不会立即生效，所以配置都是这样，要注意。

- `memory`：从命令行界面进行修改时会更改。(**会话级别**)
- `runtime`：ProxySQL将其用作有效配置。（**运行时级别**）
- `disk`：用于使配置在重新启动后保持不变。（**持久化**）

在命令行修改后的配置存储在`memory`中。为了使更改生效，必须将`memory`设置复制到`runtime`，然后将它们保存到`disk`以使其持久。

```sql
LOAD ADMIN VARIABLES TO RUNTIME;
SAVE ADMIN VARIABLES TO DISK;
```

- `ADMIN`变量处理与管理命令行界面有关的变量
- `MYSQL`变量处理其配置的其他部分

### 配置监控

配置已经在MySQL示例中创建好的监控账户，返回ProxySQL管理界面，将`mysql-monitor_username`变量更新为新帐户的用户名，并将配置持久化，注意此时使用的是`MYSQL`变量。

```sql
UPDATE global_variables SET variable_value='monitor' WHERE variable_name='mysql-monitor_username';
LOAD MYSQL VARIABLES TO RUNTIME;
SAVE MYSQL VARIABLES TO DISK;
```

### 添加节点

在Proxysql中配置主机组，每个主机组都由一个正数标识，例如1或2。使用ProxySQL查询路由时，主机组可以将不同的SQL查询路由到不同的主机集群。

在静态副本设置中，主机组可以任意设置。但是，ProxySQL的组副本机制支持自动将组中的所有节点划分为四个逻辑状态：

- `writer`：MySQL节点可以接受更改数据的查询。ProxySQL确保将所有主节点保存在该组中的最大定义数量之内。
- `backup writer`：MySQL节点以接受更改数据的查询。但是，这些节点并没有被设置为`writer`；超过`writer`定义数量的主节点将保留在该组中，如果`writer`中有节点发生故障，则将`backup writer`中的一个节点升级到`writer`。
- `reader`：MySQL节点不能接受更改数据的查询，将其用作只读节点。ProxySQL仅在此处放置从属节点。
- `offline`：放置由于缺乏连通性或流量缓慢等问题而出现异常的节点。

这四个状态中的每一个都有对应的主机组，但是不会自动分配数字组标识符。

我们需要告诉ProxySQL每个状态应使用哪些标识符。比如：

- 将1用于`offline`主机组
- 将2用于`writer`主机组
- 将3用于`reader`主机组
- 将4用于`backup writer`主机组

要设置这些标识符，就在`mysql_group_replication_hostgroups`表中使用这些变量和值创建一个新行。

```sql
INSERT INTO mysql_group_replication_hostgroups (writer_hostgroup, backup_writer_hostgroup, reader_hostgroup, offline_hostgroup, active, max_writers, writer_is_also_reader, max_transactions_behind) VALUES (2, 4, 3, 1, 1, 3, 1, 100);
```

字段说明：

- 如果将`active`设置为1，则ProxySQL可以监视这些主机组。
- `max_writers`定义`writer`节点的最大数量。在这里使用3是因为在多主节点设置中，所有节点都可以被视为相等，因此在这里使用3（节点总数）。
- 将`writer_is_also_reader`设置为1会指示ProxySQL将`writer`也认为是`reader`。
- `max_transactions_behind`设置定义节点为`offline`状态的最大延迟事务数。

这样ProxySQL就知道如何在主机组之间分配节点，然后将MySQL服务器添加到池中。为此，需要将每个服务器的IP地址和初始主机组插入`mysql_servers`表，该表包含ProxySQL可以与之交互的服务器列表。

添加MySQL服务器，并确保替换以下命令中的示例IP地址。

```sql
INSERT INTO mysql_servers(hostgroup_id, hostname, port) VALUES (2, '203.0.113.1', 3306);

LOAD MYSQL SERVERS TO RUNTIME;
SAVE MYSQL SERVERS TO DISK;
```

ProxySQL现在应按指定在主机组之间分布我们的节点。我们通过对`runtime_mysql_servers`表执行SELECT查询来进行检查，该表显示了ProxySQL正在使用的服务器的当前状态。

```sql
SELECT hostgroup_id, hostname, status FROM runtime_mysql_servers;
```

### 配置数据库用户凭据

ProxySQL充当负载均衡；用户连接到ProxySQL，然后ProxySQL将该连接依次传递到所选的MySQL节点。为了连接到单个节点，ProxySQL会重用其访问的凭据。

为了允许访问位于复制节点上的数据库，需要创建一个具有与ProxySQL相同的凭据的用户帐户，并向该用户授予必要的特权。

创建一个名为`PlaygroundUser`的新用户，密码是`playgroundpassword`。

```sql
CREATE USER 'playgrounduser'@'%' IDENTIFIED BY 'playgroundpassword';
GRANT ALL PRIVILEGES on playground.* to 'playgrounduser'@'%';
FLUSH PRIVILEGES;
EXIT;
```

可以通过直接在节点上尝试使用新配置的凭据访问数据库来验证是否已正确创建用户。

```bash
mysql -u playgrounduser -p
```

```sql
SHOW TABLES FROM playground;

+----------------------+
| Tables_in_playground |
+----------------------+
| equipment            |
+----------------------+
1 row in set (0.00 sec)
```

### 创建Proxysql用户

最后的配置步骤是允许`playgrounduser`用户与ProxySQL建立连接，并将这些连接传递给节点。

为此，我们需要在`mysql_users`表中设置配置变量，该表包含用户凭据信息。将用户名，密码和默认主机组添加到配置数据库（对于`writer`主机组，该名称为2）。

```sql
INSERT INTO mysql_users(username, password, default_hostgroup) VALUES ('playgrounduser', 'playgroundpassword', 2);
LOAD MYSQL USERS TO RUNTIME;
SAVE MYSQL USERS TO DISK;
```

ProxySQL在端口6033上监听传入的客户端连接，因此请尝试使用`Playgrounduser`和端口6033连接到真实数据库（而非管理界面）。在这里，我们将提示设置为`ProxySQLClient>`，以便将其与管理界面提示区分开。

```bash
mysql -u playgrounduser -p -h 127.0.0.1 -P 6033 --prompt='ProxySQLClient> '
```

让我们执行一条简单的语句来验证ProxySQL是否将连接到其中一个节点。此命令在数据库中查询正在运行的服务器的主机名，并返回服务器的主机名作为唯一输出。

```sql
SELECT @@hostname;

+------------+
| @@hostname |
+------------+
| member1    |
+------------+
1 row in set (0.00 sec)
```

根据我们的配置，此查询应由ProxySQL定向到分配给编写者主机组的三个节点之一。其中`member1`是一个MySQL节点的主机名。

### 验证ProxySQL配置

ProxySQL和MySQL节点之间的连接正常，测试确保数据库权限允许从ProxySQL读取和写入语句，并确保当其中的某些节点宕机仍能成功执行这些语句。

```sql
-- 读
SELECT * FROM playground.equipment;

-- 写
INSERT INTO playground.equipment (type, quant, color) VALUES ("drill", 5, "red");

-- 再读
SELECT * FROM playground.equipment;
```

完成。
