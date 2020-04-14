---
title: 63-temp-file-dir.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 63-temp-file-dir.go
showInMenu: false

---

// 在整个程序执行过程中，经常需要创建程序退出后不需要的数据
// 临时文件和目录对此非常有用，因为它们不会随着时间的推移而污染文件系统

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 创建临时文件的最便捷的方式就是调用ioutil.TempFile函数
	// 它创建一个文件并打开它进行读写
	// 给ioutil.TempFile函数传递的一个参数是空字符串，所以会操作系统的默认位置创建文件
	f, err := ioutil.TempFile("", "sample")
	check(err)

	// 在类Unix操作系统上，目录可能是/tmp
	// 文件名以ioutil.TempFile的第二个参数作为前缀，
	// 其余部分自动选择以确保并发调用将始终创建不同的文件名。
	fmt.Println("Temp file name", f.Name())

	// 文件使用完成后清理文件，
	// 操作系统也可能在一段时间后自动清理临时文件
	// 显式的执行清理操作是一个最佳实践
	defer os.RemoveAll(f.Name())

	// 向临时文件中写入一些数据
	_, err = f.Write([]byte{1, 2, 3, 4})
	check(err)

	// 如果要创建多个临时文件，可以先创建一个临时目录
	// ioutil.TempDir的参数与ioutil.TempFIle的参数相同
	// 只是返回的是目录名而不是一个文件名
	dname, err := ioutil.TempDir("", "sampledir")
	fmt.Println("Temp dir name", dname)

	defer os.RemoveAll(dname)

	// 可以通过在临时目录前添加前缀来合成临时文件名。
	fname := filepath.Join(dname, "file1")
	err = ioutil.WriteFile(fname, []byte{1, 2}, 0666)
	check(err)
}
