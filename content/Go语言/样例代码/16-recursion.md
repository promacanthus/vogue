---
title: 16-recursion
date: 2020-01-10T19:58:47.857511+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 16-recursion
showInMenu: false

---


```go
//  Go支持递归函数

package main

import "fmt"

func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

func main() {
	fmt.Println(fact(7))
}

```