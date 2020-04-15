---
title: 22-goroutines
date: 2020-01-10T20:00:59.67767+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 22-goroutines
showInMenu: false

---


```go
// Goroutine是一个轻量级的执行线程

package main

import "fmt"

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {
	f("direct") // 调用f()函数，并以同步的方式运行

	go f("goroutine") // 在goroutine中调用f()函数，这个新的goroutine将与发起调用的goroutine同时执行

	go func(msg string) { // 以匿名函数调用启动goroutine
		fmt.Println(msg)
	}("going")

	// 上述两个函数调用在不同的goroutine中异步运行，执行到此为止
	fmt.Scanln() // Scanln()函数需要在程序退出钱按一个按键
	fmt.Println("done")
}

// 程序运行后，首先看到阻塞调用的输出，然后是两个goroutine的交错输出
// 这种交错输出反应了Go运行时启动的并发运行的goroutine

```