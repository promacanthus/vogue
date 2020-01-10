# 28-timeout

```go
// timeouts 对于需要连接到外部资源或者需要绑定执行时间的程序非常重要
//  使用通道和select在Go中实现timeout简单切优雅
package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string, 1) // 创建一个缓冲区为1的通道，在不读取的情况下，这是防止goroutine泄露的常见模式
	go func() {                // 另外的goroutine将在两秒后向通道c1中写入“result1”
		time.Sleep(2 * time.Second)
		c1 <- "result1"
	}()

	select { // select语句实现超时选择
	case res := <-c1: // 等待上面的goroutine返回结果
		fmt.Println(res)
	case <-time.After(1 * time.Second): // 等待1秒后将时间写入返回的通道中
		fmt.Println("timeout1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result2"
	}()

	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second): // 等待3秒后将时间写入返回的通道中
		fmt.Println("timeout2")
	}
}

```