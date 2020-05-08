---
title: package endpoint
date: 2020-04-14T10:09:14.254627+08:00
draft: false
---

## overview

endpoint包为RPCs定义了一种抽象概念。端点（Endpoints）是许多Go-kit组件的基本构建块，它由服务端实现并有客户端调用。

## func Nop

```go
func Nop(context.Context, interface{}) (interface{}, error)
```

Nop是一个不执行任何操作的endpoint，它会返回一个nil错误。在测试时非常有用。

## type Endpoint

```go
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
```

Endpoint是服务端和客户端的基本构建块，它表示单个RPC方法。

## type Failer

```go
type Failer interface {
    Failed() error
}
```

可以通过包含业务逻辑错误详情的 Go kit 响应类型来实现故障恢复器（Failer）。 如果 Failed 方法返回一个non-nil错误，那么Go kit 传输层可能将其解释为业务逻辑错误，并且可能将其编码为与常规的成功响应不同的错误。

对于响应类型来说，不需要实现 Failer，但是对于更复杂的用例来说，它可能会有所帮助。 [Addsvc](https://github.com/go-kit/kit/tree/master/examples/addsvc) 示例显示了一个完整的应用程序应该如何使用 Failer。

## type Middleware

```go
type Middleware func(Endpoint) Endpoint
```

中间件（Middleware）是端点的链式行为修改器。

## func Chain

```go
func Chain(outer Middleware, others ...Middleware) Middleware
```

Chain 是组成中间件的辅助函数，请求将按声明的顺序遍历它们， 也就是说，第一个中间件被视为最外层的中间件。