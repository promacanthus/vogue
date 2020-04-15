

//  在Go中写入文件的方式和读取文件的方式类似。

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 将字符串或者字节写入到文件中
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("./59-write-file1", d1, 0644)
	check(err)

	// 要进行更精细的写入，那么就打开一个文件然后写入
	f, err := os.Create("./59-write-file2")
	check(err)

	// 根据惯例打开一个文件后，立即使用defer关闭它
	defer f.Close()

	// 将字节切片写入到文件中
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	// 使用WriteString函数写入字符串到文件中
	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	// 启动一个Sync()将写入数据刷新到持久化存储中
	f.Sync()

	// bufio也提供缓冲写入器
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)

	// 使用Flush()函数确保将所有缓冲操作都已经应用到底层的写入器
	w.Flush()
}
