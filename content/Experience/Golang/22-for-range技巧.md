---
title: "22 for range 技巧"
date: 2020-08-06T18:15:21+08:00
draft: true
---

- [0.1. 技巧](#01-技巧)
  - [0.1.1. 示例一](#011-示例一)
  - [0.1.2. 示例二](#012-示例二)
- [0.2. 更多](#02-更多)
- [0.3. for-range与goroutine](#03-for-range与goroutine)
  - [0.3.1. 问题代码](#031-问题代码)
  - [0.3.2. 原因](#032-原因)
  - [0.3.3. 解决方案](#033-解决方案)

## 0.1. 技巧

1. 在for range开始之前，就先获取slice的大小，在后面的迭代中不会改变
2. 在for range开始之前，就先声明两个全局变量`index`和`value`

### 0.1.1. 示例一

```golang
func main() {
    v := []int{1, 2, 3}
    for i := range v {
        v = append(v, i)
    }
}
```

1. 先初始化了一个内容为1、2、3的slice
2. 然后遍历这个slice
3. 然后给这个切片追加元素

随着遍历的进行，数组v也在逐渐增大，但是for循环并不会死循环。只会遍历三次，v的结果是`[0, 1, 2]`。原因就在于for range实现的时候用到了语法糖。

对于切片的for range，它的底层代码就是：

```golang
//   for_temp := range
//   len_temp := len(for_temp)
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = for_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

第二行，在遍历之前就获取切片的长度`len_temp := len(for_temp)`，遍历的次数不会随着切片的变化而变化。

### 0.1.2. 示例二

```golang
func main() {
 slice := []int{0, 1, 2, 3}
 myMap := make(map[int]*int)
 for index, value := range slice {
  fmt.Println(&index, &value)
  myMap[index] = &value
 }
 fmt.Println("=====new map=====")
 for k, v := range myMap {
  fmt.Printf("%d => %d\n", k, *v)
 }
}

// 输出
0xc0000140e0 0xc0000140e8
0xc0000140e0 0xc0000140e8
0xc0000140e0 0xc0000140e8
0xc0000140e0 0xc0000140e8
=====new map=====
0 => 3
1 => 3
2 => 3
3 => 3
```

循环切片时，`index`和`value`这两个变量的地址在一开始是就分配好，之后一直没变过,只是被赋予的值不断变化。

`myMap[index] = &value`语句把`value`变量的地址保存到myMap中，for range迭代结束后，map的值存储的都是`value`变量在for range 开始时申请的内存地址，所以他们的值都是最后一次赋予`value`变量的值3。

> 理解技巧：`for index, value := range slice`其实是在开始之前先声明了两个全局变量，而不是在每次循环中声明局部变量（临时变量），这样也是更为合理的操作。

## 0.2. 更多

map：

```golang
// Lower a for range over a map.
// The loop we generate:
//   var hiter map_iteration_struct
//   for mapiterinit(type, range, &hiter); hiter.key != nil; mapiternext(&hiter) {
//           index_temp = *hiter.key
//           value_temp = *hiter.val
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

channel：

```golang
// Lower a for range over a channel.
// The loop we generate:
//   for {
//           index_temp, ok_temp = <-range
//           if !ok_temp {
//                   break
//           }
//           index = index_temp
//           original body
//   }
```

array：

```golang
// Lower a for range over an array.
// The loop we generate:
//   len_temp := len(range)
//   range_temp := range
//   for index_temp = 0; index_temp < len_temp; index_temp++ {
//           value_temp = range_temp[index_temp]
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

string：

```golang
// Lower a for range over a string.
// The loop we generate:
//   len_temp := len(range)
//   var next_index_temp int
//   for index_temp = 0; index_temp < len_temp; index_temp = next_index_temp {
//           value_temp = rune(range[index_temp])
//           if value_temp < utf8.RuneSelf {
//                   next_index_temp = index_temp + 1
//           } else {
//                   value_temp, next_index_temp = decoderune(range, index_temp)
//           }
//           index = index_temp
//           value = value_temp
//           // original body
//   }
```

## 0.3. for-range与goroutine

### 0.3.1. 问题代码

```golang
package main

import (
    "fmt"
    "sync"
)

func mockSendToServer(url string) {
    fmt.Printf("server url: %s\n", url)
}

func main() {
    urls := []string{"0.0.0.0:5000", "0.0.0.0:6000", "0.0.0.0:7000"}
    wg := sync.WaitGroup{}
    for _, url := range urls {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mockSendToServer(url)
        }()
    }
    wg.Wait()
}

// output
$ go run main.go
server url: 0.0.0.0:7000
server url: 0.0.0.0:7000
server url: 0.0.0.0:7000
```

### 0.3.2. 原因

goroutine的启动需要准备时间。

当主goroutine中的for循环逻辑已经走完并阻塞于`wg.Wait()`一段时间后，go func的goroutine才启动准备（准备资源，挂载M线程等）完毕。

此时url局部变量中的值是最后一次for循环的url的内容，三个goroutine准备完毕开始启动读取url局部变量时都读取到同样的内容，因此就造成了上面的bug。

### 0.3.3. 解决方案

```golang
package main

import (
    "fmt"
    "sync"
)

func mockSendToServer(url string) {
    fmt.Printf("server url: %s\n", url)
}

func main() {
    urls := []string{"0.0.0.0:5000", "0.0.0.0:6000", "0.0.0.0:7000"}
    wg := sync.WaitGroup{}
    for _, url := range urls {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()
            mockSendToServer(url)
        }(url)
    }
    wg.Wait()
}
```

将每次遍历的url所指向值，通过函数入参，作为数据资源赋予给go func,这样不管goroutine启动会有多耗时，其url已经作为goroutine的私有数据保存，后续运行就用上了正确的url，那么，上文bug也相应解除。
