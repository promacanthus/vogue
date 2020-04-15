

// 行过滤器是一种常见的程序类型，它读取stdin的输入，处理它，然后将一些派生结果打印到stdout
// 常见的行过滤器有：grep和sed。

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 使用缓冲扫描器包装无缓冲的os.Stdin,
	// 以此来提供一种方便的扫描方法，将扫描器推进到下一个token（即默认扫描其中的下一行）
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ucl := strings.ToUpper(scanner.Text()) // Text函数返回当前token（在此处为输入的下一行）
		fmt.Println(ucl)
	}

	// 检查扫描中的错误
	// 遇到文件结尾是很常见的，但在扫描时不会报错
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
