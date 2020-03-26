// 使用os.Exit根据定状态码立即退出

package main

import (
	"fmt"
	"os"
)

func main() {
	// 使用os.Exit时不会运行defers，因此永远不会调用此fmt.Println
	defer fmt.Println("!")
	os.Exit(3)
}

// 请注意，不像C语言，Go不使用main的整数返回值来指示退出状态。
// 如果想以非零状态退出，则应使用os.Exit。

// 如果您使用go run运行73-exit.go，则将通过go并打印退出。

// 通过构建和执行二进制文件，可以在终端中查看状态。
//  go build exit.go
// ./exit
// echo $?
// 3
// 注意程序永远不会打印感叹号（!） 。
