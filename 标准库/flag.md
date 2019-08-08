# package flag

```go
import "flag"
```

flag包实现了命令行标识的解析。

## 使用方式

1. 使用`flag.String()`,`Bool()`,`Int()`等来定义命令行标识。

### 方式一

定义一个整数标识`-flagname`，并将值保存在`*int`类型的指针`ip`中，其中`1234`为标识的默认值，`”help message for flagname“`是该标识的帮助信息：

```go
var ip = flag.Int("flagname", 1234, "help message for flagname")
```

### 方式二

使用`Var()`函数将该标识绑定到一个**变量**上，如下所示：

```go
var flagvar int     // 定义标识要绑定的变量
func init() {
    flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
}
```

### 方式三

创建一个自定义的标识（使用指针接收器），只要该标识满足Value接口，然后将它绑定到标识解析器，如果使用这样的变量来自定义标识，那么默认值就是**变量的初始值**。

```go
flag.Var(&flagVal, "name", "help message for flagname")
```

2. 在所有的标识都定义完成之后，执行下面的调用,将命令行解析到定义的标识中：

```go
flag.Parse()
```

然后可以直接使用这些标识：

- 如果直接使用标识本身，那么它们都是**指针**
- 如果将标识绑定到变量,那么它们就是**值**

```go
fmt.Println("ip has value ", *ip)
fmt.Println("flagvar has value ", flagvar)
```

解析后，所有标识对应的参数可用用切片`flag.Args()`或用`flag.Arg(i)`一个个独立的获取到。参数从0开始通过`flag.NArg()-1`索引。

## 命令行标识语法

以下形式的命令行参数是允许的：

```bash
-flag
-flag=x
-flag x  // 仅限非布尔标识
```

使用`-`或者`--`都是可以，它们是等价的。由于命令可能包含的含义，布尔型的标识不允许使用第三种形式。

```go
cmd  -x *
```

其中`*`是Unix shell的通配符，如果有一个文件的名字是`0`、`false`等等时，会使命令产生歧义，因此必须使用`-flag=false`来关闭一个布尔型的标识。

停止解析标识的位置：

1. 第一个非标识参数前（`-`是一个非标识参数）
2. 在终止符`--`后

整型标识接收`1234`、`0666`、`0x1234`或者可能是负整数。

布尔型的标识可以如下：

- 1, t, T, true,  TRUE， True
- 0, f, F, false, FALSE, False

持续时间(`Duration`)标识接受任何`time.ParseDuration`的有效输入。

默认的命令行标识集由顶级函数控制。`FlagSet`类型允许定义独立的标识集，例如在命令行界面中实现子命令。`FlagSet`的方法类似于命令行标识集的顶级函数。

## 示例

```go
// 这些示例演示了flag包的更复杂用法
package main

import (
    "errors"
    "flag"
    "fmt"
    "strings"
    "time"
)

// 示例1：名为“species”的单字符串标识，默认值为“gopher”。
var species = flag.String("species", "gopher", "the species we are studying")

//示例2：共享变量的两个标识，因此我们可以使用简写。
//初始化顺序未定义，因此请确保两者都使用相同的默认值。必须使用init函数设置它们。
var gopherType string

func init() {
    const (
        defaultGopher = "pocket"
        usage         = "the variety of gopher"
    )
    flag.StringVar(&gopherType, "gopher_type", defaultGopher, usage)
    flag.StringVar(&gopherType, "g", defaultGopher, usage+" (shorthand)")
}

// 示例3:用户自定义的标识类型,是一个持续时间的切片
type interval []time.Duration

// String方法用于格式户标识的值,是flag.Value接口的一部分。
// String方法的输出将用于诊断。
func (i *interval) String() string {
    return fmt.Sprint(*i)
}

// Set方法用于设置标识的值，也是flag.Value接口的一部分。
// Set方法的参数是string类型的被解析后用于设置该标识。
// 这是一个逗号分隔的列表，所以我们拆分它。
func (i *interval) Set(value string) error {
    // 如果我们允许多次设置标识来累积值，就需要将下面的if语句删除。
    // 这将运行诸如此类的用法： -deltaT 10s -deltaT 15s 和其他的组合。
    if len(*i) > 0 {
        return errors.New("interval flag already set")
    }
    for _, dt := range strings.Split(value, ",") {
        duration, err := time.ParseDuration(dt)
        if err != nil {
            return err
        }
        *i = append(*i, duration)
    }
    return nil
}

// 定义一个标识来累积持续时间。因为它有一个特殊的类型，
// 我们需要使用Var函数，并在init期间创建该标识。

var intervalFlag interval

func init() {

    // 将命令行标识绑定到intervalFlag变量并设置使用方式消息。
    flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
}

func main() {
    // 有趣的是使用上面声明的变量，为了使flag包能够看到上面定义的标识，必须执行，通常在main（而不是init）的开头执行：
    flag.Parse()
    // 我们不在这里运行它，因为这不是main函数，并且测试套件已经解析了标识。
}
```

## Variables

```go
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)
```

CommandLine是默认的命令行标识集，从`os.Args`解析。顶级函数（如BoolVar，Arg等）是CommandLine方法的包装器。

```go
var ErrHelp = errors.New("flag: help requested")
```

ErrHelp是在调用`-help`或`-h`标识但未定义此类标识时返回的错误。

```go
var Usage = func() {
    fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
    PrintDefaults()
}
```

`Usage`打印一条用法消息，该消息中记录所有已定义的命令行标识并输出到`CommandLine.output()`，默认情况下是输出到`os.Stderr`。只有当解析标识发生错误时才调用`Usage`，[见类型部分中Flag结构体](##类型)。

`Usage`变量是一个函数，可以更改为指向其他的自定义函数。默认情况下，它会打印一个简单的标题并调用`PrintDefaults`；有关输出格式及其控制方法的详细信息，请参阅`PrintDefaults`的文档或者查看下面的函数说明。

自定义的`Usage`函数会导致程序退出；默认情况下也是这样，因为命令行的错误处理策略设置为`ExitOnError`。

## 函数

```go
func Arg（i int） string
```

Arg返回第i个命令行参数。`Arg(0)`是处理完标识后的第一个剩余参数。如果请求的元素不存在，Arg将返回一个空字符串。

---

```go
func Args() []string
```

Args返回非标识命令行参数。

---

```go
// Bool定义了一个带有指定名称，默认值和用法字符串的bool标识。返回值是存储标识值的bool变量的地址（指针）
func Bool(name string, value bool, usage string) *bool{
}

// BoolVar定义了一个带有指定名称，默认值和用法字符串的bool标识。参数p指向一个bool变量，用于存储标识的值。
func BoolVar(p *bool, name string, value bool, usage string){
}
```

```go
// Duration定义具有指定名称，默认值和用法字符串的time.Duration标识。返回值是存储标识值的time.Duration变量的地址。该标识接受time.ParseDuration可接受的值。

func Duration(name string, value time.Duration, usage string) *time.Duration{
}

// DurationVar定义具有指定名称，默认值和用法字符串的time.Duration标识。参数p指向time.Duration变量，用于存储标识的值。该标识接受time.ParseDuration可接受的值。

func DurationVar(p *time.Duration, name string, value time.Duration, usage string){
}
```

```go
// Float64定义了一个带有指定名称，默认值和用法字符串的float64标识。返回值是存储标识值的float64变量的地址。
func Float64(name string, value float64, usage string) *float64{
}

// Float64Var定义了一个带有指定名称，默认值和用法字符串的float64标识。参数p指向一个float64变量，用于存储该标识的值。
func Float64Var(p *float64, name string, value float64, usage string){
}
```

```go
// Int定义具有指定名称，默认值和用法字符串的int标识。返回值是存储标识值的int变量的地址。
func Int(name string, value int, usage string) *int{
}

// IntVar定义了一个带有指定名称，默认值和用法字符串的int标识。参数p指向一个int变量，用于存储标识的值。
func IntVar(p *int, name string, value int, usage string){

}

// Int64定义了一个带有指定名称，默认值和用法字符串的int64标识。返回值是存储标识值的int64变量的地址。
func Int64(name string, value int64, usage string) *int64{
}

// Int64Var定义了一个带有指定名称，默认值和用法字符串的int64标识。参数p指向一个int64变量，用于存储该标识的值。
func Int64Var(p *int64, name string, value int64, usage string){
}
```

```go
// NArg返回的是处理标识后剩余的参数个数。
func NArg() int{
}

// NFlag返回已设置的命令行标识的数量。
func NFlag() int{
}
```

```go
// Parse从os.Args[1:]中解析命令行标识。必须在定义所有标识之后并且在程序访问标识之前调用。
func Parse(){
}

// Parsed报告是否已解析命令行标识。
func Parsed() bool{
}
```

---

```go
// 除非另外配置了输出位置，否则PrintDefaults会将内容打印在标准错误中。输出的内容为一个包含所有已经定义的命令行标识的默认设置的用法消息。
func PrintDefaults()

// 对于整数值标识x，默认输出如下所示：
-x int
    usage-message-for-x (default 7)
```

除了带有单字节名称的bool标识外，用法消息将另起一行显示。对于bool标识，省略类型，如果标识名称是一个字节，则用法消息将显示在同一行。

如果标识的默认值为零值，则会省略用法消息后面括号中的内容。在定义标识时，可以用反引号将标识的用法字符串中的某个名称引起来，那么，用法消息中的第一个被反引号引起来的内容将被视为要在消息中显示的默认值，并且在显示时将从消息中删除反引号。

如下所示：

```go
flag.String("I", "", "search `directory` for include files")

// 用法消息输出如下
-I directory
    search directory for include files.
```

要更改标识的用法消息的输出位置，请调用`CommandLine.SetOutput()`。

---

```go
func Set(name, value string) error{
}
```

Set设置已经被命名的命令行标识的值。

---

```go
// String定义具有指定名称，默认值和用法字符串的字符串标识。返回值是存储标识值的字符串变量的地址。
func String(name string, value string, usage string) *string{
}

// StringVar定义具有指定名称，默认值和用法字符串的字符串标识。参数p指向一个字符串变量，用于存储标识的值。
func StringVar(p *string, name string, value string, usage string){
}
```

```go
// Uint定义了一个带有指定名称，默认值和用法字符串的uint标识。返回值是存储标识值的uint变量的地址。
func Uint(name string, value uint, usage string) *uint{
}

// UintVar定义了一个带有指定名称，默认值和用法字符串的uint标识。参数p指向一个uint变量，用于存储标识的值。
func UintVar(p *uint, name string, value uint, usage string){
}

// Uint64定义了一个带有指定名称，默认值和用法字符串的uint64标识。返回值是存储标识值的uint64变量的地址。
func Uint64(name string, value uint64, usage string) *uint64{
}

// Uint64Var定义了一个带有指定名称，默认值和用法字符串的uint64标识。参数p指向一个uint64变量，用于存储该标识的值。
func Uint64Var(p *uint64, name string, value uint64, usage string){
}
```

```go
// UnquoteUsage从标识的用法字符串中提取被反引号引起来的名称并返回它,同时返回去掉反引号后的用法消息。
func UnquoteUsage(flag *Flag) (name string, usage string){
}

// 给定如下用法消息
"a `name` to show"

// UnquoteUsage输出的内容如下
 ("name", "a name to show")

// 如果没有反引号，则该名称是对标识值类型的有根据的猜测，如果该标识是布尔值，则为空字符串。
```

---

```go
func Var(value Value, name string, usage string)
```

Var定义具有指定名称和用法字符串的标识。

- 标识的类型和值由第一个参数表示，类型为Value，它通常包含**用户定义的Value实现**。

例如，调用者可以创建一个标识，通过为切片提供Value方法，将逗号分隔的字符串转换为字符串切片; 尤其是，Set会将逗号分隔的字符串分解为切片。

---

```go
// Visit为每个fn调用以字典顺序访问命令行标识，它只访问已设置的那些标识。
func Visit(fn func(*Flag)){
}

// VisitAll为每个fn调用以字典顺序访问命令行标识，它访问所有标识，甚至那些未设置的标识。
func VisitAll(fn func(*Flag)){
}
```

## 类型

```go
// ErrorHandling定义了如果解析失败，FlagSet.Parse的行为方式。
type ErrorHandling int

// 如果解析失败，这些常量会反应了FlagSet.Parse的行为。
const (
    ContinueOnError ErrorHandling = iota // 返回描述性错误
    ExitOnError                          // 调用os.Exit(2).
    PanicOnError                         // 调用运行时恐慌并返回描述性错误
)
```

```go
// Flag表示标识的状态。
type Flag struct {
    Name     string // 命令行上显示的名称
    Usage    string // 帮助信息/用法信息
    Value    Value  // 设定的值
    DefValue string // 用法信息的默认值（文本形式）
}

// Lookup返回指定命令行标识的Flag结构，如果不存在则返回nil。
func Lookup(name string) *Flag
```

```go
// FlagSet表示一组定义的标识。FlagSet的零值没有名称，并且具有ContinueOnError错误处理。
type FlagSet struct {
    // Usage是在解析标识发生错误时调用的函数。该字段是一个函数（不是方法），
    // 可以更改为指向自定义错误处理程序。调用Usage后会发生什么取决于ErrorHandling的设置;
    // 对于命令行，默认是ExitOnError，它将在调用Usage后退出程序。
    Usage func(){}    // 为了显示正常，此处没有大括号
    // 包含已过滤或未导出的字段
}

// NewFlagSet返回一个带有指定名称和错误处理属性的新的空标识集。如果名称不为空，则将在默认用法消息和错误消息中打印。
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet{
}

// Arg返回第i个参数。Arg(0)是处理完标识后剩余的第一个参数。如果请求的元素不存在，Arg将返回一个空字符串。
func (f *FlagSet) Arg(i int) string{
}

// Args返回非标识的参数。
func (f *FlagSet) Args() []string{
}

// Bool定义了一个带有指定名称，默认值和用法字符串的bool标识。返回值是存储标识值的bool变量的地址。
func (f *FlagSet) Bool(name string, value bool, usage string) *bool{
}

// BoolVar定义了一个带有指定名称，默认值和用法字符串的bool标识。参数p指向一个bool变量，用于存储标识的值。
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string){
}

// Duration定义具有指定名称，默认值和用法字符串的time.Duration标识。返回值是存储标识值的time.Duration变量的地址。该标识接受time.ParseDuration可接受的值。
func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration{
}

// DurationVar定义具有指定名称，默认值和用法字符串的time.Duration标识。参数p指向time.Duration变量，用于存储标识的值。该标识接受time.ParseDuration可接受的值。
func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string){
}

// ErrorHandling返回标识集的错误处理行为。
func (f *FlagSet) ErrorHandling() ErrorHandling

// Float64定义了一个带有指定名称，默认值和用法字符串的float64标识。返回值是存储标识值的float64变量的地址。
func (f *FlagSet) Float64(name string, value float64, usage string) *float64{
}

// Float64Var定义了一个带有指定名称，默认值和用法字符串的float64标识。参数p指向一个float64变量，用于存储该标识的值。
func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string){
}

// Init为标识集设置名称和错误处理属性。默认情况下，为零值的FlagSet使用空名称和ContinueOnError错误处理策略。
func (f *FlagSet) Init(name string, errorHandling ErrorHandling){
}

// Int定义具有指定名称，默认值和用法字符串的int标识。返回值是存储标识值的int变量的地址。
func (f *FlagSet) Int(name string, value int, usage string) *int{
}

// Int64定义了一个带有指定名称，默认值和用法字符串的int64标识。返回值是存储标识值的int64变量的地址。
func (f *FlagSet) Int64(name string, value int64, usage string) *int64{
}

// Int64Var定义了一个带有指定名称，默认值和用法字符串的int64标识。参数p指向一个int64变量，用于存储该标识的值。
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string){
}

// IntVar定义了一个带有指定名称，默认值和用法字符串的int标识。参数p指向一个int变量，用于存储标识的值。
func (f *FlagSet) IntVar(p *int, name string, value int, usage string){
}

// Lookup返回指定标识的Flag结构，如果不存在则返回nil。
func (f *FlagSet) Lookup(name string) *Flag{
}

// NArg是处理标识后剩余的参数个数。
func (f *FlagSet) NArg() int{
}

// NFlag返回已设置的标识数。
func (f *FlagSet) NFlag() int{
}

// Name返回标识集的名称。
func (f *FlagSet) Name() string{
}

// Output返回用法消和错误消息的输出目的地。如果未设置output或设置为nil，则返回os.Stderr。
func (f *FlagSet) Output() io.Writer{
}

// Parse从参数列表中解析标识定义，该列表不包含命令的名称。必须在定义FlagSet中的所有标识之后并且在程序访问标识之前调用。如果未定义-help或-h，则返回值为ErrHelp。
func (f *FlagSet) Parse(arguments []string) error{
}

// Parsed报告是否已调用f.Parse。
func (f *FlagSet) Parsed() bool{
}

// 除非另外配置，否则PrintDefaults会把整个标识集中定义的标识的默认值输出到标准错误中，有关详细信息，请参阅上面的全局函数PrintDefaults的文档。
func (f *FlagSet) PrintDefaults(){
}

// Set设置指定标识的值。
func (f *FlagSet) Set(name, value string) error{
}

// SetOutput设置用法消息和错误消息的目的地。如果output设置为nil，则使用os.Stderr。
func (f *FlagSet) SetOutput(output io.Writer){
}

// String定义具有指定名称，默认值和用法字符串的字符串标识。返回值是存储标识值的字符串变量的地址。
func (f *FlagSet) String(name string, value string, usage string) *string{
}

// StringVar定义具有指定名称，默认值和用法字符串的字符串标识。参数p指向一个字符串变量，用于存储标识的值。
func (f *FlagSet) StringVar(p *string, name string, value string, usage string){
}

// Uint定义了一个带有指定名称，默认值和用法字符串的uint标识。返回值是存储标识值的uint变量的地址。
func (f *FlagSet) Uint(name string, value uint, usage string) *uint{
}

// Uint64定义了一个带有指定名称，默认值和用法字符串的uint64标识。返回值是存储标识值的uint64变量的地址。
func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64{
}

// Uint64Var定义了一个带有指定名称，默认值和用法字符串的uint64标识。参数p指向一个uint64变量，用于存储该标识的值。
func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string){
}

// UintVar定义了一个带有指定名称，默认值和用法字符串的uint标识。参数p指向一个uint变量，用于存储标识的值。
func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string){
}

// Var定义具有指定名称和用法字符串的标识。标识的类型和值由第一个参数表示，类型为Value，它通常包含用户定义的Value实现。例如，调用者可以创建一个标识，通过为切片提供Value方法，将逗号分隔的字符串转换为字符串切片; 尤其是，Set会将逗号分隔的字符串分解为切片。
func (f *FlagSet) Var(value Value, name string, usage string)

// Visit为每个fn调用以字典顺序访问标识，它只访问已设置的那些标识。
func (f *FlagSet) Visit(fn func(*Flag))

// VisitAll为每个fn调用以字典顺序访问标识，它访问所有标识，甚至那些未设置的标识。
func (f *FlagSet) VisitAll(fn func(*Flag))
```

```go
type Getter interface {
    Value
    Get() interface{}
}
```

Getter是一个允许检索Value内容的接口。因为它出现在Go1及其兼容性规则之后，所以它组合了Value接口，而不是成为value接口的一部分，此flag包提供的所有Value类型都满足Getter接口。

```go
type Value interface {
    String() string
    Set(string) error
}
```

Value是存储在一个标识中的动态值的接口，它的默认值表示是字符串。

如果某个Value接口的`IsBoolFlag()`方法返回true，则命令行解析器使`-name`等效于`-name=true`，而不是使用下一个命令行参数。

为每个存在的标识，以命令行顺序调用且仅调用一次Set。flag包会使用零值接收器调用String方法，例如nil指针。

## 样例

```go
package main

import (
    "flag"
    "fmt"
    "net/url"
)

type URLValue struct {
    URL *url.URL
}

func (v URLValue) String() string {
    if v.URL != nil {
        return v.URL.String()
    }
    return ""
}

func (v URLValue) Set(s string) error {
    if u, err := url.Parse(s); err != nil {
        return err
    } else {
        *v.URL = *u
    }
    return nil
}

var u = &url.URL{}

func main() {
    fs := flag.NewFlagSet("ExampleValue", flag.ExitOnError)
    fs.Var(&URLValue{u}, "url", "URL to parse")

    fs.Parse([]string{"-url", "https://golang.org/pkg/flag/"})
    fmt.Printf(`{scheme: %q, host: %q, path: %q}`, u.Scheme, u.Host, u.Path)

}
```
