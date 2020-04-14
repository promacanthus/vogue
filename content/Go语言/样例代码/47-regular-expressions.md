---
title: 47-regular-expressions.md
date: 2020-01-10T20:09:10.954787+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 47-regular-expressions.md
showInMenu: false

---

# 47-regular-expressions

```go
// Go内置支持正则表达式

package main

import (
	"bytes"
	"fmt"
	"regexp"
)

func main() {
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach") // 判断正则表达式是否与字符串匹配
	fmt.Println(match)

	// 编译一个优化的正则表达式结构，这个结构有许多方法可用
	r, _ := regexp.Compile("p([a-z]+)ch")
	fmt.Println(r.MatchString("peach"))                                // 判断正则表达式是否与字符串匹配
	fmt.Println(r.FindString("peach punch"))                           // 查找正则表达式的匹配项
	fmt.Println(r.FindStringIndex("peach punch"))                      // 查找第一个与正则表达式匹配的项，返回匹配的开始和结束索引而不是匹配到的文本
	fmt.Println(r.FindStringSubmatch("peach punch"))                   // 返回整个匹配模式和子匹配的信息，如p([a-z]+)ch和([a-z]+)
	fmt.Println(r.FindStringSubmatchIndex("peach punch"))              // 返回整个匹配模式和子匹配的索引
	fmt.Println(r.FindAllString("peach punch pinch", -1))              // 返回正则表达式的所有匹配项
	fmt.Println(r.FindAllStringSubmatchIndex("peach punch pinch", -1)) // 返回正则表达式的所有匹配项的整个匹配模式和子匹配模式的索引
	fmt.Println(r.FindAllString("peach punch pinch", 2))               // 提供非负整数作为第二个参数，来限制正则表达式的匹配数量
	fmt.Println(r.Match([]byte("peach")) // 以字节形式匹配正则表达式

	// 使用正则表达式常见常量时，可以使用Compile的变体MustCompile
	// 普通的Compile不适用于常量，因为它有2个返回值
	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)
	fmt.Println(r.ReplaceAllString("a peach", "<fruit>")) // 将匹配的字符串子集替换为其他值

	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)	// 使用给定函数转换匹配的文本
	fmt.Println(string(out))
}

```