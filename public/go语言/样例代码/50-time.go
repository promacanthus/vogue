---
title: 50-time.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 50-time.go
showInMenu: false

---

// Go为times和durations提供广泛支持

package main

import (
	"fmt"
	"time"
)

func main() {
	p := fmt.Println
	now := time.Now() // 获取当前时间
	p(now)

	// 可以通过提供年月日等来构建时间结构，TImes总是与位置有关，如时区
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	p(then)

	// 可以提取时间值的各个组件
	p(then.Year())
	p(then.Month())
	p(then.Day())
	p(then.Hour())
	p(then.Minute())
	p(then.Second())
	p(then.Nanosecond())
	p(then.Location())

	// 获取星期几
	p(then.Weekday())

	// 比较两个时间，比较前一个值与后一个值的前、后、相等关系
	p(then.Before(now))
	p(then.After(now))
	p(then.Equal(now))

	// Stub方法返回两个时间之间的时间间隔，默认以小时为单位
	diff := now.Sub(then)
	p(diff)

	// 可以以各种单位计算持续时间的长度
	p(diff.Hours())
	p(diff.Minutes())
	p(diff.Seconds())
	p(diff.Nanoseconds())

	// 使用Add根据给定的时间间隔向来前推进时间
	// 时间间隔前面加上负号（-）来向后推移时间
	p(then.Add(diff))
	p(then.Add(-diff))
}
