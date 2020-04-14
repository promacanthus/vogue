---
title: 55-URL-parse.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 55-URL-parse.go
showInMenu: false

---

// URLs提供了一种统一的资源定位方式。

package main

import (
	"fmt"
	"net"
	"net/url"
)

func main() {

	p := fmt.Println

	// 解析此URL，其中包含scheme、认证信息、主机、端口、路径、请求参数和请求片段
	s := "postgres://user:pass@host.com:5432/path?k=v#f"

	// 解析URL并检查是否存在错误
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	p(u.Scheme) // 获取scheme
	p(u.User)   // User中包含全部的认证信息

	// 使用Password()和Username()来获取单独的值
	pw, _ := u.User.Password()
	p(pw)

	p(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	p(host)
	p(port)

	// 获取路径和#后面的请求片段
	p(u.Path)
	p(u.Fragment)

	p(u.RawQuery) // 以k=v的字符串形式获取请求参数
	// 将请求参数解析为map
	// 解析后的请求参数的map是map[string]slice形式
	m, _ := url.ParseQuery(u.RawQuery)
	p(m)
	// 如果只需要第一个值，则索引0即可
	p(m["k"][0])
}
