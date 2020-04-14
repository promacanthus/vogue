---
title: 67-env-var.md
date: 2020-01-10T20:16:15.444294+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 67-env-var.md
showInMenu: false

---

# 67-env-var

```go
// 环境变量是将配置信息传递给Unix程序的通用机制。

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	os.Setenv("FOO", "1")                // 使用os.Setenv来设置一个键值对
	fmt.Println("FOO", os.Getenv("FOO")) // 使用os.Getenv获取给定键的值
	fmt.Println("BAR", os.Getenv("BAR")) // 如果系统中没有设置这个值则会返回一个空的字符串

	fmt.Println("")
	// os.Environ获取系统中全部环境变量的键值对
	// 这将会返回一个key=value形式的字符串切片
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=") // 使用strings.Split将键和值分开
		fmt.Println(pair[0])
	}
}

```