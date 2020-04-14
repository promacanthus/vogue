---
title: 36-rate-limit.md
date: 2020-01-10T20:05:11.246162+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 36-rate-limit.md
showInMenu: false

---

# 36-rate-limit

```go
// 速率限制是控制资源利用和维持服务质量的重要机制。
//  Golang优雅的支持 goroutine、通道和断续器（ticker）的速率限制

package main

import (
	"fmt"
	"time"
)

func main() {
	//  基本的速率限制
	requests := make(chan int, 5) //  通过同名的通道来提供对传入请求的处理限制
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond) // limiter通道每200毫秒接收一个值，这是速率限制机制中的调节器

	//	在处理每个请求之前，通过limiter通道来阻塞请求，并限制每200毫秒一个请求
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	fmt.Println()

	// 在速率限制方案中允许短时间的突发请求，同时保留总体的速率限制
	// 可以通过缓冲 limiter 通道来达到这个目的
	// burstyLimiter 通道允许最多3个突发事件
	burstyLimiter := make(chan time.Time, 3)

	// 填充通道以表示允许的突发请求
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		// 每200毫秒向burstyLimiter添加一个值，它的最大限制是3
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// 模拟传入了5个请求，其中前3个将受益于burstyLimiter的突发处理能力
	burstyRequests := make(chan int, 5)
	for i := 1; i < =5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

//  在第二批的请求中，有突发事件的处理机制，所以前3个会立即被处理，然后在每200毫米处理后2个请求

```