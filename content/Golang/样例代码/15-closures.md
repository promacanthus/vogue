---
title: 15-closures
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---


```go
//  Go支持匿名函数，可用于形成闭包
//  当需要定义一个没有名字的内部函数时，匿名函数很有用

package main

import "fmt"

func intSeq() func() int {
	//  函数intSeq返回一个定义在intSeq函数体中的匿名函数
	// 返回的函数关闭变量i以形成闭包
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	nextInt := intSeq()
	// 调用intSeq函数，将结果（函数）分配给nextInt
	// 此函数值捕获其自身的i值，每次调用nextInt都会更新该值

	// 多次调用nextInt来查看闭包的效果
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	//  创建一个新的测试函数，来确定该状态对于特定函数是唯一的
	newInts := intSeq()
	fmt.Println(newInts())
}

```