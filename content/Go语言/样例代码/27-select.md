---
title: 27-select.md
date: 2020-01-10T20:02:38.461838+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 27-select.md
showInMenu: false

---

# 27-select

```go
// Go中的select等待多个通道的操作
//  将goroutine\通道和select结合是Go中很强大的功能

package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()

	c1 := make(chan string)
	c2 := make(chan string)

	// 每个通道将在一段时间后收到一个值，以此来模拟在并发的goroutine中的阻塞RPC操作
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select { // 使用select来同事等待这两个值，并在它们到达时打印每个值
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}

	endTime := time.Since(startTime)
	fmt.Println(endTime)
}

//  由于1秒和2秒的睡眠是同时进行的，所以程序执行的总时间是2秒出头

```