---
title: 01-hello-world.go
date: 2019-11-25T11:15:47.530182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 01-hello-world.go
showInMenu: false

---

package main

import "fmt"

func main() {
	fmt.Println("hello world")
}

// 使用go run hello-world.go来执行本程序

// 使用go build 将本程序编译成二进制文件，然后使用./hello-world 执行编译后的二进制文件
