---
title: 10-maps
date: 2020-01-10T19:54:58.201479+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 10-maps
showInMenu: false

---

```go
//  map是Go内置的关联数据类型（在其他语言中有时成为hash或者字典）

package main

import "fmt"

func main() {
	m := make(map[string]int) // 使用内置函数make(map[key-type]value-type)创建空的map

	// 使用name[key]=val句法设置键值对
	m["k1"] = 7
	m["k2"] = 13
	fmt.Println("map: ", m) // 使用fmt.Println将会输出map中全部的键值对

	v1 := m["k1"] // 使用name[key]来获取key所对应的值
	fmt.Println("v1: ", v1)
	fmt.Println("len: ", len(m)) // 在map上调用内置函数len时返回的是键值对的数量

	delete(m, "k2") // 内置函数delete删除map中的键值对
	fmt.Println("map: ", m)

	// 不需要获取值本身，所以使用空标识符_忽略它
	_, prs := m["k2"] // 从map中获取值时，可选的第二个返回值标识该键是否在map中，这可以用于消除缺失键和具有零值（如 “”或者0）的键之间的歧义
	fmt.Println("prs: ", prs)

	n := map[string]int{"foo": 1, "bar": 2} // 使用这种句法在一行内声明和初始化map
	fmt.Println("map: ", n)
}

// 注意，当使用fmt.Println输出map时，输出的形式如map[k:v k:v]

```