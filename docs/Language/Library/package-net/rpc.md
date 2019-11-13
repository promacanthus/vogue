# package rpc

```go
import "net/rpc"
```

rpc包提供对通过网络或其他I/O连接的对象导出方法的访问。服务端注册一个对象，使它作为具有对象类型名称的服务而可见。在注册之后，对象的导出方法可以被远程访问。服务端可以注册不同类型的多个对象，但是如果注册同一类型的多个对象是错误的。

只有满足这些标准的方法才可用于远程访问，其他方法将被忽略：

- 方法的类型是导出的
- 方法是导出的
- 方法有两个参数，都是导出或内置类型
- 方法的第二个参数是一个指针
- 方法有返回error类型

实际上，该方法看起来应该像这样：

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

> 其中T1和T2可以通过`encoding/gob`编码、即使使用不同的编解码器，这些要求也是适用的。（在将来，这些要求可能会对自定义编解码器弱化）。

- 方法的第一个参数：表示调用方提供的参数
- 方法的第二个参数：表示要返回给调用方的结果参数

该方法的返回值（如果不是nil）将作为字符串传递回来，客户端将认为它是由`errors.New`创建的。如果有错误被返回，则不会有响应参数发送给客户端。

服务端可以通过调用`ServerConn`来处理单个连接上的请求。更常见的情况是，服务端创建一个网络监听器，然后调用`Accept`，或者创建一个HTTP监听器，然后调用`HandleHTTP`和`http.Serve`。

想要使用该服务的客户端会建立一个连接，然后在连接上调用`NewClient`。`Dial(DialHTTP)`函数为原始网络连接（HTTP连接）执行这两个步骤，然后得到的客户端对象有两个方法`Call`和`Go`，它们指定要调用的服务和方法，一个传递参数的指针和一个接收结果参数的指针。

- Call方法：等待远程调用完成
- Go方法：异步启动调用，并使用`Call`结构体中的`Done`通道发出完成信号

> 除非设置了显式编解码器，否则将使用`encoding/gob`包来传输数据。

下面是一个示例，服务端希望导出Arith类型的对象：

```go
package server

import "errors"

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}
```

服务端调用HTTP服务：

```go
arith := new(Arith)
rpc.Register(arith)
rpc.HandleHTTP()
l, e := net.Listen("tcp", ":1234")
if e != nil {
    log.Fatal("listen error:", e)
}
go http.Serve(l, nil)
```

此时，客户端可以服务“Arith”带有方法“Arith.Multiply”和“Arith.Divide”。要调用这些方法，客户端首先要连接服务端：

```go
client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
if err != nil {
    log.Fatal("dialing:", err)
}
```

然后，客户端可以执行一次远程调用：

```go
// Synchronous call
args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
    log.Fatal("arith error:", err)
}
fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
```

或者：

```go
// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.
```

服务端实现通常会为客户端提供简单的类型安全包装器，`net/rpc`包已经暂停，不会在接受新特性。

## Constants

```go
const (
    // Defaults used by HandleHTTP
    DefaultRPCPath   = "/_goRPC_"
    DefaultDebugPath = "/debug/rpc"
)
```

## Variables

```go
var DefaultServer = NewServer()
//  DefaultServer是 *Server的默认实例
var ErrShutdown = errors.New("connection is shut down")
```

## func Accpet

```go
func Accept(lis net.Listener)
```

Accept接受监听器上的连接，并为每个传入连接向DefaultServer提供请求。Accept可能会被阻塞，通常调用则在go语句中调用它。

## func HandleHTTP

```go
func HandleHTTP()
```

HandleHTTP在DefaultRPCPath上为RPC消息注册一个HTTP处理器到DefaultServer上，并在DefaultDebugPath上注册一个调试器。仍然需要调用`http.Serve()`，通常也是在go语句中。

## func Register

```go
func Register(rcvr interface{}) error
```

Register 在 DefaultServer 中发布接收方的方法。

## func RegisterName

```go
func RegisterName(name string, rcvr interface{}) error
```

Registername 类似于 Register，但使用提供的类型名称而不是接收方的具体类型。

## func ServeCodec

```go
func ServeCodec(codec ServerCodec)
```

Servecodec 类似于 ServeConn，但使用指定的编解码器对请求进行解码并对响应进行编码。

## func ServeConn

```go
func ServeConn(conn io.ReadWriteCloser)
```

ServeConn在单个连接上运行DefaultServer，ServeConn可能会被阻塞，在客户端挂起之前为连接提供服务。通常情况下，调用者在go语句中调用ServeConn。ServeConn在连接上使用gob编码格式。若要使用另外的编解码器，使用ServeCodec。有关并发访问的信息，查看NewClient的注释。

## func ServeRequest

```go
func ServeRequest(codec ServeCodec) error
```

Serverequest 类似于 ServeCodec，但同步服务于单个请求。 它不会在完成时关闭编解码器。

## type call

```go
type Call struct {
    ServiceMethod string      // 要调用的服务和方法的名称
    Args          interface{} // function (*struct)的传入参数
    Reply         interface{} // function (*struct)的返回结果
    Error         error       // 完成后的错误状态
    Done          chan *Call  // 调用结束后发出消息
}
```

Call 表示一个活跃的RPC。

## type Client

```go
type Client struct {
    codec ClientCodec

    reqMutex sync.Mutex // protects following
    request  Request

    mutex    sync.Mutex // protects following
    seq      uint64
    pending  map[uint64]*Call
    closing  bool // user has called Close
    shutdown bool // server has told us to stop
}
```

Client表示一个RPC客户端，一个Client可能有多个未完成的调用，一个Client也可能同时被多个goroutine使用。

## func Dial

```go
func Dial(network, address string) (*Client, error)
```

Dial按照到指定网络地址连接到RPC服务器。

## func DialHTTP

```go
func DialHTTP(network, address string) (*Client, error)
```

DialHTTP按照指定网络地址连接到HTTP RPC服务器，并监听默认的HTTP RPC 路径。

## func DialHTTPPath

```go
func DialHTTPPath(network, address, path string) (*Client, error)
```

DialHTTPPath 按照指定的网络地址和路径连接到 HTTP RPC 服务器。

## func NewClient

```go
func NewClient(conn io.ReadWriteCloser) *Client
```

Newclient 返回一个新 Client 来处理连接另一端的服务集发出的请求。 它向连接的写入端添加一个缓冲区，以便将头和有效负载作为一个单元发送。

连接的读取和写入部分是独立串行化的，因此不需要内联锁。 但是，每一部分都可以被并发的访问，因此 conn 的实现应该能够防止并发读或并发写。

## func NewClientWithCodec

```go
func NewClientWithCodec(codec ClientCodec) *Client
```

Newclientwithcodec 类似于 NewClient，但是使用指定的 codec 对请求进行编码并解码响应。

## func (*Client) Call

```go
func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error
```

Call调用指定名字的函数，等待它完成，并返回它的错误状态。

## func (*Client) Close

```go
func (client *Client) Close() error
```

Close 调用底层编解码器的 Close 方法。 如果连接已经关闭，则返回 ErrShutdown。

## func (*Client) Go

```go
func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
```

Go 以异步方式调用函数。 它返回表示调用的 Call 结构体。 通过返回相同的 Call 对象，done 通道将在调用完成时发出信号。 如果 done 为 nil，则 Go 将分配一个新的通道，如果不是nil，done必然已经被填充否则Go将会奔溃。

## type ClientCodec

```go
type ClientCodec interface {
    WriteRequest(*Request, interface{}) error
    ReadResponseHeader(*Response) error
    ReadResponseBody(interface{}) error

    Close() error
}
```

Clientcodec 实现在一个RPC会话的客户端侧写入RPC请求和读取RPC响应。客户端调用 WriteRequest 向连接写入请求，并成对调用 ReadResponseHeader 和 ReadResponseBody 来读取响应。 当连接结束时，客户端调用 Close。 可以使用 nil 参数调用 ReadResponseBody，以强制读取响应主体，然后将其丢弃。 有关并发访问的信息，请参阅 NewClient 的注释。

## type Request

```go
type Request struct {
    ServiceMethod string   // format: "Service.Method"
    Seq           uint64   // 客户端选择的序列号
    next          *Request // 服务器中的空闲列表
}
```

Request 是在每个 RPC 调用之前写入的请求头。 它在内部使用，但在这里作为调试的帮助而记录在案，例如在分析网络流量时。

## type Response

```go
type Response struct {
    ServiceMethod string    // echoes that of the Request
    Seq           uint64    // echoes that of the request
    Error         string    // error, if any.
    next          *Response // for free list in Server
}
```

Response 是在每个 RPC 返回之前写入的响应头。 它在内部使用，但在这里作为调试的帮助而记录在案，例如在分析网络流量时。

## type Server

```go
type Server struct {
    serviceMap sync.Map   // map[string]*service
    reqLock    sync.Mutex // protects freeReq
    freeReq    *Request
    respLock   sync.Mutex // protects freeResp
    freeResp   *Response
}
```
Server 代表RPC服务器。

## func (*Server) Accept

```go
func (server *Server) Accept(lis net.Listener)
```

Accept接受监听器上的连接，并未每个传入的连接提供请求。当监听器返回一个非nil错误时，Accetp会被阻塞。通常在Go语句中调用Accept。

## func (*Server) HandleHTTP

```go
func (server *Server) HandleHTTP(rpcPath, debugPath string)
```

Handlehttp 在 rpcPath 上为 RPC 消息注册 HTTP 处理器，在 debugPath 上注册调试器。 仍然需要调用`http.Serve()`，通常在 go 语句中。

## func (*Server) Register

```go
func (server *Server) Register(rcvr interface{}) error
```

Register 在服务器中发布满足以下条件的接收器值的方法集:

- 导出类型的导出方法
- 两个参数都是导出类型
- 第二个参数值指针类型
- 一个error类型的返回值

如果接收器不是导出类型或者没有合适的方法，则返回一个错误。 它还使用log包记录错误。 客户端使用“ Type.Method”形式的字符串访问每个方法。 其中 Type 是接收方的具体类型。

## func (*Server) RegisterName

```go
func (server *Server) RegisterName(name string, rcvr interface{}) error
```

Registername 类似于 Register，但使用给定的类型名称而不是接收方的具体类型。

## func (*Server) ServeCodec

```go
func (server *Server) ServeCodec(codec ServerCodec)
```

Servecodec 类似于 ServeConn，但使用给定的编解码器对请求进行解码并对响应进行编码。

## func (*Server) ServeConn

```go
func (server *Server) ServeConn(conn io.ReadWriteCloser)
```

Serveconn 在单个连接上运行服务器。 Serveconn 可能会被阻塞，在客户端挂断之前为连接提供服务。 调用方通常在 go 语句中调用 ServeConn。 Serveconn 在连接上使用gob编码格式。 若要使用另外的编解码器，请使用 ServeCodec。 有关并发访问的信息，请参阅 NewClient 的注释。

## func (*Server) ServeHTTP

```go
func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

ServeHTTP实现了一个`http.Handler`来响应RPC请求。

## func (*Server) ServeRequest

```go
func (server *Server) ServeRequest(codec ServerCodec) error
```

Serverequest 类似于 ServeCodec，但同步服务于单个请求。 它不会在完成时关闭编解码器。

## type ServerCodec

```go
type ServerCodec interface {
    ReadRequestHeader(*Request) error
    ReadRequestBody(interface{}) error
    WriteResponse(*Response, interface{}) error

    // Close 可以被多次调用，并且必须是幂等的
    Close() error
}
```

Servercodec 为 RPC 会话的服务器端实现对 RPC 请求的读取和 RPC 响应的写入。 服务器成对调用 readrequesttheader 和 ReadRequestBody 来读取来自连接的请求，并调用 WriteResponse 来写回应。 服务器在连接结束时调用 Close。 可以使用 nil 参数调用 ReadRequestBody，以强制读取和丢弃请求的主体。 有关并发访问的信息，请参阅 NewClient 的注释。

## type ServerError

```go
type ServerError string
```

Servererror 表示从 RPC 连接的远程端返回的错误。

## func (ServerError) Error

```go
func (e ServerError) Error() string
```