# 06-if-else

```go
//  用Golang中的分支（if/else）非常简洁明了

package main

import "fmt"

func main() {
	// Go语言中的条件不需要括号,但需要大括号
	if 7%2 == 0 {
		fmt.Println("7 is even")
	} else {
		fmt.Println("7 is odd")
	}

	// 没有else语句的if语句
	if 8%4 == 0 {
		fmt.Println("8 is divisible by 4")
	}

	// 变量声明可以优先于条件,变量声明后在分支中都可以使用
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}

// 在Go语言中没有三元组,所以即使是最基本的条件也需要完整的if语句

```