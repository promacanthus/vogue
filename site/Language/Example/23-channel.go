// 通道是连接并发goroutine的管道
// 可以通过管道在一个goroutine与另一个goroutine之间发送和接收值

package main

import "fmt"

func main() {
	message := make(chan string) // 使用make(chane val-type)语法来创建一个新的通道，通道根据传入的值来确定类型

	// 使用channel <- 语法向通道中传入一个值，此处在一个新的goroutine中向message通道中传入一个字符串ping
	go func() {
		message <- "ping"
	}()

	msg := <-message
	fmt.Println(msg)
}

// 程序运行后，字符串ping将会通过通道成功的从一个goroutine传递到另一个goroutine中

// 默认情况下，直到发送和接收都准备好，才会进行发送和接收操作
// 这种特性使得的我们可以在程序结束时等到ping消息，而不需要任何其他同步
