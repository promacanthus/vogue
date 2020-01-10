# 34-worker-pools

```go
//  使用goroutine和通道实现worker pool

package main

import (
	"fmt"
	"time"
)

//  运行多个并发的实例
//  从jobs通道中获取工作,并将结果写入result通道中
// 每个工作都会睡眠一秒钟来模拟执行任务
func worker(id int, jobs <-chan int, result chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		result <- j * 2
	}
}

func main() {

	//  创建两个通道用于发送jobs和接收结果
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// 启动3个worker,最后会被阻塞,因为jobs通道为空
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 向jobs通道中传入5个值,然后将通道关闭
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// 收集所有工作的结果,这也能确保worker的goroutine已经完成
	for a := 1; a <= 5; a++ {
		<-results
	}
}

//  等待多个goroutine的另一种方法是使用WaitGroup

```