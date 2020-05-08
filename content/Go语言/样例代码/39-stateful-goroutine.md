---
title: 39-stateful-goroutine
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---
```go
//  在前面的例子中，使用显示定义的互斥锁来跨多个goroutine同步对共享状态的访问。

// 另一种方式是使用goroutine和通道内置的同步功能来实现相同的效果，
// 这种基于通道的方式和Go共享内存的想法一致，通过通信使每块数据只被一个goroutine拥有。

package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// 在这个例子中，状态只会被一个goroutine拥有，
// 这将保证数据永远不会因为并发访问而损坏。
// 为了读取或者写入状态，其他的goroutine需要发送消息给状态的持有goroutine并收到相应的响应

// readOp和WriteOp结构体封装了这些请求和持有状态的goroutine的一个响应方法
type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	// 记录总共执行的读取和写入次数
	var readOps uint64
	var writeOps uint64

	// 读取和写入通道将被其他的goroutine使用来分别发出读取或写入请求
	reads := make(chan readOp)
	writes := make(chan writeOp)

	// 这是持有状态的goroutine
	go func() {
		// goroutine私有的状态map
		var state = make(map[int]int)
		for {
			// 这个goroutine不断的从读取和写入通道中进行选择，并响应获取到的结果
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		// 启动100个goroutine通过reads通道向持有状态的goroutine发送读取请求
		go func() {
			// 每次读取请求都会构造一个readOp结构，并通过reads通道进行发送，然后通过resp通道获取结果
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for w := 0; w < 10; w++ {
		// 启动10个goroutine进行写入请求
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// 等待上述goroutine执行一秒钟
	time.Sleep(time.Second)

	// 捕获并返回操作的总次数
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}

// 程序运行结果表明，给予goroutine的状态管理完成了10万次操作

// 这种情况下，基于goroutine的方法比基于互斥锁的方法更复杂一些，
// 在某些情况下它可能很有用，例如涉及多个管道或管理多个互斥锁时容易出错。

```