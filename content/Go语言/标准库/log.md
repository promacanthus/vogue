---
title: packages log
date: 2020-04-14T10:09:14.274627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 标准库
summary: packages log
showInMenu: false

---

```go
import "log"
```

log包实现了一个简单的日志包，它定义了一种类型`Logger`，其中包含格式化输出的方法。

它还有一个预定义的“标准”Logger，可以通过辅助函数`Print[f|ln]`，`Fatal[f|ln]`和`Panic[f|ln]`访问，使用这些辅助函数比手动创建Logger更容易使用。Logger写入标准错误输出并打印每个日志记录的时间和日期。

每条日志记录都在单独的行中输出：如果正在打印的消息未以换行符结尾，则Logger将会在行尾添加一个换行符。

- Fatal函数在写入日志消息后调用`os.Exit(1)`。
- Panic函数在写入日志消息后调用运行时恐慌（panic）。

## 常量

```go
const (
    Ldate         = 1 << iota     // 当地时区的日期: 2009/01/23
    Ltime                         // 当地时区的时间: 01:23:23
    Lmicroseconds                 // 微秒精度: 01:23:23.123123.  假设是Ltime
    Llongfile                     // 完整的文件名和行号: /a/b/c/d.go:23
    Lshortfile                    // 最终文件名元素和行号: d.go:23. 覆盖Llongfile
    LUTC                          // 如果设置了Ldate或者Ltime，那么使用UTC而不是本地时区
    LstdFlags     = Ldate | Ltime // 标准Logger的初始值
)
```

这些标志定义了Logger生成的每个日志条目的**前缀**，（`|`） 或运算来控制哪些需要输出。无法控制它们出现的顺序（如，此处列出的顺序）或它们呈现的格式（如，注释中所述）。仅当指定`Llongfile`或`Lshortfile`时，前缀后跟冒号。查看如下示例：

```go
// 设置Ldate | Ltime（或LstdFlags）则产生如下效果
2009/01/23 01:23:23 message

// 设置Ldate | Ltime | Lmicroseconds | Llongfile 则产生如下效果
2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
```

## 函数Fatal

```go
func Fatal(v ...interface{})
```

`Fatal()`等效于`Print()`,然后调用`os.Exit(1)`。

## 函数Fatalf

```go
func Fatalf(format string, v ...interface{})
```

`Fatal()`等效于`Printf()`,然后调用`os.Exit(1)`。

## 函数Fatalln

```go
func Fatalln(v ...interface{})
```

`Fatal()`等效于`Println()`,然后调用`os.Exit(1)`。

## 函数Flags

```go
func Flags() int
```

`Flags()`返回标准`Logger`的输出标志。

## 函数Output

```go
func Output(calldepth int, s string) error
```

`Output()`将记录事件的输出写出来。

- 字符串`s`包含要在`Logger`的`flags`中指定的前缀之后打印的文本内容。如果`s`的最后一个字符不是换行符，则会附加换行符。
- `Calldepth`是设置了`Llongfile`或`Lshortfile`后，计算文件名和行号时要跳过的帧数; 如果`Calldepth`的值为1将打印`Output()`调用者的详细信息。

## 函数Panic

```go
func Panic(v ...interface{})
```

`Panic()`等效于`Print()`，然后调用`panic()`。

## 函数Panicf

```go
func Panicf(format string, v ...interface{})
```

`Panic()`等效于`Printf()`，然后调用`panic()`。

## 函数Panicln

```go
func Panicln(v ...interface{})
```

`Panic()`等效于`Println()`，然后调用`panic()`。

## 函数Prefix

```go
func Prefix() string
```

`Prefix()`返回标准`Logger`的输出前缀。

## 函数Print

```go
func Print(v ...interface{})
```

`Print()`调用`Output()`来打印到标准`Logger`，参数以`fmt.Print()`的方式处理。

## 函数Printf

```go
func Printf(format string, v ...interface{})
```

`Print()`调用`Output()`来打印到标准`Logger`，参数以`fmt.Printf()`的方式处理。

## 函数Println

```go
func Println(v ...interface{})
```

`Print()`调用`Output()`来打印到标准`Logger`，参数以`fmt.Println()`的方式处理。

## 总结几个方法

![image](/标准库../images/log.png)

## 函数SetFlags

```go
func SetFlags(flag int)
```

`SetFlags()`为标准`Logger`设置输出标识。

## 函数SetOutput

```go
func SetOutput(w io.Writer)
```

`SetOutput()`设置标准`longer`的输出目的地。

## 函数SetPrefix

```go
func SetPrefix(prefix string)
```

`SetPrefix()`设置标准`Logger`的输出前缀。

## 函数New

```go
func New(out io.Writer, prefix string, flag int) *Logger
```

`New()`创建了一个新的`Logger`。

- `out`设置日志数据的目的地
- `prefix`每行日志记录的前缀
- `flag`定义日志记录操作的属性



## 类型Logger

```go
type Logger struct {
    // contains filtered or unexported fields
}
```

Logger表示一个活跃的日志记录对象，它向`io.Writer`生成输出。每个日志记录操作只对`Writer`的`Write()`方法进行一次调用。Logger可以同时被多个goroutine使用; 它保证对Writer的序列化访问。

### 例子

```go
var (
    buf    bytes.Buffer
    logger = log.New(&buf, "logger: ", log.Lshortfile)
)

logger.Print("Hello, log file!")

fmt.Print(&buf)

// 输出
logger: example_test.go:19: Hello, log file!
```

### 方法 (*Logger)Fatal

```go
func (l *Logger) Fatal(v ...interface{})
```

`Fatal()`等效于`l.Print()`，然后调用`os.Exit(1)`。

### 方法 (*Logger)Fatalf

```go
func (l *Logger) Fatalf(format string, v ...interface{})
```

`Fatalf()`等效于`l.Printf()`，然后调用`os.Exit(1)`。

### 方法 (*Logger)Fatalln

```go
func (l *Logger) Fatalln(v ...interface{})
```

`Fatalln()`等效于`l.Println()`，然后调用`os.Exit(1)`。

### 方法 (*Logger)Flags

```go
func (l *Logger) Flags() int
```

Flags返回Logger的输出标志。

### 方法 (*Logger)Output

```go
func (l *Logger) Output(calldepth int, s string) error
```

Output写入记录事件的输出。

- 字符串`s`包含要在`Logger`中`flags`指定的前缀之后打印的文本。如果`s`的最后一个字符不是换行符，则会附加换行符
- `Calldepth`提供通用的用于恢复PC(程序计数器)，当前所有预定义路径上它都是2

#### 样例

```go
var (
    buf    bytes.Buffer
    logger = log.New(&buf, "INFO: ", log.Lshortfile)

    infof = func(info string) {
        logger.Output(2, info)
    }
)

infof("Hello world")

fmt.Print(&buf)

// 输出
INFO: example_test.go:36: Hello world
```

### 方法(*Logger) Panic

```go
func (l *Logger) Panic(v ...interface{})
```

`Panic()`等效于`l.Print()`，然后调用`panic()`。

### 方法(*Logger) Panicf

```go
func (l *Logger) Panicf(format string, v ...interface{})
```

`Panicf()`等效于`l.Printf()`，然后调用`panic()`。

### 方法(*Logger) Panicln

```go
func (l *Logger) Panicln(v ...interface{})
```

`Panicln()`等效于`l.Println()`，然后调用`panic()`。

### 方法(*Logger) Prefix

```go
func (l *Logger) Prefix() string
```

`Prefix()`返回`Logger`的输出前缀。

### 方法(*Logger) Print

```go
func (l *Logger) Print(v ...interface{})
```

`Print()`调用`l.Output()`来输出到`Logger`。参数以`fmt.Println()`的方式处理。

### 方法(*Logger) Printf

```go
func (l *Logger) Printf(format string, v ...interface{})
```

`Printf()`调用`l.Output()`来输出到`Logger`。参数以`fmt.Printf()`的方式处理。

### 方法(*Logger) Println

```go
func (l *Logger) Println(v ...interface{})
```

`Println()`调用`l.Output()`来输出到`Logger`。参数以`fmt.Println()`方式处理。

### 方法(*Logger) SetFlags

```go
func (l *Logger) SetFlags(flag int)
```

`SetFlags()`设置`Logger`的输出`flags`。

### 方法(*Logger) SetOutput

```go
func (l *Logger) SetOutput(w io.Writer)
```

`SetOutput()`设置Logger的输出目的地。

### 方法(*Logger) SetPrefix

```go
func (l *Logger) SetPrefix(prefix string)
```

`SetPrefix()`设置Logger的输出前缀。

### 方法(*Logger) Writer

```go
func (l *Logger) Writer() io.Writer
```

`Writer()`返回Logger的输出目的地。
