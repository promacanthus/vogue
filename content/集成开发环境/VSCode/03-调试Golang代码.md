---
title: 03-调试Golang代码
date: 2020-04-14T10:09:14.278627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- VSCode
summary: 03-调试Golang代码
showInMenu: false

---

## 安装`delve`

有两种方法可以安装`delve`:

1. 执行指令`Go: Install/Update Tools`，然后选择`dlv`,点击`ok`来安装或者更新`delve`
2. 参照[安装指南](https://github.com/derekparker/delve/tree/master/Documentation/installation)手动安装

## 在`settings`中进行设置

调试器使用以下的设置项，你可能不需要添加/更改它们中的任何一个，以便在简单的情况下进行调试，但有时会给它们读取一些配置：

- `go.gopath`，查看扩展程序中设置`GOPATH`]
- `go.inferGopath`，查看扩展程序中设置`GOPATH`
- `go.delveConfig`：
  - `apiVersion`：控制将要在无头服务中被使用的`delve apis`的版本，默认值为2。
  - `dlvLoadConfig`：当`apiVersion`为1时不适用。这个配置会传递给`delve`。控制`delve`的[各种功能](https://github.com/Microsoft/vscode-go/blob/0.6.85/package.json#L431-L468)，这些功能会影响调试窗格中显示的变量。
    - `maxStringLen`：从字符串中读取的最大字节数
    - `maxArrayValues`：从数组，切片或`map`中读取的最大元素数
    - `maxStructFields`：从结构体中读取的最大字段数，`-1`表示将读取所有字段
    - `maxVariableRecurse`：评估嵌套类型时递归的程度
    - `followPointers`：请求指针自动解除引用

一些常见的情况，可能想要调整传递给delve的配置：

- 在调试模式中查看变量时，更改字符串和数组长度的默认上限64。
- 在调试模式中检查和评估嵌套的变量。

## 在`launch.json`中设置配置项

安装`delve`后，运行命令`Debug：Open launch.json`。如果还没有`launch.json`文件，则会创建一个具有以下默认配置的文件，该配置可用于调试当前包。

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": {},
            "args": []
        }
    ]
}
```

以下是可以在调试配置中调整的一些常见属性的更多信息：

|属性|描述|
|---|---|
|name|Debug视图下拉列表中显示的配置的名称
|type|始终设置为`“go”`，VS Code使用它来确定应该使用哪个扩展程序来调试代码
|request|`launch`或者`attach`，使用`attach`当需要附加到一个已经在运行中的程序
|mode|对于上个属性中的`launch` ，可选择`auto`,`debug`,`remote`,`test`,`exec`；对于上个属性中的`attach`，可选择`local`，`remote`
|program|在`debug`和`test`模式下，指定要调试的包或文件的绝对路径，或者在`exec`模式下调试的预构建二进制文件的绝对路径；不适用于`attach`
|env|调试时要使用的环境变量，如`{ "ENVNAME": "ENVVALUE" }` 
|envFile|包含环境变量定义的文件的绝对路径，`env`属性中传递的环境变量会覆盖此文件中的变量
|args|将传递给正在调试的程序的命令行参数数组
|showLogs|布尔值，指示是否应在调试控制台中打印来自`delve`的日志
|logOutput|以逗号分隔的`delve`组件列表（`debugger`，`gdbwire`，`lldbout`，`debuglineerr`，`rpc`），当设置为true时应生成调试输出
|buildFlags|要传递给Go编译器的构建(`build`)标志
|remotePath|远程调试时远程计算机上正在调试的文件的绝对路径，即`mode`设置为`remote`。有关详细信息，请参阅“[远程调试](#远程调试)”一节
|processeid|仅在使用`local`mode的`attach`request时适用，这是运行需要调试的可执行文件的进程的id

## 在调试设置中使用`VSCode`变量

调试配置中采用文件夹/文件路径的任何属性都可以使用以下`VSCode`变量：

- `${workspaceFolder}`：调试在`VSCode`中打开的工作空间的根目录下的程序包
- `${file}`：调试当前文件
- `${fileDirname}`：调试当前文件所属的包

### 使用构建标签

如果构建时需要构建标记（例如`go build -tags=whatever_tag`），那么添加参数`buildFlags`，内容为“`-tags=whatever_tag`”。

支持多个标记，将它们用双引号括起来，如下所示：`“-tags='first_tag second_tag third_tag'”`。

## 调试配置的代码段

编辑`launch.json`文件时，可以使用代码段进行调试配置。输入`“Go”`，将获得用于调试当前文件/包，测试功能等的片段。

### 用于调试当前文件的示例配置

```json
{
    "name": "Launch file",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${file}"
}
```

### 用于调试单个测试的示例配置

```json
{
    "name": "Launch test function",
    "type": "go",
    "request": "launch",
    "mode": "test",
    "program": "${workspaceFolder}",
    "args": [
        "-test.run",
        "MyTestFunction"
    ]
}
```

### 用于调试包中所有测试的示例配置

```json
{
    "name": "Launch test package",
    "type": "go",
    "request": "launch",
    "mode": "test",
    "program": "${workspaceFolder}"
}
```

### 用于调试预构建二进制文件的示例配置

```json
{
    "name": "Launch executable",
    "type": "go",
    "request": "launch",
    "mode": "exec",
    "program": "absolute-path-to-the-executable"
}
```

### 使用`processId`附加到已在运行的本地进程的示例配置

```json
{
    "name": "Attach to local process",
    "type": "go",
    "request": "attach",
    "mode": "local",
    "processId": 0
}
```

## 远程调试

要使用`VSCode`进行远程调试，必须首先在目标计算机上运行无头的`Delve`服务器。例如：

```bash
dlv debug --headless --listen=:2345 --log --api-version=2
```

要传递给正在调试的程序的任何参数都必须传递给在目标计算机上运行的`Delve`服务器。例如：

```bash
dlv debug --headless --listen=:2345 --log -- -myArg=123
```

然后，在VS Code `launch.json`中创建远程调试配置。

```json
{
    "name": "Launch remote",
    "type": "go",
    "request": "launch",
    "mode": "remote",
    "remotePath": "absolute-path-to-the-file-being-debugged-on-the-remote-machine",
    "port": 2345,
    "host": "127.0.0.1",
    "program": "absolute-path-to-the-file-on-the-local-machine",
    "env": {}
}
```

- 上面的示例在同一台机器上本地运行无头`dlv`服务器和VScode调试器。根据在远程计算机上的设置更新**端口**和**主机**。
- `remotePath`应指向远程计算机中正在调试的文件（在源代码中）的绝对路径
- `program`应指向本地计算机上与`remotePath`中对应的文件的绝对路径

当选择此新的`Launch remote`目标启动调试器时，`VS Code`会将调试命令发送到之前启动的`dlv`服务器，而不是针对应用启动它自己的`dlv`实例。

请参阅[这篇文章](https://github.com/lukehoban/webapp-go/tree/debugging)中调试在docker主机中运行的进程的示例。

## 故障排除

如果在调试Go代码时遇到问题，请首先尝试更新的`delve`版本，以确保使用的是最新版本，并且已使用当前的Go版本进行编译。要执行此操作，请运行命令`Go：Install/Update Tools`，选择`dlv`，然后点击`ok`。

### 启动调试日志

- 将调试配置中的`showLog`属性设置为`true`。将从`delve`中看到调试控制台中的日志。
- 在调试配置中设置`trace`属性以进行记录。将从Go扩展程序的调试适配器中看到调试控制台中的日志。这些日志将保存到一个文件，该文件的路径将在调试控制台的开头打印。
- 将调试配置中的`logOutput`属性设置为`rpc`。将看到对应于在`VS Code`和`delve`之间来回传递的`RPC`消息的日志。请注意，首先需要将`showLog`设置为`true`。
  - `logOutput`属性对应于`delve`使用的`--log-output`标志，可以是一个逗号分隔的组件列表，它应该生成调试输出。

### 使用源代码调试调试器

如果想深入挖掘并使用此扩展程序的源代码调试调试器，请参阅[build-and-debugging-the-extension](https://github.com/Microsoft/vscode-go/wiki/Building,-Debugging-and-Sideloading-the-extension-in-Visual-Studio-Code#building-and-debugging-the-extension)。

### 常见问题

1. 调试二进制文件时未验证断点或变量

> 确保正在调试的二进制文件是在没有优化的情况下构建的。在构建二进制文件时使用标志`-gcflags="all = -N -l"`。

2. Cannot find package ".." in any of ...

调试器没有使用正确的`GOPATH`。这不应该发生，如果确实如此，请记录一个错误。

> 在记录的错误被解决之前，解决方法是在`launch.json`文件的`env`属性中将`GOPATH`添加为`env var`。

3. Failed to continue: "Error: spawn EACCES"

`dlv`从命令行运行得很好，但`VS Code`提供此访问相关的错误。如果扩展程序试图从错误的位置运行`dlv`二进制文件，则会发生这种情况。Go扩展首先尝试在`$GOPATH/bin`中找到`dlv`，然后在`$PATH`中找到。

> 解决方案：在命令行中运行`dlv`。如果这与`GOPATH/bin`不匹配，则删除`GOPATH/bin`中的`dlv`文件

4. could not launch process: stat ***/debug.test: no such file or directory

可以在调试控制台中看到此信息，同时尝试在测试模式下运行。当程序属性指向没有测试文件的文件夹时，会发生这种情况。

> 解决方案：确保程序属性指向包含要运行的测试文件的文件夹

5. could not launch process: could not fork/exec
   1. OSX

    由于签名问题，这通常发生在OSX中。请参阅讨论，请参阅[＃717](https://github.com/Microsoft/vscode-go/issues/717)，[＃269](https://github.com/Microsoft/vscode-go/issues/269)和[derekparker/delve/357](https://github.com/derekparker/delve/issues/357)

    > 解决方案：可能必须卸载`dlv`并按照[说明](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md#manual-install)手动安装它

    2. Linux/Docker

    Docker具有安全设置，可以在容器内默认阻止`ptrace(2)`操作。

    > 解决方案：要不安全地运行容器，请在启动时将`--security-opt=seccomp:unconfined`传递给`docker`。参考：[derekparker/delve/515](https://github.com/derekparker/delve/issues/515)

6. could not launch process: exec: "lldb-server": executable file not found in $PATH

> 对于使用版本`0.12.2`或更高版本的版本的`Mac`用户，此错误可能会显示。不知道为什么，但是做`xcode-select --install`已经解决了看到这个问题的用户的问题。

7. 远程调试未验证的断点

> 检查`remote delve process`使用的`delve api`的版本，即检查标志`-api-version`的值。这需要匹配Go扩展程序使用的版本，默认情况下使用版本2。可以通过编辑`launch.json`文件中的调试配置来更改扩展使用的`api`版本。

8. 尝试使用终端/命令行中的`dlv`

> 将`"trace"："log"`添加到调试配置并在`VS Code`中调试。这会将日志发送到调试控制台，可以在其中查看对`dlv`的实际调用。可以复制并在终端中运行它
