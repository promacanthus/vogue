---
title: 45-string-func
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---


```go
//   标准库strings包提供了很多有用的字符串相关函数。

package main

import (
	"fmt"
	s "strings"
)

// 为 fmt.Println 创建别名 p
var p = fmt.Println

func main() {
	//  这些都是strings包的函数，而不是字符串对象自身的方法
	// 需要将待处理的字符串作为第一个参数传递给函数
	p("Contains:\t", s.Contains("test", "es"))
	p("Count:\t\t", s.Count("test", "t"))
	p("HasPrefix:\t", s.HasPrefix("test", "te"))
	p("HasSuffix:\t", s.HasSuffix("test", "st"))
	p("Index:\t\t", s.Index("test", "e"))
	p("Join:\t\t", s.Join([]string{"a", "b"}, "-"))
	p("Repeat:\t\t", s.Repeat("a", 5))
	p("Replace:\t", s.Replace("foo", "o", "0", -1))
	p("Split:\t\t", s.Split("a-b-c-d-e", "-"))
	p("ToLower:\t", s.ToLower("TEST"))
	p("ToUpper:\t", s.ToUpper("test"))
	p()

	// 以下函数不是strings包中的函数

	// 以字节的形式获取字符串的长度
	p("Len:\t\t", len("hello"))
	// 以字节的形式获取字节在字符串中的索引
	p("Char:\t\t", "hello"[1])
	//  需要注意的是，上面两个函数都是工作中字节级别的

	// Go使用UTF-8编码字符串，所以这两个函数通常很有用
	// 如果正在使用可能的多字节字符，则需要使用编码感知操作，查看这片博文 https://blog.golang.org/strings
}

```