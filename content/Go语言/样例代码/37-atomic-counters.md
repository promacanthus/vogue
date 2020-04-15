---
title: 37-atomic-counters
date: 2020-01-10T20:05:30.450207+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 37-atomic-counters
showInMenu: false

---


```go
//  Golang中管理状态的首要机制是通过通道进行通信
//  使用 sync/atomic 包中的原子计数器实现多个goroutine的访问

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var ops uint64        // 使用无符号整数表示计数器
	var wg sync.WaitGroup // WaitGroup将会等待所以goroutine完成任务

	for i := 0; i < 50; i++ { // 启动50个goroutine，每个都会将计数器增加1000
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				// 使用AddUint64函数以原子的方式增加计数器的值
				// 使用&语法将ops变量的内存地址传递给它
				atomic.AddUint64(&ops, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait() // 等待所以goroutine完成

	// 现在访问ops是安全的，因为已经没有其他的goroutine在写它
	fmt.Println("ops:", ops)
	// 使用atomic.LoadUint64等函数可以在更新数据时进行安全的读取
}

```