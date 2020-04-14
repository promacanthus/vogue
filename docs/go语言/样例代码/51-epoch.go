---
title: 51-epoch.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 51-epoch.go
showInMenu: false

---

// 程序中的一个常见要求是获取自Unix时代以来的秒数，毫秒数或纳秒数。

package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()       // 获取当前时间
	secs := now.Unix()      // 将当前时间修改为从Unix时代以来的秒数
	nanos := now.UnixNano() // 将当前时间修改为从Unix时代以来的纳秒数
	fmt.Println(now)

	// 没有UnixMillis,所以需要手动除以纳秒来获取
	millis := nanos / 1000000
	fmt.Println(secs)
	fmt.Println(millis)
	fmt.Println(nanos)

	// 将自Unix时代以来的整秒数或纳秒数转换为相应的时间
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))
}
