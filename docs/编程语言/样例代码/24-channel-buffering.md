# 24-channel-buffering

```go
// 默认情况下，通道是不带缓冲区的,这意味着只有接收端和发送端同时准备好才能发送数据
// 带缓冲区的通道在没有相应的接收端时,可以接收有限数量的值

package main

import "fmt"

func main() {
	message := make(chan string, 2) // 创建字符串通道，最多缓冲两个值

	//  message通道是带有缓存的，因此可以将值发送到通道中，而不需要相应的并发接收
	message <- "buffered"
	message <- "channel"

	// 正常接收到通道中的两个值
	fmt.Println(<-message)
	fmt.Println(<-message)
}

```