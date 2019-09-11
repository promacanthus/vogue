// Go标准库为net/http包中的HTTP客户端和服务器提供了出色的支持

package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {
	// 向HTTP服务器发出HTTP GET请求
	// http.Get是创建http.Client对象并调用其Get方法的便捷方式
	// 它使用http.DefaultClient对象，该对象具有有用的默认设置
	resp, err := http.Get("http://gobyexample.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 输出HTTP响应状态
	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	// 输出响应体的前五行
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
