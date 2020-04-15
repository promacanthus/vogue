---
title: 65-cmd-line-flags
date: 2020-01-10T20:15:37.509103+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go语言
- 样例代码
summary: 65-cmd-line-flags
showInMenu: false

---

```go
// 命令行标识是制定命令行程序选项的常用方法
// 例如： wc -l 其中 -l 就是一个命令行标识

package main

import (
	"flag" // Go提供flag包支持基本的命令行标识解析
	"fmt"
)

func main() {
	// 基本的标识声明可以用字符串、整数和布尔型

	// 声明一个字符串型标识word，默认值为foo并带有一个简短的描述
	// flag.String函数返回一个字符串指针（而不是一个字符串值）
	wordPtr := flag.String("word", "foo", "a string")
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	// 也可以使用一个程序中已经存在的变量来声明一个命令行标识
	// 注意，西药传递一个指针给标识的声明函数
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// 一旦所以的标识都声明好，就可以调用flag.Parse函数来执行命令行解析
	flag.Parse()

	// 请注意在输出时，需要取消指针引用，即获取指针的实际值
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
}

//请注意在命令行执行程序时，如果省略标识，则会自动采用其默认值

// 请注意，flag包要求所有命令行标识要在位置参数之前出现，否则标识会被解析为位置参数
// go  run  65-cmd-line-flags.go -word=opt a1 a2 a3 -numb=7
// word: opt
// numb: 42
// fork: false
// svar: bar
// tail: [a1 a2 a3 -numb=7]

// 使用-h或--help标志可以获得命令行程序的自动生成的帮助文本
// go run 65-cmd-line-flags.go -h
// Usage of /tmp/go-build690548213/b001/exe/65-cmd-line-flags:
//   -fork
//         a bool
//   -numb int
//         an int (default 42)
//   -svar string
//         a string var (default "bar")
//   -word string
//         a string (default "foo")
// exit status 2

// 如果使用未提供给flag包的标识，程序将打印错误消息并再次显示帮助文本

```