---
title: 30-closing-channel.md
date: 2020-01-10T20:03:23.761926+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 30-closing-channel.md
showInMenu: false

---

# 30-closing-channel

```go
//  关闭通道表示不再会有值传递给它,这对于与通道的接收器进行通信完成非常有用。

package main

import "fmt"

func main() {
	jobs := make(chan int, 5) // 通道jobs在main goroutine与worker goroutine之间通讯，当没有jobs之后将关闭通道
	done := make(chan bool)

	// 这里是worker goroutine
	go func() {
		for {
			// 不断的从通道中获取值
			j, more := <-jobs // 在这种特殊的二值形式中，如果通道已经关闭且其中的值都已经被接收，那么more将会接收到false值
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true // 完成所有操作后通过done通道来发出通知
				return
			}
		}
	}()

	for j := 0; j < 3; j++ { // 发送3个job到jobs通道中，然后把通道关闭
		jobs <- j
		fmt.Println("send job", j)
	}

	close(jobs)
	fmt.Println("send all jobs")

	<-done
}

```