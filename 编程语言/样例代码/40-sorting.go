// Go的sort包实现了内置和用户定义类型的排序。

package main

import (
	"fmt"
	"sort"
)

func main() {
	// sort函数对于内置类型有特定方法
	// 对于切片使用就地排序，不会生成新切片
	strs := []string{"c", "b", "a"}
	sort.Strings(strs)
	fmt.Println("Strings:", strs)

	ints := []int{7, 2, 5}
	sort.Ints(ints)
	fmt.Println("Ints:", ints)

	// 判断给定切片是否 已经升序排列
	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted:", s)
}
