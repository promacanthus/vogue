# Go并发模式： Context

## 简介

在Go服务器中，每个传入请求都在其自己的goroutine中进行处理。请求处理程序通常会启动其他goroutine来访问后端（例如数据库和RPC服务）。处理请求的goroutine集合通常需要访问特定于请求的值（例如，最终用户的身份，授权令牌和请求的期限）。当一个请求被取消或超时时，处理该请求的所有goroutine应该迅速退出，以便系统可以回收他们正在使用的任何资源。

在Google，我们开发了一个context包，可以轻松地跨API边界将请求范围值，取消信号和截止时间传递给处理请求的所有goroutine。该软件包可作为[context](../Library/context.md)公开使用。本文介绍了如何使用该程序包，并提供了一个完整的工作示例。

## Context

Context包的核心是Context类型：

```go
// 上下文在API边界上包含截止日期，取消信息和请求范围值，
// 它的方法可以安全地被多个goroutine同时使用。
type Context interface {
    // Done返回一个通道，当上下文取消或超时，该通道会被关闭
    Done() <-chan struct{}

    //  Err说明在Done通道被关闭后，为何取消此上下文
    Err() error

    // Deadline 返回一个时间，表示该上下文将在何时被取消（如果有）。
    Deadline() (deadline time.Time, ok bool)

    //  Value 返回与key相关的值，如果没有则返回nil。
    Value(key interface{}) interface{}
}
```

Done方法将一个通道作为取消信号返回给代表上下文运行的函数：当该通道被关闭后，这些函数应放弃工作并返回。 Err方法返回一个错误，说明为何取消此上下文。 [管道与取消](https://blog.golang.org/pipelines)一文更详细地讨论了Done 通道的用法。

Context没有Cancel方法，因为Done通道是唯一的接收通道，接收取消信号的功能通常不是发送信号的功能。特别是，当父操作为子操作启动goroutine时，这些子操作应该不能取消父操作。相反，WithCancel函数（如下所述）提供了一种取消新的上下文值的方法。

由多个goroutine同时使用上下文也是安全的。代码可以将单个上下文传递给任意数量的goroutine，并可以取消该上下文来通知所有持有该上下文的goroutine。

Deadline方法允许函数确定它们是否应该开始工作，如果剩余时间过少，那可能是不值得的。代码还可以使用deadline来设置I/O操作的超时时间。

Value允许上下文携带请求范围的数据。该数据必须安全的被多个goroutine同时使用。

## Derived（派生） context

Context包提供了从现有值派生新的Context值的函数。这些值形成一棵树：取消context时，从该context派生的所有context也会被取消。

Background是任何context树的根；它永远不会被取消：

```go
// Background返回一个空的Context。 它永远不会被取消，没有截止时间，也没有值。
// Background 通常用于main，init和test中，并用作传入请求的顶层Context。
func Background() Context
```

WithCancel和WithTimeout返回派生的Context值，该Context可以比先于父Context被取消。请求处理程序返回时，通常会取消与传入请求关联的context。使用多个副本时，WithCancel对于取消冗余请求很有用。 WithTimeout用于设置对后端服务器请求的截止时间：

```go
// WithCancel返回 parent的副本，该副本的通道会在parent.Done关闭或取消后立即关闭。
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// CancelFunc 取消 Context
type CancelFunc func()

// WithTimeout returns a copy of parent whose Done channel is closed as soon as
// parent.Done is closed, cancel is called, or timeout elapses. The new
// Context's Deadline is the sooner of now+timeout and the parent's deadline, if
// any. If the timer is still running, the cancel function releases its
// resources.

// WithTimeout返回 parent 的副本，该副本的通道会在parent.Done关闭、取消或超时后立即关闭。
// 新 context的截止时间是 now + timeout 与 父级截止时间（如果有）中的较早者。 
// 如果计时器仍在运行，那么cancel 函数将释放其资源。
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

WithValue提供了一种将请求范围的值与Context相关联的方法：

```go
// WithValue返回父项的副本，该副本的 Value方法返回key 的 val。
func WithValue(parent Context, key interface{}, val interface{}) Context
```

查看如何使用上下文包的最佳方法是通过一个有效的示例。

## 示例：Google Web Search

我们的示例是一个HTTP服务器，该服务器转发对“golang”的查询到[Google Web Search API](https://developers.google.com/web-search/docs/)并渲染结果来处理类似`/search?q=golang&timeout=1s`的URL。 超时参数告诉服务器在该持续时间过去之后取消请求。

代码被拆分为三个包：

- [server](https://blog.golang.org/context/server/server.go)提供`/search`的主要功能和处理程序。
- [userip](https://blog.golang.org/context/userip/userip.go)提供用于从请求中提取用户IP地址并将其与上下文关联的功能。
- [google](https://blog.golang.org/context/google/google.go)提供用于向Google发送查询的搜索功能。

### server

server程序通过为golang提供前几个Google搜索结果来处理`/search?q=golang`之类的请求。 它注册handleSearch来处理`/ search` 端点（Endpoint）。 处理程序将创建一个称为ctx的初始Context，并安排在处理程序返回时将其取消。 如果请求中包含超时URL参数，则在超时后会自动取消Context：

```go
func handleSearch(w http.ResponseWriter, req *http.Request) {

    // ctx是这个处理程序的Context，调用cancel来关闭ctx.Done通道
    // 这是此处理程序启动的请求的取消信息
    var (
        ctx    context.Context
        cancel context.CancelFunc
    )

    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    if err == nil {
        // 该请求具有超时，因此创建一个在超时后自动取消的Context
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        ctx, cancel = context.WithCancel(context.Background())
    }

    defer cancel() // handleSearch返回后立即取消ctx
```

处理程序从请求中提取查询，并通过调用userip包来提取客户端的IP地址。后端请求需要客户端的IP地址，因此handleSearch会将其附加到ctx：

```go
    // 检查search查询
    query := req.FormValue("q")
    if query == "" {
        http.Error(w, "no query", http.StatusBadRequest)
        return
    }

    // 将用户IP存储在ctx中，以供其他包中的代码使用。
    userIP, err := userip.FromRequest(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ctx = userip.NewContext(ctx, userIP)
```

处理程序使用ctx和查询调用google.Search：

```go
  // 运行Google搜索并打印结果
    start := time.Now()
    results, err := google.Search(ctx, query)
    elapsed := time.Since(start)
```

如果搜索成功，则处理程序将呈现结果：

```go
 if err := resultsTemplate.Execute(w, struct {
        Results          google.Results
        Timeout, Elapsed time.Duration
    }{
        Results: results,
        Timeout: timeout,
        Elapsed: elapsed,
    }); err != nil {
        log.Print(err)
        return
    }
```

### userip

userip包提供用于从请求中提取用户IP地址并将其与Context相关联的功能。Context提供了键值映射，其中键和值均为`interface {}`类型。 键类型必须可以判等，并且值必须是安全的，以便多个goroutine同时使用。 像userip这样的软件包隐藏了此映射的详细信息，并提供了对特定Context值的强类型访问。

为了避免键冲突，userip定义了一个未导出的类型key，并将此类型的值作为Context的key：

```go
// key类型是未导出的，以防止与其他程序包中定义的context的key冲突。
type key int

// userIPkey是用户IP地址context的key，其值为零是任意的（arbitrary）。
// 如果此程序包定义了其他context key，则它们将具有不同的整数值。
const userIPKey key = 0
```

FromRequest从`http.Request`中提取一个userIP值：

```go
func FromRequest(req *http.Request) (net.IP, error) {
    ip, _, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
    }
```

NewContext返回一个新的Context，该Context带有提供的userIP值：

```go
func NewContext(ctx context.Context, userIP net.IP) context.Context {
    return context.WithValue(ctx, userIPKey, userIP)
}
```

FromContext从Context中提取用户IP：

```go
func FromContext(ctx context.Context) (net.IP, bool) {
    // 如果ctx没有 key，则ctx.Value返回nil； 对于nil，net.IP类型断言将返回 ok = false。
    userIP, ok := ctx.Value(userIPKey).(net.IP)
    return userIP, ok
}
```

### google

`google.Search`函数向Google Web Search API发出HTTP请求，并解析JSON编码的结果。它接受一个Context类型的参数ctx，如果`ctx.Done`在请求过程中被关闭，那么它将立即返回。

Google Web Search API请求将搜索请求和用户IP作为查询参数：

```go
func Search(ctx context.Context, query string) (Results, error) {
    // 准备Google Search API请求
    req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
    if err != nil {
        return nil, err
    }
    q := req.URL.Query()
    q.Set("q", query)

    // 如果ctx带有用户IP地址，则将其转发到服务器。
    // Google API使用用户IP来区分服务器发起的请求和最终用户发起的请求。
    if userIP, ok := userip.FromContext(ctx); ok {
        q.Set("userip", userIP.String())
    }
    req.URL.RawQuery = q.Encode()
```

Search使用辅助函数httpDo来发出HTTP请求，如果`ctx.Done`在处理请求或响应过程中被关闭则直接取消它。Search将闭包传递给httpDo处理HTTP响应：

```go
var results Results

err = httpDo(ctx, req, func(resp *http.Response, err error) error {
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 解析JSON搜索结果。
    // https://developers.google.com/web-search/docs/#fonje
    var data struct {
        ResponseData struct {
            Results []struct {
                TitleNoFormatting string
                URL               string
            }
        }
    }

    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return err
    }

    for _, res := range data.ResponseData.Results {
        results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
    }

    return nil
})

// httpDo等待我们提供的闭包返回，因此可以在此处安全读取结果。
return results, err
```

httpDo函数运行HTTP请求并在新的goroutine中处理其响应。如果在goroutine退出之前`ctx.Done`关闭，它将立即取消请求：

```go
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
    // 在goroutine中运行HTTP请求，并将响应传递给f。
    c := make(chan error, 1)
    req = req.WithContext(ctx)

    go func() { c <- f(http.DefaultClient.Do(req)) }()

    select {
    case <-ctx.Done():
        <-c // 等待 f 返回。
        return ctx.Err()
    case err := <-c:
        return err
    }
}
```

## 为Context适配代码

许多服务器框架提供用于携带请求范围值的包和类型。我们可以定义Context接口的新实现，以使用现有框架在代码与需要Context参数的代码之间进行桥接。

例如，Gorilla 的 `github.com/gorilla/context` 包允许处理程序通过提供从HTTP请求到键值对的映射，将数据与传入请求相关联。 在[gorilla.go](https://blog.golang.org/context/gorilla/gorilla.go)中，我们提供了一个Context实现，其Value方法返回与Gorilla包中的特定HTTP请求关联的值。

其他软件包提供了与Context类似的取消支持。 例如，Tomb提供了一种Kill方法，该方法通过关闭Dying通道来发出取消信号。 Tomb还提供了等待这些goroutine退出的方法，类似于sync.WaitGroup。 在[tomb.go](https://blog.golang.org/context/tomb/tomb.go)中，我们提供了一个Context实现，当其父Context被取消或提供的Tomb被杀死时，该实现将被取消。

## 总结

在Google，我们要求Go程序员将Context参数作为传入和传出请求之间的调用路径上每个函数的第一个参数传递。 这使许多不同团队开发的Go代码可以很好地进行互操作。 它提供了对超时和取消的简单控制，并确保安全凭证之类的关键值正确地传递Go程序。

希望基于Context构建的服务器框架应提供Context的实现，以在其程序包和需要Context参数的程序包之间架起桥梁。 然后，他们的客户端库将从调用代码中接受上下文。 通过为请求范围的数据和取消建立通用接口，Context使程序包开发人员更容易共享用于创建可伸缩服务的代码。
