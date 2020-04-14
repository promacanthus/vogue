---
title: 02-values.md
date: 2020-01-10T19:50:04.610136+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 02-values.md
showInMenu: false

---

# 02-values

```go
//  Golang有多种值类型，包括字符串、整数、浮点数、布尔值等

package main

import "fmt"

func main() {

	// 使用加号（+）可以将string合并到一起
	fmt.Println("go" + "lang")

	fmt.Println("1+1=", 1+1)
	fmt.Println("7.0/3.0=", 7.0/3.0)

	// 布尔运算符是短路运算符
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}

```