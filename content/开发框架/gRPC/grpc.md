---
title: grpc.md
date: 2020-04-14T10:09:14.262627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 开发框架
- gRPC
summary: grpc.md
showInMenu: false

---

# package grpc

grpc包实现了一个称为gRPC的RPC系统。关于gRPC的更多详细访问 [grpc.io](https://grpc.io/)。

## Constants

```go
const (
    SupportPackageIsVersion3 = true
    SupportPackageIsVersion4 = true
    SupportPackageIsVersion5 = true
)
```

从生成的`protocol buffer`文件中引用`SupportPackageIsVersion`变量，以确保与所使用的gRPC版本兼容。最新的支持包版本为5。

保留较旧的版本是为了兼容性.。如果无法保持兼容性，则可以将其删除。

这些常量不应从任何其他代码中引用。

```go
const PickFirstBalancerName = "pick_first"
```

`PickFirstBalancerName`是`pick_first`负载均衡器的名字。

```go
const Version = "1.26.0-dev"
```

`Version`是当前gRPC版本。

## Variables

```go
var DefaultBackoffConfig = BackoffConfig{
    MaxDelay: 120 * time.Second,
}
```

`DefaultBackoffConfig`使用[本文中的内容](https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md)为`backoff`指定值。

弃用，改用`ConnectParams`代替。在整个1.x系列中都会支持它。

```go
var EnableTracing bool
```

`EnableTracing`控制是否使用`golang.org/x/net/trace`包跟踪RPC。仅应在此程序发送或接收任何RPC之前设置此项。

```go
var (
    // ErrClientConnClosing 表示此操作是非法的，因为ClientConn正在被关闭
    //
    // 弃用: 用户不应该依赖此错误，使用Cancled状态码代替它
    ErrClientConnClosing = status.Error(codes.Canceled, "grpc: the client connection is closing")
)
```

```go
var ErrClientConnTimeout = errors.New("grpc: timed out when dialing")
```

`ErrClientConnTimeout`表示`ClientConn`无法在指定的时间内建立底层连接。

弃用：grpc永远不会返回此错误，用户不应该引用它。

```go
var ErrServerStopped = errors.New("grpc: the server has been stopped")
```

`ErrServerStopped`表示此操作现在是非法的，因为服务器被停止了。

## func Code

```go
func Code(err error) codes.Code
```

如果错误是由rpc系统产生的，那么Code返回错误代码，否则返回`codes.Unknown`。

弃用：使用`status.Code`代替。

## func ErrorDesc

```go
func ErrorDesc(err error) string
```

如果错误是有rpc系统产生的，那么ErrorDesc返回错误的描述信息，否则返回`err.Error()`，当err为nil的时候返回空字符串。

弃用：使用`status.Convert`和Method方法代替它。

## func Errorf

```go
func Errorf(c codes.Code, format string, a ...interface{}) error
```

Errorf返回错误其中包含错误代码和描述信息，如果c的值是ok，那么返回nil。

弃用：使用`status.Errorf`代替。

## func Invoke

```go
func Invoke(ctx context.Context, method string, args, reply interface{}, cc *ClientConn, opts ...CallOption) error
```

Invoke在线路上发送RPC请求，并在接收到响应后返回，这通常由生成的代码调用。

弃用：使用`ClientConn.Invoke`代替它。

## func Method

```go
func Method(ctx context.Context) (string, bool)
```

Metthod返回服务端上下文中`method`字符串，返回的字符串格式为"/service/method"。

## func MethodFromServerStream

```go
func MethodFromServerStream(stream ServerStream) (string, bool)
```

MethodFromServerStream返回输入流的`method`字符串，返回的字符串格式为"/service/method"。

## func NewContextWithServerTransportStream

```go
func NewContextWithServerTransportStream(ctx context.Context, stream ServerTransportStream) context.Context
```

NewContextWithServerTransportStream 从ctx创建一个新的上下文并将流附加到该上下文。此API是实验性的。

## func SendHeader

```go
func SendHeader(ctx context.Context, md metadata.MD) error
```

SendHeader发送标头元数据，做多只能调用一次，将发送由SetHeader()设置的提供的md和headers。

## func SetHeader

```go
func SetHeader(ctx context.Context, md metadata.MD) error
```

SetHeader设置标头元数据，当多次调用的时候，所有提供的元数据将会被合并。发生以下情况之一，将发送所有元数据：

- `grpc.SendHeader()`被调用
- 第一个响应被发送出去
- RPC状态被发送（错误或成功）

## func SetTrailer

```go
func SetTrailer(ctx context.Context, md metadata.MD) error
```

SetTrailer设置在RPC返回时将发送的trailer元数据，当多次调用时，所有提供的元数据将会被合并。

## type Address

```go
type Address struct {
    // Addr是建立连接的服务器地址
    Addr string
    // 元数据是与Addr相关的信息，用于做负载均衡决策
    Metadata interface{}
}
```

Address代表客户端连接到的服务器。

弃用：使用balancer包

## type BackoffConfig

```go
type BackoffConfig struct {
    // MaxDelay 是 backoff 延迟的上界
    MaxDelay time.Duration
}
```

BackoffConfig定义缺省gRPCbackoff策略的参数。

弃用：使用ConnectParams代替，在1.x系列中都会支持。

## type Balancer

```go
type Balancer interface {
    // 负责启动balancer的初始化工作。
    //  例如，它可以启动名称解析并观察更新。它将会在dial时调用。
    Start(target string, config BalancerConfig) error
    // 通知balancergRPC已经在addr地址连接到服务器。
    // 一旦与addr的连接关闭或者丢失将会返回down。
    // TODO: 目前还不清楚如何构造和利用down的有意义的错误参数，
    // 需要现实的需求来指导。
    Up(addr Address) (down func(error))
    // 获取与 ctx 对应的 RPC 的服务器地址。
    // i) 如果它返回一个连接的地址，gRPC内部会在这个连接上发出到这个地址的RPC。
    // ii) 如果它返回一个正在构建（有Notify(...)启动）但是还没有连接的地址，
    //  如果gRPC处于fail-fast连接处于TransientFailure或关闭的状态 ，那么gRPC内部构建失败，
    //  否则将在这个连接上发送RPC。
    // iii) 如果它返回一个连接不存在的地址，gRPC内部将任务是一个错误，
    // 并且将使响应的RPC失败。

    // 因此，在编写自定义balancer时，推荐使用以下规则。
    //  如果ops.BlockingWait是true，它应该返回一个已经连接的地址或块（如果没有已连接的地址）
    //  在阻塞时，它应该尊重ctx的超时或取消。
    //  如果opts.BlockingWait是false（对于fail-fast的RPC），
    // 它应该立即返回一个由Notify(...)通知的地址而不是块。

    // 函数返回put，这个函数在rpc完成或者失败后被调用。
    // put可以收集RPC统计数据并将其报告给远程负载均衡器。
    //
    // 这个函数应该只返回错误的balancer而不能恢复它。
    // 如果返回错误，gRPC 内部将使 RPC 失败。
    Get(ctx context.Context, opts BalancerGetOptions) (addr Address, put func(), err error)
    // 通知返回一个通道，gRPC内部使用该通道来监视gRPC需要连接的地址。
    // 这些地址可能来自名称解析器或远程负载平衡器。
    //  gRPC内部会将其与现有的连接地址进行比较。
    //  如果平衡器通知的地址不在现有的连接地址中，则gRPC将开始连接该地址。
    // 如果现有连接地址中的地址不在通知列表中，则相应的连接会正常关闭。
    // 否则，将不执行任何操作。 请注意，地址片必须是应连接的地址的完整列表。
    // 它不是delta.
    Notify() <-chan []Address
    // 关闭 balancer
    Close() error
}
```

balancer为 rpc 选择网络地址。

弃用：使用balancer包，可能会在未来的1.x 版本中被删除。

## func RoundRobin

```go
func RoundRobin(r naming.Resolver) Balancer
```

RoundRobin返回一个Balancer，它选择地址循环，它使用r监视名称解析更新并相应更新可用的地址。

弃用：使用`balancer/doundrobin`，可能在将来的1.x版本中删除。

## type BalancerConfig

```go
type BalancerConfig struct {
    // DialCreds是Balancer实施可用于拨号到远程负载均衡器服务器的传输凭据。
    // 如果Balancer实现不需要安全地与另一方通话，则可以忽略此设置。
    DialCreds credentials.TransportCredentials
    // Dialer是Balancer实现可用于拨打到远程负载均衡器服务器的自定义拨号器。
    // Balancer实现如果不需要与远程平衡器通信，可以忽略此操作。
    Dialer func(context.Context, string) (net.Conn, error)
}
```

BalancerConfig指定Balancer的配置。

弃用：使用balancer包，可能在将来的1.x版本中删除。

## type BalancerGetOptions

```go
type BalancerGetOptions struct {
    //   Blockingwait 指定当没有连接地址时 Get 是否应该阻塞。
    BlockingWait bool
}
```

Balancergetutions 配置 Get 调用。

弃用：使用balancer包，可能在将来的1.x版本中删除。

## type CallOption

```go
type CallOption interface {
    // contains filtered or unexported methods
}
```

CallOption在调用开始前配置它，或在完成之后从调用中提取信息。

## func CallContentSubtype

```go
func CallContentSubtype(contentSubtype string) CallOption
```

Callcontentsubtype 返回一个 CallOption，该 CallOption 将为调用设置 content-subtype。 例如，如果 content-subtype 是“ json” ，那么连接上的 Content-Type 将是`application / grpc + json`。 Content-type 中包含内容子类型之前，会将其转换为小写形式。更多细节请参见[这里](https://github.com/grpc/grpc/blob/master/doc/protocol-http2.md#requests)。

如果不使用 ForceCodec，那么将使用 content-subtype 在 RegisterCodec 控制的注册表中查找要使用的 Codec。 有关注册的详细信息，请参阅 RegisterCodec 上的文档。 内容子类型的查找不区分大小写。 如果没有找到这样的编解码器，调用将导致代码错误。

如果还使用了 ForceCodec，那么该 Codec 将用于所有请求和响应消息，对于请求，其内容子类型设置为给定的 contentSubtype。

## func CallCustomCodec

```go
func CallCustomCodec(codec Codec) CallOption
```

Callcustomcodec 的行为类似于 ForceCodec，但是接受 `grpc.Codec` 而不是 `encoding.Codec`。

弃用: 使用 ForceCodec代替它。

## func FailFast

```go
func FailFast(failFast bool) CallOption
```

FailFast 是 WaitForReady 的反义词。

弃用：使用WaitForReady

## func ForceCodec

```go
func ForceCodec(codec encoding.Codec) CallOption
```

Forcecodec 返回一个 CallOption，它将设置给定的 Codec 用于调用的所有请求和响应消息。 调用 String ()的结果将以不区分大小写的方式用作 content-subtype。

更多关于Content-Type的细节看[这里](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests)。有关 Codec 和 content-subtype 之间交互的更多细节，请参阅 RegisterCodec 和 CallContentSubtype 的文档。

这个函数是为高级用户提供的; 更喜欢只使用 CallContentSubtype 来选择已注册的编解码器。

这是一个实验性的 API。

## func Header

```go
func Header(md *metadata.MD) CallOption
```

Header 返回一个 CallOptions，用于检索普通 RPC 的 Header 元数据。

## func MaxCallRecvMsgSize

```go
func MaxCallRecvMsgSize(s int) CallOption
```

