---
title: "19 开启pprof"
date: 2020-07-22T09:10:58+08:00
draft: true
---

## 导入pprof

```go
import "net/http/pprof"
```

`pprof`包通过程序的HTTP运行时提供服务，使用pprof可视化工具性能分析数据并提供期望的输出格式。

只要导入这个包来注册它的HTTP处理程序就可以了，处理路径都是以`/debug/pprof/`开头的。要使用`pprof`主要在程序中导入`import _ "net/http/pprof"`

如果应用程序不启动http服务器，那么需要导入`net/http`和`log`包，如下所示：

```go
go func() {
 log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

如果不使用`DefaultServeMux`，那么在路由器上注册使用的处理程序即可。

## 调用pprof

查看所有可用的概要信息，在浏览器其中打开`http://localhost:6060/debug/pprof/`，这里的端口根据实际运行情况确定。

|案例|调用|
---|---
查看堆信息|`go tool pprof http://localhost:6060/debug/pprof/heap`
查看30秒内CPU信息|`go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30`
查看`goroutine`的阻塞信息（在程序中调用`runtime.SetBlockProfileRate`）|`go tool pprof http://localhost:6060/debug/pprof/block`
手收集5秒内执行的跟踪信息|`wget http://localhost:6060/debug/pprof/trace?seconds=5`
查看锁的持有者（在程序中调用`runtime.SetMutexProfileFraction`）|`go tool pprof http://localhost:6060/debug/pprof/mutex`

## 更多

了解更多，可查看golang官方博客的文章，点击[这里](https://blog.golang.org/pprof)。
