---
title: 29-nonblocking-channel
date: 2020-01-10T20:03:07.929894+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 29-nonblocking-channel
showInMenu: false

---


```go
//  通道上的发送和接收操作都是阻塞的
//  使用带有default子句的select来实现非阻塞发送,接收甚至是非阻塞多路select

package main

import "fmt"

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	//  非阻塞接收，如果message中有值，那么select将会使用<-message 这个case，否则执行default
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	// 非阻塞发送，这里msg不能发送到message通道中，因为这是一个非缓冲区通道，同时通道没有接收器，所以执行default
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	//  在default之上实现多路非阻塞select
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

```