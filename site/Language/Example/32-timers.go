//  我们通常需要在未来的某些时间点执行Go代码或者在某个时间段内重复
//  Golang内置了timer（计时器）和ticker（断续器），使得上述目的很容易实现

package main

import (
	"fmt"
	"time"
)

func main() {
	// Timer代表未来的某个时间
	timer1 := time.NewTimer(2 * time.Second) //告诉timer等待的时长，然后提供一个时间到达后进行通知的通道
	<-timer1.C                               // timer的通道C被<-timer1.C阻塞，直到通道内被传入一个表示时间到期的值
	fmt.Println("Timer 1 expired")

	//  如果只是想要等待，可以使用time.Sleep
	//  使用Timer的好处是，可以在时间过期之前取消它
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()

	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}

// 第一个Timer将在程序运行2秒后过期，第二个Timer将会被停止而没有机会过期
