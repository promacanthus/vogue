---
title: package httputil
date: 2020-04-14T10:09:14.274627+08:00
draft: false
---

```go
import "net/http/httputil"
```

httputil包提供HTTP实用函数，补充`net/http`包中更常见的函数。

## Variables

```go
var ErrLineTooLong = internal.ErrLineTooLong
```

当读取格式不正确的带有过长行的块数据时，将返回`ErrLineTooLong`。

## func DumpRequest

```go
func DumpRequest(req *http.Request, body bool) ([]byte, error)
```

DumpRequest以`HTTP/1.x`连接形式返回给定的请求。 它只应该被服务端用来调试客户端请求。 返回值只是一个近似值; 在将初始请求解析为 `http.Request` 时，会丢失初始请求的一些细节。 特别是，头字段名称的顺序和大小写都会丢失。 多值标头中值的顺序保持不变。 `HTTP/2`请求以 `HTTP/1.x` 的形式被转储，而不是原始的二进制表示形式。

如果 body 为 true，则 DumpRequest 也会返回 body。 为此，它获取 `req.Body` 然后用一个新的 `io.ReaderCloser` 替换它来生成相同字节的。 如果 DumpRequest 返回一个错误，那么说明 req 的状态是未定义的。

`http.Request.Write`文档详细说明了转储中包括哪些req字段。

### DumpRequest Example

```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    dump, err := httputil.DumpRequest(r, true)
    if err != nil {
        http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "%q", dump)
}))
defer ts.Close()

const body = "Go is a general-purpose language designed with systems programming in mind."
req, err := http.NewRequest("POST", ts.URL, strings.NewReader(body))
if err != nil {
    log.Fatal(err)
}
req.Host = "www.example.org"
resp, err := http.DefaultClient.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

b, err := ioutil.ReadAll(resp.Body)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%s", b)

// output
"POST / HTTP/1.1\r\nHost: www.example.org\r\nAccept-Encoding: gzip\r\nContent-Length: 75\r\nUser-Agent: Go-http-client/1.1\r\n\r\nGo is a general-purpose language designed with systems programming in mind."
```

## func DumpRequestOut

```go
func DumpRequestOut(resp *http.Response, body bool)([]byte, error)
```

DumpRequestOut类似于DumpRequest，但用于传出客户端请求。它包括标准`http.Transport`添加的任何标头，例如User-Agent。

### DumpRequestOut Example

```go
const body = "Go is a general-purpose language designed with systems programming in mind."
req, err := http.NewRequest("PUT", "http://www.example.org", strings.NewReader(body))
if err != nil {
    log.Fatal(err)
}

dump, err := httputil.DumpRequestOut(req, true)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%q", dump)

// Output
"PUT / HTTP/1.1\r\nHost: www.example.org\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 75\r\nAccept-Encoding: gzip\r\n\r\nGo is a general-purpose language designed with systems programming in mind."
```

## func DumpResponse

```go
func DumpResponse(resp *http.Response, body bool) ([]byte, error)
```

DumpResponse类似于DumpRequest，但转储响应。

### DumpResponse Example

```go
const body = "Go is a general-purpose language designed with systems programming in mind."
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Date", "Wed, 19 Jul 1972 19:00:00 GMT")
    fmt.Fprintln(w, body)
}))
defer ts.Close()

resp, err := http.Get(ts.URL)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

dump, err := httputil.DumpResponse(resp, true)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%q", dump)
// Output
"HTTP/1.1 200 OK\r\nContent-Length: 76\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 19 Jul 1972 19:00:00 GMT\r\n\r\nGo is a general-purpose language designed with systems programming in mind.\n"
```

## func NewChunkedReader

```go
func NewChunkedReader(r io.Reader) io.Reader
```

NewChunkedReader返回一个新的chunkedReader，该块将从r读取的数据转换为HTTP“分块”格式，然后再返回。读取最后的`0-length`的块时，chunkedReader返回`io.EOF`。

普通应用程序不需要NewChunkedReader。读取响应正文时，http包会​​自动解码分块。

## func NewChunkedWriter

```go
func NewChunkedWriter(w io.Writer) io.WriteCloser
```

NewChunkedWriter 返回一个新的 chunkedWriter，它将写入转换为 HTTP“分块”格式，然后写入到 w。关闭返回的chunkedWriter将发送标记为流结束的最后一个`0-length`的块，但不发送出现在trailer的最后一个CRLF；trailer和最后的CRLF必须分开写。

普通的应用程序不需要 NewChunkedWriter。 如果处理程序没有设置 Content-Length 标头，http 包会自动添加 chunking。 在处理程序中使用 NewChunkedWriter 将导致双重块或者Content-Length长度的分块，这两种情况都是错误的。

## type BufferPool

```go
type BufferPool interface{
    Get() []byte
    Put([]byte)
}
```

BufferPool是用于获取和返回由`io.CopyBuffer`使用的临时字节切片的接口。

## type ReverseProxy

```go
type ReverseProxy struct {
    // Director必须具有将请求修改为要使用Transport发送的新请求的函数。
    // 然后将其响应复制回未经修改的原始客户端。Director 返回后不得访问提供的请求。
    Diretor func(*http.Request)

    // Transport用于代理请求，如果是nil，那么使用 http.DefaultTransport
    Transport http.RoundTripper

    // FlushInterval指定在复制响应主体时要刷新到客户端的刷新间隔。
    // 如果为零，则不执行定期刷新。负值表示每次写入客户端后立即刷新。
    // 当ReverseProxy将响应识别为流响应时，将忽略FlushInterval。
    // 对于此类响应，立即将写操作刷新到客户端
    FlushInterval time.Duration

    // ErrorLog为尝试代理请求时发生的错误指定一个可选的记录器。
    // 如果为nil，则通过日志包的标准记录器完成记录。
    ErrorLog *log.Logger

    // BufferPool可以选择指定一个缓冲池，
    // 以在复制HTTP响应正文时获取io.CopyBuffer使用的字节切片。
    BufferPool BufferPool

    // ModifyResponse是一个可选函数，用于修改后端的响应。
    // 如果后端完全返回带有任何HTTP状态代码的响应，则调用该方法。
    // 如果后端不可达，则调用可选的ErrorHandler而不调用ModifyResponse。

    // 如果ModifyResponse返回错误，则会使用错误值调用ErrorHandler。
    // 如果ErrorHandler为nil，则使用其默认实现。
    ModifyResponse func(*http.Response) error

    // ErrorHandler是一个可选函数，用于处理到达后端的错误或ModifyResponse中的错误。
    // 如果为nil，则默认是记录所提供的错误并返回 502 Bad Gateway 响应。
    ErrorHandler func(http.ResponseWriter, *http.Request, error)
}
```

ReverseProxy是一个HTTP处理程序，它接收传入的请求并将其发送到另一台服务器，并将响应代理回客户端。

### ReverseProxy Example

```go
backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "this call was relayed by the reverse proxy")
}))
defer backendServer.Close()

rpURL, err := url.Parse(backendServer.URL)
if err != nil {
    log.Fatal(err)
}
frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(rpURL))
defer frontendProxy.Close()

resp, err := http.Get(frontendProxy.URL)
if err != nil {
    log.Fatal(err)
}

b, err := ioutil.ReadAll(resp.Body)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%s", b)

// Output

this call was relayed by the reverse proxy
```

## func NewSingleHostReverseProxy

```go
func NewSingleHostReverseProxy(target *url.URL) *ReverseProxy
```

NewSingleHostReverseProxy返回一个新的ReverseProxy，将URL路由到target中提供的scheme，host和base path。 如果目标的路径是`"/base"`，而传入的请求是`"/dir"`，则目标请求将是`"/base/dir"`。 NewSingleHostReverseProxy不会重写Host标头。 要重写Host标头，请直接将ReverseProxy与自定义Director一起使用。

## func(*ReverseProxy) ServeHTTP

```go
func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)
```
