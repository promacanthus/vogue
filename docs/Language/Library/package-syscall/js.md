# package js

```go
import "syscall/js"
```

使用`js/wasm`体系结构时，使用js包可以访问WebAssembly主机环境。其API基于JavaScript语义。

这个软件包是实验性的。当前的范围仅是允许运行测试，而尚未为用户提供全面的API。它不受Go兼容性承诺的约束。

## func CopyBytesToGo

```go
func CopyBytesToGo(dst []byte, src Value) int
```

CopyBytesToGo将字节从Uint8Array src复制到dst。它返回复制的字节数，这将是src和dst的最小长度。如果src不是Uint8Array，则CopyBytesToGo会出现panic。

## func CopyBytesToJS

```go
func CopyBytesToJS(dst Value, src []byte) int
```

CopyBytesToJS将字节从src复制到Uint8Array dst。它返回复制的字节数，这将是src和dst的最小长度。如果dst不是Uint8Array，则CopyBytesToJS会出现panic。

## type Error

```go
type Error struct {
    // Value 是基础的JavaScript错误值。
    Value
}
```

Error包装了JavaScript错误。

### func (Error) Error

```go
func (e Error) Error() string
```

Error实现了error接口。

## type Func

```go
type Func struct {
    Value // 调用Go函数的JavaScript函数
    id    uint32
}
```

Func是包装的Go函数，将由JavaScript调用。

### func FuncOf

```go
func FuncOf(fn func(this Value, args []Value) interface{}) Func
```

FuncOf返回包装的函数。

调用JavaScript函数将使用JavaScript的“this”关键字的值和调用的参数来同步调用Go函数fn。 调用的返回值是根据ValueOf将Go函数映射回JavaScript的结果。

- 在从Go到JavaScript的调用期间触发的包装函数将在同一goroutine上执行。
- 由JavaScript的事件循环触发的包装函数将在额外的goroutine上执行。

包装函数中的阻塞操作将阻塞事件循环。 如果一个包装函数被阻止，其他包装函数将不被处理。 因此，阻塞函数应显式启动新的goroutine。

当不再使用该函数时，必须调用Func.Release以释放资源。

### example

```go
var cb js.Func
cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    fmt.Println("button clicked")
    cb.Release() // release the function if the button will not be clicked again
    return nil
})
js.Global().Get("document").Call("getElementById", "myButton").Call("addEventListener", "click", cb)
```

### func (Func) Release

```go
func (c Func) Release()
```

Release释放为该函数分配的资源。调用Release之后不得调用该函数。

## type Type

```go
type Type int
```

Type代表Value的JavaScript类型。

```go
const (
    TypeUndefined Type = iota
    TypeNull
    TypeBoolean
    TypeNumber
    TypeString
    TypeSymbol
    TypeObject
    TypeFunction
)
```

### func (Type) String

```go
func (t Type) String() string
```

## type Value

```go
type Value struct {
    ref ref
}
```

Value表示JavaScript值。零值为JavaScript值“undefined”。

### func Global

```go
func Global() Value
```

Global返回JavaScript全局对象，通常是“window”或“global”。

### func Null

```go
func Null() Value
```

Null返回JavaScript值“ null”。

### func Undefined

```go
func Undefined() Value
```

Undefined返回JavaScript值“undefined”。

### func ValueOf

```go
func ValueOf(x interface{}) Value
```

ValueOf返回x作为JavaScript值：

| Go                     | JavaScript             |
| ---------------------- | ---------------------- |
| js.Value               | [its value]            |
| js.Func                | function               |
| nil                    | null                   |
| bool                   | boolean                |
| integers and floats    | number                 |
| string                 | string                 |
| []interface{}          | new array              |
| map[string]interface{} | new object             |

如果x不是期望的类型之一，则发生panic。

### func (Value) Bool

```go
func (v Value) Bool() bool
```

Bool返回值v作为布尔。如果v不是JavaScript布尔值则返回panic。

### func (Value) Call

```go
func (v Value) Call(m string, args ...interface{}) Value
```

Call使用给定的参数对值v的方法m进行JavaScript调用。如果v没有方法m则返回panic。根据ValueOf函数将参数映射到JavaScript值。

### func (Value) Float

```go
func (v Value) Float() float64
```

Float返回值v作为float64。如果v不是JavaScript的数字，则返回panic。

### func (Value) Get

```go
func (v Value) Get(p string) Value
```

Get返回值v的JavaScript属性p。如果v不是JavaScript对象，则返回panic。

### func (Value) Index

```go
func (v Value) Index(i int) Value
```

Index返回值为v的JavaScript索引i。如果v不是JavaScript对象，则返回panic。

### func (Value) InstanceOf

```go
func (v Value) InstanceOf(t Value) bool
```

InstanceOf根据JavaScript的instanceof运算符报告v是否是类型t的实例。

### func (Value) Int

```go
func (v Value) Int() int
```

Int返回截断为int的值v。如果v不是JavaScript的数字，则返回panic。

### func (Value) Invoke

```go
func (v Value) Invoke(args ...interface{}) Value
```

Invoke使用给定的参数对值v进行JavaScript调用。如果v不是JavaScript函数，则返回panic。根据ValueOf函数将参数映射到JavaScript值。

### func (Value) JSValue

```go
func (v Value) JSValue() Value
```

JSValue实现Wrapper接口。

### func (Value) Length

```go
func (v Value) Length() int
```

Length返回v的JavaScript属性“ length”。如果v不是JavaScript对象，则返回panic。

### func (Value) New

```go
func (v Value) New(args ...interface{}) Value
```

New使用JavaScript的“ new”运算符，将值v作为构造函数和给定参数。如果v不是JavaScript函数，则返回panic。根据ValueOf函数将参数映射到JavaScript值。

### func (Value) Set

```go
func (v Value) Set(p string, x interface{})
```

Set将值v的JavaScript属性p设置为ValueOf(x)。如果v不是JavaScript对象，则返回panic。

### func (Value) SetIndex

```go
func (v Value) SetIndex(i int, x interface{})
```

SetIndex将值v的JavaScript索引i设置为ValueOf(x)。如果v不是JavaScript对象，则返回panic。

### func (Value) String

```go
func (v Value) String() string
```

String返回值v作为字符串。 由于Go的String方法约定，所以String是一种特殊情况。 与其他getter不同，如果v的Type不是TypeString，不会返回panic。 而是返回格式为`“<T>”`或`“<T:V>”`的字符串，其中T是v的类型，V是v的值的字符串表示形式。

### func (Value) Truthy

```go
func (v Value) Truthy() bool
```

Truthy返回值v的JavaScript“真实性”。在JavaScript中，false，0，“”，null，undefined和NaN为“falsy”，其他所有内容均为“truthy”，更多内容点击[这里](https://developer.mozilla.org/en-US/docs/Glossary/Truthy)。

### func (Value) Type

```go
func (v Value) Type() Type
```

Type返回值v的JavaScript类型。它类似于JavaScript的typeof运算符，不同之处在于它返回TypeNull而不是TypeObject表示null。

## type ValueError

```go
type ValueError struct {
    Method string
    Type   Type
}
```

在不支持Value的Value方法上调用Value方法时，会发生ValueError。在每种方法的说明中都记录了这种情况。

### func (*ValueError) Error

```go
func (e *ValueError) Error() string
```

## type Wrapper

```go
type Wrapper interface {
    // JSValue returns a JavaScript value associated with an object.
    JSValue() Value
}
```

Wrapper由JavaScript值支持的类型实现。
