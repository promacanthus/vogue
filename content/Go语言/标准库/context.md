---
title: package context
date: 2020-04-14T10:09:14.270627+08:00
draft: false
---

```go
import "context"
```

`context`包定义了`Context`类型，该类型自带

- 截止时间
- 取消信号
- 跨API边界或在两个进程之间的其他请求值

对服务端的传入请求会创建一个`Context`，服务端的传出调用会携带一个`Context`。客户端与服务端之间的函数调用链必须传递`Context`，其他可选的替代`Context`是一个使用`WithCancel`，`WithDeadline`，`WithTimeout`或`WithValue`创建的派生`Context`。当某个`Context`被取消后，基于它派生的所有`Context`也都会被取消。

`WithCancel`，`WithDeadline`和`WithTimeout`函数接受父`Context`并返回派生的子`Context`和一个`CancelFunc`函数。调用`CancelFunc`函数会取消子派生及子派生的子项，删除父项对子项的引用，并停止所有关联的计时器。调用`CancelFunc`函数失败会导致泄漏子派生及子派生子项，直到最开始的父项被取消或者计时器到达而被触发。`go vet`工具检查`CancelFunc`函数是否在所有控制流路径上使用。

使用`Context`的程序应遵循这些规则，以使各个包之间的接口保持一致，并启用静态分析工具来检查`Context`的传播：

- 不要将`Context`存储在结构类型中；而应该将`Context`明确传递给需要它的每个函数。
- Context应该是第一个参数，通常命名为ctx：

```go
func DoSomething(ctx context.Context, arg Arg) error {
    // ... use ctx ...
}
```

即使函数允许，也不要传递`nil Context`。如果不确定要使用哪个`context`，请传递`context.TODO`。

仅将`Context`值用于转换进程和API之间的请求数据，而不是将可选参数传递给函数。

可以将相同的`Context`传递给在不同`goroutine`中运行的函数；`Context`对于多个`goroutine`同时使用是安全的。

有关使用`Context`的服务示例代码，请参阅[示例](../官方博客/context.md)

## 变量

```go
// Canceled是context取消时Context.Err返回的错误
var Canceled = errors.New("context canceled")

// DeadlineExceeded是Context.Err在context截止时间过后返回的错误
var DeadlineExceeded error = deadlineExceededError{}
```

## 函数`WithCancel`

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

`WithCancel`返回带有新`Done`通道的父`context`的副本。返回的`context`的`Done`通道将会在调用返回的`cancel`函数或父`context`的`Done`通道关闭时被关闭(这两种情况中先发生的为准）。

**取消此`context`会释放与其关联的资源，因此代码应在此`context`中运行的操作完成后立即调用`cancel`函数。**

### WithCancel示例

此示例演示了使用可取消的context来防止goroutine泄漏。在示例函数结束时，gen启动的goroutine将返回而不会泄漏。

```go


// gen在单独的goroutine中生成整数将它们发送到返回的channel。
// gen的调用者需要取消一次context，他们消费生成的整数而不泄漏，内部的goroutine由gen开始

gen := func(ctx context.Context) <-chan int {
    dst := make(chan int)
    n := 1
    go func() {
        for {
            select {
            case <-ctx.Done():
                return // returning not to leak the goroutine
            case dst <- n:
                n++
            }
        }
    }()
    return dst
}

ctx, cancel := context.WithCancel(context.Background())
defer cancel() // cancel when we are finished consuming integers

for n := range gen(ctx) {
    fmt.Println(n)
    if n == 5 {
        break
    }
}

// 输出
1
2
3
4
5
```

## 函数`WithDeadline`

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

WithDeadline返回`context`的副本，且该副本的截止时间被调整为不迟于`d`。如果父`context`的截止时间早于`d`，则`WithDeadline（parent，d）`在语义上等同于父`context`。返回的`context`的`Done`通道在截止时间到期、调用返回的`cancel`函数时或父`context`的`Done`通道关闭时关闭（以先发生的情况为准）。

**取消此`context`会释放与其关联的资源，因此代码应在此`context`中运行的操作完成后立即调用`cancel`函数。**

### WithDeadline示例

这个例子传递一个带有任意截止时间的`context`来告诉阻塞函数它应该在截止时间到达时放弃它的工作。

```go
d := time.Now().Add(50 * time.Millisecond)
ctx, cancel := context.WithDeadline(context.Background(), d)

// 尽管ctx将会过期，最佳实践是在任何情况下都要调用它的取消函数。
// 如果不这样做可能会导致保留context及其父context的活动时间超过我们的预期

defer cancel()

select {
case <-time.After(1 * time.Second):
    fmt.Println("overslept")
case <-ctx.Done():
    fmt.Println(ctx.Err())
}
```

## 函数`WithTimeout`

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithTimeout`函数返回`WithDeadline(parent,time.Now().Add(timeout))`。

**取消此`context`会释放与其关联的资源，因此代码应在此`context`中运行的操作完成后立即调用`cancel`函数。**

```go
func slowOperationWithTimeout(ctx context.Context) (Result, error) {
    ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
    defer cancel()  // releases resources if slowOperation completes before timeout elapses
    return slowOperation(ctx)
}
```

### WithTimeout示例

此示例传递具有超时的context，以告知阻塞函数在超时过后它应该放弃其工作。

```go
// 传递一个带有超时时长的context来告诉一个正在阻塞的函数，
// 它应该在超时时长过后放弃正在进的工作。
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()

select {
case <-time.After(1 * time.Second):
    fmt.Println("overslept")
case <-ctx.Done():
    fmt.Println(ctx.Err()) // prints "context deadline exceeded"
}
```

## 类型`CancelFunc`

```go
type CancelFunc func()
```

`CancelFunc`告诉某个操作放弃它正在进行的工作。`CancelFunc`不会等待工作停止。第一次调用后，对`CancelFunc`的后续调用都不执行任何操作。

## 类型`Context`

```go
type Context interface {
    // Deadline函数返回完成工作的时间（该时间代表此context应该被取消）。
    // 如果没有设置deadline，Deadline函数返回ok==false。
    // 对Deadline函数的连续调用都返回相同的结果。
    Deadline() (deadline time.Time, ok bool)

    // Done函数返回一个通道，当工作完成时通道将会被关闭，(这代表了这个context应该被取消)。
    // 如果这个context不能被取消，Done函数也可能返回nil。
    // 对Done函数的连续调用会返回相同的结果。

    // 当取消被调用时，WithCancel安排Done通道关闭
    // 当截止时间过期时，WithDeadline安排Done通道关闭
    // 当超时时长到期时，WithTimeout安排Done通道关闭
    //
    // Done提供用于select语句:
    //
    //  // Stream使用DoSomething来生成值并将这些值发送出去，
    //  // 直到DoSomething返回一个错误或者ctx.Done被关闭。
    //  func Stream(ctx context.Context, out chan<- Value) error {
    //      for {
    //          v, err := DoSomething(ctx)
    //          if err != nil {
    //              return err
    //          }
    //          select {
    //          case <-ctx.Done():
    //              return ctx.Err()
    //          case out <- v:
    //          }
    //      }
    //  }
    //
    // 查阅 https://blog.golang.org/pipelines 获取跟多如何使用Done通道来取消context的例子。
    Done() <-chan struct{}

    // 如果Done没有关闭，Err函数返回nil。
    // 如果Done已经关闭，Err函数返回non-nil来解释：
    //      如果是取消，为什么context被取消
    //      如果是截止时间超时，为什么context的截止时间过了
    // 在Err函数返回一个non-nil错误后, 持续的调用Err返回的都是同一个错误。
    Err() error

    // Value函数返回与此context的key关联的value，如果没有value与key关联，则返回nil。
    // 使用相同的key连续调用Value函数会返回相同的结果。

    // 仅将context值用于切换进程和API的请求数据，而不是将可选参数传递给函数。

    // key标识context中特定的值。想要在Context中存储值的函数通常在全局变量中分配一个key，
    // 然后使用该key作为context.WithValue和Context.Value的参数。
    // key可以是支持判等的任何类型; 包应该将key定义为未导出类型以避免冲突。

    // 定义Context key的包应该为使用该key存储的值提供类型安全的访问器：
    //  // 包使用者定义一个User类型存储在Context中。
    //  package user
    //
    //  import "context"
    //
    //  // User是存储在context中的值的类型
    //  type User struct {...}
    //
    //  // 定义在本包中的key是一个非导出类型。
    //  // 这避免了与定义在其他包中的key产生冲突。
    //  type key int
    //
    //  // userKey是context中user.User值的key。它的未导出的。
    //  // 客户端使用user.NewContext和user.FromContext而不是直接使用这个key。
    //  var userKey key
    //
    //  // NewContext返回一个新的Context其中携带值u。
    //  func NewContext(ctx context.Context, u *User) context.Context {
    //      return context.WithValue(ctx, userKey, u)
    //  }
    //
    //  // FromContext返回User的值，该值存储在ctx中(如果有)。
    //  func FromContext(ctx context.Context) (*User, bool) {
    //      u, ok := ctx.Value(userKey).(*User)
    //   return u, ok
    //  }
    Value(key interface{}) interface{}
}
```

`Context`跨越API边界携带截止时间，取消信号和其他跨API边界的值。`Context`的方法可以由多个`goroutine`同时调用。

### 方法`Background`

```go
func Background() Context
```

`Background()`返回一个`non-nil`和空的`Context`。它永远不会被取消，没有值，也没有截止时间。它通常由主函数初始化和测试，并作为传入请求的顶级Context。

### 方法`TODO`

```go
func TODO() Context
```

`TODO()`返回一个`non-nil`和空的`Context`。代码应该使用`context.TODO`当不清楚使用哪个`Context`或者它还不可用时（因为周围的函数尚未扩展为接受`Context`参数）。

### 方法`WithValue`

```go
func WithValue(parent Context, key, val interface{}) Context
```

`WithValue()`返回父`context`的副本，其中与key关联的值为`val`。

仅将`context`值用于转换进程和API请求数据，而不是将可选参数传递给函数。

提供的key必须是可比较的，不应该是字符串类型或任何其他内置类型，以避免使用`context`的包之间的冲突。`WithValue()`的用户应该为key定义他们自己的类型。为了避免在指派接口时分配，`context`的key通常有具体类型的结构体。或者，导出`context` 的key的变量的静态类型应该是指针或接口。

#### WithValue示例

此示例演示如何将值传递给`context`以及如何检索它（如果存在）。

```go
type favContextKey string

f := func(ctx context.Context, k favContextKey) {
    if v := ctx.Value(k); v != nil {
        fmt.Println("found value:", v)
        return
    }
    fmt.Println("key not found:", k)
}

k := favContextKey("language")
ctx := context.WithValue(context.Background(), k, "Go")

f(ctx, k)
f(ctx, favContextKey("color"))

// 输出
found value: Go
key not found: color
```
