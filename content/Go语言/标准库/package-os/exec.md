---
title: exec.md
date: 2020-04-14T10:09:14.274627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 标准库
- package-os
summary: exec.md
showInMenu: false

---

# package exec

```go
import "os/exec"
```

`exec`包用于运行外部命令,它包装了`os.StartProcess`，使它更容易重映射到`stdin`和`stdout`，或者使用管道连接I/O并进行其他调整。

与C语言或者其他编程语言的“系统”库调用不同，`os/exec`包有意不调用系统`shell`，也不扩展任何`glob patterns`或处理一般情况下由`shell`完成的其他扩展，管道或重定向。

> `glob patterns`是一种匹配模式，运用通配符来匹配一个文件列表。
> 
> `exec`包的行为更像是C的“exec”系列函数的功能。

- 要扩展`glob patterns`或者直接调用`shell`时：
  - 要注意转义任何危险的输入
  - 也可以使用`path/filepath`包的`Glob`函数
- 要扩展环境变量，请使用`os`包的`ExpandEnv`。

请注意，此包中的示例是在Unix系统中云的，可能无法在Windows上运行，并且不能在使用`golang.org`和`godoc.org`的Go Playground中运行。

## 变量

```go
var ErrNotFound = errors.New("executable file not found in $PATH")
```

`ErrNotFound`是路径搜索未能找到可执行文件时导致的错误。

## 函数LookPath

```go
func LookPath(file string) (string, error)
```

`LookPath`在`PATH`环境变量命名的目录中搜索可执行文件`file`。如果`file`包含斜杠，则直接尝试搜索且不查询`PATH`。返回结果是**绝对路径**或**相对于当前目录的相对路径**。

### 示例

```go
path, err := exec.LookPath("fortune")
if err != nil {
    log.Fatal("installing fortune in your future")
}
fmt.Printf("fortune is available at %s\n", path)
```

## 类型Cmd

```go
type Cmd struct {
    // Path是所要运行命令的路径，这是唯一必须设置为非零值的字段。
    // 如果是Path是相对的，则相对于{Dir}进行评估。
    Path string

    // Args保存命令行参数，包括命令为Args[0]。
    // 如果Args字段为空或为nil，则直接运行上述{Path}字段。
    // 通常情况下，Path和Args都是通过调用Command来设置的。
    Args []string

    // Env指定该进程的环境变量。每个条目都是“key=value”的格式。
    // 如果Env为nil，那么新创建的进程使用当前进程的环境变量。
    // 在此Env切片中，如果包含的环境变量中有重复的键值，
    // 那么每个重复的键值中只有最新的值会被使用。
    Env []string

    // Dir指定命令的工作目录。
    // 如果Dir是空字符串，那么{Run}会在发起调用的进程的当前目录中运行该命令。
    Dir string

    // Stdin指定进程的标准输入。
    // 如果Stdin为nil，则进程从null设备（os.DevNull）读取。
    // 如果Stdin是*os.File，则进程的标准输入直接连接到该文件。
    // 否则，在执行命令期间，单独的goroutine从Stdin读取并通过管道将数据传递给该命令。
    // 在这种情况下，{Wait}将会一直等待直到goroutine停止拷贝数据，
    // 即，已经读取到Stdin的末尾（EOF或读取错误），或数据写入管道时出错。
    Stdin io.Reader

    // Stdout和Stderr指定进程的标准输出和标准错误输出。
    // 如果其中一个为nil，则{Run}将相应的文件描述符连接到空设备（os.DevNull）。
    // 如果其中一个是*os.File，则进程相应的输出将直接连接到该文件。
    // 否则，在命令执行期间，单独的goroutine通过管道从进程读取并将数据传递给对应的{Writer}。
    // 在这种情况下，{Wait}将会一直等待直到goroutine读取到EOF或遇到错误。
    // 如果Stdout和Stderr是同一个writer，并且具有可以判等（==）的类型，则一次最多只有一个goroutine调用写入函数。
    Stdout io.Writer
    Stderr io.Writer

    // ExtraFiles指定新进程要继承的其他已打开文件。
    // 它不包括标准输入，标准输出或标准错误输出。
    // 如果非nil，那么条目i变为文件描述符3+i。
    // Windows系统不支持ExtraFiles。
    ExtraFiles []*os.File

    // SysProcAttr包含可选的操作系统特定属性。
    // {Run}将它作为os.ProcAttr的Sys字段传递给os.StartProcess。
    SysProcAttr *syscall.SysProcAttr

    // 一旦启动，Process就是底层进程。
    Process *os.Process

    // ProcessState包含已退出进程的信息，在一次调用后无论处于{Wait}或{Run}都可用。
    ProcessState *os.ProcessState
    // 包含已过滤或未导出的字段
}
```

Cmd表示正在准备或运行的外部命令，在调用其`Run`，`Output`或`CombinedOutput`方法后，无法被重用。

### 函数Command

```go
func Command(name string, arg ...string) *Cmd
```

`Command`返回`Cmd`结构以使用给定的参数执行指定的程序，它仅在返回的Cmd结构中设置`Path`和`Args`。

如果`name`不包含路径分隔符，则`Command`使用`LookPath`将`name`解析为完整路径（如果可能）。否则，它直接使用`name`作为`Path`。

返回的`Cmd结构`中`Args`字段是从命令`name`后紧跟的`arg`元素构造而来，因此在`arg`中不需要再包含命令`name`。例如，`Command（“echo”，“hello”）`。`Args[0]`始终是`name`，而不可能是已解析的路径。

在`Windows`上，进程将整个命令行作为单个字符串接收并执行自己的解析。`Command`使用与应用程序兼容的算法`CommandLineToArgvW`（这是最常用的方式）将`Args`组合并引用到命令行字符串中。值得注意的例外是`msiexec.exe`和`cmd.exe`（包括所有批处理文件`batch files`），它们具有不同的反引用算法。在这些或其他类似情况下，可以自己进行引用并在`SysProcAttr.CmdLine`中提供完整的命令行，让Args空着。

```go
cmd := exec.Command("tr", "a-z", "A-Z")
cmd.Stdin = strings.NewReader("some input")
var out bytes.Buffer
cmd.Stdout = &out
err := cmd.Run()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("in all caps: %q\n", out.String())
```

```go
cmd := exec.Command("prog")
cmd.Env = append(os.Environ(),
    "FOO=duplicate_value", // 重复值将被忽略
    "FOO=actual_value",    // 最新值将被使用
)
if err := cmd.Run(); err != nil {
    log.Fatal(err)
}
```

### 函数CommandContext

```go
func CommandContext(ctx context.Context, name string, arg ...string) *Cmd
```

`CommandContext`与`Command`类似，但包含上下文。

如果在命令自行完成之前`context`已经结束，则提供的`context`用于终止进程（通过调用`os.Process.Kill`）。

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
    // 这将在100毫秒后失败。
    // 5秒的sleep进程将会被中断。
}
```

### 方法 (*Cmd)CombinedOutput

```go
func (c *Cmd) CombinedOutput() ([]byte, error)
```

`CombinedOutput`运行命令并返回标准输出和标准错误错误输出的组合。

```go
cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")
stdoutStderr, err := cmd.CombinedOutput()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s\n", stdoutStderr)
```

### 方法 (*Cmd)Output

```go
func (c *Cmd) Output() ([]byte, error)
```

`Output`运行命令并返回其标准输出。任何返回的错误通常都是`*ExitError`类型。如果`c.Stderr`为`nil`，则`Output`填充`ExitError.Stderr`。

```go
out, err := exec.Command("date").Output()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("The date is %s\n", out)
```

### 方法 (*Cmd)Run

```go
func (c *Cmd) Run() error
```

`Run`启动指定的命令并等待它完成。

- 如果命令运行，在标准输入、标准输出和标准错误输出之间复制数据没有问题并且退出状态为零，那么返回的错误为`nil`。
- 如果命令启动但未成功完成，则返回的错误类型为`*ExitError`。对于其他情况，可能会返回其他错误类型。

如果正在被调用的goroutine已使用`runtime.LockOSThread`锁定操作系统线程并修改了任何可继承的OS级别线程状态（例如，Linux或Plan 9名称空间），则新进程将继承调用方的线程状态。

```go
cmd := exec.Command("sleep", "1")
log.Printf("Running command and waiting for it to finish...")
err := cmd.Run()
log.Printf("Command finished with error: %v", err)
```

### 方法 (*Cmd)Start

```go
func (c *Cmd) Start() error
```

`Start`启动指定的命令，但不等待它完成。

一旦命令退出，`Wait`方法将返回退出状态码并释放相关资源。

```go
cmd := exec.Command("sleep", "5")
err := cmd.Start()
if err != nil {
    log.Fatal(err)
}
log.Printf("Waiting for command to finish...")
err = cmd.Wait()
log.Printf("Command finished with error: %v", err)
```

### 方法 (*Cmd)StderrPipe

```go
func (c *Cmd) StderrPipe() (io.ReadCloser, error)
```

`StderrPipe`返回一个管道，该管道将在命令启动时连接到命令的标准错误输出。

`Wait`将在看到命令退出后关闭管道，因此大多数调用者不需要自己关闭管道; 但是，需要注意的是在管道的所有读取完成之前调用`Wait`是不正确的。同样的，使用`StderrPipe`时使用`Run`也是不正确的。

有关习惯用法，请参阅StdoutPipe示例。

```go
cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")
stderr, err := cmd.StderrPipe()
if err != nil {
    log.Fatal(err)
}

if err := cmd.Start(); err != nil {
    log.Fatal(err)
}

slurp, _ := ioutil.ReadAll(stderr)
fmt.Printf("%s\n", slurp)

if err := cmd.Wait(); err != nil {
    log.Fatal(err)
}
```

### 方法 (*Cmd)StdinPipe

```go
func (c *Cmd) StdinPipe() (io.WriteCloser, error)
```

`StdinPipe`返回一个管道，该管道将在命令启动时连接到命令的标准输入。`Wait`查看到命令退出后将自动关闭管道。调用者只需要调用`Close`来强制管道更快关闭。例如，如果正在运行的命令在标准输入关闭之后不会退出，则调用者必须关闭管道。

```go
cmd := exec.Command("cat")
stdin, err := cmd.StdinPipe()
if err != nil {
    log.Fatal(err)
}

go func() {
    defer stdin.Close()
    io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
}()

out, err := cmd.CombinedOutput()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%s\n", out)
```

### 方法 (*Cmd)StdoutPipe

```go
func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
```

`StdoutPipe`返回一个管道，该管道将在命令启动时连接到命令的标准输出。

`Wait`将在看到命令退出后关闭管道，因此大多数调用者不需要自己关闭管道; 但是，值得注意的是在管道的所有读取完成之前调用`Wait`是不正确的。同样的，使用`StdoutPipe`时调用Run是不正确的。有关习惯用法，请参阅示例。

```go
cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
stdout, err := cmd.StdoutPipe()
if err != nil {
    log.Fatal(err)
}
if err := cmd.Start(); err != nil {
    log.Fatal(err)
}
var person struct {
    Name string
    Age  int
}
if err := json.NewDecoder(stdout).Decode(&person); err != nil {
    log.Fatal(err)
}
if err := cmd.Wait(); err != nil {
    log.Fatal(err)
}
fmt.Printf("%s is %d years old\n", person.Name, person.Age)
```

### 方法 (*Cmd)Wait

```go
func (c *Cmd) Wait() error
```

`Wait`等待命令退出并等待任何复制到标准输入或者从标准输出或标准错误输出复制完成。

该命令必须由`Start`启动。

- 如果命令运行，在标准输入、标准输出和标准错误输出之间复制数据没有问题且退出状态码为零，那么，返回的错误为`nil`。
- 如果命令无法运行或未成功完成，则错误类型为`*ExitError`。对于I/O问题，可能会返回其他错误类型。

如果`c.Stdin`，`c.Stdout`或`c.Stderr`中的任何一个不是`*os.File`，那么`Wait`也会等待相应的I/O循环复制到进程或从进程中复制完成。

`Wait`释放与`Cmd`相关的任何资源。

## 类型Error

```go
type Error struct {
    // Name是发生错误的文件名。
    Name string
    // Err是潜在的错误。
    Err error
}
```

Error是`LookPath`无法将文件归类为可执行文件时返回的。

### 方法 (*Error)Error

```go
func (e *Error) Error() string
```

## 类型ExitError

```go
type ExitError struct {
    *os.ProcessState

    // 如果没有收集标准错误输出，Stderr会保留Cmd.Output方法的标准错误输出的子集。
    // 如果错误输出很长，Stderr可能只包含输出的前缀和后缀，中间替换为相关省略字节数的文本。
    // 提供Stderr是用于代码调试，来包含在错误消息中。具有其他需求的用户应根据需要重定向Cmd.Stderr。
    Stderr []byte
}
```

ExitError报告命令退出失败。

### 方法 (*ExitError)Error

```go
func (e *ExitError) Error() string
```
