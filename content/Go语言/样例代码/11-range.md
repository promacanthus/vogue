---
title: 11-range.md
date: 2020-01-10T19:56:01.177451+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 11-range.md
showInMenu: false

---

# 11-range

```go
// range遍历各种数据结构中的元素

package main

import "fmt"

func main() {
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums { // 使用range来对切片中的数字求和，数组也是这样使用
		sum += num
	}
	fmt.Println("sum: ", sum)

	for i, num := range nums { // 在数组和切片上使用range都会返回每个条目的索引和值，上面的例子中，不需要返回的索引时，使用空白标识符来忽略它
		if num == 3 {
			fmt.Println("index: ", i)
		}
	}

	kvs := map[string]string{"a": "apple", "b": "banana"} // 在map上使用range将会迭代其中的键值对
	for k, v := range kvs {
		fmt.Printf("%s --> %s\n", k, v)
	}

	for k := range kvs { // range也可以只遍历map中的键值
		fmt.Println("key: ", k)
	}

	for i, c := range "go" { // 在字符串上使用range将会遍历unicode代码点，返回的第一个值是rune的索引，第二个值是rune自身
		fmt.Println(i, c)
	}
}

// go中byte就是uint8，rune就是int32

```