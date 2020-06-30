---
title: "12 Slice作为参数的坑"
date: 2020-06-30T15:15:19+08:00
draft: true
---

关于数组和切片，golang官方博客有文章详细说明，点击[这里](https://blog.golang.org/slices)。其实这里说的已经很清楚了，论好好阅读官方说明的重要性。

## 问题

```go
package main

import (
 "fmt"
)

func main() {
 slice := []int{0, 1, 2, 3}

 fmt.Printf("slice: %v slice addr %p \n", slice, &slice)
 // slice: [0 1 2 3] slice addr 0xc00000c080

 ret := changeSlice(slice)
 fmt.Printf("slice: %v slice addr %p | ret: %v ret addr %p \n", slice, &slice, ret, &ret)
 // slice: [0 111 2 3] slice addr 0xc00000c080 | ret: [0 111 2 3] ret addr 0xc00000c0c0

 res := appendSlice(slice)
 fmt.Printf("slice: %v slice addr %p | res: %v ret addr %p \n", slice, &slice, res, &res)
 //  slice: [0 111 2 3] slice addr 0xc00000c080 | res: [0 111 2 3 1] ret addr 0xc00000c120
}

func changeSlice(s []int) []int {
 s[1] = 111
 return s
}

func appendSlice(s []int) []int {
 s = append(s, 1)
 return s
}

```

从上面代码和输出结果（注释部分）可以看出：

1. `changeSlice()`函数对外部slice生效了
2. `appendSlcie()`函数对外部没有生效

## 分析

### 值传递和引用传递

golang中只有**值传递**，所有的**引用传递**都是直接把对应的指针拷贝过去了，所以修改能直接在原对象生效。

### slice

很多地方都说slice是引用类型（这是相对于slice底层的数组而言的），其实slice是一个结构体类型（也就是值类型）。

```go
type slice struct {
 array unsafe.Pointer
 len   int
 cap   int
}
```

### 为啥changeSlice生效了

因为，`slice`是一个结构体且参数传递是值传递，所以`changeSlice()`函数中的`s`是`slice`的一个副本，所以`changeSlice()`函数的返回值`ret`的地址与`slice`不同，他们是内存中的两个对象。

在`slice`中`array`是一个指针，指向底层数组的开头，所以在`changeSlice()`函数中`s[1] = 111`是对底层数组的修改。那么在`main()`函数中不论是读取`slice`还是读取`ret`，他们都指向同一个底层数组，所以看起来就是`changeSlice()`函数修改了传入的切片对象的预原始值。

### 为啥appendSlice没有生效

根据上面的分析，在`appendSlice()`函数中的`append()`操作是作用在`res`上而不是`slice`上。

### 让`appendSlice`生效

因为slice其实是一个结构体而不是一个引用。要让`appendSlice`生效，只要传入引用就可以，代码修改如下：

```go
res := appendSlice(&slice)

func appendSlice(s *[]int) *[]int {
 *s = append(*s, 1)
 return s
}
```
