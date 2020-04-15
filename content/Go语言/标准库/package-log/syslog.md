---
title: package syslog
date: 2020-04-14T10:09:14.274627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 标准库
summary: package syslog
showInMenu: false

---

```go
import "log/syslog"
```

syslog包提供了一个简单的系统日志服务接口，可以通过使用`UNIX domain sockets`，`UDP`或`TCP`将消息发送到syslog守护程序。

只需调用一次`Dial`，因为在写入失败时，syslog客户端将尝试重新连接到服务器并再次写入。

syslog包已归档，不再接受新功能。一些外部包提供更多功能。看到：https://godoc.org/?q=syslog

## 函数NewLogger

```go
func NewLogger(p Priority, logFlag int) (*log.Logger, error)
```

`NewLogger()`创建一个`log.Logger`，其输出将写入具有指定优先级的系统日志服务，即syslog工具和severity（优先级别）的组合。

- logFlag是传递给`log.New()`以创建`Logger`的标志集(flag set)。

## 类型Priority

```go
type Priority int
```

优先级是syslog工具和serverity（优先级别）的组合。例如，`LOG_ALERT`|`LOG_FTP`从FTP工具发送警报serverity（优先级别）消息。默认serverity（优先级别）为`LOG_EMERG`；默认工具是`LOG_KERN`。

```go
const (

    // From /usr/include/sys/syslog.h.
    // These are the same on Linux, BSD, and OS X.
    LOG_EMERG Priority = iota
    LOG_ALERT
    LOG_CRIT
    LOG_ERR
    LOG_WARNING
    LOG_NOTICE
    LOG_INFO
    LOG_DEBUG
)
```

```go
const (

    // From /usr/include/sys/syslog.h.
    // These are the same up to LOG_FTP on Linux, BSD, and OS X.
    LOG_KERN Priority = iota << 3
    LOG_USER
    LOG_MAIL
    LOG_DAEMON
    LOG_AUTH
    LOG_SYSLOG
    LOG_LPR
    LOG_NEWS
    LOG_UUCP
    LOG_CRON
    LOG_AUTHPRIV
    LOG_FTP

    LOG_LOCAL0
    LOG_LOCAL1
    LOG_LOCAL2
    LOG_LOCAL3
    LOG_LOCAL4
    LOG_LOCAL5
    LOG_LOCAL6
    LOG_LOCAL7
)
```

## 函数Dail

```go
func Dial(network, raddr string, priority Priority, tag string) (*Writer, error)
```

Dial通过连接到指定网络的`raddr`地址来建立与日志守护程序的连接。每次写入返回的writer都会发送一条日志消息，其中包含工具和serverity（优先级别）和标记。

- 如果tag为空，则使用`os.Args[0]`。
- 如果网络为空，Dial将连接到本地`syslog`服务器。

否则，请参阅`net.Dial`的文档以获取`network`和`raddr`的有效值。

```go
sysLog, err := syslog.Dial("tcp", "localhost:1234",
    syslog.LOG_WARNING|syslog.LOG_DAEMON, "demotag")
if err != nil {
    log.Fatal(err)
}
fmt.Fprintf(sysLog, "This is a daemon warning with demotag.")
sysLog.Emerg("And this is a daemon emergency with demotag.")
```

## 函数New

```go
func New(priority Priority, tag string) (*Writer, error)
```

`New()`建立与系统日志守护程序的新连接。每次写入返回的writer都会发送一条日志消息，其中包含给定的优先级（syslog工具和serverity的组合）和前缀标记。如果tag为空，则使用`os.Args[0]`。

## 类型Writer

```go
type Writer struct {
    // 包含已过滤或未导出的字段
}
```

`Writer`是一个`syslog`服务的连接。

### 方法 (*Writer) Alert

```go
func (w *Writer) Alert(m string) error
```

Alert()记录serverity（优先级别）为LOG_ALERT的消息，忽略传递给New的severity（优先级别）。

### 方法 (*Writer) Close

```go
func (w *Writer) Close() error
```

Close关闭与syslog守护程序的连接。

### 方法 (*Writer) Crit

```go
func (w *Writer) Crit(m string) errors
```

Crit记录serverity（优先级别）为`LOG_CRIT`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Debug

```go
func (w *Writer) Debug(m string) error
```

Debug会记录serverity（优先级别）为`LOG_DEBUG`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Emerg

```go
func (w *Writer) Emerg(m string) error
```

Emerg记录serverity（优先级别）为`LOG_EMERG`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Err

```go
func (w *Writer) Err(m string) error
```

Err记录serverity（优先级别）为`LOG_ERR`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Info

```go
func (w *Writer) Info(m string) error
```

Info会记录serverity（优先级别）为`LOG_INFO`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Notice

```go
func (w *Writer) Notice(m string) error
```

注意记录serverity（优先级别）为`LOG_NOTICE`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Warning

```go
func (w *Writer) Warning(m string) error
```

Warning会记录serverity（优先级别）为`LOG_WARNING`的消息，忽略传递给New的serverity。

### 方法 (*Writer) Write

```go
func (w *Writer) Write(b []byte) (int, error)
```

Write将日志消息发送到syslog守护程序。
