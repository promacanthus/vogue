---
title: 42-panic
date: 2020-01-10T20:07:07.61445+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 42-panic
showInMenu: false

---

```go
// panic通常表示一些意料之外的错误。
// 大多数情况下，我们使用它来快速应对在正常操作期间不应该发生的错误或者不准备优雅地处理的错误。

package main

import (
	"os"
)

func main() {
	// 使用panic检查意料之外的错误
	panic("a problem")

	// panic的常见用法是一个函数返回了一个不知道如何处理或者不想处理的错误值时直接终止
	_, err := os.Create("/tmp/file")
	if err != nil {
		panic(err)
	}
}

// 运行程序将会出发panic，打印错误消息和goroutine追踪信息，然后以非零状态退出

```