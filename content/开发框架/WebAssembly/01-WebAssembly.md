---
title: 01-WebAssembly.md
date: 2020-04-14T10:09:14.258627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 开发框架
- WebAssembly
summary: 01-WebAssembly.md
showInMenu: false

---

# 01-WebAssembly

## 简介

Go 1.11向WebAssembly添加了一个实验端口。 Go 1.12对它的某些部分进行了改进，并有望在Go 1.13中进行进一步的改进。

WebAssembly在其[主页](https://webassembly.org/)上描述为：

> WebAssembly（缩写为Wasm）是基于堆栈的虚拟机的二进制指令格式。 Wasm被设计为可移植目标，用于编译高级语言（如C/C ++ / Rust），从而可以在Web上为客户端和服务器应用程序进行部署。

如果不熟悉WebAssembly，请阅读下面的“[入门](##入门)”部分，观看下面的一些“[Go WebAssembly讲座](##GoWebAssembly讲座)”，然后查看下面的更多“[示例](##示例)”。

## 入门

此页面假定Go 1.11或更高版本可以正常运行。有关故障排除，请参阅“[安装故障排除](https://github.com/golang/go/wiki/InstallTroubleshooting)”页面。

为Web编译一个基本的Go语言包：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, WebAssembly!")
}
```

设置环境变量`GOOS=js`和`GOARCH=wasm`来针对WebAssembly进行编译：

```bash
GOOS=js GOARCH=wasm go build -o main.wasm
```

这将生成程序包并生成一个名为`main.wasm`的可执行WebAssembly模块文件。`.wasm`文件扩展名将使以后通过带有正确的`Content-Type`标头的HTTP服务更加容易。

请注意，只能编译主软件包。否则，将获得无法在WebAssembly中运行的目标文件。如果具有要与WebAssembly一起使用的软件包，请将其转换为主软件包并生成二进制文件。

要在浏览器中执行`main.wasm`，我们还需要一个`JavaScript`支持文件和一个`HTML`页面来将所有内容连接在一起。

复制JavaScript支持文件：

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

新建一个`index.html`文件：

```html
<html>
    <head>
        <meta charset="utf-8"/>
        <script src="wasm_exec.js"></script>
        <script>
            const go = new Go();
            WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
                go.run(result.instance);
            });
        </script>
    </head>
    <body></body>
</html>
```

如果浏览器尚不支持`WebAssembly.instantiateStreaming`，则可以使用[polyfill](https://github.com/golang/go/blob/b2fcfc1a50fbd46556f7075f7f1fbf600b5c9e5d/misc/wasm/wasm_exec.html#L17-L22)。

然后从Web服务器提供三个文件（`index.html`、`wasm_exec.js`、`main.wasm`）。例如使用goexec：

```go
# install goexec: go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```

或者使用[基本的HTTP服务器命令](https://play.golang.org/p/pZ1f5pICVbV)。

将`/usr/local/go/bin`添加到PATH环境变量中。可以通过将以下行添加到`/etc/profile`（用于系统范围的安装）或`$HOME/ .profile`中来完成此操作：

```bash
export PATH=$PATH:/usr/local/go/bin
```

注意：对配置文件的更改可能要等到下一次登录计算机后才能应用，或者执行source命令。

最后，浏览至`http//localhost:8080/index.html`，打开JavaScript调试控制台，应该会看到输出。可以修改程序，重建`main.wasm`，然后刷新以查看新的输出。

## 使用Node.js执行WebAssembly

可以使用`Node.js`而非浏览器执行已编译的WebAssembly模块，这对于测试和自动化非常有用。

安装`Node.js`并在PATH中设置好，在执行`go run`或`go test`时，将`-exec`标志设置为`go_js_wasm_exec`所在的目录。默认情况下，`go_js_wasm_exec`位于Go安装的`misc/wasm`目录中。

```bash
GOOS=js GOARCH=wasm go run -exec="$(go env GOROOT)/misc/wasm/go_js_wasm_exec" .
Hello, WebAssembly!
GOOS=js GOARCH=wasm go test -exec="$(go env GOROOT)/misc/wasm/go_js_wasm_exec"
PASS
ok      example.org/my/pkg      0.800s
```

在PATH中添加`go_js_wasm_exec`，可以在执行时不必每次都手动设置`-exec`标志，`go run`和`go test`能直接对`js/wasm`生效。

```shell
export PATH="$PATH:$(go env GOROOT)/misc/wasm"
GOOS=js GOARCH=wasm go run .
Hello, WebAssembly!
GOOS=js GOARCH=wasm go test
PASS
ok      example.org/my/pkg  0.800s
```

## 在浏览器中运行测试

也可以使用[wasmbrowsertest](https://github.com/agnivade/wasmbrowsertest)在浏览器中运行测试。它可以自动完成网络服务器的工作，并使用无头的Chrome浏览器在其中运行测试，并将日志中继到控制台。

就像上一节那样，只需执行`go get github.com/agnivade/wasmbrowsertest`即可获取二进制文件。将其重命名为`go_js_wasm_exec`并将其放置到PATH中。

```shell
mv $GOPATH/bin/wasmbrowsertest $GOPATH/bin/go_js_wasm_exec
export PATH="$PATH:$GOPATH/bin"
GOOS=js GOARCH=wasm go test
PASS
ok      example.org/my/pkg  0.800s
```

或者，使用`-exec`测试标志：

```shell
GOOS=js GOARCH=wasm go test -exec="$GOPATH/bin/wasmbrowsertest"
```

## GoWebAssembly讲座

- [Building a Calculator with Go and WebAssembly](https://www.youtube.com/watch?v=4kBvvk2Bzis) [源码地址](https://tutorialedge.net/golang/go-webassembly-tutorial/)
- [Get Going with WebAssembly](https://www.youtube.com/watch?v=iTrx0BbUXI4)
- [Go&WebAssembly简介 - by chai2010](https://talks.godoc.org/github.com/chai2010/awesome-go-zh/chai2010/chai2010-golang-wasm.slide)

## 与DOM交互

详细查阅[`syscall/js`](https://godoc.org/syscall/js)包。

其他项目：

- [gas](https://github.com/gascore/gas)-WebAssembly应用程序的基于组件的框架
- [app](https://github.com/maxence-charriere/app)-基于React兼容PWA的自定义工具框架
- [Vugu](https://github.com/vugu/vugu)-wasm Web UI库，具有HTML布局，带有Go格式的应用程序逻辑，单个文件组件，快速开发和原型工作流程
- [vue](https://github.com/norunners/vue) -WebAssembly应用程序的渐进框架
- [dom](https://github.com/dennwc/dom)-用于简化DOM操作的库
- [webapi](https://gowebapi.github.io/)-绑定生成器
- [vert](https://github.com/norunners/vert)-Go和JS值之间的WebAssembly互操作
- [GoWebian](https://github.com/bgokden/gowebian)-用于纯粹地构建页面并添加WebAssembly绑定的库

## Canvas

- [go-canvas](https://github.com/markfarnan/go-canvas)-Canvas绘图库，[Simple demo](https://markfarnan.github.io/go-canvas/)

## 使用net/http时配置fetch options

可以使用`net/http`库从Go发出HTTP请求，这些请求将转换为[fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)调用。 但是，fetch [options](https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters) 和 http [客户端](https://golang.org/pkg/net/http/#Client)选项之间没有直接映射。 为此，我们有一些特殊的标头值，这些标头值被认为是fetch options。如下表所示：

|选项|描述|默认有效值|其他有效值|
|---|---|---|---|
|`js.fetch.mode`|设置Fetch API 模式|`same-origin`|`cors`、`no-cors`、`navigation`|
|`js.fetch.credentials`|设置Fetch API 证书|`same-origin`|`omit`、`include`|
|`js.fetch.redirect`|设置Fetch API重定向|`follow`|`error`、`manual`|

因此，如果想在发出请求时将模式设置为`cors`，则将类似于：

```go
req, err := http.NewRequest("GET", "http://localhost:8080", nil)
req.Header.Add("js.fetch:mode", "cors")
if err != nil {
  fmt.Println(err)
  return
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
  fmt.Println(err)
  return
}
defer resp.Body.Close()
// handle the response
```

请随时订阅[＃26769](https://github.com/golang/go/issues/26769)，以获取更多信息和可能的更新。

## 编辑器配置

在Goland和Intellij Ultimate中设置WebAssembly的具体步骤，参考[这里](https://github.com/golang/go/wiki/Configuring-GoLand-for-WebAssembly)。

## Chrome中的WebAssembly

如果运行较新版本的Chrome，则会有一个标志（`chrome://flags＃enable-webassembly-baseline`）启用新编译器`Liftoff`，这将大大缩短加载时间，更多信息在[这里](https://chinagdg.org/2018/08/liftoff-a-new-baseline-compiler-for-webassembly-in-v8/)。

## 调试

WebAssembly尚不支持调试器，因此，现在需要使用良好的 `println()`方法在JavaScript控制台上显示输出。

已经创建了一个官方的[WebAssembly调试子组](https://github.com/WebAssembly/debugging)来解决此问题，并且正在进行一些初步调查和提议：

- [WebAssembly Debugging Capabilities Living Standard](https://fitzgen.github.io/wasm-debugging-capabilities/)[源代码](https://github.com/fitzgen/wasm-debugging-capabilities)
- [DWARF for WebAssembly Target](https://yurydelendik.github.io/webassembly-dwarf/)[源代码](https://github.com/yurydelendik/webassembly-dwarf/)

如果对调试器方面有兴趣，请参与并帮助实现这一目标。

### 分析WebAssembly文件的结构

[WebAssembly代码资源管理器](https://wasdk.github.io/wasmcodeexplorer/)对于可视化WebAssembly文件的结构很有用。

- 单击左侧的十六进制值将突出显示其所属的部分，并在右侧显示相应的文本表示形式
- 单击右侧的一行将在左侧突出显示它的十六进制字节表示

## 已知错误

1.11.2之前的Go版本存在一个[错误](https://github.com/golang/go/issues/27961)，该错误可能在某些（罕见）情况下生成错误的wasm代码。

如果您的Go代码可以毫无问题地编译为wasm，但是在浏览器中运行时会产生如下错误：`CompileError: wasm validation error: at offset 1269295: type mismatch: expression has type i64 but expected f64`，那么可能会遇到此错误。

**解决方案是升级到Go 1.11.2或更高版本**。

## 示例

更多示例点击[这里](https://github.com/golang/go/wiki/WebAssembly#further-examples)。

## 减少Wasm文件的大小

目前，Go会生成大型的Wasm文件，可能的最小大小约为2MB。如果您的Go代码导入了库，则此文件的大小可能会急剧增加。 10MB +是常见的。

目前有两种主要方法来减小此文件的大小：

### 1. 手动压缩`.wasm`文件

1. 使用gz压缩将示例WASM文件的2MB（最小文件大小）减小到大约500kB。使用[Zopfli](https://github.com/google/zopfli)进行gzip压缩可能会更好，因为它提供的结果比`gzip --best`更好，但是运行时间要长得多。
2. 使用[Brotli](https://github.com/google/brotli)进行压缩，文件大小明显优于Zopfli和`gzip --best`，并且压缩时间也在两者之间。这种[（新的）Brotli压缩工具](https://github.com/andybalholm/brotli)看起来很合理。

#### 压缩对比

例子1

|大小|命令|压缩时间|
|---|---|---|
|16M|未压缩大小|N/A|
|2.4M|brotli -o test.wasm.br test.wasm|53.6s
|3.3M|go-zopfli test.wasm|3m 2.6s
|3.4M|gzip --best test.wasm|2.5s
|3.4M|gzip test.wasm|0.8s

例子2

|大小|命令|压缩时间|
|---|---|---|
|2.3M|未压缩大小|N/A
|496K|brotli -o main.wasm.br main.wasm|5.7s
|640K|go-zopfli main.wasm|16.2s
|660K|gzip --best main.wasm|0.2s
|668K|gzip main.wasm|0.2s

使用 `https://github.com/lpar/gzipped` 之类的东西来自动提供带有正确标题的压缩文件（如果可用）。

### 2. 使用[TinyGo](https://github.com/tinygo-org/tinygo)生成Wasm文件

TinyGo支持面向嵌入式设备的Go语言的子集，并具有WebAssembly输出目标。

尽管确实有局限性（尚无完整的Go实现），但它仍然相当强大，并且生成的Wasm文件很小。 10kB并不罕见。 “ Hello world”示例为575字节。如果您将`gz -6`设为gz，它会下降到408个字节。

该项目也非常积极地开发，因此其功能正在迅速扩展。有关将WebAssembly与TinyGo结合使用的更多信息，请参见[这里](https://tinygo.org/webassembly/webassembly/)。

## 其他关于WebAssembly的资源

- [Awesome-Wasm](https://github.com/mbasso/awesome-wasm)-大量其他Wasm资源的清单。不是具体针对Go语言的。
