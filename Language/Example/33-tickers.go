//  Timer（计时器）用于在未来的执行一次某件事
//  Ticker（断续器）用于定期重复做某件事

package main

import (
	"fmt"
	"time"
)

func main() {
	// 类似Timer的机制，使用一个通道来传递值
	ticker := time.NewTicker(500 * time.Millisecond) // 使用通道内置的range每个500毫秒遍历一下进入通道中的值
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(3600 * time.Millisecond)
	ticker.Stop() //  ticker可以像timer那样被停止，一旦ticker被停止就不会在接收任何值传入通道中
	done <- true
	fmt.Println("Ticker stopped")
}
