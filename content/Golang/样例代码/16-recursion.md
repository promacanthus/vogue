---
title: 16-recursion
date: 2020-01-10T20:37:47.913857+08:00
draft: false
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