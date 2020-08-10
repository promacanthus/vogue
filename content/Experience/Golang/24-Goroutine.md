---
title: "24 Goroutine退出、泄露"
date: 2020-08-07T12:07:13+08:00
draft: true
---

Go中，goroutine是否结束执行（退出）是由其自身决定，其他goroutine只能通过消息传递的方式通知其关闭，而并不能在外部强制结束一个正在执行的goroutine。

> 有一种特殊情况会导致正在运行的goroutine会因为其他goroutine的结束而终止，即main函数退出。

## 常见goroutine退出方式

### main函数结束

```golang
package main

import (
 "fmt"
 "time"
)

func main() {
 go func() {
  time.Sleep(time.Second)
  fmt.Println("goroutine exit")
 }()
 fmt.Println("main exit")
}

// output
main exit
```

如上所示，程序未等`goroutine`执行完毕，即随着`main`函数的退出而停止执行。

### context通知退出

```golang
package main

import (
 "context"
 "fmt"
 "time"
)

func main() {
 ctx, cancel := context.WithCancel(context.Background())
 go func(ctx context.Context) {
  num := 0
  for {
   select {
   case <-ctx.Done():
    fmt.Println("goroutine exit")
    return
   case <-time.After(time.Second):
    num++
    fmt.Printf("goroutine wait times: %d\n", num)
   }
  }
 }(ctx)

 time.Sleep(time.Second * 3)
 cancel()
 time.Sleep(time.Second)
 fmt.Println("main exit")
}

// output
goroutine wait times: 1
goroutine wait times: 2
goroutine exit
main exit
```

### panci异常退出

```golang
package main

import (
 "fmt"
 "os"
 "time"
)

func main() {
 go func() {
  defer func() {
   if err := recover(); err != nil {
    fmt.Printf("goroutine exit by panic: %v\n", err)
   }
  }()

  _, err := os.Open("notExistFile.txt")
  if err != nil {
   panic(err)
  }
  fmt.Println("goroutine exit naturally")
 }()

 time.Sleep(time.Second)
 fmt.Println("main exit")
}


// output
goroutine exit by panic: open notExistFile.txt: no such file or directory
main exit
```

上面自定义函数中defer函数使用了`recover`来捕获`panic`，当`panic`发生时可使`goroutine`拿回控制权，确保程序不会将`panic`传递到`goroutine`调用栈顶部后引起崩溃。

### 执行完毕退出

```golang
package main

import (
 "fmt"
 "time"
)

func main() {
 go func() {
  for i := 0; i < 10000; i++ {
   // TODO: do some thing
  }
  fmt.Println("goroutine exit")
 }()

 time.Sleep(time.Second)
 fmt.Println("main exit")
}

// output
goroutine exit
main exit
```

goroutine里的任务执行完毕，即结束。

## goroutine泄露

如果启动了一个goroutine，但并没有按照预期的一样退出，等到程序结束，此goroutine才结束，这种情况就是 goroutine 泄露。

> 当 goroutine 泄露发生时，该 goroutine 的**栈一直被占用而不能释放**，goroutine 里的函数**在堆上申请的空间也不能被垃圾回收器回收**。这样，在程序运行期间，内存占用持续升高，可用内存越来也少，最终将导致系统崩溃。

大多数情况下，引起goroutine泄露的原因有两类：

- channel阻塞
- goroutine陷入死循环

### channel阻塞

```golang
// 从channel中读取，但是没有向channel中写入
package main

import (
 "fmt"
 "runtime"
 "time"
)

func main() {
 go func() {
  c := make(chan int)
  go func() {
   <-c
  }()
  time.Sleep(time.Second * 2)
  fmt.Println("goroutine exit")
 }()
 c := time.Tick(time.Second)
 for range c {
  fmt.Printf("goroutine [nums]: %d\n", runtime.NumGoroutine())
 }
}

// output
goroutine [nums]: 3
goroutine exit
goroutine [nums]: 3
goroutine [nums]: 2
goroutine [nums]: 2
goroutine [nums]: 2
...

// 向已满的channel中写入，但是没有读取
package main

import (
 "flag"
 "fmt"
 "runtime"
 "time"
)

var size = flag.Int("c", 0, "define channel size")

func main() {
 flag.Parse()

 go func(size int) {
  c := make(chan int, size)
  go func() {
   <-c
  }()

  go func() {
   for i := 0; i < 10; i++ {
    c <- i
   }
  }()
  fmt.Println("goroutine exit")
 }(*size)

 c := time.Tick(time.Second)
 for range c {
  fmt.Printf("goroutine [nums]: %d\n", runtime.NumGoroutine())
 }
}

// output
go run main.go -c 2
goroutine exit
goroutine [nums]: 2
goroutine [nums]: 2
goroutine [nums]: 2
...

go run main.go -c 11
goroutine exit
goroutine [nums]: 1
goroutine [nums]: 1
goroutine [nums]: 1
...
```

### 死循环

当代码里循环的退出条件不可达时，会令该goroutine进入死循环中，进而导致资源一直无法释放，引起泄露。

**在实际项目中，往往死循环会发生在一些后台的常驻服务中**。

## goroutine泄露的预防和检测

### 预防

1. **在创建goroutine时，就应该知道goroutine啥时能结束**。

2. `channel`引起的`goroutine`泄露问题，主要是看在`channel`阻塞`goroutine`时，该`goroutine`的阻塞是正常的，还是可能会导致`goroutine`永远没有机会执行(极大可能会造成协程泄露)。

    > `channel`的实际使用中，常用的两种模型：**生产者-消费者模型**；**master-worker模型**。一般的解决方案是：当主线程结束时，告知`worker goroutine`，`worker goroutine`得到通知后，进行清理工作然后退出；为每个worker任务制定超时，当超时触发，返回给master超时信息，并结束该`worker goroutine`。

3. 实现循环语句时必须清晰地知道退出循环的条件，避免死循环。

### 检测

1. Go提供的pprof工具。
2. 利用`runtime.NumGoroutine`接口，实时查看程序中运行的`goroutine`数。
3. 开源三方`profiling`库，如：[gops](https://github.com/google/gops)或[goleak](https://github.com/uber-go/goleak)。


