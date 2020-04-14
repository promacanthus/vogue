---
title: 61-file-path.md
date: 2020-01-10T20:14:29.886698+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 61-file-path.md
showInMenu: false

---

# 61-file-path

```go
// filepath 包提供可以在操作系统之间移植的函数来解析或者构建file path
// 如 : 1. Linux中的dir/file
//         2. Windows中的dir\file。

package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	// 使用Join函数来快捷的构造路径
	// 它接受任意数量的参数并从中构造分层路径
	p := filepath.Join("dir1", "dir2", "filename")
	fmt.Println("p:", p)

	// 应该始终使用Join而不是手动连接 \s 或者 /s
	// 除了提供可移植性之外，Join还将通过删除多余的分隔符和目录来规范路径
	fmt.Println(filepath.Join("dir1//", "filename"))
	fmt.Println(filepath.Join("dir1/../dir1", "filename"))

	// Dir和Base函数可以拆分目录和文件的路径
	// split函数将在一次调用中把目录和文件两者的路径都返回
	fmt.Println("Dir(p):", filepath.Dir(p))
	fmt.Println("Base(p):", filepath.Base(p))

	// 判断一个路径是否是绝对路径
	fmt.Println(filepath.IsAbs("dir/file"))
	fmt.Println(filepath.IsAbs("/dir/file"))

	// 有些文件名在点后面后扩展名
	// 可以使用Ext函数将这些名称拆开
	filename := "config.json"
	ext := filepath.Ext(filename)
	fmt.Println(ext)

	// 要删除扩展名以查找文件名，使用strings.TrimSuffix函数
	fmt.Println(strings.TrimSuffix(filename, ext))

	// Rel返回基础路径与目标路径之间的相对路径
	// 如果无法相对于基础路径生成目标路径那么返回错误
	rel, err := filepath.Rel("a/b", "a/b/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)

	rel, err = filepath.Rel("a/b", "a/c/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)
}

```