---
title: 48-json
date: 2020-01-10T20:09:27.902835+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 48-json
showInMenu: false

---

```go
// Go提供对JSON编码和解码的内置支持，包括内置和自定义类型

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// 使用下面的两个结构体来演示编码和解码自定义类型
type response1 struct {
	Page   int
	Fruits []string
}

// 只有导出字典才会被编码/解码为JSON
// 导出字段首字母大写
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {
	// 一些原子类型编码为JSON字符串
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))

	// 将切片和map编码为JSON数组和对象
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	// JSON包可以自动编码自定义数据类型
	// 它只包含编码输出中的导出字段，默认情况下将这些名称作为JSON的键
	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

	// 可以在结构体的字段声明上使用标记来自定义编码成的JSON的键的名称
	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

	// 将JSON数据解码为Go中的值
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)

	// 提供一个变量用于放置JSON包解码出的数据
	// 此map将存储键为string值为任意类型的数据
	var dat map[string]interface{}

	// 进行解码并检查相关的错误
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	// 为了使用解码后的map中的值，需要将它们转换为合适的类型
	num := dat["num"].(float64)
	fmt.Println(num)

	//访问嵌套数据需要进行一些列转换
	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

	//  将JSON解码为自定义数据类型
	// 这样做的好处是可以为程序增加额外的类型安全性
	//  并且在访问解码数据时不需要类型断言
	str := `{"page":1,"fruits":["apple","peach"]}`
	res := response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

	// 上面的例子中，使用字节和字符串作为标准输出和JSON之间的媒介
	// 还可以将JSON编码直接流式传输到os.Writers
	//  如： 1. os.Stdout
	//           2. HTTP响应体
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}

```