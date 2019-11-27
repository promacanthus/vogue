# package http

```go
import "net/http"
```

http包提供HTTP客户端和服务端实现。

Get、Head、Post和PostForm发起HTTP（或HTTPS）请求：

```go
resp , err := http.Get("http://example.com/")
...
resp, err := http.Post("http://example.com/upload","image/jpeg",&buf)
...
resp, err := http.PostForm("http://example.com/form",
    url.Values{"Key":{"Value"}, "id":{"123"}})
```

结束后，客户端必须关闭响应体：

```go
resp, err := http.Get("http://example.com")
if err != nil {
    // handle error
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
```

想要控制HTTP客户端的头、重定向策略和其他设置，需要创建一个`Client`：

```go
client := &http.Client{
    CheckRedirect: redirectPolicyFunc,
}

resp, err := client.Get("http://example.com")
// ...

req, err := http.NewRequest("GET","http://example.com",nil)
// ...
req.Header.Add("if-None-Match",`W/"wyzzy"`)
resp, err := client.Do(req)
// ...
```

想要控制代理，TLS配置，keep-alive，压缩和其他设置，需要创建一个`Transport`：

```go
tr := &http.Transport{
    MaxldleConns: 10,
    ldleConnTimeout: 30 * time.Second,
    DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```

`Client`和`Transport`可以安全地被多个`goroutine`并发使用，为了高效应该只创建一次并重复使用。

`ListenAndServe`使用给定地址和处理程序启动一个HTTP服务端。这个处理程序通常是`nil`，这表示使用`DefaultServeMux`。`Handle`和`HandleFunc`将处理程序添加到`DefaultServeMux`:

```go
http.Handle("/foo"，fooHandler)

http.HandleFunc("/bar",func(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello,%q", html.EscapeString(r.URL, Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))
```

创建一个自定义的Server来更好的控制server的行为：

```go
s := &http.Server{
    Addr :              ":8080",
    Handler:            myHandler,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1<< 20,
    }

log.Fatal(s.ListenAndServe())
```

从Go 1.6开始，使用HTTPS时，http包对HTTP/2协议具有透明支持。

必须禁用HTTP/2的程序可以通过设置`Transport.TLSNextProto`(客户端设置)或者`Server.TLSNextProto`(服务端设置)的值为非nil的空map。或者，当前支持如下的`GODEBUG`环境变量：

```go
GODEBUG=http2client=0 // 禁用HTTP/2客户端支持
GODEBUG=http2server=0 // 禁用HTTP/2服务端支持
GODEBUG=http2debug=1 // 启动详细的HTTP/2 调试日志
GODEBUG=http2debug=2 // 更详细的日志，包含 frame dumps
```

Go的API兼容性保证不涵盖GODEBUG变量。在禁用HTTP/2支持之前，请报告所有问题：`https://golang.org/s/http2bug`

为了简化配置， http 包中的`Transport`和`Server`都自动启动HTTP/2支持。要为更复杂的配置启用HTTP/2，以使用低级别HTTP/2功能或使用Go的http2软件包的新版本，请直接导入 `golang.org/x/net/http2` 并使用其`ConfigureTransport`或`ConfigureServer`功能。通过`golang.org/x/net/http2`包手动配置HTTP/2优先于`net/http`包的内置HTTP/2支持。

## Constants

```go
const(
    MethodGet       =       "GET"
    MethodHead     =       "HEAD"
    MethodPost      =       "POST"
    MethodPut       =        "PUT"
    MethodPath      =       "PATCH"     // RFC 5789
    MethodDelete    =       "DELETE"
    MethodConnect   =       "CONNECT"
    MethodOptions   =       "OPTIONS"
    MethodTrace        =        "TRACE"
)
```

常见的HTTP方法。

除非另有说明，否则它们在[RFC 7231第4.3节](http://tools.ietf.org/html/rfc7231#section-4.3)中定义。

```go
const (
    StatusContinue           = 100 // RFC 7231, 6.2.1
    StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
    StatusProcessing         = 102 // RFC 2518, 10.1
    StatusEarlyHints         = 103 // RFC 8297

    StatusOK                   = 200 // RFC 7231, 6.3.1
    StatusCreated              = 201 // RFC 7231, 6.3.2
    StatusAccepted             = 202 // RFC 7231, 6.3.3
    StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
    StatusNoContent            = 204 // RFC 7231, 6.3.5
    StatusResetContent         = 205 // RFC 7231, 6.3.6
    StatusPartialContent       = 206 // RFC 7233, 4.1
    StatusMultiStatus          = 207 // RFC 4918, 11.1
    StatusAlreadyReported      = 208 // RFC 5842, 7.1
    StatusIMUsed               = 226 // RFC 3229, 10.4.1

    StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
    StatusMovedPermanently = 301 // RFC 7231, 6.4.2
    StatusFound            = 302 // RFC 7231, 6.4.3
    StatusSeeOther         = 303 // RFC 7231, 6.4.4
    StatusNotModified      = 304 // RFC 7232, 4.1
    StatusUseProxy         = 305 // RFC 7231, 6.4.5

    StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
    StatusPermanentRedirect = 308 // RFC 7538, 3

    StatusBadRequest                   = 400 // RFC 7231, 6.5.1
    StatusUnauthorized                 = 401 // RFC 7235, 3.1
    StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
    StatusForbidden                    = 403 // RFC 7231, 6.5.3
    StatusNotFound                     = 404 // RFC 7231, 6.5.4
    StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
    StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
    StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
    StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
    StatusConflict                     = 409 // RFC 7231, 6.5.8
    StatusGone                         = 410 // RFC 7231, 6.5.9
    StatusLengthRequired               = 411 // RFC 7231, 6.5.10
    StatusPreconditionFailed           = 412 // RFC 7232, 4.2
    StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
    StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
    StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
    StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
    StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
    StatusTeapot                       = 418 // RFC 7168, 2.3.3
    StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
    StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
    StatusLocked                       = 423 // RFC 4918, 11.3
    StatusFailedDependency             = 424 // RFC 4918, 11.4
    StatusTooEarly                     = 425 // RFC 8470, 5.2.
    StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
    StatusPreconditionRequired         = 428 // RFC 6585, 3
    StatusTooManyRequests              = 429 // RFC 6585, 4
    StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
    StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

    StatusInternalServerError           = 500 // RFC 7231, 6.6.1
    StatusNotImplemented                = 501 // RFC 7231, 6.6.2
    StatusBadGateway                    = 502 // RFC 7231, 6.6.3
    StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
    StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
    StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
    StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
    StatusInsufficientStorage           = 507 // RFC 4918, 11.5
    StatusLoopDetected                  = 508 // RFC 5842, 7.2
    StatusNotExtended                   = 510 // RFC 2774, 7
    StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
```

注册在IANA上的HTTP状态码，查看： `https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml`

```go
const DefaultMaxHeaderBytes = 1 << 20 // 1 MB
```

`DefaultMaxHeaderBytes`是HTTP请求头的最大允许大小。可以通过设置`Server.MaxHeaderBytes`来覆盖它。

```go
const DefaultMaxIdleConnsPerHost = 2
```

`DefaultMaxIdleConnsPerHost`是`Transport`的`MaxIdleConnsPerHost`的默认值。

```go
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
```

`TimeFormat`是在HTTP标头中生成时间时要使用的时间格式。就像`time.RFC1123`一样，但是将GMT硬编码为时区。格式化时间必须采用UTC格式才能生成正确的格式。

有关解析此时间格式的信息，请参见`ParseTime`。

```go
const TrailerPrefix = "Trailer:"
```

`TrailerPrefix`是`ResponseWriter.Header`的map的键的前缀（如果存在的话），表示map条目实际上是用于响应尾部的，而不是响应头的。`ServeHTTP`调用完成后，前缀将被删除，并且值将在尾部中发送。

此机制仅适用于在写入标头之前未知的尾部。如果在写标头之前固定的或已知的尾部组，则首选普通的Go尾部机制：

[<https://golang.org/pkg/net/http/#ResponseWriter>](https://golang.org/pkg/net/http/#ResponseWriter)

[<https://golang.org/pkg/net/http/#example_ResponseWriter_trailers>](https://golang.org/pkg/net/http/#example_ResponseWriter_trailers)

## Variables

```go
var (
    // ErrNotSupported 是由Pusher的Push方法返回的，表示不可获得对HTTP/2的支持
    ErrNotSupported = &ProtocolError{"feature not supported"}

    // 弃用
    // net/http包中的任何方法都不会再返回 ErrUnexpectedTrailer
    // 调用者不应该将错误与这个变量进行比较
    ErrUnexpectedTrailer = &ProtocolError{"trailer header without chunked transfer encoding"}

    // 当请求的 Content-Type 不包含一个“boundary”参数时，Request.MultipartReader 返回 ErrMissingBoundary
    ErrMissingBoundary = &ProtocolError{"no multipart boundary param in Content-Type"}

    // 当请求的 Content-Type 不是 multipart/form-data 时，Request.MultipartReader 返回 ErrNotMultipart
    ErrNotMultipart = &ProtocolError{"request Content-Type isn't multipart/form-data"}

    // 弃用
    // net/http包中的任何方法都不会再返回 ErrHeaderTooLong
    // 调用者不应该将错误与这个变量进行比较
    ErrHeaderTooLong = &ProtocolError{"header too long"}

    // 弃用
    // net/http包中的任何方法都不会再返回 ErrShortBody
    // 调用者不应该将错误与这个变量进行比较
    ErrShortBody = &ProtocolError{"entity body too short"}

    // 弃用
    // net/http包中的任何方法都不会再返回 ErrMissingContentLength
    // 调用者不应该将错误与这个变量进行比较
    ErrMissingContentLength = &ProtocolError{"missing ContentLength in HEAD response"}
)
```

```go
var (
    // HTTP方法或响应代码不允许body时，ResponseWriter.Write调用返回ErrBodyNotAllowed
    ErrBodyNotAllowed = errors.New("http: request method or response status code does not allow body")

    // 当使用Hijacker接口时基础连接被劫持，ResponseWriter.Write 返回 ErrHijacked
    // 在被劫持的连接上进行零字节写入将返回ErrHijacked，而不会产生任何其他副作用
    ErrHijacked = errors.New("http: connection has been hijacked")

    // 当 Handler 为Content-Length响应头设置一个长度，然后尝试写入比声明更多的字节时，
    // ResponseWriter.Write调用返回ErrContentLength
    ErrContentLength = errors.New("http: wrote more than the declared Content-Length")

    // 弃用
    // net/http包中的任何方法都不会再返回 ErrWriteAfterFlush
    // 调用者不应该将错误与这个变量进行比较
    ErrWriteAfterFlush = errors.New("unused")
)
```

HTTP服务端使用的错误。

```go
var (
    // ServerContextKey是 context key。可以在带有 Context.Value 的 HTTP handler 中
    // 使用它来访问启动 handler 的服务器，关联的值将是*Server类型
    ServerContextKey = &contextKey{"http-server"}

    // LocalAddrContextKey是context key。 可以在带有Context.Value的HTTP handler中使用它来访问连接的本地地址
    // 关联的值将为net.Addr类型。
    LocalAddrContextKey = &contextKey{"local-addr"}
)
```

```go
var DefaultClient = &Client{}
```

DefaultClient是默认客户端，由Get，Head和Post使用。

```go
var DefaultServeMux = &defaultServeMux
```

DefaultServeMux是Serve使用的默认ServeMux。

```go
var ErrAbortHandler = errors.New("net/http: abort Handler")
```

`ErrAbortHandler`是一个中止panic值，用于中止handler。来自ServeHTTP的任何panic都会中止对客户端的响应，但是使用`ErrAbortHandler`进行panic也会抑制将堆栈跟踪记录到服务器的错误日志中。

```go
var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")
```

在关闭请求或响应主体后，读取请求或响应主体时，将返回`ErrBodyReadAfterClose`。通常在HTTP处理程序在其`ResponseWriter`上调用`WriteHeader`或`Write`后读取正文时，会发生这种情况。

```go
var ErrHandlerTimeout = errors.New("http: Handler timeout")s
```

在已超时的处理程序上进行`ResponseWriter`的Write调用会返回`ErrHandlerTimeout`。

```go
var ErrLineTooLong = internal.ErrLineTooLong
```

读取分块编码异常的请求或响应体时，将返回`ErrLineTooLong`。

```go
var ErrMissingFile = errors.New("http: no such file")
```

当提供的文件字段不存在于请求中或者它不是一个文件字段中时，FormFile返回`ErrMissingFile`。

```go
var ErrNoCookie = errors.New("http: named cookie not present")
```

未找到Cookie时，Request的Cookie方法将返回`ErrNoCookie`。

```go
var ErrNoLocation = errors.New("http: no Location header in response")
```

如果不存在Location标头，则由Response的Location方法返回`ErrNoLocation`。

```go
var ErrServerClosed = errors.New("http: Server closed")
```

调用Shutdown或Close之后，服务器的Serve，ServeTLS，ListenAndServe和ListenAndServeTLS方法返回`ErrServerClosed`。

```go
var ErrSkipAltProtocol = errors.New("net/http: skip alternate protocol")
```

`ErrSkipAltProtocol`是由`Transport.RegisterProtocol`定义的标记错误值。

```go
var ErrUseLastResponse = errors.New("net/http: use last response")
```

`Client.CheckRedirect` hooks 可以返回`ErrUseLastResponse`，以控制如何处理重定向。如果返回，则不发送下一个请求，并且返回最近的响应，且其主体未关闭。

```go
var NoBody = noBody{}
```

`NoBody`是一个没有字节的`io.ReadCloser`。读取始终返回EOF，而关闭始终返回nil。可以在传出客户端请求中使用它来明确表示请求的字节数为零。也可以很容易的将`Request.Body`设置为nil。

## func CanonicalHeaderKey

```go
func CanonicalHeaderKey(s string) string
```

`CanonicalHeaderKey`返回header key `s` 的规范格式。规范化将第一个字母和连字符后的任何字母转换为大写；其余的将转换为小写。例如，"accept-encoding"的规范格式是"Accept-Encoding"。如果 `s` 包含空格或无效的header field bytes，则返回它而无需进行任何修改。

## func DetectContentType

```go
func DetectContentType(data []byte) string
```

`DetectContentType`实现在`https://mimesniff.spec.whatwg.org/`上描述的算法，以确定给定数据的`Content-Type`。它最多考虑前512个字节的数据。 `DetectContentType`始终返回有效的`MIME`类型：如果无法确定更具体的类型，则返回`"application/octet-stream"`。

## func Error

```go
func Error(w ResponseWriter, error string, code int)
```

Error回复带有特定错误消息和HTTP代码的请求。否则，它不会结束请求；调用者应确保不再对`w`进行写操作。错误消息应为纯文本。

## func Handle

```go
func Handle(pattern string, handler Handler)
```

`Handle`在`DefaultServeMux`中注册给定模式的处理程序。`ServeMux`的文档说明了如何匹配模式。

### Handle Example

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
)

type countHander struct{
    mu sync.Mutex // guards n
    n int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    h.mu.Lock()
    defer h.mu.Unlock()
    h.n++
    fmt.Fprintf(w, "count is %d\n", h.n)
}

func main(){
    http.Handle("/count",new(countHandler))
    log.Fatal(http.ListenAndServe(":8080",nil))
}
```

## func HandleFunc

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

`HandleFunc`在`DefaultServeMux`中注册给定模式的处理函数。`ServeMux`的文档说明了如何匹配模式。

### HandleFunc Example

```go
h1 := func(w http.ResponseWriter, _ *http.Request){
    io.WriteString(w, "Hello from a HandleFunc #1!\n")
}

h2 := func(w http.ResponseWriter, _ *http.Request){
    io.WriteString(w,"Hello from a HandleFunc #2!\n")
}

http.HandleFunc("/", h1)
http.HandleFunc("/endpoint", h2)

log.Fatal(http.ListenAndServe(":8080",nil))
```

## func ListenAndServe

```go
func ListenAndServe(addr string, handler Handler) error
```

`ListenAndServe`监听TCP网络地址`addr`，然后调用带有处理程序的Serve来处理传入连接上的请求。能够接受的连接需要配置为启用TCP keep-alives。

该处理程序通常为nil，在这种情况下，将使用`DefaultServeMux`。

`ListenAndServe`始终返回非nil错误。

### ListenAndServe Example

```go
helloHandler := func(w http.ResponseWriter, req *http.Request){
    io.WriteString(w,"Hello,world!\n")
}

http.HandlerFunc("/hello",helloHandler)
log.Fatal(http.ListenAndServe(":8080",nil))
```

## func ListenAndServeTLS

```go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
```

`ListenAndServeTLS`的行为与`ListenAndServe`相同，不同之处在于它需要`HTTPS`连接。此外，必须提供包含服务器证书和匹配私钥的文件。如果证书是由证书颁发机构签名的，则`certFile`应该是服务器证书，任何中间件和CA证书的串联。

### ListenAndServeTLS Example

```go
http.HandlerFunc("/",func(w http.ResponseWriter, req *http.Request){
    io.WriteString(w,"Hello,TLS!\n")
})

// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem
log.Printf("About to listen on 8443. Go to https://127.0,0.1:8443/")
err := http.ListenAndServeTLS(":8443","cert.pem","key.pem",nil)
log.Fatal(err)
```

## func MaxBytesReader

```go
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser
```

`MaxBytesReader`与`io.LimitReader`相似，但旨在限制传入请求体的大小。与`io.LimitReader`相比，`MaxBytesReader`的结果为`ReadCloser`，对于超出限制的Read返回非EOF错误，并在调用其Close方法时关闭底层reader。

`MaxBytesReader`可以防止客户端意外或恶意发送大请求并浪费服务器资源。

## func NotFound

```go
func NotFound(w ResponseWriter, r *Request)
```

`NotFound`回复请求，并显示HTTP 404 not found错误。

## func ParseHTTPVersion

```go
func ParseHTTPVersion(vers string) (major, minor int, ok bool)
```

`ParseHTTPVersion`解析HTTP版本字符串。 `"HTTP/1.0"`返回`(1,0,true)`。

## funcParseTime

```go
func ParseTime(text string) (t time.Time, err error)
```

`ParseTime`解析时间标头（例如`Date: header`），尝试使用HTTP/1.1允许的三种格式：`TimeFormat`，`time.RFC850`和`time.ANSIC`。

## func ProxyFromEnvironment

```go
func ProxyFromEnvironment(req *Request) (*url.URL, error)
```

`ProxyFromEnvironment`返回用于给定请求的代理的URL，如环境变量`HTTP_PROXY`，`HTTPS_PROXY`和`NO_PROXY`（或其小写版本）所指示。**对于HTTPS请求，HTTPS_PROXY优先于HTTP_PROXY**。

环境值可以是完整的URL或`"host [:port]"`，在这种情况下，假定使用"http"协议。如果值的格式不同，则返回错误。

如果环境中未定义任何代理，则返回nil URL和nil错误，或者不应将代理用于这样的请求（如`NO_PROXY`所定义）。

在特殊情况下，如果`req.URL.Host`为`"localhost"`（有或没有端口号），则将返回nil URL和nil错误。

## func ProxyURL

```go
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error)
```

`ProxyURL`返回一个代理函数（在Transport中使用），该函数始终返回相同URL。

## func Redirect

```go
func Redirect(w ResponseWriter, r *Request, url string, code int)
```

`Redirect`通过重定向到URL来响应请求，URL可能是相对于请求路径的路径。

返回的状态码应在3xx范围内，通常为`StatusMovedPermanently`（301），`StatusFound`（302）或`StatusSeeOther`(303）。

如果尚未设置`Content-Type`标头，`Redirect`会将其设置为 `"text / html; charset = utf-8"` 并编写一个小的HTML。将`Content-Type`标头设置为任何值（包括nil）将禁用该行为。

## func Serve

```go
func Serve(l net.Listener, handler Handler) error
```

`Serve`接受侦听器`l`上传入的HTTP连接，从而为每个侦听器创建一个新的服务`goroutine`。服务`goroutine`读取请求，然后调用处理程序以回复请求。该处理程序通常为nil，在这种情况下，将使用`DefaultServeMux`。

仅当侦听器返回`* tls.Conn`连接并且在`TLS Config.NextProtos`中将它们配置为`"h2"`时，才启用HTTP/2支持。`Serve`始终返回非nil错误。

## func ServeContent

```go
func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)
```

`ServeContent`使用提供的`ReadSeeker`中的内容回复请求。`ServeContent`优于`io.Copy`的主要好处是它可以正确处理`Range`请求，设置`MIME`类型并处理`If-Match`，`If-Unmodified-Since`，`If-None-Match`，`If-Modified-Since`和`If-Range` 请求。

如果未设置响应的`Content-Type`标头，则`ServeContent`首先尝试从名称的文件扩展名中推断出类型，如果失败，则退回到读取内容的第一块并将其传递给`DetectContentType`。否则该名称未使用。特别是它可以为空，并且永远不会在响应中发送。

如果`modtime`不是零时间或Unix时期，则`ServeContent`会将其包含在响应的`Last-Modified`头中。如果请求中包含`If-Modified-Since`标头，则`ServeContent`使用`modtime`来决定是否需要发送内容。

内容的Seek方法必须起作用：ServeContent使用对内容结尾的查找来确定其大小。

如果调用方设置了按照[RFC 7232第2.3节](http://tools.ietf.org/html/rfc7232#section-2.3)格式化的w的ETag标头，则ServeContent使用它来处理使用If-Match，If-None-Match或If-Range的请求。

请注意，`*os.File`实现了`io.ReadSeeker`接口。
