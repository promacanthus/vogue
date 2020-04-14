---
title: 72-signals.md
date: 2020-01-10T20:17:32.778816+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 72-signals.md
showInMenu: false

---

# 72-signals

```go
// 有时希望Go程序能够智能地处理Unix信号。
// 例如，希望服务器在收到SIGTERM时正常关闭，或者命令行工具在收到SIGINT时停止处理输入

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 通过在通道上发送os.Signal值来发送信号通知
	// 创建一个通道来接收这些通知（还会在程序退出时发出通知）
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// signal.Notify注册​​给定通道以接收指定信号的通知
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 该goroutine执行信号的阻塞接收
	// 当它获得一个信号，将会打印出来该信号，然后发出完成通知
	go func() {
		sig := <-sigs
		fmt.Println("")
		fmt.Println(sig)
		done <- true
	}()

	// 程序将在此处等待，直到获得预期的信号
	// 如上面的goroutine所示，发送完成后的值，然后退出
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

```