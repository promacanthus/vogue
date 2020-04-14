---
title: 09-slice.md
date: 2020-01-10T19:54:36.101497+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 09-slice.md
showInMenu: false

---

# 09-slice

```go
// 切片是Go中的关键数据类型，为序列提供了比数组更强大的接口。

package main

import "fmt"

func main() {
	// 与数组不同,切片的类型由它包含的元素(而不是元素数量)决定,使用内置的make函数创建长度非零的空切片
	s := make([]string, 3) // 在这里创建一个长度为3的空的字符串切片,其中默认零值为空字符串
	fmt.Println("emp:", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	// 可以像数组那样设置和获取切片的值
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])

	fmt.Println("len:", len(s)) // len函数返回切片的长度

	// 除了上面的基本操作,切片还有更多其他操作,使得它比数组更丰富
	s = append(s, "d")      // 内置函数append返回包含一个或多个新值的切片
	s = append(s, "e", "f") // 需要接收来自append函数的返回值，因为可能会得到一个新的切片
	fmt.Println("apd:", s)

	c := make([]string, len(s)) // 创建一个与s长度相同的空的切片c用于存储s的拷贝
	copy(c, s)                  // 切片也可以被拷贝
	fmt.Println("cpy:", c)

	// 切片支持slice[low:high]语法的切片操作
	l := s[2:5] // 得到包含元素s[2],s[3]和s[4]的切片
	fmt.Println("sl1:", l)

	l = s[2:] //从s[2]开始到结尾
	fmt.Println("sl2:", l)

	l = s[:5] //从头开始到s[5]结束
	fmt.Println("sl3:", l)

	t := []string{"g", "h", "i"} // 可以在一行中声明并初始化一个切片变量
	fmt.Println("dcl:", t)

	twoD := make([][]int, 3) // 切片可以被组合成多维数据结构
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen) // 内部切片的长度是可变的，这与多维数组不同
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d", twoD)
}

// 注意，虽然切片的类型与数组不同，但是在fmt.Println中的输出形式是相似的

// Go语言团队在切片设计与实现中的更对细节：https://blog.golang.org/go-slices-usage-and-internals

```