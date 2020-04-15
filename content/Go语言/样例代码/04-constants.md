---
title: 04-constants
date: 2020-01-10T19:50:56.185945+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 04-constants
showInMenu: false

---


```go
//  Golang支持的常量包括字符、字符串、布尔和数值

package main

import (
	"fmt"
	"math"
)

// const 用于声明常量，const语句可以出现在任何var语句可以出现的地方
const s string = "constant"

func main() {
	fmt.Println(s)

	// 常量表达式以任意精度执行算术
	const n = 500000000
	const d = 3e20 / n
	fmt.Println(d)

	// 数值常量在被设置之前是无类型的,(如 上面例子中的常量 d 通过显式转换)
	fmt.Println(int64(d))

	// 在上下文中使用数字时可以自动给它设置类型,(如 变量赋值或者函数调用)
	fmt.Println(math.Sin(n))
	// 此处的math.Sin()函数需要一个float64类型的数值
}

```