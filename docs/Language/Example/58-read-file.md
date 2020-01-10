# 58-read-file

```go
// 读取和写入文件是Go程序要完成的基本任务。

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// 读取文件需要检查大部分调用错误，
// check函数有助于我们精简代码
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 读取文件最基本的操作是将整个文件个内容都放在内存中
	dat, err := ioutil.ReadFile("./58-read-file")
	check(err)
	fmt.Println(string(dat))

	// 控制文件的读取方式和读取位置
	// 使用os.Open函数获取一个os.FIle返回值
	f, err := os.Open("./58-read-file")
	check(err)

	// 从文件开头读取一些字节
	// 创建一个byte切片来运行读取最多5个字节，但也需要注意实际读取了多少字节
	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes:%s\n", n1, string(b1[:n1]))

	// 搜索文件中的指定位置，并从那里读取
	o2, err := f.Seek(6, 0)
	check(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes @ %d:", n2, o2)
	fmt.Printf("%v\n", string(b2[:n2]))

	// io包提供一些读取文件的有用函数
	// 例如使用ReadAtLast函数实现上述功能
	o3, err := f.Seek(6, 0)
	check(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	check(err)
	fmt.Printf("%d bytes @ %d:%s\n", n3, o3, string(b3))

	// 没有内置的后退功能，使用Seek(0,0)来实现后退
	_, err = f.Seek(0, 0)
	check(err)

	// bufio包实现了一个缓冲读取器，对于多个小的读取很高效，而且可以提供额外的读取方法
	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5) // 根据参数值读取文件中的字节数而不修改读取器的位置
	check(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	// 文件操作完成后关闭文件
	// 通常使用defer在打开文件操作后立即添加上关闭操作
	f.Close()
}

```