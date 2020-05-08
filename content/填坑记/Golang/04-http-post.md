---
title: 04-http POST 请求
date: 2020-04-14T10:09:14.238627+08:00
draft: false
---

通过HTTP POST方法传输数据的时候，有两种不同的方式需要区分。

- form Data：url的参数形式
- payload（有效负载）：json串

不同形式的数据传输，需要对应不同的请求头。

- `Content-Type: application/x-www-form-urlencoded`

```json
POST /some-path HTTP/1.1
Content-Type: application/x-www-form-urlencoded

foo=bar&name=John
```

- `Content-Type: application/json`

```json
POST /some-path HTTP/1.1
Content-Type: application/json

{ "foo" : "bar", "name" : "John" }
```

## RFC的定义

如果不受请求方法和响应状态码的限制，则HTTP消息可以传输payload。payload由header字段形式的元数据和消息主体中的八位组序列形式的数据组成，在对任何传输编码进行解码之后。

HTTP中的“有效负载”始终是某些资源的部分或完整表示。 我们对payload使用单独的术语，因为某些消息仅包含关联的表示的header字段（例如，对HEAD的响应）或仅表示的某些部分（例如206状态代码）。

专门定义payload而不是相关表示的HTTP header字段称为“payload header 字段”。 payload header 字段由`HTTP/1.1`定义。

仅当存在消息body时，payload body 才出现在消息中。 payload body 是通过对可能已应用以确保安全正确地传输消息的任何传输编码进行解码而从消息主体获得的。
