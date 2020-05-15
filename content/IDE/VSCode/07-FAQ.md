---
title: 07-FAQ
date: 2020-04-14T10:09:14.278627+08:00
draft: false
---

1. 安装了扩展程序，但是没有任何功能正常工作？

> 检查是否安装了所有依赖的Go工具。执行`Go: Install/Update Tools`命令。如果希望只是安装指定的工具，查看[插件依赖的Go工具](../06-插件依赖的go工具/)，然后手动安装需要的工具。

2. 安装工具时显示`git pull --ff-only`错误？

> 有可能时因为正在安装的工具库被强制push，删除相关工具在`$GOPATH/src`（如果设置了`go.toolsGopath`,相应的路径也需要检查）中的文件夹然后重新安装。

3. 为什么导入的软件包会被删除或者重新排序？

> 默认情况下，扩展程序会格式化代码并重新组织每个文件中导入的包（添加缺失的包、删除未用的包、重新排序导入包）。可以通过下面的设置来关闭整个功能。

- 如果没有使用`gopls`语言服务器，可以添加`"go.formatTools": "gofmt`到设置中来选择一个格式化工具或者禁用导入时格式化功能：

    ```json
    "[go]": {
    "editor.formatOnSave": false
    }
    ```

- 如果使用`gopls`语言服务器，那么重新组织导入的功能在格式化进程之外执行。要在这种情况下禁用整个功能，添加如下的设置：

    ```json
    "[go]": {
        "editor.codeActionsOnSave": {
        "source.organizeImports": false
        }
    }
    ```

4. 为什么在保存文件时，空格被替换为制表符？

> 因为此扩展程序使用的格式化工具`goreturns`，`goimports`或`gofmt`都遵循使用制表符而不是空格的规则。

5. 格式化工具中一个制表符的大小是否为8？

> 在VSCode中制表符的默认大小是4，要在每个Go文件中将制表符大小修改为8，那么添加如下配置：

```json
"[go]": {
  "editor.tabSize": 8
}
```

6. 扩展程序如何确定要使用的`GOPATH`?

> 查看扩展程序中设置GOPATH。

7. VSCode 是否支持 Go Modules?

> 查看VSCode支持Go-Models。

8. 使用Go Modules时，为什么代码导航和代码完成速度会变慢？

> 查看VSCode支持Go-Models。

9. 使用Go Modules时可以受用语言服务器吗？

> 查看VSCode中使用Go语言服务器。

10.  如何仅仅运行代码而不是调试代码？

> 使用快捷键`Ctrl+F5`或者执行`Debug: Start without Debugging`命令。

在这种场景下，我们使用`go run`指定一个文件路径。因此：

- 如果已经有一个带有默认配置的`launch.json`文件，请将其更新为使用`${file}`而不是`${fileDirname}`
- 如果使用自定义的调试配置，请配置的程序属性指向文件为不是目录

如果程序属性中没有指定文件路径，那么`Start without Debugging`命令将会恢复到正常的调试。

11. 为什么在终端中设置的`GOPATH`没有被插件使用？为什么程序没有获取在终端中设置的环境变量？

> 扩展程序运行在一个独立的进程中而不是在终端或者VSCode窗口的其他部分。因此在终端中设置的环境变量对于扩展程序时不可见的。

12. 在哪里可以查看到扩展程序的日志？

> 在`View`菜单中选择`Output`,将显示输出面板。在此面板的右上角，从下拉列表中选择`Log(Extension Host)`。**提示：如果在执行特定操作后查找日志，请先清除日志并再次尝试以降低其他日志的干扰**。

13. 自动补全功能不生效了，咋整？

> 首先检查日志中的错误。

- 如果正在使用Go语言服务器，那么上一步中的输出窗口面板的下拉带单中有与一个关于语言服务器的条目，选择并查看它的输出。

- 如果没有没有使用Go语言服务器，
  1. 并且这是第三方的包，首先确定是否安装了。可以通过执行`Go: Build Current Package`命令，这样将会安装所有的依赖，后者手动安装所有的依赖通过命令`go install`。
  2. 如果依然不起作用，在终端中执行`gocode close`或者`gocode exit`然后再试一次。如果使用`Go Modules`，那么使用`gocode-gomod`而不是`gocode`。
  3. 如果依然不起作用，执行`Go: Install/Updata Tools`,选择`gocode`来更新这些工具。如果使用`Go Modules`，那么使用`gocode-gomod`而不是`gocode`。
  4. 如果依然不起作用，在终端中执行`gocode -s -debug`后立即执行`gocode close`或者`gocode exit`然后再试一次。`gocode`的结果将输出再终端中。如果使用`Go Modules`，那么使用`gocode-gomod`而不是`gocode`。
  5. 如果再终端中看到预期的结果，但是在VSCode中没有预期的结果，那么在[vscode-go的仓库](https://github.com/Microsoft/vscode-go)和[gocode的仓库](https://github.com/mdempsky/gocode)中开一个`issue`。如果使用`Go Modules`在[这里](https://github.com/stamblerre/gocode)中开一个`issue`。

14. 为什么在文件保存时格式化不生效？

> 查看日志（操作见问题12），具体的消息如`"Formatting took too long"`或者`Format On Save feature could be aborted`。如果看到此类消息，很大概率因为格式化花太长时间而被中止，这影响到了保存体验。可以设置`editor.formatOnSaveTimeout`来控制超时参数。

15. 导入的包有红色的下划线显示“包未找到”？

> 这些是构建错误，点击`View`->`Output`->从面板右上角的下拉菜单中选择`go`。然后就可以看到`go build`的输出（或者是`go test`的输出，如果当前文件时测试文件）。将`go build`命令和参数一起拷贝后尝试在终端中运行。如果仍然看到相同的错误，那么问题在于`GOPATH`的设置。如果它运行正常，那么提交一个`issue`。

16. 如何获取已实现但尚未发布的功能/错误修复程序？如何获得Go扩展的测试版？

> 查看安装测试版本
