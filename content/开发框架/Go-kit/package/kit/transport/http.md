---
title: package http
date: 2020-04-14T10:09:14.254627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Go-kit
summary: http
showInMenu: false

---

## Overview

http包为端点提供通用的HTTP绑定。

## Constants

```go
const (
    // ContextKeyRequestMethod由 PopulateRequestContext在 context中填充，
    // 它的值是  r.Method。
    ContextKeyRequestMethod contextKey = iota

    // ContextKeyRequestURI 由 PopulateRequestContext 在 context 中填充，
    // 它的值是 r.RequestURI。
    ContextKeyRequestURI

    // ContextKeyRequestPath 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.URL.Path。
    ContextKeyRequestPath

    // ContextKeyRequestProto由 PopulateRequestContext 在 context 中填充，
    // 它的值是 r.Proto。
    ContextKeyRequestProto

    // ContextKeyRequestHost 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Host。
    ContextKeyRequestHost

    // ContextKeyRequestRemoteAddr 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.RemoteAddr。
    ContextKeyRequestRemoteAddr

    // ContextKeyRequestXForwardedFor由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("X-Forwarded-For")。
    ContextKeyRequestXForwardedFor

    // ContextKeyRequestXForwardedProto 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("X-Forwarded-Proto")。
    ContextKeyRequestXForwardedProto

    // ContextKeyRequestAuthorization 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("Authorization")。
    ContextKeyRequestAuthorization

    // ContextKeyRequestReferer 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("Referer")。
    ContextKeyRequestReferer

    // ContextKeyRequestUserAgent 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("User-Agent")。
    ContextKeyRequestUserAgent

    // ContextKeyRequestXRequestID 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("X-Request-Id")。
    ContextKeyRequestXRequestID

    // ContextKeyRequestAccept 由 PopulateRequestContext在 context 中填充，
    // 它的值是 r.Header.Get("Accept")。
    ContextKeyRequestAccept

    // 每当指定ServerFinalizerFunc，就会在context中填充 ContextKeyResponseHeaders，
    // 它的值是  http.Header类型, 并且仅在写入整个响应后才捕获。
    ContextKeyResponseHeaders

    // 每当指定ServerFinalizerFunc，就会在 context 中填充 ContextKeyResponseSize，
    // 它的值是 int64 类型。
    ContextKeyResponseSize
)
```

## func DefaultErrorEncoder

```go
func DefaultErrorEncoder(_ context.Context, err error, w http.ResponseWriter)
```

`Defaulterrorencoder` 将错误写入 `ResponseWriter`，默认情况下是 `text/plain` 的内容类型、纯文本体的错误和500状态码。

1. 如果`error`实现 `Headerer`，则提供的头将应用于响应。
2. 如果`error`实现 `json.Marshaler`，那么解析器在处理成功后，将使用`application/json` 的内容类型和以 JSON 编码的错误形式。
3. 如果`error`实现了 `StatusCoder`，那么提供的 `StatusCode` 将被使用而不是500。

## func EncodeJSONRequest

```go
func EncodeJSONRequest(c context.Context, r *http.Request, request interface{}) error
```

`Encodejsonrequest` 是一个 `EncodeRequestFunc`，它将请求序列化为 `ResponseWriter` 的JSON 对象。 许多 JSON-over-HTTP 服务可以将其作为一种合理的默认值。

1. 如果请求实现 `Headerer`，则提供的头将应用于请求。

## func EncodeJSONResponse

```go
func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error
```

`Encodejsonresponse` 是一个 `EncodeResponseFunc`，它将响应序列化为 `ResponseWriter` 的 JSON 对象。 许多 JSON-over-HTTP 服务可以将其作为一种合理的默认值。

1. 如果响应实现 `Headerer`，则提供的头将应用于响应。
2. 如果响应实现了 `StatusCoder`，则将使用提供的 `StatusCode` 而不是200。

## func EncodeXMLRequest

```go
func EncodeXMLRequest(c context.Context, r *http.Request, request interface{}) error
```

`Encodexmlrequest` 是一个 `EncodeRequestFunc`，它将请求序列化为请求体的 XML 对象。 

1. 如果请求实现 `Headerer`，则提供的头将应用于请求。

## func NopRequestDecoder

```go
func NopRequestDecoder(ctx context.Context, r *http.Request) (interface{}, error)
```

`Noprequestdecoder` 是一个`DecodeRequestFunc` ，可用于不需要解码的请求，并简单地返回 nil，nil。

## func PopulateRequestContext

```go
func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context
```

`Populaterequestcontext` 是一个 `RequestFunc`，它将多个值从 HTTP 请求填充到`context`中。 这些值可以使用此包中相应的 `ContextKey` 类型提取。

## 样例

```go
handler := NewServer(
    func(ctx context.Context, request interface{}) (response interface{}, err error) {
        fmt.Println("Method", ctx.Value(ContextKeyRequestMethod).(string))
        fmt.Println("RequestPath", ctx.Value(ContextKeyRequestPath).(string))
        fmt.Println("RequestURI", ctx.Value(ContextKeyRequestURI).(string))
        fmt.Println("X-Request-ID", ctx.Value(ContextKeyRequestXRequestID).(string))
        return struct{}{}, nil
    },
    func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
    func(context.Context, http.ResponseWriter, interface{}) error { return nil },
    ServerBefore(PopulateRequestContext),
)

server := httptest.NewServer(handler)
defer server.Close()

req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/search?q=sympatico", server.URL), nil)
req.Header.Set("X-Request-Id", "a1b2c3d4e5")
http.DefaultClient.Do(req)

// 输出
Method PATCH
RequestPath /search
RequestURI /search?q=sympatico
X-Request-ID a1b2c3d4e5
```

## type Client

```go
type Client struct {
    client         HTTPClient
    method         string
    tgt            *url.URL
    enc            EncodeRequestFunc
    dec            DecodeResponseFunc
    before         []RequestFunc
    after          []ClientResponseFunc
    finalizer      []ClientFinalizerFunc
    bufferedStream bool
}
```

Clinet封装URL并提供实现`endpoint.Endpoint`的方法。

## func NewClient

```go
func NewClient(method string, tgt *url.URL, enc EncodeRequestFunc, dec DecodeResponseFunc, options ...ClientOption,) *Client
```

Newclient 为单个远程方法构造一个可用的 Client。

## func (Client) Endpoint

```go
func (c Client) Endpoint() endpoint.Endpoint
```

Endpoint 返回调用远程endpoint的可用endpoint。

## type ClientFinalizerFunc

```go
type ClientFinalizerFunc func(ctx context.Context, err error)
```

Clientfinalizerfunc 可以用于在客户端 HTTP 请求已经结束，响应返回后执行工作。 主要用途是用于错误日志记录。 在context中，带 ContextKeyResponse 前缀的 key 下提供了额外的响应参数。 注意: err 可能为 nil。 根据错误发生的时间，可能也没有额外的响应参数。

## type ClientOption

```go
type ClientOption func(*Client)
```

ClientOption 为客户机设置一个可选参数。

## func BufferedStream

```go
func BufferedStream(buffered bool) ClientOption
```

BufferedStream设置`Response.Body`是否依然处于打开状态，并允许它之后可以被读取。 将文件作为缓冲流传输时很有效。 在适当的时候不得不关闭Body来结束请求。

## func ClientAfter

```go
func ClientAfter(after ...ClientResponseFunc) ClientOption
```

Clientafter 在传入的 HTTP 请求被解码之前，设置应用于该请求的 ClientResponseFuncs。 这对于从响应中获取任何信息并在解码之前添加到上下文中非常有用。

## func ClientBefore

```go
func ClientBefore(before ...RequestFunc) ClientOption
```

Clientbefore 设置在调用传出 HTTP 请求之前应用于该请求的 RequestFuncs。

## func ClientFinalizer

```go
func ClientFinalizer(f ...ClientFinalizerFunc) ClientOption
```

Clientfinizer 在每个 HTTP 请求的末尾执行。 默认情况下，没有注册finalizer。

## func SetClient

```go
func SetClient(client HTTPClient) ClientOption
```

Setclient 设置用于请求的基础 HTTP 客户端。 默认情况下，使用`http.DefaultClient`。

## type ClientResponseFunc

```go
type ClientResponseFunc func(context.Context, *http.Response) context.Context
```

ClientResponseFunc 可以从 HTTP 请求中获取信息，并使响应可用于使用。 ClientResponseFunc只在clients中执行，在一个请求已经到达，但是还没被解码之前执行。

## type DecodeRequestFunc

```go
type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)
```

DecodeRequestFunc 从 HTTP 请求对象中提取用户域请求对象。 它被设计用于 HTTP 服务器，用于服务器侧端点。 一个简单的 DecodeRequestFunc 可以是从请求体解码JSON到具体请求类型。

## type DecodeResponseFunc

```go
type DecodeResponseFunc func(context.Context, *http.Response) (response interface{}, err error)
```

DecodeResponseFunc 从 HTTP 响应对象中提取用户域响应对象。 它被设计用于 HTTP 客户端，用于客户端侧端点。 一个简单的 DecodeResponseFunc 可以是从响应体解码JSON到具体响应类型。

## type EncodeRequestFunc

```go
type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error
```

EncodeRequestFunc 将传递的请求对象编码为 HTTP 请求对象。 它被设计用于 HTTP 客户端，用于客户端侧端点。 一个简单的 EncodeRequestFunc 可以是将 JSON对象直接编码到请求体。

## type EncodeResponseFunc

```go
type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error
```

EncodeResponseFunc 将传递的响应对象编码到 HTTP 响应写入器。 它被设计用于 HTTP 服务器，用于服务器侧端点。 一个简单的 EncodeResponseFunc 可以是将 JSON 对象直接编码到响应主体。

## type ErrorEncoder

```go
type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)
```

ErrorEncoder 负责将错误编码到 ResponseWriter。 鼓励用户使用自定义 ErrorEncoder 向其客户端编码 HTTP 错误，并可能希望传递和检查自定义的错误类型。 请参阅`shaipping/handling`服务示例。

## type HTTPClient

```go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}
```

HTTPClient 是一个建模 `*http.Client` 的接口。

## type Headerer

```go
type Headerer interface {
    Headers() http.Header
}
```

Headerer 由 DefaultErrorEncoder 检查 。 如果一个错误值实现了 Headerer，那么在设置了 Content-Type 之后，提供的头将应用于响应写入器。

## type RequestFunc

```go
type RequestFunc func(context.Context, *http.Request) context.Context
```

Requestfunc 可以从 HTTP 请求中获取信息并将其放入请求context中。

1. 在服务器中，在调用端点之前执行 RequestFuncs。
2. 在客户端中，RequestFuncs 在请求创建完成后，调用 HTTP 客户端之前执行。

## func SetRequestHeader

```go
func SetRequestHeader(key, val string) RequestFunc
```

SetRequestHeader 返回一个设置给定头的 RequestFunc。

## type Server

```go
type Server struct {
    e            endpoint.Endpoint
    dec          DecodeRequestFunc
    enc          EncodeResponseFunc
    before       []RequestFunc
    after        []ServerResponseFunc
    errorEncoder ErrorEncoder
    finalizer    []ServerFinalizerFunc
    errorHandler transport.ErrorHandler
}
```

服务器封装了端点并实现 `http.Handler`。

## func NewServer

```go
func NewServer(e endpoint.Endpoint, dec DecodeRequestFunc, enc EncodeResponseFunc, options ...ServerOption,) *Server
```

构造了一个新的server，它实现了 `http.Hander`并封装提供的端点。

## func (Server) ServeHTTP

```go
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

ServeHTTP 实现 `http.Handler`。

## type ServerFinalizerFunc

```go
type ServerFinalizerFunc func(ctx context.Context, code int, r *http.Request)
```

ServerFinalizerFunc 可以用于响应已经被写入客户端，在HTTP 请求结束时执行工作。 主要用途是用于请求日志记录。 除了函数签名中提供的响应代码之外，在带 ContextKeyResponse 前缀的 key 下面的 context 中还提供了其他响应参数。

## type ServerOption

```go
type ServerOption func(*Server)
```

ServerOption 为服务器设置一个可选参数。

## func ServerAfter

```go
func ServerAfter(after ...ServerResponseFunc) ServerOption
```

ServerAfter 函数在端点被调用后，向客户端写入任何内容前，在 HTTP 响应写入器上执行。

## func ServerBefore

```go
func ServerBefore(before ...RequestFunc) ServerOption
```

在请求被解码之前，对 HTTP 请求对象执行 ServerBefore 函数。

## func ServerErrorEncoder

```go
func ServerErrorEncoder(ee ErrorEncoder) ServerOption
```

每当在处理请求时遇到错误，就将它们编码到 `http.ResponseWriter` 。 客户端可以使用这个来提供自定义的错误格式和响应代码。 默认情况下，错误将使用 DefaultErrorEncoder 编写。

## func ServerErrorHandler

```go
func ServerErrorHandler(errorHandler transport.ErrorHandler) ServerOption
```

ServerErrorHandler 用于处理非终端(non-terminal)错误。 默认情况下，会忽略非终端错误。 这是作为一种诊断措施。 更细粒度的错误处理控制，包括更详细的日志记录，应该在定制的 ServerErrorEncoder 或 ServerFinizer 中执行，这两者都可以访问context。

## func ServerErrorLogger

```go
func ServerErrorLogger(logger log.Logger) ServerOption
```

ServerErrorLogger 用于记录非终端错误。 默认情况下，不会记录任何错误。 这是作为一种诊断措施。 更细粒度的错误处理控制，包括更详细的日志记录，应该在定制的 ServerErrorEncoder 或 ServerFinizer 中执行，这两者都可以访问上下文。 **不推荐: 改用 ServerErrorHandler**。

## func ServerFinalizer

```go
func ServerFinalizer(f ...ServerFinalizerFunc) ServerOption
```

ServerFinizer 在每个 HTTP 请求结束时执行。 默认情况下，没有注册Finalizer。

## type ServerResponseFunc

```go
type ServerResponseFunc func(context.Context, http.ResponseWriter) context.Context
```

ServerResponseFunc 可以从请求 context 中获取信息，并使用它来操作 ResponseWriter。 ServerResponseFuncs 只在服务器中执行，在写响应之前，调用端点之后执行。

## func SetContentType

```go
func SetContentType(contentType string) ServerResponseFunc
```

SetContentType 返回一个 ServerResponseFunc，该 ServerResponseFunc 将 Content-Type 头设置为给定的值。

## func SetResponseHeader

```go
func SetResponseHeader(key, val string) ServerResponseFunc
```

SetResponseHeader 返回设置为给定头的 ServerResponseFunc。

## type StatusCoder

```go
type StatusCoder interface {
    StatusCode() int
}
```

StatusCoder 由 DefaultErrorEncoder 检查, 如果一个错误值实现了 StatusCoder，则在编码错误时将使用 StatusCode。 默认情况下，使用 StatusInternalServerError (500)。
