---
title: 56-SHA1-hashes
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---

```go
// SHA1哈希经常用于计算二进制或者文本blob的短标识，
// 例如Git就使用SHA1来标识版本化的文件和目录。

package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	s := "sha1 this string"
	// 生成散列的模式是：
	//  1. sha1.New()
	// 	2. sha1.Write(bytes)
	//  3. sha1.Sum([]byte{})
	h := sha1.New()
	h.Write([]byte(s)) //将字符串强制类型转换为字节，[]byte(s)，其中s为字符串

	// 将最终的哈希结果作为字节切片
	// 可以将Sum()方法的参数添加到一个已经存在的字节切片中，通常不需要使用它
	bs := h.Sum(nil)

	// SHA1的值通常以十六进制的形式打印出来
	// 使用%x格式化谓词将哈希的结果转换为十六进制的字符串
	fmt.Println(s)
	fmt.Printf("%x\n", bs)
}

// 还有其他的计算哈希的方式，如MD5哈希，
// 导入crypto/md5包，然后使用md5.New()

// 如果需要加密安全的hash，那就要仔细研究一下哈希强度了。

```