//  将通道作为函数的参数时,可以指定通道为单向通道
//  这种特性增加了程序的类型安全

package main

import "fmt"

func ping(pings chan<- string, msg string) { // ping函数只接收输入通道，尝试在输入通道上进行接收数据将发生编译时错误
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) { // pong函数接收一个输出通道和一个输入通道
	msg := <-pings
	pongs <- msg
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
