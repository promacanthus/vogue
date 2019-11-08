# Go Modules

VS Code使用大量[Go工具](/VSCode-go/06-扩展程序依赖的Go工具.md)来提供代码导航，代码完成，构建，分析等功能。

> 这些工具还没有为[Go模块](https://blog.golang.org/modules2019)提供良好的支持。

新的语言服务器[`gopls`](https://github.com/golang/go/wiki/gopls)支持`Go Modules`。在`settings`中添加以下内容以使用它。

```json
"go.useLanguageServer": true
```

> 注意：使用`Go：Install/Update Tools`定期更新`gopls`版本，以便不断获得对语言服务器所做的改进。

使用`gopls`时`VS Code`中的已知问题:

- 无法自动添加未导入的包
- 查找引用并重命名仅在单个包中生效
- 更多查看这里：
  - [`golps`已知问题](https://github.com/golang/go/wiki/gopls#known-issues)
  - 在这个仓库中的已知问题都带有`go-modules`标签

要对语言服务器进行故障排除，请参阅[故障排除`gopls`](https://github.com/golang/go/wiki/gopls#troubleshooting)。

如果不想使用语言服务器，请了解并非此扩展所依赖的所有Go工具都支持Go模块。

https://github.com/golang/go/issues/24661 是Go工具团队用来跟踪各种Go工具中Go模块支持更新的问题列表。

## FAQ

1. 使用Go模块时可以使用语言服务器吗？

> 可以，这是`VS Code`中模块支持的前进之路。请注意语言服务器本身处于alpha模式，正在进行活动开发。有关详细信息，请参阅此问题上方的部分。

2. 使用Go模块时，为什么代码导航和代码补全速度会变慢？

> 使用Google的语言服务器时，代码导航和代码补全肯定会更好。所以，请试一试。如果不使用语言服务器，那么这主要由于`godef`和`gocode`的限制。Google的Go工具团队正在开发一种[语言服务器](https://godoc.org/golang.org/x/tools/cmd/gopls)，它将成为所有语言功能的长期解决方案。请按照本页第一部分的说明试用语言服务器。

如果不想使用语言服务器:

- 如果代码补全缓慢，请在[`gocode`仓库](https://github.com/stamblerre/gocode)中记录一个问题。
- 如果代码导航速度慢，请在`[godef仓库](https://github.com/rogpeppe/godef)`中记录问题，或者如果在设置中选择了`gogetdoc`，请在[`gogetdoc仓库](https://github.com/zmb3/gogetdoc)`中记录问题。

3. 文件保存时不再自动导入。为什么？

> 如果不使用语言服务器，则此扩展程序默认使用`goreturns`工具格式化文件并自动导入缺失的包。由于此工具不支持模块，因此文件保存中的自动导入功能不再有效。添加设置"`go.formatTool`"："`goimports`"，然后使用`Go：Install/Update Tools`来安装/更新`goimports`，因为它最近添加了对模块的支持。