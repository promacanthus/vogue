---
title: 03-error处理
date: 2020-04-14T10:09:14.238627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- 填坑记
- Golang
summary: 03-error处理
showInMenu: false

---

## 问题

报错信息如下：

```bash
http: named cookie not present
```

引起问题的代码:

```golang
cookie, err := r.Cookie("user")
    if err != nil {
        log.Fatal(err)
    }
```

问题的现象就是，浏览器在发送http请求的时候，在Request中没有我们需要的`Cookie：user=XXX`，因此代码执行`log.Fatal(err)`,然后退出了。

开始以为是Cookie的问题，其实是golang错误处理的理解不到位。

## 分析

- `Fatal()`：等价于`Print()`然后会执行`os.Exit()`调用。

> 也就是说，上面的程序执行过程中遇到err就会直接打印日志信息，并退出了。

- `Panic()`：等价于`Println()`然后会执行内置函数panic()。

> 同样的，panic之后程序也会退出，不过会打印出更详细的信息。但是panic可以用recover来避免程序崩溃。

关于，panic、recover、defer的具体说明和用法，参考[这里](../../编程语言/基础/16-panic&recover&defer.md)。

大部分情况下，都希望记录并打印异常，但是程序不崩溃继续运行。

所以上面的代码可以这样修改：

```golang

defer func() {
    // p!=nil 判断确实发生了panic
    if p := recover(); p != nil {
        log.Printf("panic: %s\n", p)
    }
}()

cookie, err := r.Cookie("user")
    if err != nil {
        panic(err)
    }
```

## 另外

其实有很多比官方log库更好用的一些记录日志的库，比如[`go-kit`](https://pkg.go.dev/github.com/go-kit/kit/log?tab=doc)的log库。
