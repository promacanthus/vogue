---
title: 44-collection-func
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---


```go
// 我们经常需要我们的程序对数据集合执行操作，
// 例如 	1.选择满足给定谓词的所有项目
//         		2.使用自定义函数将所有项目映射到新集合。

// 在某些编程语言中,使用generic数据结构和算法
// Go不支持generics,在Go中,如果程序或者数据类型特别需要,通常会提供集合函数

// 以下是一些字符串切片的集合函数示例,
// 可以使用以下示例构建自己的函数
// 请注意,在某些情况下,最简单的方法是直接内联集合操作代码,而不是创建和调用辅助函数

package main

import (
	"fmt"
	"strings"
)

// Index 返回目标字符串t的第一个索引,如果没有匹配到则返回 -1
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include 返回 true 如果目标字符串 t 存在于切片中
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

// Any 返回 true 如果切片中有一个字符串满足函数 f
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All 返回 true 如果切片中所有的字符串都满足函数 f
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter 返回一个新的切片,新切片中包含原切片中所有满足函数 f 的字符串
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map 返回一个新的切片，其中包含将函数 f 应用与原始切片中每个字符串的结果
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func main() {
	var strs = []string{"peach", "apple", "pear", "plum"}

	// 尝试各种集合函数
	fmt.Println(Index(strs, "pear"))
	fmt.Println(Include(strs, "grape"))
	fmt.Println(Any(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println(All(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println(Filter(strs, func(v string) bool {
		return strings.Contains(v, "e")
	}))
	fmt.Println(Map(strs, strings.ToUpper))

	// 上述例子中使用的都是匿名函数，也可以使用正确类型的命名函数
}

```