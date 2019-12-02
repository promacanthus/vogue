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

`ListenAndServe`使用给定地址和handler启动一个HTTP服务端。这个handler通常是`nil`，这表示使用`DefaultServeMux`。`Handle`和`HandleFunc`将handler添加到`DefaultServeMux`:

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

在关闭请求或响应主体后，读取请求或响应主体时，将返回`ErrBodyReadAfterClose`。通常在HTTPhandler在其`ResponseWriter`上调用`WriteHeader`或`Write`后读取正文时，会发生这种情况。

```go
var ErrHandlerTimeout = errors.New("http: Handler timeout")s
```

在已超时的handler上进行`ResponseWriter`的Write调用会返回`ErrHandlerTimeout`。

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

`Handle`在`DefaultServeMux`中注册给定模式的handler。`ServeMux`的文档说明了如何匹配模式。

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

`ListenAndServe`监听TCP网络地址`addr`，然后调用带有handler的Serve来处理传入连接上的请求。能够接受的连接需要配置为启用TCP keep-alives。

该handler通常为nil，在这种情况下，将使用`DefaultServeMux`。

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

如果尚未设置`Content-Type`标头，`Redirect`会将其设置为 `"text/html; charset = utf-8"` 并编写一个小的HTML。将`Content-Type`标头设置为任何值（包括nil）将禁用该行为。

## func Serve

```go
func Serve(l net.Listener, handler Handler) error
```

`Serve`接受侦听器`l`上传入的HTTP连接，从而为每个侦听器创建一个新的服务`goroutine`。服务`goroutine`读取请求，然后调用handler以回复请求。该handler通常为nil，在这种情况下，将使用`DefaultServeMux`。

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

## func ServeFile

```go
func ServeFile(w ResponseWriter, r *Request, name string)
```

ServeFile使用命名文件或目录的内容答复请求。

如果提供的文件或目录名称是相对路径，则会相对于当前目录进行解释，并且可能会升至父目录。 如果提供的名称是根据用户输入构造的，则应在调用ServeFile之前对其进行清理。

作为预防措施，ServeFile将拒绝`r.URL.Path`包含“..”路径元素的请求； 这样可以防止调用者可能不安全地使用`filepath.Join`而不清理它，然后使用该`filepath.Join`结果作为name参数。

作为另一种特殊情况，ServeFile将`r.URL.Path`以“/index.html”结尾的任何请求重定向到同一路径，而没有最终的“index.html”。 为避免此类重定向，请修改路径或使用ServeContent。

除了这两种特殊情况外，ServeFile不使用`r.URL.Path`来选择要提供的文件或目录。 仅使用name参数中提供的文件或目录。

## func ServeTLS

```go
func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error
```

ServeTLS在侦听器`l`上接受传入的HTTPS连接，从而为每个侦听器创建一个新的服务goroutine。 服务goroutine读取请求，然后调用handler以回复请求。 

该handler通常为nil，在这种情况下，将使用DefaultServeMux。

此外，必须提供包含服务器证书和匹配私钥的文件。 如果证书是由证书颁发机构签名的，则certFile应该是服务器证书，任何中间件和CA证书的串联。 

ServeTLS始终返回非nil错误。

## func SetCookie

```go
func SetCookie(w ResponseWriter, cookie *Cookie)
```

SetCookie将Set-Cookie标头添加到提供的ResponseWriter的标头中。提供的cookie必须具有有效的名称。无效的cookie可能会被静默删除。

## func StatusText

```go
func StatusText(code int) string
```

StatusText返回HTTP状态代码的文本。如果状态码未知，它将返回空字符串。

## type Client

```go
type Client struct {
    // Transport 指定发出单个HTTP请求的机制。
    // 如果为nil，则使用DefaultTransport。
    Transport RoundTripper

    // CheckRedirect指定用于处理重定向的策略。 如果CheckRedirect不为nil，则客户端将在执行 HTTP重定向之前调用它。
    // 参数req和via是即将到来的请求和已发出的请求，先到达的先发出。
    // 如果CheckRedirect返回错误，则客户端的Get方法将返回先前的Response（关闭该响应体）和CheckRedirect的错误（包装在url.Error中），而不是发出Request请求。
    // 作为一种特殊情况，如果CheckRedirect返回ErrUseLastResponse，则返回最近的响应，且其响应体未关闭，并返回nil错误。 
    // 如果CheckRedirect为nil，则客户端使用其默认策略，该策略将在连续10个请求后停止。
    CheckRedirect func(req *Request, via []*Request) error

    // Jar 指定cookie jar。Jar用于将相关cookie插入每个出站请求，并使用每个入站Response的cookie值进行更新。
    // 客户端遵循的每个重定向都会咨询Jar。 如果Jar为nil，则仅当在Request上显式设置cookie时，才发送cookie。
    Jar CookieJar

    // Timeout 指定此客户端发出的请求的时间限制。 超时包括连接时间，任何重定向和读取响应正文。
    // 在Get，Head，Post或Do返回之后，计时器保持运行状态，并且将中断Response.Body的读取。
    // Timeout为零表示没有超时。 客户端取消对基础传输的请求，就像请求的context结束一样。
    // 为了兼容性，如果找到，客户端还将在Transport上使用已经弃用的CancelRequest方法。
    // 新的RoundTripper实现应使用请求的context进行取消，而不是实现CancelRequest。
    Timeout time.Duration
}
```

Client是HTTP客户端。它的零值（DefaultClient）是使用DefaultTransport的可用客户端。

客户端的传输通常具有内部状态（缓存的TCP连接），因此应重用客户端，而不是根据需要创建客户端。客户端可以安全地被多个goroutine并发使用。

客户端比RoundTripper（例如Transport）更高级别，并且还处理HTTP详细信息，例如cookie和重定向。

执行重定向时，客户端将转发在初始请求上设置的所有标头，但以下情况除外：

- 将诸如“Authorization”，“ WWW-Authenticate”和“ Cookie”之类的敏感标头转发到不受信任的目标时。当重定向到与子域不匹配或与初始域不完全匹配的域时，将忽略这些标头。例如，从“ foo.com”重定向到“ foo.com”或“ sub.foo.com”将转发敏感标头，但重定向到“ bar.com”则不会。
- 用非零值的Cookie Jar转发“ Cookie”标头时。由于每个重定向可能会更改Cookie Jar的状态，因此重定向可能会更改初始请求中设置的Cookie。当转发“ Cookie”标头时，任何突变的cookie都将被省略，并期望Jar将插入具有更新值的那些突变的cookie（假设原点匹配）。如果Jar为零，则将转发原始cookie，而不进行任何更改。

## func (*Client) CloseIdleConnections

```go
func (c *Client) CloseIdleConnections()
```

CloseIdleConnections关闭其Transport上先前请求建的现在处于“keep-alive”状态的所有连接。 它不会中断当前正在使用的任何连接。

如果客户端的Transport没有CloseIdleConnections方法，则此方法不执行任何操作。

## func (*Client) Do

```go
func (c *Client) Do(req *Request) (*Response, error)
```

Do 按照客户端上配置的策略（例如重定向，Cookie，身份验证）发送HTTP请求并返回HTTP响应。

如果是由客户端策略（例如CheckRedirect）或无法连接HTTP（例如网络连接问题）引起的，则返回错误。非2xx状态代码不会引起错误。

如果返回的错误为nil，则响应将包含一个非nil的响应体，用户希望将其关闭。如果未同时将主体读入EOF并关闭，则客户端的底层RoundTripper（通常是Transport）可能无法将与服务器的持久TCP连接重新用于后续的“keep-alive”请求。

请求主体（如果非nil）将被底层Transport关闭，即使发生错误也是如此。

出错时，任何响应都可以忽略。仅当CheckRedirect失败并且返回的`Response.Body`已关闭时，才会出现具有非nil错误的非nil响应。

通常，将使用Get，Post或PostForm代替Do。

如果服务器回复重定向，则客户端首先使用CheckRedirect函数来确定是否应遵循重定向。如果允许，则301、302或303重定向会导致后续请求使用HTTP方法GET（如果原始请求为HEAD，则为HEAD），而没有请求体。如果定义了`Request.GetBody`函数，则307或308重定向将保留原始的HTTP方法和主体。 NewRequest函数自动为常见的标准库主体类型设置GetBody。

返回的任何错误均为`*url.Error`类型。如果请求超时或被取消，则`url.Error`值的Timeout方法将报告true。

## func (*Client) Get

```go
func (c *Client) Get(url string) (resp *Response, err error)
```

Get将GET发送到指定的URL。如果响应是以下重定向代码之一，则Get在调用客户端的CheckRedirect函数后执行重定向：

```go
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
```

如果客户端的CheckRedirect函数失败或存在HTTP协议错误，则返回错误。 非2xx响应不会导致错误。 返回的任何错误均为`*url.Error`类型。 如果请求超时或被取消，则`url.Error`值的Timeout方法将报告true。

当err为nil时，resp始终包含非nil的 `resp.Body` 。调用者完成读取后，应关闭`resp.Body`。

要使用自定义标头发出请求，请使用NewRequest和`Client.Do`。

## func (*Client) Head

```go
func (c *Client) Head(url string) (resp *Response, err error)
```

Head向指定的URL发出HEAD。如果响应是以下重定向代码之一，则Head在调用客户端的CheckRedirect函数后执行重定向：

```go
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
```

## func (*Client) Post

```go
func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error)
```

Post发布POST到指定的URL。 调用者完成读取后，应关闭`resp.Body`。 如果提供的body是`io.Closer`，则在请求后将其关闭。 若要设置自定义标头，请使用NewRequest和`Client.Do`。 有关如何处理重定向的详细信息，请参阅`Client.Do`方法文档。

## func (*Client) PostForm

```go
func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)
```

PostForm向指定的URL发出POST，并将数据的键和值URL编码为请求正文。 Content-Type标头设置为`application/x-www-form-urlencoded`。 若要设置其他标头，请使用NewRequest和`Client.Do`。 当err为nil时，resp始终包含非nil `resp.Body`，调用者完成读取后，应关闭`resp.Body`。 有关如何处理重定向的详细信息，请参阅`Client.Do`方法文档。

## type CloseNotifier

```go
type CloseNotifier interface {
    // 当客户端连接断开时，CloseNotify返回一个通道，该通道最多接收单个值（true）。
    // 在完全读取Request.Body之前，CloseNotify可能会等待通知。
    // handler返回后，不能保证该通道会收到一个值。
    // 如果协议是HTTP/1.1，并且在使用HTTP/1.1管道处理幂等请求（例如GET）时调用了CloseNotify，则后续管道请求的到来可能会导致在返回的通道上发送值。
    // 实际上，HTTP/1.1管道未在浏览器中启用，并且很不常见。 如果这是一个问题，请使用HTTP/2或仅对诸如POST之类的方法使用CloseNotify。
    CloseNotify() <-chan bool
}
```

## type ConnState

```go
type ConnState int
```

ConnState表示客户端与服务器的连接状态。由可选的`Server.ConnState`挂钩使用。

```go
const (
    // StateNew表示一个新连接，该连接应立即发送请求。 连接从此状态开始，然后过渡到StateActive或StateClosed。
    StateNew ConnState = iota

    // StateActive表示已读取一个或多个请求字节的连接。 StateActive的Server.ConnState钩子在请求进入handler之前触发，直到请求被处理后才再次触发。
    // 处理请求后，状态将转换为StateClosed，StateHijacked或StateIdle。 对于HTTP / 2，StateActive在从零到一个活动请求的转换上触发，并且仅在所有活动请求完成后才转换。 
    // 这意味着ConnState不能用于执行每个请求的预处理工作； ConnState仅记录连接的整体状态。
    StateActive

    // StateIdle表示已完成处理请求并处于keep-alive状态的连接，正在等待新请求。 连接从StateIdle转换为StateActive或StateClosed。
    StateIdle

    // StateHijacked 表示被劫持的连接， 这是状态的终态，它不会过渡到StateClosed。
    StateHijacked

    // StateClosed表示已关闭的连接。 这是状态的终态， 被劫持的连接不会过渡到StateClosed。
    StateClosed
)
```

## func (ConnState) String

```go
func (c ConnState) String() string
```

## type Cookie

```go
type Cookie struct {
    Name  string
    Value string

    Path       string    // optional
    Domain     string    // optional
    Expires    time.Time // optional
    RawExpires string    // for reading cookies only

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int
    Secure   bool
    HttpOnly bool
    SameSite SameSite
    Raw      string
    Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

Cookie代表在HTTP响应的Set-Cookie标头或HTTP请求的Cookie标头中发送的HTTP Cookie，更多详情看[这里](https://tools.ietf.org/html/rfc6265)。

## func (*Cookie) String

```go
func (c *Cookie) String() string 
```

String返回用于Cookie头（如果仅设置了Name和Value）或Set-Cookie响应头（如果设置了其他字段）的cookie的序列化。如果`c`为nil或`c.Name`无效，则返回空字符串。

## type CookieJar

```go
type CookieJar interface {
    // SetCookies在给定URL的回复中处理cookie的接收。
    // 它可能会或可能不会选择保存cookie，具体取决于jar的策略和实现。
    SetCookies(u *url.URL, cookies []*Cookie)

    // Cookies返回cookie，以发送对给定URL的请求。
    // 具体实现取决于标准的cookie使用限制，例如RFC 6265。
    Cookies(u *url.URL) []*Cookie
}
```

CookieJar管理HTTP请求中cookie的存储和使用。 CookieJar的实现必须安全，可以被多个goroutine并发使用。 `net/http/cookiejar`软件包提供了CookieJar实现。

## type Dir

```go
type Dir string
```

Dir使用限于特定目录树的原生文件系统来实现FileSystem。

尽管`FileSystem.Open`方法采用`'/'`分隔的路径，但Dir的字符串值是原生文件系统上的文件名，而不是URL，因此它由`filepath.Separator`分隔，不一定是`'/'`。

请注意，Dir将允许以句点开头（隐藏文件）的文件和目录访问，这可能会公开敏感目录（如`.git`目录）或敏感文件（如`.htpasswd`）。 要排除隐藏文件，请从服务器中删除文件/目录或创建自定义FileSystem实现。

空目录将被视为`“.”`。

## func (Dir) Open

```go
func (d Dir) Open(name string) (File, error)
```

Open使用`os.Open`实现FileSystem，打开文件以读取root以及相对于目录d的文件。

## type File

```go
type File interface {
    io.Closer
    io.Reader
    io.Seeker
    Readdir(count int) ([]os.FileInfo, error)
    Stat() (os.FileInfo, error)
}
```

File由FileSystem的Open方法返回，并且可以由FileServer实现提供服务。 

这些方法的行为应与`*os.File`上的行为相同。

## type FileSystem

```go
type FileSystem interface {
    Open(name string) (File, error)
}
```

FileSystem实现对命名文件集合的访问。无论主机操作系统的约定如何，文件路径中的元素都用斜杠（`'/'`，U+002F）字符分隔。

## type Flusher

```go
type Flusher interface {
    // Flush 刷新会将所有缓冲的数据发送到客户端。
    Flush()
}
```

Flusher接口由ResponseWriters实现，它允许HTTPhandler将缓冲的数据刷新到客户端。

默认的`HTTP/1.x`和`HTTP/2 ResponseWriter`实现支持Flusher，但ResponseWriter包装器可能不支持。 handler应始终在运行时测试此功能。

请注意，即使对于支持Flush的ResponseWriters，如果客户端通过HTTP代理连接，在响应完成之前，缓冲的数据也可能无法到达客户端。

## type Handler

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Handler响应HTTP请求。

ServeHTTP应将响应标头和数据写入ResponseWriter，然后返回。返回信号表明请求已完成；在ServeHTTP调用完成后或与之同时使用ResponseWriter或从`Request.Body`中读取是无效的。

根据HTTP客户端软件，HTTP协议版本以及客户端和Go服务器之间的任何中介，可能无法在写入ResponseWriter之后从`Request.Body`中读取。谨慎的handler应先读取`Request.Body`，然后再进行回复。

除读取正文外，handler不应修改提供的请求。

如果出现ServeHTTP出现运行时恐慌（painc），则服务端（ServeHTTP的调用方）将假定panic的影响与处于活跃状态的请求无关。它将会恢复panic，将堆栈跟踪记录到服务器错误日志中，然后关闭网络连接或发送`HTTP/2 RST_STREAM`，具体取决于HTTP协议。要中止handler，以便客户端看到中断的响应，而服务器不会记录错误，那么panic应该带有ErrAbortHandler值。

## func FileServer

```go
func FileServer(root FileSystem) Handler
```

FileServer返回一个handler，该handler以根目录为根文件系统内容为HTTP请求提供服务。

要使用操作系统的文件系统实现，请使用`http.Dir`：

```go
http.Handle("/", http.FileServer(http.Dir("/tmp")))
```

作为一种特殊情况，返回的文件服务器会将以`“/index.html”`结尾的所有请求重定向到同一路径，而没有最终的`“ index.html”`。

### Example

```go
// 简单的静态Web服务器
log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))
```

### Example(DotFileHiding)

```go
package http_test

import (
    "log"
    "net/http"
    "os"
    "strings"
)

// containsDotFile报告name是否包含以句点开头的path元素。
// 假定name由正斜杠分隔，由http.FileSystem接口保证。
func containsDotFile(name string) bool {
    parts := strings.Split(name, "/")
    for _, part := range parts {
        if strings.HasPrefix(part, ".") {
            return true
        }
    }
    return false
}

// dotFileHidingFile是dotFileHidingFileSystem中的http.File使用。
// 它用于包装http.File的Readdir方法，以便我们可以从其输出中删除以句点开头的文件和目录。
type dotFileHidingFile struct {
    http.File
}

// Readdir是嵌入式File的Readdir方法的包装，可过滤掉名称中以句点开头的所有文件。
func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
    files, err := f.File.Readdir(n)
    for _, file := range files { // 过滤点文件
        if !strings.HasPrefix(file.Name(), ".") {
            fis = append(fis, file)
        }
    }
    return
}

// dotFileHidingFileSystem是一个http.FileSystem，可隐藏隐藏的“点文件”，使其不再提供服务。
type dotFileHidingFileSystem struct {
    http.FileSystem
}

// Open是嵌入式FileSystem的Open方法的包装，
// 当name具有名称以其路径开头的文件或目录时，它将包装403权限错误。
func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
    if containsDotFile(name) { // 如果是点文件, 返回 403
        return nil, os.ErrPermission
    }

    file, err := fs.FileSystem.Open(name)
    if err != nil {
        return nil, err
    }
    return dotFileHidingFile{file}, err
}

func ExampleFileServer_dotFileHiding() {
    fs := dotFileHidingFileSystem{http.Dir(".")}
    http.Handle("/", http.FileServer(fs))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## func NotFoundHandler

```go
func NotFoundHandler() Handler
```

NotFoundHandler返回一个简单的handler，该handler以``404页面未找到''答复来回复每个请求。

### Example（NotFoundHandler）

```go
mux := http.NewServeMux()

// 创建示例handler以返回404
mux.Handle("/resources", http.NotFoundHandler())

// 创建示例handler以返回200
mux.Handle("/resources/people/", newPeopleHandler())

log.Fatal(http.ListenAndServe(":8080", mux))
```

## func RedirectHandler

```go
func RedirectHandler(url string, code int) Handler
```

RedirectHandler返回一个handler，该handler使用给定的状态码将收到的每个请求重定向到给定的url。 

提供的代码应在3xx范围内，通常为StatusMovedPermanently，StatusFound或StatusSeeOther。

## func StripPrefix

```go
func StripPrefix(prefix string, h Handler) Handler
```

StripPrefix通过从请求URL的路径中删除给定前缀并调用handler h，返回处理HTTP请求的handler。 StripPrefix通过回复HTTP 404 not found错误来处理不以前缀开头的路径的请求。

### Example（StripPrefix）

```go
// 要在备用URL路径（/tmpfiles/）下为磁盘（/tmp）上的目录提供服务，请在FileServer看到之前使用StripPrefix修改请求URL的路径：
http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))
```

## func TimeoutHandler

```go
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler
```

TimeoutHandler返回一个在给定时间限制下运行h的Handler。

新的handler调用`h.ServeHTTP`来处理每个请求，但是如果调用的运行时间超过其时间限制，则该handler将以503 Service Unavailable错误和响应体中的给定消息作为响应。 （如果msg为空，则将发送适当的默认消息。）在这样的超时之后，用`h`对其ResponseWriter进行写操作将返回ErrHandlerTimeout。

TimeoutHandler支持Flusher和Pusher接口，但不支持Hijacker接口。

## type HandlerFunc

```go
type HandlerFunc func(ResponseWriter, *Request)
```

HandlerFunc类型是一个适配器，允许将普通函数用作HTTP处理程序。如果`f`是具有适当签名的函数，则HandlerFunc(f) 是调用f的处理程序。

## func (HandlerFunc) ServeHTTP

```go
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request)
```

ServeHTTP 调用 f(w, r)。

## type Header

```go
type Header map[string][]string
```

Header代表HTTP标头中的键/值对。 key应采用CanonicalHeaderKey返回的规范形式。

## func (Header) Add

```go
func (h Header) Add(key, value string)
```

Add将键值对添加到Header。它附加到与键关联的任何现有值。key不区分大小写；它由CanonicalHeaderKey规范化。

## func (Header) Clone

```go
func (h Header) Clone() Header
```

如果`h`为nil，则Clone返回`h`或nil的副本。

## func (Header) Del

```go
func (h Header) Del(key string)
```

Del删除与键关联的值。key不区分大小写；它由CanonicalHeaderKey规范化。

## func (Header) Get

```go
func (h Header) Get(key string) string
```

Get获取与给定键关联的第一个值。 如果没有与键关联的值，则Get返回`""`。 不区分大小写； `textproto.CanonicalMIMEHeaderKey`用于规范化提供的key。 要访问键的多个值或使用非规范键，请直接访问map。

## func (Header) Set

```go
func (h Header) Set(key, value string)
```

Set将与key关联的header条目设置为单个元素值。 它替换了与键关联的任何现有值。 key不区分大小写； 它由`textproto.CanonicalMIMEHeaderKey`规范化。 要使用非规范键，请直接分配给map。

## func (Header) Write

```go
func (h Header) Write(w io.Writer) error
```

Write以传输格式写header。

## func (Header) WriteSubset

```go
func (h Header) WriteSubset(w io.Writer, exclude map[string]bool) error
```

WriteSubset以传输格式写入header。如果exclude不为nil，则不会写入exclude [key] == true的键。

## type Hijacker

```go
type Hijacker interface {
    // Hijack使调用者可以接管连接。
    // 调用Hijack之后，HTTP服务器库将不会对该连接执行任何其他操作。

    // 管理和关闭连接成为调用者的责任。

    // 根据服务器的配置，返回的net.Conn可能已设置了读取或写入的期限。
    // 调用者有责任根据需要设置或清除这些截止时间。

    // 返回的bufio.Reader可能包含来自客户端的未处理的缓冲数据。

    // 调用Hijack之后，不得使用原始的Request.Body。
    // 原始请求的context仍然有效，并且直到该请求的ServeHTTP方法返回后才被取消。
    Hijack() (net.Conn, *bufio.ReadWriter, error)
}
```

Hijacker接口是由ResponseWriters实现的，它允许HTTP hander接管连接。

`HTTP/1.x`连接的默认ResponseWriter支持Hijacker，但`HTTP/2`连接有意设置为不支持。 ResponseWriter包装器也可能不支持Hijacker。 handler应始终在运行时测试此功能。

### Example（Hijacker）

```go
http.HandleFunc("/hijack", func(w http.ResponseWriter, r *http.Request) {
    hj, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
        return
    }
    conn, bufrw, err := hj.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // 不要忘记关闭连接
    defer conn.Close()
    bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
    bufrw.Flush()
    s, err := bufrw.ReadString('\n')
    if err != nil {
        log.Printf("error reading string: %v", err)
        return
    }
    fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
    bufrw.Flush()
})
```

## type PushOptions

```go
type PushOptions struct {
    // Method为承诺的请求指定HTTP方法。
    // 如果设置，则必须为“ GET”或“ HEAD”。空表示“ GET”。
    Method string

    // Header指定其他承诺的请求标头。
    // 这不能包括HTTP/2伪标头字段，例如“:path”和“:scheme”，它们会自动添加。
    Header Header
}
```

PushOptions描述了Pusher.Push的选项。

## type Pusher

```go
type Pusher interface {
    // Push启动HTTP/2服务器来推送。 这将使用给定的目标和选项构造一个综合请求，
    // 将该请求序列化为PUSH_PROMISE帧，然后使用服务器的handler分派该请求。 
    // 如果opts为nil，则使用默认选项。

    // target必须是绝对路径（例如“ /path”）或包含有效host和与父请求相同的scheme的绝对URL。
    // 如果目标是路径，它将继承父请求的scheme和host。

    // HTTP/2规范不允许递归推送和交叉授权推送。
    // 推送可能会或可能不会检测到这些无效的推送； 但是，合格的客户端将检测并取消无效的推送。

    // 希望推送URL X的handler应在发送任何可能触发URL X请求的数据之前调用Push。
    // 这避免了客户端在收到X的PUSH_PROMISE之前发出X请求的竞争。

    // Push将在单独的goroutine中运行，从而使到达顺序不确定。 调用者需要实现任何所需的同步。

    // 如果客户端禁用了Push或底层连接不支持Push，则Push返回ErrNotSupported。
    Push(target string, opts *PushOptions) error
}
```

Pusher是由ResponseWriters实现的支持`HTTP/2`服务器推送的接口。有关更多背景信息，请参见[这里](https://tools.ietf.org/html/rfc7540#section-8.2)。

## type Request

```go
type Request struct {
    // Method指定HTTP方法（GET，POST，PUT等）。
    // 对于客户端请求，空字符串表示GET。

    // Go的HTTP客户端不支持使用CONNECT方法发送请求。 
    // 有关详细信息，请参见Transport文档。
    Method string

    // URL指定要请求的URI（对于服务端请求）或要访问的URL（对于客户端请求）。
    // 对于服务端请求，将从存储在RequestURI中的请求行上提供的URI解析URL。
    // 对于大多数请求，Path和RawQuery以外的其他字段将为空。 （请参阅RFC 7230，第5.3节）

    // 对于客户端请求，URL的host指定要连接的服务器，
    // 而请求的host字段则可选地指定要在HTTP请求中发送的host标头值。
    URL *url.URL

    // 传入服务器请求的协议版本。
    // 对于客户请求，将忽略这些字段。HTTP客户端代码始终使用HTTP/1.1或HTTP/2。
    // 有关详细信息，请参见Transport文档。
    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0

    // Header包含服务端接收或客户端发送的请求标头字段。

    // 如果服务端收到带有标头行的请求：
    // Host: example.com
    // accept-encoding: gzip, deflate
    // Accept-Language: en-us
    // fOO: Bar
    // foo: two

    //  Header = map[string][]string{
    //      "Accept-Encoding": {"gzip, deflate"},
    //      "Accept-Language": {"en-us"},
    //      "Foo": {"Bar", "two"},
    //  }

    // 对于传入的请求，Host Header被提升为Request.Host字段并从 Header映射中删除。

    // HTTP定义Header名称不区分大小写。 请求解析器通过使用CanonicalHeaderKey来实现此目的，
    // 使第一个字符以及连字符后的所有字符都变为大写，其余的都变为小写。

    // 对于客户端请求，某些header（例如Content-Length和Connection）会在需要时自动写入，
    // 并且标头中的值可能会被忽略。 请参阅文档中的Request.Write方法。
    Header Header

    // Body 是请求体。

    // 对于客户端请求，nil body表示该请求没有请求体，例如GET请求。
    // HTTP客户端的transport负责调用Close方法。

    // 对于服务端请求，请求体始终为non-nil，但在不存在请求体时将立即返回EOF。
    // 服务端将关闭请求体。 ServeHTTP handler不需要这么做。
    Body io.ReadCloser

    // GetBody定义了一个可选的func来返回Body的新副本。
    // 它用于当客户端请求重定向需要多次读取正文时。使用GetBody仍然需要设置Body。 

    // 对于服务端请求，它是不使用的。

    GetBody func() (io.ReadCloser, error)


    // ContentLength记录关联内容的长度。
    // 值-1表示长度未知。
    // 值>=0表示可以从请求体读取给定的字节数。
    // 对于客户请求，non-nil Body 的值为0也会被认为是未知。
    ContentLength int64

    // TransferEncoding列出从最外层到最内层的传输编码。
    // 空列表表示“身份”编码。 通常可以忽略TransferEncoding。
    // 发送和接收请求时，将根据需要自动添加和删除分块编码。
    TransferEncoding []string

    // Close 指示是在回复此请求后（对于服务器）还是在发送此请求并读取其响应（对于客户端）之后关闭连接。

    //  对于服务端请求，HTTP服务器会自动处理此请求，并且处理程序不需要此字段。
    // 对于客户端请求，设置此字段可防止在相同主机的请求之间重复使用TCP连接，
    // 就像设置了Transport.DisableKeepAlives一样。
    Close bool

    // 对于服务器请求，Host指定在其上搜索URL的host。
    //  根据RFC 7230，第5.4节，这是“host”标头的值或URL本身中提供的主机名。
    // 它的形式可能是“ host:port”。 对于国际域名，Host可以采用Punycode或Unicode形式。
    // 如果需要，请使用golang.org/x/net/idna将其转换为两种格式。

    // 为了防止DNS重新绑定攻击，服务器handler应验证主机标头是否具有其认为自己具有权威性的值。
    // 随附的ServeMux支持注册到特定host名的模式，从而保护其注册的handler。
    // 对于客户端请求，host可以选择覆盖要发送的主机头。 如果为空，则Request.Write方法使用URL.Host的值。 
    // 主机可能包含国际域名。
    Host string

    // 表单包含已解析的表单数据，包括URL字段的查询参数和PATCH，POST或PUT表单数据。
    // 该字段仅在调用ParseForm之后可用。 HTTP客户端会忽略Form，而使用Body。
    Form url.Values

    // PostForm包含从PATCH，POST或PUT正文参数解析的表单数据。
    // 该字段仅在调用ParseForm之后可用。 HTTP客户端会忽略PostForm并改用Body。
    PostForm url.Values

    // MultipartForm是已解析的多部分表单，包括文件上传。
    // 仅在调用ParseMultipartForm之后，此字段才可用。 HTTP客户端会忽略MultipartForm并改用Body。
    MultipartForm *multipart.Form

    // Trailer指定在请求正文之后发送的其他标头。
    // 对于服务器请求，Trailer map 最初仅包含 Trailer 键，其值为nil。
    // （客户端声明它将稍后发送的Trailer。）handler从Body读取时，它不得引用 Trailer。
    // 从Body读取返回EOF后，Trailer可以再次读取，并且包含非null值（如果它们是由客户端发送的）。
    //  对于客户端请求，必须将Trailer初始化为包含Trailer key 的 map，以便以后发送。
    // 这些值可以为nil或它们的最终值。 ContentLength必须为0或-1，才能发送分块的请求。
    // 发送HTTP请求后，可以在读取请求正文的同时更新map值。 一旦正文返回EOF，调用方就不得使Trailer发生改变。 
    // 很少有HTTP客户端，服务器或代理支持HTTP Trailer。
    Trailer Header

    // RemoteAddr允许HTTP服务器和其他软件记录发送请求的网络地址，通常用于日志记录。
    // ReadRequest不会填写此字段，并且没有定义的格式。
    // 在调用handler之前，此程序包中的HTTP服务器将RemoteAddr设置为“ IP:Port”地址。 HTTP客户端将忽略此字段。
    RemoteAddr string

    // RequestURI是客户端发送到服务器的请求行（RFC 7230，第3.1.1节）的未经修改的请求目标。
    // 通常应改用URL字段。 在HTTP客户端请求中设置此字段是错误的。
    RequestURI string

    // TLS允许HTTP服务器和其他软件记录有关在其上接收到请求的TLS连接的信息。ReadRequest不会填写此字段。
    // 此程序包中的HTTP服务器在调用handler之前为启用TLS的连接设置字段。
    // 否则，该字段将为零。 HTTP客户端会忽略此字段。
    TLS *tls.ConnectionState

    // Cancel是一个可选通道，当它关闭则表示客户端的请求应被视为已取消。 
    // 并非RoundTripper的所有实现都支持取消。
    // 对于服务器请求，此字段不适用。
    // 弃用：而是使用NewRequestWithContext设置请求的context。
    // 如果同时设置了请求的“cancel”字段和context，则不确定是否遵守“cancel”。
    Cancel <-chan struct{}

    // Response是导致创建此请求的重定向响应。仅在客户端重定向期间填充此字段。
    Response *Response
}
```

Request表示服务器接收或客户端发送的HTTP请求。

客户端和服务器使用情况之间的字段语义略有不同。 除了以下字段上的注释外，请参阅`Request.Write`和RoundTripper的文档。

## func NewRequest

```go
func NewRequest(method, url string, body io.Reader) (*Request, error)
```

NewRequest使用background context包装NewRequestWithContext。

## func NewRequestWithContext

```go
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error)
```

给定Method，URL和可选Body，NewRequestWithContext返回一个新的请求。

如果提供的body同时也是`io.Closer`，则返回的`Request.Body`设置为body，并且将通过客户端方法Do，Post和PostForm以及`Transport.RoundTrip`关闭。

NewRequestWithContext返回适合与`Client.Do`或`Transport.RoundTrip`一起使用的请求。若要创建用于测试服务器handler的请求，请使用`net/http/httptest`程序包中的NewRequest函数，使用ReadRequest，或手动更新Request字段。对于传出的客户端请求，context控制请求的整个生命周期及其响应：获取连接，发送请求以及读取响应标头和正文。有关入站和出站请求字段之间的区别，请参见Request type的文档。

如果body的类型为`*bytes.Buffer`，`*bytes.Reader`或`*strings.Reader`，则返回的请求的ContentLength设置为其确切值（而不是-1），将填充GetBody（因此307和308重定向可以重播body），如果ContentLength为0，则将Body设置为NoBody。

## func ReadRequest

```go
func ReadRequest(b *bufio.Reader) (*Request, error)
```

ReadRequest读取并解析来自`b`的传入请求。

ReadRequest是一个底层功能，仅应用于专业应用程序。大多数代码应使用服务器读取请求并通过Handler接口处理请求。 ReadRequest仅支持`HTTP/1.x`请求。 对于`HTTP/2`，请使用`golang.org/x/net/http2`。

## func (*Request) AddCookie

```go
func (r *Request) AddCookie(c *Cookie)
```

AddCookie将cookie添加到请求中。根据RFC 6265第5.4节，AddCookie不会附加多个Cookie header字段。这意味着所有cookie（如果有的话）都写在同一行中，并用分号分隔。

## func (*Request) BasicAuth

```go
func (r *Request) BasicAuth() (username, password string, ok bool)
```

如果请求使用HTTP基本认证，则BasicAuth返回请求的授权header中提供的用户名和密码。请参阅RFC 2617，第2节。

## func (*Request) Clone

```go
func (r *Request) Clone(ctx context.Context) *Request
```

Clone 返回`r`的深层副本，其context更改为ctx。 提供的ctx必须为非零。

对于传出的客户端请求，context控制请求的整个生命周期及其响应：获取连接，发送请求以及读取响应header和body。

## func (*Request) Context

```go
func (r *Request) Context() context.Context
```

Context返回请求的context。 若要更改context，请使用WithContext。

返回的context始终为非零； 它默认为 background context。 对于传出的客户端请求，context控制取消。

对于传入的服务器请求，当客户端的连接关闭，请求被取消（使用`HTTP/2`）或ServeHTTP方法返回时，context将被取消。

## func (*Request) Cookie

```go
func (r *Request) Cookie(name string) (*Cookie, error)
```

Cookie返回请求中提供的命名cookie，如果未找到，则返回ErrNoCookie。如果多个Cookie与给定名称匹配，则仅返回一个Cookie。

## func (*Request) Cookies

```go
func (r *Request) Cookies() []*Cookie
```

Cookies解析并返回与请求一起发送的HTTP cookie。

## func (*Request) FormFile

```go
func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
```

FormFile返回提供的表单key的第一个文件。如果需要，FormFile调用ParseMultipartForm和ParseForm。

## func (*Request) FormValue

```go
func (r *Request) FormValue(key string) string
```

FormValue返回查询的命名组件的第一个值。 POST和PUT请求体的参数优先于URL查询字符串值。 FormValue在必要时调用ParseMultipartForm和ParseForm，并忽略这些函数返回的任何错误。 如果key不存在，则FormValue返回空字符串。 要访问同一key的多个值，请调用ParseForm，然后直接检查`Request.Form`。

## func (*Request) MultipartReader

```go
func (r *Request) MultipartReader() (*multipart.Reader, error)
```

如果这是`multipart/form-data`或`multipart/mixed `POST请求，则MultipartReader返回MIME分段阅读器，否则返回nil和错误。使用此函数代替ParseMultipartForm将请求体作为流处理。

## func (*Request) ParseForm

```go
func (r *Request) ParseForm() error
```

ParseForm填充`r.Form`和`r.PostForm`。

对于所有请求，ParseForm都会从URL解析原始查询并更新`r.Form`。

对于POST，PUT和PATCH请求，它还将请求体解析为一种形式，并将结果放入`r.PostForm`和`r.Form`中。 请求体参数优先于`r.Form`中的URL查询字符串值。

对于其他HTTP方法，或者当Content-Type不是`application/x-www-form-urlencoded`时，不会读取请求体，并且`r.PostForm`初始化为非null的空值。

如果请求体的大小尚未受到MaxBytesReader的限制，则大小上限为10MB。

ParseMultipartForm自动调用ParseForm。 ParseForm是幂等的。

## func (*Request) ParseMultipartForm

```go
func (r *Request) ParseMultipartForm(maxMemory int64) error
```

ParseMultipartForm将请求体解析为`multipart/form-data`。 整个请求体将被解析，其文件部分的总计maxMemory字节最多存储在内存中，其余部分存储在磁盘上的临时文件中。

如果需要，ParseMultipartForm调用ParseForm。 一次调用ParseMultipartForm之后，后续调用无效。

## func (*Request) PostFormValue

```go
func (r *Request) PostFormValue(key string) string
```

PostFormValue返回POST，PATCH或PUT请求体的命名组件的第一个值。 URL查询参数将被忽略。 如果需要，PostFormValue调用ParseMultipartForm和ParseForm，并忽略这些函数返回的任何错误。 如果键不存在，则PostFormValue返回空字符串。

## func (*Request) ProtoAtLeast

```go
func (r *Request) ProtoAtLeast(major, minor int) bool
```

ProtoAtLeast报告请求中使用的HTTP协议是否至少为`major.minor`。

## func (*Request) Referer

```go
func (r *Request) Referer() string
```

如果请求中发送了Referer，则返回Referering URL。

Referer来源网址拼写错误，就像请求本身一样，这是HTTP最早的错误。 该值也可以从Header map中提取为`Header [“Referer”]`; 将其作为方法使用的好处是，编译器可以诊断使用备用（正确的英语）拼写`req.Referrer()`的程序，但不能诊断使用`Header [“Referrer”]`的程序。

## func (*Request) SetBasicAuth

```go
func (r *Request) SetBasicAuth(username, password string)
```

SetBasicAuth将请求的Authorization header设置为使用HTTP基本身份验证以及提供的用户名和密码。

使用HTTP基本认证时，提供的用户名和密码未加密。

某些协议可能会对预转义用户名和密码提出其他要求。 例如，当与OAuth2一起使用时，两个参数都必须首先使用`url.QueryEscape`进行URL编码。

## func (*Request) UserAgent

```go
func (r *Request) UserAgent() string
```

如果在请求中发送了UserAgent，则返回客户端的User-Agent。

## func (*Request) WithContext

```go
func (r *Request) WithContext(ctx context.Context) *Request
```

WithContext返回`r`的浅表副本，其context更改为ctx。 提供的ctx必须为非零。

对于传出的客户端请求，context控制请求的整个生命周期及其响应：获取连接，发送请求以及读取响应头和响应体。 

要使用context创建新请求，请使用NewRequestWithContext。 要更改请求（例如传入）的context，您还需要修改以发送出去，请使用`Request.Clone`。 在这两种用途之间，很少需要WithContext。

## func (*Request) Write

```go
func (r *Request) Write(w io.Writer) error
```

Write以传输格式写入`HTTP/1.1`请求，即请求头和请求体。此方法查询请求的以下字段：

```go
Host
URL
Method (defaults to "GET")
Header
ContentLength
TransferEncoding
Body
```

如果存在请求体，则Content-Length为<= 0，并且TransferEncoding尚未设置为“identity”，Write将“ Transfer-Encoding: chunked”添加到请求头，发送后将其关闭。

## func (*Request) WriteProxy

```go
func (r *Request) WriteProxy(w io.Writer) error
```

WriteProxy类似于Write，但是以HTTP代理期望的形式写入请求。 特别是，WriteProxy根据RFC 7230的5.3节（包括scheme和host）用绝对URI写入请求的初始Request-URI行。无论哪种情况，WriteProxy都会使用`r.Host`或`r.URL.Host`写入Host header。

## type Response

```go
type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0

    // Header将header key映射到value。
    // 如果响应中的多个header具有相同的键，则可以使用逗号分隔符将它们连接在一起。
    // （RFC 7230第3.2.2节要求，多个header在语义上等效于逗号分隔的序列。）
    // 当header值被该结构中的其他字段（例如ContentLength，TransferEncoding，Trailer）复制时，
    // 这些字段值是权威的。

    // map中的key已规范化（请参见CanonicalHeaderKey）。
    Header Header

    // Body代表响应体。

    // 在读取“body”字段时，将按需流式传输响应体。 如果网络连接失败或服务器终止响应，则Body.Read调用将返回错误。

    // http 的Client 和Transport 保证即使在没有响应体的响应或长度为零的响应中，响应体也始终为非零。
    // 关闭响应体是调用者的责任。 如果未读取并关闭响应体，则默认的HTTP客户端的传输可能不会重用HTTP/1.x“keep-alive” TCP连接。

    // 如果服务器回复了“分块”传输编码，则响应体将自动分块。

    // 从Go 1.12开始，Body还将在成功的“101交换协议”响应上实现io.Writer，
    // 该响应由WebSocket和HTTP/2的“h2c”模式使用。
    Body io.ReadCloser

    // ContentLength记录关联内容的长度。 值-1表示长度未知。
    // 除非Request.Method为“HEAD”，否则值>=0表示可以从Body中读取给定的字节数。
    ContentLength int64

    // 包含从最外部到最内部的传输编码。 值为nil，表示使用“identity”编码。
    TransferEncoding []string

    // Close记录在读取响应体后，header是否指示关闭连接。 该值是给客户的建议：ReadResponse和Response.Write都不会关闭连接。
    Close bool

    // Uncompressed报告响应是否已压缩发送，并已被http包解压缩。
    // 如果为true，则从Body读取将产生未压缩的内容，而不是从服务器实际设置的压缩内容，
    // ContentLength设置为-1，并且从responseHeader中删除“Content-Length”和“Content-Encoding”字段。
    // 要从服务器获取原始响应，请将Transport.DisableCompression设置为true。
    Uncompressed bool

    // Trailer将trailer键映射到与header相同格式的值。
    // Trailer最初仅包含nil值，服务器“Trailer”header值中指定的每个键对应一个值。这些值不会添加到Header中。
    // 不得在响应体上与Read调用同时访问Trailer。 在Body.Read返回io.EOF之后，Trailer将包含服务器发送的所有trailer值。
    Trailer Header
    // Request是为获取此响应而发送的请求。 请求体为零（已被消耗）。 仅针对客户端请求填充。
    Request *Request

    // TLS包含有关在其上接收到响应的TLS连接的信息。 对于未加密的响应，它为nil。
    // 指针在响应之间共享，不应修改。
    TLS *tls.ConnectionState
}
```

响应表示来自HTTP请求的响应。 一旦收到响应头，Client和Transport将从服务器返回响应。在读取“Body”字段时，将按需流式传输响应体。

## func Get

```go
func Get(url string) (resp *Response, err error)
```

Get将GET发送到指定的URL。如果响应是以下重定向代码之一，则Get将遵循该重定向，最多10个重定向：

```go
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
```

如果重定向太多或HTTP协议错误，则返回错误。 非2xx响应不会导致错误。 返回的任何错误均为`*url.Error`类型。 如果请求超时或被取消，则`url.Error`值的Timeout方法将报告true。

当err为nil时，resp始终包含一个非nil `resp.Body`。 调用者完成读取后，应关闭`resp.Body`。 Get是`DefaultClient.Get`的包装。 要使用自定义header发出请求，请使用NewRequest和`DefaultClient.Do`。

## Example（Get）

```go
res, err := http.Get("http://www.google.com/robots.txt")
if err != nil {
    log.Fatal(err)
}
robots, err := ioutil.ReadAll(res.Body)
res.Body.Close()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s", robots)
```

## func Head

```go
func Head(url string) (resp *Response, err error)
```

Head向指定的URL发出HEAD。如果响应是以下重定向代码之一，则Head遵循该重定向，最多10个重定向：

```go
301 (Moved Permanently)
302 (Found)
303 (See Other)
307 (Temporary Redirect)
308 (Permanent Redirect)
```

Head是`DefaultClient.Head`的包装器。

## func Post

```go
func Post(url, contentType string, body io.Reader) (resp *Response, err error)
```

Post发布POST到指定的URL。 调用者完成读取后，应关闭`resp.Body`。 如果提供的body是`io.Closer`，则在请求后将其关闭。

Post是`DefaultClient.Post`的包装。 要设置自定义header，请使用NewRequest和`DefaultClient.Do`。 有关如何处理重定向的详细信息，请参阅`Client.Do`方法文档。

## func PostForm

```go
func PostForm(url string, data url.Values) (resp *Response, err error)
```

PostForm向指定的URL发出POST，并将数据的键和值URL编码为请求体。 Content-Type header设置为`application/x-www-form-urlencoded`。 要设置其他header，请使用NewRequest和`DefaultClient.Do`。

当err为nil时，resp始终包含一个非nil `resp.Body`。 调用者完成读取后，应关闭`resp.Body`。 PostForm是`DefaultClient.PostForm`的包装。 有关如何处理重定向的详细信息，请参阅`Client.Do`方法文档。

## func ReadResponse

```go
func ReadResponse(r *bufio.Reader, req *Request) (*Response, error)
```

ReadResponse从`r`读取并返回HTTP响应。 req参数可选地指定与此响应相对应的请求。 如果为零，则假定为GET请求。

读取完`resp.Body`后，客户端必须调用`resp.Body.Close`。 调用之后，客户端可以检查trailer，以查找响应tailer中包含的键/值对。

## func (*Response) Cookies`

```go
func (r *Response) Cookies() []*Cookie
```

Cookies解析并返回Set-Cookie header中设置的cookie。

## func (*Response) Location

```go
func (r *Response) Location() (*url.URL, error)
```

Location返回响应的“location”header（如果存在）的URL。相对重定向是相对于响应的请求来解决的。如果不存在Location header，则返回ErrNoLocation。

## func (*Response) ProtoAtLeast

```go
func (r *Response) ProtoAtLeast(major, minor int) bool
```

ProtoAtLeast报告响应中使用的HTTP协议是否至少为`major.minor`。

## func (*Response) Write

```go
func (r *Response) Write(w io.Writer) error
```

Write以`HTTP/1.x`服务器响应格式将`r`写入`w`，包括状态行，header，body和可选的trailer。 此方法查询响应`r`的以下字段：

```go
StatusCode
ProtoMajor
ProtoMinor
Request.Method
TransferEncoding
Trailer
Body
ContentLength
Header, values for non-canonical keys will have unpredictable behavior
```

发送后，响应体将关闭。

## type ResponseWriter

```go
type ResponseWriter interface {
    // Header返回将由WriteHeader发送的header map。header map也是handler可以用来设置HTTP trailer的机制。

    // 除非修改后的header是trailer，否则在调用WriteHeader（或Write）后更改header map无效。

    // 设置trailer有两种方法。 首选方法是在header中预先声明稍后将发送的trailer，
    // 方法是将“Trailer”header设置为稍后将出现的键的名称。
    // 在这种情况下，Header映射的那些键被视为trailer。参见示例。
    // 第二种方法是，对于直到第一次处理之后才被handler所知的trailer键，
    // 在Header映射键之前加上TrailerPrefix常量值。 请参阅TrailerPrefix。

    // 要取消自动响应header（例如“Date”），请设置它们的值为nil。
    Header() Header

    // Write将数据作为HTTP回复的一部分写入连接。

    // 如果尚未调用WriteHeader，则Write在写入数据之前会调用WriteHeader（http.StatusOK）。
    // 如果header不包含Content-Type行，则Write将Content-Type集添加到将初始512字节的写入数据传递到DetectContentType的结果中。
    // 此外，如果所有写入数据的总大小小于几KB，并且没有Flush调用，则会自动添加Content-Length header。
    // 根据HTTP协议版本和客户端，调用Write或WriteHeader可能会阻止将来对Request.Body进行读取。
    // 对于HTTP/1.x请求，handler应在写入响应之前读取所有需要的请求体数据。
    // 刷新header后（由于显式的Flusher.Flush调用或写入足够的数据以触发刷新），请求体可能不可用。
    // 对于HTTP/2请求，Go HTTP服务器允许handler在同时写入响应的同时继续读取请求体。
    // 但是，并非所有的HTTP/2客户端都支持这种行为。
    // 如果可能，handler应在写入之前先进行读取，以最大程度地实现兼容性。
    Write([]byte) (int, error)

    // WriteHeader发送带有提供的状态代码的HTTP响应头。 如果未显式调用WriteHeader，则对Write的第一次调用将触发一个隐式WriteHeader（http.StatusOK）。
    // 因此，对WriteHeader的显式调用主要用于发送错误代码。 提供的代码必须是有效的HTTP 1xx-5xx状态代码。 只能写入一个header。
    // Go当前不支持发送用户定义的1xx信息性header，但服务器会在读取Request.Body时自动发送的100-xx响应头，但不支持发送该header。
    WriteHeader(statusCode int)
}
```

HTTP handler使用ResponseWriter接口构造HTTP响应。 返回`Handler.ServeHTTP`方法后，不得使用ResponseWriter。

## Example（Trailer）

HTTP Trailer是一组键/值对，例如header，位于HTTP响应之后而不是之前。

```go
mux := http.NewServeMux()
mux.HandleFunc("/sendstrailers", func(w http.ResponseWriter, req *http.Request) {
    // 在对WriteHeader或Write的任何调用之前，声明将在HTTP响应期间设置的Trailer。
    // 这三个header实际上是在Trailer中发送的。
    w.Header().Set("Trailer", "AtEnd1, AtEnd2")
    w.Header().Add("Trailer", "AtEnd3")

    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    w.WriteHeader(http.StatusOK)

    w.Header().Set("AtEnd1", "value 1")
    io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
    w.Header().Set("AtEnd2", "value 2")
    w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.
})
```

## type RoundTripper

```go
type RoundTripper interface {
    // RoundTrip执行一个HTTP事务，为提供的请求返回响应。

    // RoundTrip不应尝试解释响应。
    // 特别是，如果RoundTrip获得响应，则必须返回err == nil，
    // 而不管响应的HTTP状态代码如何。如果失败，则应保留非null错误。
    // 同样，RoundTrip不应尝试处理更高级别的协议详细信息，例如重定向，身份验证或cookie。

    // 除了消费和关闭请求体之外，RoundTrip不应修改请求。 RoundTrip可以在单独的goroutine中读取请求的字段。
    // 在响应体关闭之前，调用者不应更改或重用请求。

    // RoundTrip必须要关闭响应体，包括发生错误时，但根据实现的不同，
    // 即使在RoundTrip返回之后，也可能在单独的goroutine中关闭它。
    // 这意味着希望重用响应体以用于后续请求的调用者必须安排在等待Close调用之后再这样做。

    // 请求的URL和标头字段必须初始化。
    RoundTrip(*Request) (*Response, error)
}
```

RoundTripper是表示执行单个HTTP事务，获取给定请求的响应的能力的接口。RoundTripper必须安全，可以同时被多个goroutine使用。

```go
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}
```

DefaultTransport是Transport的默认实现，由DefaultClient使用。 它根据需要建立网络连接，并缓存它们以供后续调用重用。 它按照`$HTTP_PROXY`和`$NO_PROXY`（或`$http_proxy`和`$no_proxy`）环境变量的指示使用HTTP代理。

## func NewFileTransport

```go
func NewFileTransport(fs FileSystem) RoundTripper
```

NewFileTransport返回一个新的RoundTripper，服务于提供的FileSystem。 返回的RoundTripper会忽略传入请求中的URL host以及该请求的大多数其他属性。 NewFileTransport的典型用例是在TRansport中注册“file”协议，如下所示：

```go
t := &http.Transport{}
t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
c := &http.Client{Transport: t}
res, err := c.Get("file:///etc/passwd")
...
```

## type SameSite

```go
type SameSite int
```

SameSite允许服务器定义cookie属性，从而使浏览器无法将cookie与跨站点请求一起发送。 主要目标是减轻跨域信息泄漏的风险，并提供针对跨站点请求伪造攻击的某种保护。

更多信息点击[这里](https://tools.ietf.org/html/draft-ietf-httpbis-cookie-same-site-00)。

```go
const (
    SameSiteDefaultMode SameSite = iota + 1
    SameSiteLaxMode
    SameSiteStrictMode
    SameSiteNoneMode
)
```

## type ServeMux

```go
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    es    []muxEntry // slice of entries sorted from longest to shortest.
    hosts bool       // whether any patterns contain hostnames
}
```

ServeMux是一个HTTP请求多路复用器。它根据注册的模式列表将每个传入请求的URL匹配，并为与URL最匹配的模式调用处理程序。

模式命名固定的，有根的路径（例如`“/favicon.ico”`）或有根的子树（例如`“/images/”`）（请注意结尾的斜杠）。较长的模式优先于较短的模式，因此，如果同时为`“/images/”`和`“/images/thumbnails/”`注册了handler，则将为从`“/images/thumbnails/”`开始的路径调用后handler，将在`“/images/”`子树中接收对任何其他路径的请求。

请注意，由于以斜杠结尾的模式命名了一个有根的子树，因此模式`“/”`与所有其他已注册模式不匹配的路径匹配，而不仅仅是`Path==“/”`的URL。

如果已经注册了一个子树，并且接收到一个命名该子树根的请求而没有其后斜杠，则ServeMux将该请求重定向到该子树根（添加后斜杠）。可以用单独的路径注册来覆盖此行为，而不必使用斜杠。例如，注册`“/images/”`会使ServeMux将对`“/images”`的请求重定向到`“/images/”`，除非已单独注册了`“/images”`。

模式可以选择以主机名开头，仅将匹配项限制在该主机上。特定于主机的模式优先于常规模式，因​​此处理程序可以注册两个模式`“/codesearch”`和`“codesearch.google.com/”`，而不必同时接收对`“http://www.google.com/”`的请求”。

ServeMux还负责清理URL请求路径和Host header，清除端口号并重定向任何包含的请求`.`或`..`元素，或重复的斜杠表示为等效的，更简洁的URL。

## func NewServeMux

```go
func NewServeMux() *ServeMux
```

NewServeMux分配并返回一个新的ServeMux。

## func (*ServeMux) Handle

```go
func (mux *ServeMux) Handle(pattern string, handler Handler)
```

Handle注册给定模式的处理程序。如果已经存在用于模式的处理程序，则Handle发出运行时恐慌（panic）。

## Example

```go
mux := http.NewServeMux()
mux.Handle("/api/", apiHandler{})
mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    // “/”模式匹配所有内容，因此我们需要检查我们是否在root这里。
    if req.URL.Path != "/" {
        http.NotFound(w, req)
        return
    }
    fmt.Fprintf(w, "Welcome to the home page!")
})
```

## func (*ServeMux) HandleFunc

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

HandleFunc注册给定模式的处理函数。

## func (*ServeMux) Handler

```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
```

Handler通过访问`r.Method`，`r.Host`和`r.URL.Path`返回用于给定请求的处理程序。 它总是返回一个非nil handler。 如果路径的格式不规范，则该handler将是内部生成的handler，该handler将重定向到规范路径。 如果主机包含端口，则在匹配处理程序时将忽略该端口。

path和host不做变更直接用于CONNECT请求。

Handler还会返回与请求匹配的已注册模式，如果是内部生成的重定向，则返回在跟随重定向之后将匹配的模式。 如果没有适用于该请求的注册handler，则处理程序返回一个`“找不到页面'”`处理程序和一个空模式。

## func (*ServeMux) ServeHTTP

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

ServeHTTP将请求调度到其模式与请求URL最匹配的处理程序。

## type Server

```go
type Server struct {
    Addr    string  // 监听的TCP地址 , 如果为空，则为":http"
    Handler Handler // 要调用的处理程序，如果为空则为 http.DefaultServeMux

    // TLSConfig可选地提供TLS配置，以供ServeTLS和ListenAndServeTLS使用。
    // 请注意，此值由ServeTLS和ListenAndServeTLS克隆，
    // 因此无法使用tls.Config.SetSessionTicketKeys之类的方法修改配置。
    // 若要使用SetSessionTicketKeys，请改为将Server.Serve与TLS侦听器一起使用。
    TLSConfig *tls.Config

    // ReadTimeout是读取整个请求（包括请求体）的最大持续时间。
    // 由于ReadTimeout不允许处理程序根据每个请求体的可接受截止日期或上载速率做出每个请求的决定，
    // 因此大多数用户将更喜欢使用ReadHeaderTimeout。 两者都使用也是有效的。
    ReadTimeout time.Duration

    // ReadHeaderTimeout是允许读取请求头的时间量。
    // 读取请求体后，将重置连接的读取截止时间，并且Handler可以确定对主体而言太慢的速度。
    // 如果ReadHeaderTimeout为零，则使用ReadTimeout的值。 如果两者均为零，则没有超时。
    ReadHeaderTimeout time.Duration

    // WriteTimeout是写入响应超时之前的最大持续时间。
    // 每当读取新请求头时，都会将其重置。
    // 与ReadTimeout一样，它也不允许处理程序根据每个请求做出决策。
    WriteTimeout time.Duration

    // IdleTimeout是启用保持活动状态后等待下一个请求的最长时间。
    // 如果IdleTimeout为零，则使用ReadTimeout的值。 如果两者均为零，则没有超时。
    IdleTimeout time.Duration

    // MaxHeaderBytes控制服务器读取的最大字节数，以解析请求头的键和值（包括请求行）。
    // 它不限制请求体的大小。 如果为零，则使用DefaultMaxHeaderBytes。
    MaxHeaderBytes int

    // TLSNextProto可以选择指定一个函数，以在进行NPN/ALPN协议升级时接管所提供的TLS连接的所有权。
    // map的key是协商的协议名称。 Handler参数应用于处理HTTP请求，
    // 并将初始化请求的TLS和RemoteAddr（如果尚未设置）。 函数返回时，连接将自动关闭。
    // 如果TLSNextProto不为nil，则不会自动启用HTTP/2支持。
    TLSNextProto map[string]func(*Server, *tls.Conn, Handler)

    // ConnState指定一个可选的回调函数，当客户端连接更改状态时调用该函数。
    // 有关详细信息，请参见ConnState类型和关联的常量。
    ConnState func(net.Conn, ConnState)

    // ErrorLog指定一个可选的记录器，用于接收连接的错误，
    // 来自处理程序的意外行为以及潜在的FileSystem错误。
    // 如果为nil，则通过日志包的标准记录器完成记录。
    ErrorLog *log.Logger


    // BaseContext可以选择指定一个函数，该函数返回此服务器上传入请求的基本context。
    // 提供的侦听器是即将开始接收请求的特定侦听器。
    // 如果BaseContext为nil，则默认值为context.Background()。 如果为非nil，则它必须返回非nil Context。
    BaseContext func(net.Listener) context.Context

    // ConnContext可选地指定一个函数，该函数修改用于新连接的上下文c。
    // 提供的ctx派生自基本上下文，并且具有ServerContextKey值。
    ConnContext func(ctx context.Context, c net.Conn) context.Context

    disableKeepAlives int32     // accessed atomically.
    inShutdown        int32     // accessed atomically (non-zero means we're in Shutdown)
    nextProtoOnce     sync.Once // guards setupHTTP2_* init
    nextProtoErr      error     // result of http2.ConfigureServer if used

    mu         sync.Mutex
    listeners  map[*net.Listener]struct{}
    activeConn map[*conn]struct{}
    doneChan   chan struct{}
    onShutdown []func()
}
```

Server定义用于运行HTTP服务器的参数。服务器的零值是有效配置。

## func (*Server) Close

```go
func (srv *Server) Close() error
```

Close立即关闭所有活动的`net.Listeners`以及状态StateNew，StateActive或StateIdle中的所有连接。 要优雅的关机，请使用Shutdown。

Close不会尝试关闭（甚至不知道）任何被劫持的连接，例如WebSockets。 Close返回从关闭服务器的基础侦听器返回的任何错误。

## func (*Server) ListenAndServe

```go
func (srv *Server) ListenAndServe() error
```

ListenAndServe侦听TCP网络地址`srv.Addr`，然后调用Serve处理传入连接上的请求。 接受的连接配置为启用TCP保持活动状态。 

如果`srv.Addr`为空，则使用`“:http”`。

ListenAndServe始终返回非nil错误。 Shutdown或Close后，返回的错误为ErrServerClosed。

## func (*Server) ListenAndServeTLS

```go
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error
```

ListenAndServeTLS侦听TCP网络地址`srv.Addr`，然后调用ServeTLS处理传入TLS连接上的请求。 接受的连接配置为启用TCP keep-alives。

如果未填充服务器的`TLSConfig.Certificates`或`TLSConfig.GetCertificate`，则必须提供包含证书和服务器匹配私钥的文件名。 如果证书是由证书颁发机构签名的，则certFile应该是服务器证书，任何中间件和CA证书的串联。

如果`srv.Addr`为空，则使用`“:https”`。

ListenAndServeTLS始终返回非nil错误。 Shutdown或Close后，返回的错误为ErrServerClosed。

## func (*Server) RegisterOnShutdown

```go
func (srv *Server) RegisterOnShutdown(f func())
```

RegisterOnShutdown注册一个函数来调用Shutdown。 这可用于正常关闭已进行NPN/ALPN协议升级或被劫持的连接。 此功能应启动特定于协议的正常关机，但不应等待关机完成。

## func (*Server) Serve

```go
func (srv *Server) Serve(l net.Listener) error
```

服务在侦听器`l`上接受传入连接，为每个连接创建一个新的服务goroutine。 服务goroutine读取请求，然后调用`srv.Handler`对其进行回复。

仅当侦听器返回`*tls.Conn`连接并且在`TLSConfig.NextProtos`中将它们配置为`“h2”`时，才启用`HTTP/2`支持。

服务始终返回非零错误并关闭`l`。 Shutdown或Close后，返回的错误为ErrServerClosed。

## func (*Server) ServeTLS

```go
func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error
```

ServeTLS在侦听器`l`上接受传入连接，为每个连接创建一个新的服务goroutine。 服务goroutine会执行TLS设置，然后读取请求，并调用`srv.Handler`对其进行回复。

如果未填充服务器的`TLSConfig.Certificates`或`TLSConfig.GetCertificate`，则必须提供包含证书和服务器匹配私钥的文件。

如果证书是由证书颁发机构签名的，则certFile应该是服务器证书，任何中间件和CA证书的串联。 ServeTLS始终返回非nil错误。 Shutdown或Close后，返回的错误为ErrServerClosed。

## func (*Server) SetKeepAlivesEnabled

```go
func (srv *Server) SetKeepAlivesEnabled(v bool)
```

SetKeepAlivesEnabled控制是否启用HTTP keep-alive。默认情况下，始终启用keep-alives状态。只有在资源非常有限的环境或正在关闭的服务器中才能禁用它们。

## func (*Server) Shutdown

```go
func (srv *Server) Shutdown(ctx context.Context) error
```

Shutdown可以优雅的关闭服务器，而不会中断任何活动的连接。Shutdown的工作方式是先关闭所有打开的侦听器，然后关闭所有空闲连接，然后无限期等待连接返回到空闲状态，然后关闭。如果提供的上下文在关闭完成之前到期，则Shutdown返回上下文的错误，否则它将返回从关闭服务器的基础侦听器返回的任何错误。

调用Shutdown时，Serve，ListenAndServe和ListenAndServeTLS立即返回ErrServerClosed。确保程序没有退出，而是等待Shutdown返回。

Shutdown不会尝试关闭也不等待被劫持的连接，例如WebSockets。如果需要，Shutdown的调用者应单独通知此类长期存在的连接，并等待它们关闭。有关注册关闭通知功能的方法，请参见RegisterOnShutdown。

一旦在服务器上调用了Shutdown，就可能无法重用它。将来对诸如Serve之类的方法的调用将返回ErrServerClosed。

### Example

```go
var srv http.Server

idleConnsClosed := make(chan struct{})
go func() {
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    <-sigint

    // We received an interrupt signal, shut down.
    if err := srv.Shutdown(context.Background()); err != nil {
        // Error from closing listeners, or context timeout:
        log.Printf("HTTP server Shutdown: %v", err)
    }
    close(idleConnsClosed)
}()

if err := srv.ListenAndServe(); err != http.ErrServerClosed {
    // Error starting or closing listener:
    log.Fatalf("HTTP server ListenAndServe: %v", err)
}

<-idleConnsClosed
```

## type Transport struct {

```go
    // Proxy指定一个函数来返回给定请求的代理。
    // 如果函数返回非零错误，则请求将中止并提供所提供的错误。
    // Proxy类型由URL scheme确定。 支持“http”，“ https”和“socks5”。
    //  如果scheme为空，则假定为“http”。 如果Proxy为nil或返回nil *URL，则不使用任何代理。
    Proxy func(*Request) (*url.URL, error)

    // DialContext指定用于创建未加密的TCP连接的拨号功能。
    // 如果DialContext为nil（并且下面弃用的Dial也为nil），则Transport使用程序包net进行调用。
    // DialContext与RoundTrip的调用同时运行。 当较早的连接在以后的DialContext完成之前变为空闲时，
    // 发起调用的RoundTrip调用可能会使用先前调用的连接结束。
    DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

    // Dial指定用于创建未加密的TCP连接的拨号功能。 Dial与RoundTrip的调用同时运行。
    // 当较早的连接在之后的Dial完成之前变为空闲时，发起Dial的RoundTrip调用可能会使用先前调动的连接结束。
    // 弃用：改用DialContext，它允许Transport在不再需要调用时立即取消Dial。 
    // 如果两者都设置，则DialContext优先。
    Dial func(network, addr string) (net.Conn, error)

    // DialTLS指定用于为非代理HTTPS请求创建TLS连接的可选Dial功能。
    // 如果DialTLS为nil，则使用Dial和TLSClientConfig。
    // 如果设置了DialTLS，则Dial挂钩不用于HTTPS请求，并且TLSClientConfig和TLSHandshakeTimeout将被忽略。
    // 假定返回的net.Conn已通过TLS握手。
    DialTLS func(network, addr string) (net.Conn, error)

    // TLSClientConfig指定要与tls.Client一起使用的TLS配置。
    // 如果为nil，则使用默认配置。 如果为非nil，则默认情况下可能不会启用HTTP/2支持。
    TLSClientConfig *tls.Config

    // TLSHandshakeTimeout指定等待TLS握手的最长时间。 零表示没有超时。
    TLSHandshakeTimeout time.Duration

    // DisableKeepAlives（如果为true）将禁用HTTP keep-alives，并且仅将与服务器的连接用于单个HTTP请求。
    // 这与类似命名的TCP keep-alives。
    DisableKeepAlives bool

    // DisableCompression如果为true，则当请求不包含现有的Accept-Encoding值时，
    // 阻止Transport使用“Accept-Encoding:gzip”请求标头请求压缩。
    // 如果Transport请求gzip并获得gzip压缩的响应，则会在Response.Body中对其进行透明解码。
    // 但是，如果用户明确请求gzip，则不会自动将其解压缩。
    DisableCompression bool

    // MaxIdleConns控制所有主机之间的最大空闲（keep-alive）连接数。 零表示无限制。
    MaxIdleConns int

    // MaxIdleConnsPerHost（如果非零）控制最大空闲（keep-alive）连接以保留每个主机。
    // 如果为零，则使用DefaultMaxIdleConnsPerHost。
    MaxIdleConnsPerHost int

    // MaxConnsPerHost可以选择限制每个主机的连接总数，
    // 包括处于拨号，活动和空闲状态的连接。 超出限制时，拨号将阻塞。 零表示无限制。
    MaxConnsPerHost int

    // IdleConnTimeout是空闲（keep-alive）连接在关闭自身之前将保持空闲状态的最长时间。
    // 零表示无限制。
    IdleConnTimeout time.Duration

    // ResponseHeaderTimeout（如果非零）指定在完全写入请求（包括其主体（如果有））
    // 之后等待服务器的响应头的时间。 该时间不包括读取响应体的时间。
    ResponseHeaderTimeout time.Duration

    // ExpectContinueTimeout（如果非零）指定如果请求具有“Expect: 100-continue”header，
    // 则在完全写入请求头之后等待服务器的第一个响应头的时间。
    // 零表示没有超时，并导致请求体立即发送，而无需等待服务器批准。 此时间不包括发送请求头的时间。
    ExpectContinueTimeout time.Duration

    // TLSNextProto指定在TLS NPN/ALPN协议协商之后，Transport如何切换到备用协议（例如HTTP/2）。
    // 如果Transport使用非空协议名称调用TLS连接，并且TLSNextProto包含该键的映射条目（例如“ h2”），
    // 则将以请求的权限（例如“example.com”或“example .com:1234“）和TLS连接。
    // 该函数必须返回RoundTripper，然后再处理请求。 如果TLSNextProto不为nil，则不会自动启用HTTP/2支持。
    TLSNextProto map[string]func(authority string, c *tls.Conn) RoundTripper

    // ProxyConnectHeader可以选择指定在CONNECT请求期间发送给代理的header。
    ProxyConnectHeader Header

    // MaxResponseHeaderBytes指定对服务器的响应头中允许的响应字节数的限制。 零表示使用默认限制。
    MaxResponseHeaderBytes int64

    // WriteBufferSize指定在写入Transport时使用的写入缓冲区的大小。 如果为零，则使用默认值（当前为4KB）。
    WriteBufferSize int

    // ReadBufferSize指定从Transport读取时使用的读取缓冲区的大小。 如果为零，则使用默认值（当前为4KB）。
    ReadBufferSize int

    // 当提供非零Dial，DialTLS或DialContext函数或TLSClientConfig时，ForceAttemptHTTP2控制是否启用HTTP/2。
    // 默认情况下，保守地使用这些字段会禁用HTTP/2。要使用自定义Dial程序或TLS配置并仍尝试HTTP/2升级，
    // 请将其设置为true。
    ForceAttemptHTTP2 bool

    idleMu       sync.Mutex
    closeIdle    bool                                // user has requested to close all idle conns
    idleConn     map[connectMethodKey][]*persistConn // most recently used at end
    idleConnWait map[connectMethodKey]wantConnQueue  // waiting getConns
    idleLRU      connLRU

    reqMu       sync.Mutex
    reqCanceler map[*Request]func(error)

    altMu    sync.Mutex   // guards changing altProto only
    altProto atomic.Value // of nil or map[string]RoundTripper, key is URI scheme

    connsPerHostMu   sync.Mutex
    connsPerHost     map[connectMethodKey]int
    connsPerHostWait map[connectMethodKey]wantConnQueue // waiting getConns

}
```

Transport是RoundTripper的实现，它支持HTTP，HTTPS和HTTP代理（对于HTTP或带CONNECT的HTTPS）。

默认情况下，Transport缓存连接以供将来重用。访问许多主机时，这可能会留下许多打开的连接。可以使用Transport的CloseIdleConnections方法以及MaxIdleConnsPerHost和DisableKeepAlives字段来管理此行为。

Transport应该被重用，而不是根据需要创建。多个goroutine并发使用Transport是安全的。Transport是用于发出HTTP和HTTPS请求的低级原语。有关Cookie和重定向之类的高级功能，请参阅Client。

Transport使用`HTTP/1.1`作为HTTP URL，`HTTP/1.1`或`HTTP/2`作为HTTPS URL，这取决于服务器是否支持`HTTP/2`，以及Transport的配置方式。 DefaultTransport支持`HTTP/2`。要在Transport上显式启用`HTTP/2`，请使用`golang.org/x/net/http2`并调用ConfigureTransport。有关HTTP/2的更多信息，请参见软件包文档。

状态代码在1xx范围内的响应将自动处理（100 expect-continue）或被忽略。一个例外是HTTP状态代码101（交换协议），它被认为是终止状态，由RoundTrip返回。若要查看被忽略的1xx响应，请使用httptrace跟踪包的`ClientTrace.Got1xxResponse`。

如果请求是幂等且没有请求体或已定义其`Request.GetBody`，则Transport仅在遇到网络错误时重试该请求。如果HTTP请求具有HTTP方法GET，HEAD，OPTIONS或TRACE，则它们被认为是幂等的。或者其header映射包含“Idempotency-Key”或“ X-Idempotency-Key”条目。如果幂等键值为零长度切片，则将请求视为幂等，但header不会在被发送出去。

## func (*Transport) Clone

```go
func (t *Transport) Clone() *Transport
```

Clone返回`t`导出字段的深层副本。

## func (*Transport) CloseIdleConnections

```go
func (t *Transport) CloseIdleConnections()
```

CloseIdleConnections关闭先前与以前的请求建立连接但现在处于“keep-alive”状态的空闲连接。它不会中断当前正在使用的任何连接。

## func (*Transport) RegisterProtocol

```go
func (t *Transport) RegisterProtocol(scheme string, rt RoundTripper)
```

RegisterProtocol使用scheme注册新协议。 Transport将使用给定scheme将请求传递给`rt`。 模拟HTTP请求语义是`rt`的责任。

其他包可以使用RegisterProtocol提供协议scheme的实现，例如“ftp”或“file”。 如果`rt.RoundTrip`返回ErrSkipAltProtocol，则Tramsport将为该请求处理RoundTrip本身，就像未注册协议一样。

## func (*Transport) RoundTrip`

```go
func (t *Transport) RoundTrip(req *Request) (*Response, error)
```

RoundTrip实现RoundTripper接口。 有关更高级别的HTTP客户端支持（例如cookie和重定向的处理），请参阅Get，Post和Client类型。 与RoundTripper界面类似，RoundTrip返回的错误类型未指定。
