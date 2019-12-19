# 01-WebAssembly

## 简介

Go 1.11向WebAssembly添加了一个实验端口。 Go 1.12对它的某些部分进行了改进，并有望在Go 1.13中进行进一步的改进。

WebAssembly在其[主页](https://webassembly.org/)上描述为：

> WebAssembly（缩写为Wasm）是基于堆栈的虚拟机的二进制指令格式。 Wasm被设计为可移植目标，用于编译高级语言（如C/C ++ / Rust），从而可以在Web上为客户端和服务器应用程序进行部署。

如果不熟悉WebAssembly，请阅读下面的“入门”部分，观看下面的一些Go WebAssembly讲座，然后查看下面的更多示例。

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

