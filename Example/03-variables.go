package main

import "fmt"

func main() {
	// var 声明一个或多个变量，可以在一次声明中同时声明多个变量
	// Go将会推断初始化变量的类型
	var a = "initial"
	fmt.Println(a)

	var b, c int = 1, 2
	fmt.Println(b, c)

	var d = true
	fmt.Println(d)

	// 在声明变量时没有对它进行任何初始化操作，那么该变量的默认值为该类型的零值。（如 int 的零值为 0）
	var e int
	fmt.Println(e)

	// := 句法是一种短变量声明和初始化的方式。（如 上面的变量f的声明和初始化）
	f := "apple"
	fmt.Print(f)
}
