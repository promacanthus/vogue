---
title: 64-cmd-line-arg
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---


```go
// 命令行参数是参数化执行程序的常用方法。
// 例如： go run hello.go 使用run和hello.go 作为go的参数

package main

import (
	"fmt"
	"os"
)

func main() {
	argsWithProg := os.Args        // os.Args提供对原始命令行参数的访问
	argsWithoutProg := os.Args[1:] // 请注意，此切片中的第一个值是程序的路径，os.Args [1:]保存程序的参数

	arg := os.Args[3] // 使用正常的索引获取单个参数

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}

```