---
title: 52-time-format-parse
date: 2020-01-10T20:10:27.191007+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 52-time-format-parse
showInMenu: false

---


```go
// Go支持通过基于模式设计的时间格式化和解析。

package main

import (
	"fmt"
	"time"
)

func main() {
	p := fmt.Println
	t := time.Now()
	// 根据RFC3339规范格式化时间
	p(t.Format(time.RFC3339))

	// 根据RFC3339规范解析时间
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	p(t1)

	// 根据自定义样式格式化时间
	p(t.Format("3:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))

	// 根据自定义样式解析时间
	form := "3 04 PM"
	t2, _ := time.Parse(form, "8 41 PM")
	p(t2)

	// 对于纯数字的表示可以使用标准格式化输出和时间提取组件
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	// 格式解析错误
	ansic := "Mon Jan _2 15:04:05 2006"
	_, e := time.Parse(ansic, "8:41PM")
	p(e) // parsing time "8:41PM" as "Mon Jan _2 15:04:05 2006": cannot parse "8:41PM" as "Mon"
}

```