---
title: 17-pointer.md
date: 2020-01-10T19:59:23.077545+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 17-pointer.md
showInMenu: false

---

# 17-pointer

```go
//  Go支持指针，允许传递对程序中的值和记录的引用

package main

import "fmt"

func zeroval(ival int) { // zeroval函数的参数是int，所以参数将按值传递
	ival = 0 // zeroval将会获得ival值的拷贝，而不是调用函数时传入的ival
}

func zeroptr(iptr *int) { // zeroptr函数的参数是*int ,这表示它需要一个int指针
	*iptr = 0 // 函数体中的*iptr解除在那个地址上的从内存地址到当前值的引用，为解除引用的指针赋值会改变引用地址的值
}

func main() {
	i := 1
	fmt.Println("initial:", i)

	zeroval(i)
	fmt.Println("zeroval:", i)

	zeroptr(&i) // &i语法给出i的内存地址，即指向i的指针
	fmt.Println("zeroptr:", i)

	fmt.Println("pointer:", &i) // 指针也可以被打印
}

// zeroval不会更改main()中的i，但是zeroptr会更改，因为它具有对该变量的内存地址的引用。

```