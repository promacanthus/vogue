---
title: 69-HTTP-server
date: 2020-01-10T20:16:49.487616+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 69-HTTP-server
showInMenu: false

---

```go
// 使用net/http包可以轻松编写基本的HTTP服务器

package main

import (
	"fmt"
	"net/http"
)

// net/http 服务器的基本概念是handlers(处理器)
// 一个处理器是一个实现http.Handler接口的对象
// 编写处理器的常用方法是在具有相应签名的函数上使用http.HandlerFunc适配器
func hello(w http.ResponseWriter, req *http.Request) {
	// 作为处理器的函数将http.ResponseWriter和http.Request作为参数
	// 响应输出器用于填写HTTP响应
	// 在这里只是响应"hello\n"
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	// 这个处理器通过读取所有HTTP请求标头并将它们回显到响应主体来做一些更复杂的事情
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v:%v\n", name, h)
		}
	}
}

func main() {
	// 使用http.HandleFunc函数在服务器路由上注册上面定义的处理器。
	// http.HandleFunc在net/http包中设置默认路由器，并将函数作为参数
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	// 使用端口和处理器调用ListenAndServe函数
	// nil表示使用刚刚设置的默认路由器
	http.ListenAndServe(":8090", nil)
}

```