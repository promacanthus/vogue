# package templage

```go
import "text/template"
```

templage包实现数据驱动的模板以生成文本输出。

要生成HTML输出，请查看`html/template`包，该软件包具有与此软件包相同的接口，但会自动保护HTML输出免受某些攻击。

Templates通过将模板应用于数据结构来运行。 模板中的注释指的是数据结构的元素（通常是结构体的字段或map中的键），以控制执行并派生要显示的值。 模板执行引擎将遍历结构并设置光标（以句点`“.”`表示）到结构中当前位置的值。

模板的输入文本是UTF-8编码的任何格式的文本。 

- 动作：数据评估或控制
- 结构：以`“ {{”`和`“}}”`界定

所有动作之外的文本将原样复制到输出中。 除了`raw string`，动作可能不会跨越换行符，尽管注释可以。

一旦开始解析，可以并行安全地执行模板，如果并行执行共享一个Writer，则输出可以交错。

这是一个简单的示例，输出“ 17 items are made of wool ”。

```go
type Inventory struct {
    Material string
    Count    uint
}

sweaters := Inventory{"wool", 17}
tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
if err != nil { panic(err) }
err = tmpl.Execute(os.Stdout, sweaters)
if err != nil { panic(err) }
```

下面展示更多复杂的示例。

## 文字和空格

默认情况下，执行模板时，将逐个复制动作之间的所有文本。例如，在上面的示例中，运行程序时，字符串`“ items made of ”`显示在标准输出上。

但是，为帮助格式化模板源代码，如果操作的左定界符（默认为`“{{”`）后紧跟一个减号和ASCII空格字符（`“{{- ”`），则会紧接在前的文本。 同样，如果右定界符（`“}}”`）前面有一个空格和减号（`“ -}}”`），则所有紧随其后的空格都将修剪掉。 在这些修饰标记中，必须存在ASCII空格； `“{{-3}}”`解析为包含数字-3的动作。

例如，执行以下模板：

```go
"{{23 -}} < {{- 45}}"
// 生成的输出
"23<45"
```

对于此修剪操作，空格字符的定义与Go中的相同：空格，水平制表符，回车符和换行符。

## 动作

这是动作列表，“参数”和“管道”是对数据的评估，将在随后的相应部分中详细定义。

```go
{{/* a comment */}}
{{- /* a comment with white space trimmed from preceding and following text */ -}}
一条注释可能包含换行符。注释不嵌套，但是必须像这里一样以界定符开头和结尾。

{{pipeline}}
管道值的默认的文字表示形式（和fmt.Print输出的一样）被复制到输出中。

{{if pipeline}} T1 {{end}}
如果管道的值为空，则不生成任何输出，否则执行T1。
空值是false，0，任何为nil的指针或者接口，任何长度为0的数组、切片、字典或字符串。
点不受影响。

{{if pipeline}} T1 {{else}} T0 {{end}}
如果管道的值为空，则执行T0，否则执行T1。点不受影响。

{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
为了简化if-else链的外观，上面的例子和下面写的效果是一样的。
        {{if pipeline}} T1 {{else}}{{if pipeline}} T0 {{end}}{{end}}

{{range pipeline}} T1 {{end}}
管道的值必须是数组，切片，字典或通道。 如果管道的值的长度为零，则不输出任何内容。
否则，将点设置为数组，切片或字典的连续元素，然后执行T1。
如果值是字典，并且键是具有定义顺序（“可比较”）的基本类型，则将按排序的键顺序访问元素。

{{range pipeline}} T1 {{else}} T0 {{end}}
管道的值必须是数组，切片，字典或通道。 如果管道的长度为零，则点不受影响并执行T0；
否则，将点设置为数组，切片或字典的连续元素，然后执行T1。

{{template "name"}}
指定名称的模板将使用nil数据执行。

{{template "name" pipeline}}
指定名称的模板通过点设置执行管道的值。

{{block "name" pipeline}} T1 {{end}}
块是定义模板的简写
        {{define "name"}} T1 {{end}}
然后执行它
        {{template "name" pipeline}}
通常的用途是定义一组根模板， 然后通过重新定义其中的块模板进行自定义。

{{with pipeline}} T1 {{end}}
如果管道的值为空，则不生成任何输出。 否则，将点设置为管道的值，并且T1被执行。

{{with pipeline}} T1 {{else}} T0 {{end}}
如果管道的值为空，则点不受影响并且T0 被执行;否则，将点设置为管道的值然后执行T1。
```

## 参数

参数是一个简单值，由以下之一表示。

- Go语法中的布尔值，字符串，字符，整数，浮点数，虚数或复数常量。 它们的行为就像Go的未类型化常量。 请注意，与Go中一样，大整数常量在赋值或传递给函数时是否溢出取决于主机的int是32位还是64位。
- 关键字nil，代表无类型的Go nil。
- 字符`'.'`，结果是点的值。
- 变量名称，它是一个字母数字字符串（可能为空），前面加一个美元符号，例如`$piOver2`或者`$`结果是变量的值。
- 数据字段的名称（必须是结构体），例如：`.Field`，结果是该字段的值，字段的调用可能是链式的：` .Field1.Field2`，还可以对变量进行评估包括链接：`$x.Field1.Field2`

## func HTMLEscape

```go
func HTMLEscape(w io.Writer, b []byte)
```

HTMLEscape将等同于纯文本数据`b`的转义HTML写入`w`。

## func HTMLEscapeString

```go
func HTMLEscapeString(s string) string
```

HTMLEscapeString返回等效于纯文本数据`s`的转义HTML。

## func HTMLEscaper

```go
func HTMLEscaper(args ...interface{}) string
```

HTMLEscaper返回与`args`的文本表示形式等效的转义HTML。

## func IsTrue

```go
func IsTrue(val interface{}) (truth, ok bool)
```

IsTrue报告`val`是否为“true”（不是其类型的零），以及该`val`是否具有有意义的真值。这是if和其他此类操作使用的`true`定义。

## func JSEscape

```go
func JSEscape(w io.Writer, b []byte)
```

JSEscape将等同于纯文本数据`b`的转义JavaScript写入`w`。

## func JSEscapeString

```go
func JSEscapeString(s string) string
```

JSEscapeString返回等效于纯文本数据`s`的转义JavaScript。

## func JSEscaper

```go
func JSEscaper(args ...interface{}) string

```

SEscaper返回等效于`args`的文本表示形式的转义JavaScript。

## func URLQueryEscaper

```go
func URLQueryEscaper(args ...interface{}) string
```

URLQueryEscaper以适合嵌入URL查询的形式返回其参数的文本表示形式的转义值。

## type ExecError

```go
type ExecError struct {
    Name string // Name of template.
    Err  error  // Pre-formatted error.
}
```

ExecError是Execute在评估其模板出错时返回的自定义错误类型。 （如果发生写错误，则返回实际错误；它不会属于ExecError类型。）

### func (ExecError) Error

```go
func (e ExecError) Error() string
```

### func (ExecError) Unwrap

```go
func (e ExecError) Unwrap() error
```

## type FuncMap

```go
type FuncMap map[string]interface{}
```

FuncMap是定义从名称到函数的Map的映射类型。

- 每个函数必须具有单个返回值，
- 或者具有两个返回值，其中第二个是error类型。在这种情况下，如果第二个（错误）返回值在执行过程中评估为non-nil，则执行终止，并且Execute返回该错误。

当模板引擎调用带有参数列表的函数时，该参数列表必须是可指派给该函数的参数类型。

- 适用任意类型参数的函数可以使用`interface{}`或`reflect.Value`类型的参数。
- 返回任意类型结果的函数可以返回`interface{}`或`reflect.Value`类型的返回值。

## type Template

```go
type Template struct {
        name string
        *parse.Tree
        *common
        leftDelim  string
        rightDelim string
}
```

Template表示一个已被解析的模板。 `*parse.Tree`字段仅导出供`html/template`使用，应被所有其他客户端视为未导出。

### 样例

#### Basic

```go
// 定义一个模板
const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.
{{- else}}
It is a shame you couldn't make it to the wedding.
{{- end}}
{{with .Gift -}}
Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

// 准备一些数据以插入模板
type Recipient struct {
    Name, Gift string
    Attended   bool
}
var recipients = []Recipient{
    {"Aunt Mildred", "bone china tea set", true},
    {"Uncle John", "moleskin pants", false},
    {"Cousin Rodney", "", false},
}

// 创建一个template对象并将常量letter解析给它
t := template.Must(template.New("letter").Parse(letter))

// 对每一个recipient执行模板
for _, r := range recipients {
    err := t.Execute(os.Stdout, r)
    if err != nil {
        log.Println("executing template:", err)
    }
}

// 输出
Dear Aunt Mildred,

It was a pleasure to see you at the wedding.
Thank you for the lovely bone china tea set.

Best wishes,
Josie

Dear Uncle John,

It is a shame you couldn't make it to the wedding.
Thank you for the lovely moleskin pants.

Best wishes,
Josie

Dear Cousin Rodney,

It is a shame you couldn't make it to the wedding.

Best wishes,
Josie
```

#### Block

```go
const (
    master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
    overlay = `{{define "list"}} {{join . ", "}}{{end}} `
)
var (
    funcs     = template.FuncMap{"join": strings.Join}
    guardians = []string{"Gamora", "Groot", "Nebula", "Rocket", "Star-Lord"}
)
masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
if err != nil {
    log.Fatal(err)
}
overlayTmpl, err := template.Must(masterTmpl.Clone()).Parse(overlay)
if err != nil {
    log.Fatal(err)
}
if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
    log.Fatal(err)
}
if err := overlayTmpl.Execute(os.Stdout, guardians); err != nil {
    log.Fatal(err)
}

// 输出
Names:
- Gamora
- Groot
- Nebula
- Rocket
- Star-Lord
Names: Gamora, Groot, Nebula, Rocket, Star-Lord
```

#### Func

本示例演示了用于处理模板文本的自定义函数。它使用`strings.Title`函数，并使用它使标题文本在模板的输出中看起来不错。

```go
// 首先，创建一个FuncMap用来注册函数
funcMap := template.FuncMap{
        // “title”是在模板文本中将调用的函数
    "title": strings.Title,
}

// 一个简单的模板定义来测试函数
// 通过几种方式打印输入文本：
// - 原版输出
// - 调用title函数
// - 调用title函数然后以%q输出
// - 以%q输出然后调用title函数
const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`

// 创建一个模板，然后添加函数，并解析文本
tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
if err != nil {
    log.Fatalf("parsing: %s", err)
}

// 运行模板来验证输出
err = tmpl.Execute(os.Stdout, "the go programming language")
if err != nil {
    log.Fatalf("execution: %s", err)
}

// 输出
Input: "the go programming language"
Output 0: The Go Programming Language
Output 1: "The Go Programming Language"
Output 2: "The Go Programming Language"
```

#### Glob

从目录加载一组模板。

```go
// 创建一个临时目录并使用示例模板定义填充该目录
// 通常，模板文件将已经存在与程序已知的某个位置
dir := createTestDir([]templateFile{
    // T0.tmpl 是一个纯模板文件仅调用T1
    {"T0.tmpl", `T0 invokes T1: ({{template "T1"}})`},
    // T1.tmpl 定义一个模板, T1 调用 T2
    {"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
    // T2.tmpl 定义一个模板 T2.
    {"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
})
// 测试后清除目录
defer os.RemoveAll(dir)

// pattern是用于查找所有模板文件的全局模式
pattern := filepath.Join(dir, "*.tmpl")

// T0.tmpl 是第一个匹配的名字 ，因此它成为起始模板即ParseGlob的返回值
tmpl := template.Must(template.ParseGlob(pattern))

err := tmpl.Execute(os.Stdout, nil)
if err != nil {
    log.Fatalf("template execution: %s", err)
}

// 输出
T0 invokes T1: (T1 invokes T2: (This is T2))
```

#### Helper

此示例演示了一种共享模板并在不同上下文中使用它们的方法。在此变体中，我们将多个驱动程序模板手动添加到现有的模板包中。

```go
// 创建一个临时目录并使用示例模板定义填充该目录
// 通常，模板文件将已经存在与程序已知的某个位置
dir := createTestDir([]templateFile{
    // T1.tmpl定义一个模板, T1 调用 T2.
    {"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
    // T2.tmpl 定义模板 T2.
    {"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
})
// 测试后清除目录
defer os.RemoveAll(dir)

// pattern是用于查找所有模板文件的全局模式
pattern := filepath.Join(dir, "*.tmpl")

// 加载帮助程序
templates := template.Must(template.ParseGlob(pattern))
// 添加一个驱动程序模板；用一个明确的模板定义来做到这一点
_, err := templates.Parse("{{define `driver1`}}Driver 1 calls T1: ({{template `T1`}})\n{{end}}")
if err != nil {
    log.Fatal("parsing driver1: ", err)
}
// 添加另一个驱动程序模板
_, err = templates.Parse("{{define `driver2`}}Driver 2 calls T2: ({{template `T2`}})\n{{end}}")
if err != nil {
    log.Fatal("parsing driver2: ", err)
}

// 在执行之前加载所有模板，text/template包不需要这种操作
// 但是html/template包的转义需要这样的操作，这是一个好习惯
err = templates.ExecuteTemplate(os.Stdout, "driver1", nil)
if err != nil {
    log.Fatalf("driver1 execution: %s", err)
}
err = templates.ExecuteTemplate(os.Stdout, "driver2", nil)
if err != nil {
    log.Fatalf("driver2 execution: %s", err)
}

// 输出
Driver 1 calls T1: (T1 invokes T2: (This is T2))
Driver 2 calls T2: (This is T2)
```

#### Share

本示例演示如何将一组驱动程序模板与不同的帮助程序模板集一起使用。

```go
// 创建一个临时目录并使用示例模板定义填充该目录
// 通常，模板文件将已经存在与程序已知的某个位置
dir := createTestDir([]templateFile{
    // T0.tmpl 是一个纯模板文件仅调用 T1.
    {"T0.tmpl", "T0 ({{.}} version) invokes T1: ({{template `T1`}})\n"},
    // T1.tmpl 定义一个模板, T1 调用 T2. 注意 T2 并没有定义
    {"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
})
// 测试后清空目录
defer os.RemoveAll(dir)

// pattern是用于查找所有模板文件的全局模式
pattern := filepath.Join(dir, "*.tmpl")

// 加载驱动程序
drivers := template.Must(template.ParseGlob(pattern))

//  必须定义T2模板的实现，首先克隆驱动程序，然后将T2的定义添加到模板名称空间

// 1. 克隆帮助程序集，以创建一个新的名称空间来运行它们
first, err := drivers.Clone()
if err != nil {
    log.Fatal("cloning helpers: ", err)
}
// 2. 定义 T2, version A, 并解析它
_, err = first.Parse("{{define `T2`}}T2, version A{{end}}")
if err != nil {
    log.Fatal("parsing T2: ", err)
}

// 使用不同版本的T2重复上述过程
// 1. 克隆驱动程序
second, err := drivers.Clone()
if err != nil {
    log.Fatal("cloning drivers: ", err)
}
// 2. 定义 T2, version B, 并解析它
_, err = second.Parse("{{define `T2`}}T2, version B{{end}}")
if err != nil {
    log.Fatal("parsing T2: ", err)
}

// 以相反的顺序执行模板，以验证第一个模板不受第二个模板的影响。
err = second.ExecuteTemplate(os.Stdout, "T0.tmpl", "second")
if err != nil {
    log.Fatalf("second execution: %s", err)
}
err = first.ExecuteTemplate(os.Stdout, "T0.tmpl", "first")
if err != nil {
    log.Fatalf("first: execution: %s", err)
}
Output:

T0 (second version) invokes T1: (T1 invokes T2: (T2, version B))
T0 (first version) invokes T1: (T1 invokes T2: (T2, version A))
```

### func Must

```go
func Must(t *Template, err error) *Template
```

Must是一个帮助程序，它包装返回（`*Template`，error）的函数的调用，如果错误为非nil，则会出现恐慌。它旨在用于变量初始化，例如

```go
var t = template.Must(template.New("name").Parse("text"))
```

### func New

```go
func New(name string) *Template
```

New使用给定名称分配新的未定义模板。

### func ParseFiles

```go
func ParseFiles(filenames ...string) (*Template, error)
```

ParseFiles创建一个新的模板，并从指定文件中解析模板定义。 返回的模板名称将具有第一个文件的基本名称和已解析的内容。 必须至少有一个文件。 如果发生错误，解析将停止并且返回的`*Template`为nil。

当在不同目录中解析具有相同名称的多个文件时，最后一个文件将是结果文件。例如，ParseFiles（`“a/foo”`，`“b/foo”`）将`“ b/foo”`存储为名为`“foo”`的模板，而`“a/foo”`不可用。

### func ParseGlob

```go
func ParseGlob(pattern string) (*Template, error)
```

ParseGlob创建一个新的模板，并从该pattern标识的文件中解析模板定义。 这些文件根据`filepath.Match`的语义进行匹配，并且该模式必须匹配至少一个文件。 返回的模板将具有该模式匹配的第一个文件的（基本）名称和（解析的）内容。 ParseGlob等效于使用模式匹配的文件列表调用ParseFiles。

当在不同目录中解析具有相同名称的多个文件时，最后一个文件将是结果文件。

### func (*Template) AddParseTree

```go
func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error)
```

AddParseTree添加具有指定名称的模板的解析树，并将其与`t`关联。如果该模板尚不存在，它将创建一个新模板。如果模板确实存在，它将被替换。

### func (*Template) Clone

```go
func (t *Template) Clone() (*Template, error)
```

Clone返回模板的副本，包括所有关联的模板。 不会复制实际的表示形式，但是会复制关联模板的命名空间，因此在副本中进一步调用Parse会将模板添加到副本中，而不是原始模板。 Clone可用于准备通用模板，并将其与其他模板的变体定义一起使用，方法是在完成克隆后添加变体。

### func (*Template) DefinedTemplates

```go
func (t *Template) DefinedTemplates() string
```

DefinedTemplates返回一个字符串，其中列出了已定义的模板，并以字符串`“;”`作为前缀。如果没有，则返回空字符串。用于在此处和`html/template`中生成错误消息。

### func (*Template) Delims

```go
func (t *Template) Delims(left, right string) *Template
```

Delims将动作定界符设置为指定的字符串，以在后续对Parse，ParseFiles或ParseGlob的调用中使用。 嵌套模板定义将继承设置。 空定界符代表相应的默认值：`{{`或`}}`。 返回值是template，因此可以链式调用。

### func (*Template) Execute

```go
func (t *Template) Execute(wr io.Writer, data interface{}) error
```

Execute将已解析的模板应用于指定的数据对象，并将输出写入`wr`。 如果执行模板或写入其输出时发生错误，则执行将停止，但是可能已将部分结果写入Writer。 模板可以安全地并行执行，如果并行执行共享一个Writer，则输出可能会交错。

如果数据是`reflect.Value`，则该模板将应用于`reflect.Value`所持有的具体值，如`fmt.Print`中所示。

### func (*Template) ExecuteTemplate

```go
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error
```

ExecuteTemplate将指定名称的与`t`关联的模板应用于指定的数据对象，并将输出写入`wr`。 如果执行模板或写入其输出时发生错误，则执行将停止，但是可能已将部分结果写入输出写入器。 模板可以安全地并行执行，如果并行执行共享一个Writer，则输出可能会交错。

### func (*Template) Funcs

```go
func (t *Template) Funcs(funcMap FuncMap) *Template
```

Funcs将参数映射的元素添加到模板的FuncMap。 必须在解析模板之前调用它。 如果映射中的值不是具有适当返回类型的函数，或者该名称不能在语法上用作模板中的函数，则它会返回恐慌。 覆盖Map元素是合法的。 返回值是Template，因此可以链式调用。

### func (*Template) Lookup

```go
func (t *Template) Lookup(name string) *Template
```

查找将返回与`t`关联的给定名称的模板。如果没有这样的模板或模板没有定义，则返回nil。

### func (*Template) Name

```go
func (t *Template) Name() string
```

Name返回模板的名称。

### func (*Template) New

```go
func (t *Template) New(name string) *Template
```

New分配与给定模板和相同定界符关联的新的未定义模板。该关联是可传递的，它允许一个模板通过`{{template}}`动作来调用另一个模板。

由于关联的模板共享基础数据，因此无法安全并行地进行模板构建。模板一旦构建，就可以并行执行。

### func (*Template) Option

```go
func (t *Template) Option(opt ...string) *Template
```

Option设置模板的选项。选项由字符串（简单字符串或`“key=value”`）描述。选项字符串中最多可以有一个等号。如果选项字符串无法识别或无效，则返回恐慌。

已知选项：

- missingkey：使用Map中不存在的键索引了Map，来控制执行期间的行为。

```go
"missingkey=default" or "missingkey=invalid"
        // default的行为：什么也不做继续执行
        // 如果输出，索引操作的值是"<no value>"字符串
"missingkey=zero"
        // 该操作返回Map类型元素的零值
"missingkey=error"
        // 执行立即停止并返回错误
```

### func (*Template) Parse

```go
func (t *Template) Parse(text string) (*Template, error)
```

Parse将文本解析为`t`的模板主体。文本中的命名模板定义(`{{define ...}}`或`{{block ...}}`语句)定义了与`t`关联的其他模板，并从`t`本身的定义中删除。
sudo
可以在连续调用Parse中重新定义模板。 具有仅包含空白和注释的主体的模板定义被认为是空的，不会替换现有模板的主体。 这允许使用Parse添加新的命名模板定义，而不会覆盖主模板主体。

### func (*Template) ParseFiles

```go
func (t *Template) ParseFiles(filenames ...string) (*Template, error)
```

ParseFiles解析指定的文件，并将生成的模板与`t`关联。 如果发生错误，则解析停止，返回的模板为nil； 否则为`t`。 必须至少有一个文件。 由于由ParseFiles创建的模板是由参数文件的基本名称命名的，因此`t`通常应具有文件（基本）名称之一的名称。 如果不是，则根据`t`的内容，在调用ParseFiles之前，`t.Execute`可能会失败。 在这种情况下，请使用`t.ExecuteTemplate`执行有效的模板。

当在不同目录中解析具有相同名称的多个文件时，最后一个将是结果文件。

### func (*Template) ParseGlob

```go
func (t *Template) ParseGlob(pattern string) (*Template, error)
```

ParseGlob解析pattern识别的文件中的模板定义，并将结果模板与`t`关联。 这些文件根据`filepath.Match`的语义进行匹配，并且该模式必须匹配至少一个文件。 ParseGlob等效于使用模式匹配的文件列表调用`t.ParseFiles`。 

当在不同目录中解析具有相同名称的多个文件时，最后一个文件将是结果文件。

### func (*Template) Templates

```go
func (t *Template) Templates() []*Template
```

Templates返回与`t`关联的已定义模板的切片。
