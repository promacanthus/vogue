# 46-string-formating

```go
//  Go在printf中对字符串的格式化输出的支持非常好。

package main

import (
	"fmt"
	"os"
)

type point struct {
	x, y int
}

func main() {
	p := point{1, 2}

	// Go提供多种打印动词来格式化一般的Go值
	fmt.Printf("%v\n", p)    // 打印point结构体
	fmt.Printf("%+v\n", p)   // 如果打印的值是结构体，%+v将会输出结构体的字段名
	fmt.Printf("%#v\n", p)   // %#v将会输出一段Go语法来表示该值，例如生成该值的源代码片段
	fmt.Printf("%T\n", p)    // 打印值的类型
	fmt.Printf("%t\n", true) // 格式化布尔值
	fmt.Printf("%d\n", 123)  // 有多种格式化输出整数的方式，%d是标准的十进制格式化输出
	fmt.Printf("%b\n", 14)   // 输出二进制形式
	fmt.Printf("%c\n", 33)   // 输出与给定整数对应的字符
	fmt.Printf("%x\n", 456)  // 输出十六进制编码
	fmt.Printf("%f\n", 78.9) //有多种格式化输出浮点数的方式，%f是标准的十进制格式化输出

	// 以科学计数法的形式输出浮点数，%e和%E略有不同
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)

	fmt.Printf("%s\n", "\"string\"") // 输出字符串使用%s
	fmt.Printf("%q\n", "\"string\"") // 像源码一样双引号引其字符串使用%q
	fmt.Printf("%x\n", "hex this")   // 像整数一样，以十六进制形式输出字符串，每个输入字节有两个输出字符
	fmt.Printf("%p\n", &p)           // %p用于输出指针

	// 格式化数字时，通常需要控制结果图的宽度和精度
	// 要指定整数的宽度，在%后加上数字
	// 默认情况下，结果右对齐，并使用空格填充
	fmt.Printf("|%6d|%6d|\n", 12, 345)
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)   // 指定浮点数宽度的同时限定小数点的精度
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45) // 使用 - 符号实现左对齐

	// 输出字符串时控制宽度，默认为右对齐
	fmt.Printf("|%6s|%6s|\n", "foo", "b")
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	// 以上都是Printf将格式化的内容输出到os.Stdout

	// Sprintf格式化并返回一个字符串而不在任何地方输出它
	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)

	// 使用Fprintf格式化并打印到os.Stdout以外的io.Writers
	fmt.Fprintf(os.Stderr, "an %s\n", "error")

}

```