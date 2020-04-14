---
title: 12-functions.md
date: 2020-01-10T19:56:34.165448+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 12-functions.md
showInMenu: false

---

# 12-functions

```go
// 函数是Go的核心

package main

import "fmt"

func plus(a, b int) int { // plus函数输入两个int并以int返回他们的和
	return a + b // Go需要显式返回，而不会自动返回最后一个表达式的值
}

func plusplus(a, b, c int) int { // 当有多个连续的相同类型的参数时，可以省略相同参数类型的参数名称，直到声明该类型的最后一个参数
	return a + b + c
}

func main() {
	res := plus(1, 2) // 使用name(args)来调用函数
	fmt.Println("1+2 =", res)

	res = plusplus(1, 2, 3)
	fmt.Println("1+2+3 =", res)
}

```