// 从字符串中解析数字是许多程序中基本且常见的任务。

package main

import (
	"fmt"
	"strconv" // 内置标准库strconv提供数字解析功能
)

func main() {
	p := fmt.Println

	f, _ := strconv.ParseFloat("1.234", 64) // 64表示要解析的精度位数
	p(f)

	i, _ := strconv.ParseInt("123", 0, 64) // 0表示从字符串的前缀推断出基数
	p(i)

	d, _ := strconv.ParseInt("0x1c8", 0, 64) // ParseIntn()能够识别十六进制的数字
	p(d)

	u, _ := strconv.ParseUint("789", 0, 64)
	p(u)

	k, _ := strconv.Atoi("135") // Atoi()是解析十进制的快捷函数
	p(k)

	_, e := strconv.Atoi("wat") //对于错误输入返回一个错误
	p(e)
}
