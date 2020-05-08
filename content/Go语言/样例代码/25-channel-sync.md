---
title: 25-channel-sync
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---


```go
// 可以使用通道来跨goroutine同步执操作

package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) { // 在goroutine中执行的函数，done通道用于通知另一个goroutine此函数已经执行完成
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true // 发送一个值来通知函数执行完成
}

func main() {
	done := make(chan bool, 1)
	go worker(done) // 启动一个goroutine，并将通知通道传递给它

	<-done // 阻塞，直到从通道中接收到worker发出的通知
	// 如果删除上面<-done这一行，程序可能在worker启动之前就已经退出了
}

```