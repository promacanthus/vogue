//  可变参数函数

package main

import "fmt"

func sum(nums ...int) { // 将任意数量的int作为参数的函数
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {
	sum(1, 2) //  可变参数函数可以像通常那样被单个参数调用
	sum(1, 2, 3)

	nums := []int{1, 2, 3, 4} // 如果在切片中有多个参数，可以使用func(slice...)句法将切片应用到可变参数函数中
	sum(nums...)
}
