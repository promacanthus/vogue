---
title: 41-sorting-by-func
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---

```go
// 如何对集合中的数据进行非自然顺序的排序。

// 假设要对string更具长度而非字母顺序排序。

package main

import (
	"fmt"
	"sort"
)

// 为了按Go中自定义函数排序，需要一个相应的类型
type byLength []string // 创建一个byLength类型，它只是内置类型[]string的别名

// 在自定义类型上实现sort.Interface中的len、Less、Swap方法
// 这样就可以使用sort包的泛型Sort函数
func (a byLength) Len() int { return len(a) }

// Swap和Less在不同的类型之间都是很相似的
func (a byLength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less中写的是实际的自定义排序逻辑
// 在我们的例子中，根据string的长度进行升序排序，所以使用len()进行长度比较
func (a byLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }

func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	// 将原始的fruits切片转换为buLength类型，
	// 然后在该类型上使用sort.Sort函数来实现自定义排序
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}

// 通过遵循这种模式来创建自定义类型，
// 并在该类型上实现sort.Interfacce的三个接口方法，
// 然后在该自定义类型的集合上调用sort.Sort函数，
// 可以按照任意方式对Go切片进行排序

```