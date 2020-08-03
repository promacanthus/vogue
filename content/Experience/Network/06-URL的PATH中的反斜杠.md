---
title: "06 URL的PATH中的反斜杠"
date: 2020-08-03T13:28:23+08:00
draft: true
---

## 示例

首先看个例子：

```bash
http://www.abc.com/abc/
http://www.abc.com/abc
```

他们两个不同的地址：

- `/abc/`：表示的是目录，
- `/abc`：表示的是文件，

## server端

一般来说，索引页面（如文章列表）作为目录，内容页面作为文件。

在开发的过程中一般也要满足上面约定俗成的规定，因此在路由中匹配Path的时候就要注意这一点。

比如，开发一个静态文件服务器，那么：

- 文件目录的Path是：`/files/datasets/`
- 具体文件的Path是：`/files/pom.xml`
- 子目录的Path是：`/files/datasets/`

如果访问目录的时候没有在Path中协商最后的斜杠`/`，那么有两种情况：

1. 在server端进行处理，自动给Path加上最后一个斜杠

2. 如果server端没有进行上述处理，那么返回301响应

```bash
curl -i localhost:10010/datasets

HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: /datasets/
Date: Mon, 03 Aug 2020 05:39:29 GMT
Content-Length: 47

<a href="/datasets/">Moved Permanently</a>.

```
