---
title: 35-waitGroup.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 35-waitGroup.go
showInMenu: false

---

//  要等待多个goroutine完成,可以使用WaitGroup

package main

import (
	"fmt"
	"sync"
	"time"
)

//  在每个goroutine中运行的函数
//  注意必须通过指针将WaitGroup传递给函数
func worker(id int, wg *sync.WaitGroup) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second) //Sleep函数模拟执行任务
	fmt.Printf("Worker %d done\n", id)

	wg.Done() // 通知WaitGroup工作已经完成
}

func main() {
	var wg sync.WaitGroup // 此WaitGroup用于等待此处启动的所有goroutine完成

	for i := 1; i < 5; i++ {
		//  启动多个goroutine,并为每个goroutine增加WaitGroup计数器
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait() // 阻塞,知道WaitGroup计数器返回0,即所以goroutine通知他们已经完成
}
