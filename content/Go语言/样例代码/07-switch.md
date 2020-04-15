---
title: 07-switch
date: 2020-01-10T19:52:36.557673+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 07-switch
showInMenu: false

---

```go
//  Switch语句表达多个分支的条件

package main

import (
	"fmt"
	"time"
)

func main() {
	i := 2
	fmt.Print("Write ", i, " as ")
	switch i {
	case 1:
		fmt.Print("one\n")
	case 2:
		fmt.Print("two\n")
	case 3:
		fmt.Print("three\n")
	}

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday: // 在同一个case语句中使用逗号分隔多个表达式
		fmt.Println("It's the weekend")
	default: // default是可选的,没有符合的情况就默认执行这个语句
		fmt.Println("It's a weekday")
	}

	t := time.Now()
	switch { // 没有表达式的switch语句可以认为是if/else语句的另一种形式
	case t.Hour() < 12: // case表达式也可以不是常量
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	whatAmI := func(i interface{}) {
		switch t := i.(type) { // switch比较类型而不是比较值,可以用于发现接口值的类型,变量t将会具有和case子句所对应的类型
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm a int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}
	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}

```