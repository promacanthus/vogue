---
title: 19-methods
date: 2020-01-10T20:00:03.473592+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 19-methods
showInMenu: false

---


```go
// Go支持在结构体类型上定义方法

package main

import "fmt"

type rect struct {
	width, height int
}

func (r *rect) area() int { // area方法有一个*rect类型的接收器
	return r.width * r.height
}

func (r rect) perim() int { // 既可以给指针接收器类型定义方法也可以给值接收器类型定义方法
	return 2*r.width + 2*r.height
}

func main() {
	// Go自动处理方法调用中值和指针之间的转换
	r := rect{width: 10, height: 5}

	fmt.Println("area:", r.area()) // 使用指针接收器类型来避免在调用方法时出现值拷贝的情况或者允许方法改变接收的结构
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area:", rp.area())
	fmt.Println("perim:", rp.perim())
}

```