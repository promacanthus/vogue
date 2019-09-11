// 对于简单的计数器状态可以通过 atomic 操作来控制
// 对于更复杂的状态，可以使用互斥锁（mutex）安全地访问多个goroutine中的数据

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var state = make(map[int]int) // 创建一个map表示状态（state）

	var mutex = &sync.Mutex{} // 互斥锁会同步对状态的访问

	// 用于追踪读写的操作数量
	var readOps uint64
	var writeOps uint64

	for r := 0; r < 100; r++ { // 启动100个goroutine,执行针对状态的重复读取
		go func() {
			total := 0
			for {
				key := rand.Intn(5)           // 对于每次读取，生成一个访问key
				mutex.Lock()                  // 锁定互斥锁，以确保对状态的独占访问
				total += state[key]           // 从map中读取所选key对应的值
				mutex.Unlock()                // 释放互斥锁
				atomic.AddUint64(&readOps, 1) // 增加readOps计数
				time.Sleep(time.Millisecond)  // 每个goroutine每毫秒执行一次操作
			}
		}()
	}

	for w := 0; w < 10; w++ { // 启动10个goroutine来模拟写入操作，处理模式与读取操作相同
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second) // 等待上述读写操作工作1秒钟

	// 获取1秒内的读取和写入的操作次数
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

	// 锁定状态并输出状态中的值
	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

//  运行结果表面，针对互斥同步状态在一秒内大概执行了9万次操作
