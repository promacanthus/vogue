# Go Modules

Go 语言中一直被人诟病的一个问题就是**没有一个比较好用的依赖管理系统**。

`GOPATH` 的设计让开发者一直有很多怨言，在 Go 语言快速发展的过程中也出现了一些比较优秀的依赖管理工具，比如:

- govendor
- dep
- glide等

有一些差不多成了半官方的工具了，但是这些工具都还是需要依赖于`GOPATH`。

随着 Go1.11 的发布，Golang 官方给我们带来了依赖管理的全新特性 `Go Modules`，这是 Golang 全新的一套依赖管理系统。

## 新建Modules

要使用 Go Modules 首先需要保证你环境中 Golang 版本大于 `1.11`：

```bash
$ go version
go version go1.12.7 linux/amd64
```

 Go Modules 主要就是为了消除 GOPATH 的，所以**新建的项目可以完全不用放在 `$GOPATH/src` 目录下面**，任何地方都可以。

### 第一步

在`bashrc`中配置环境变量以启动Go Modules:

```bash
vim ~/bashrc

export GO111MODULE=on
```

编写库函数：

```go
package stringsx
import (
 "fmt"
)
func Hello(name string) string{
    return fmt.Sprintf("Hello, %s", name), nil
}
```

### 第二步

在项目根目录初始化GO Modules:

```bash
go mod init <模块名>

# 该命令会在当前目录下生成go.mod文件
# 生成的内容就是包含一个模块名称的声明

module <模块名>
```

**注意：模块名非常重要，相当于声明了模块的名称，后面使用该模块就需要使用这个名称来获取模块**。

### 第三步

上述第二步，将当前包变成了一个Module，将代码推送到仓库，以Github为例，仓库地址为：https://github.com/sample_golang_module/example。

```bash
git init
git add .
git commit -am "init commit"
git remote add origin git@github.com:sample_golang_module/example.git
git push -u origin master
```

### 第四步

至此，完成了最简单的Go Module编写，其他任何开发者想要使用这个模块，通过`go get`命令来获取：

```bash
go get github.com/sample_golang_module/example
```

> 上面的命令是获取master分支的最新代码，这没有问题但不是最佳实践，模块可能要更新内容或者修复Bug，放在master分支会造成使用者的混乱，很可能使用者的代码在模块更新后就不兼容了，**Go Modules可以很好的解决版本问题**。

## Module版本管理

Go Modules是需要进行版本化管理的，强烈推荐使用[语义化版本控制](https://semver.org)，最主要的版本规则如下：

版本格式：主版本号.次版本号.修订号

版本号递增规则如下：

- 主版本号：当你做了不兼容的API修改
- 次版本号：当你做了向下兼容的功能性新增
- 修订号：当你做了向下兼容的问题修正

> 先行版本号及版本编译元数据可以加到“主版本号.次版本号.修订号”的后面，作为延伸。

在使用Go Modules查找版本的时候，会使用仓库中的tags并且某些版本和其他版本有一些不同之处，比如V2或者更高的版本要和V1的版本模块的导入路径是不同的，这样才能通过模块来区分使用的是不同的版本，**默认情况下，Golang会获取仓库中最新的tag版本**。

> 最重要的一点，发布模块时，要使用`git tag`来标记仓库的版本。

### 发布第一个版本

发布release包，需要给当前的包打上tag，使用语义化的版本，如v1.0.0：

```bash
git tag v1.0.0
git push --tags
```

更好的方法是创建一个名叫v1的新分支，这样可以方便以后修复当前版本代码中的Bug，也不会影响到master或者其他分支的代码：

```bash
git checkout -b v1
git push -u origin v1
```

### 模块使用

模块已经准备好，创建一个简单的程序来使用上面的模块：

```go
package main
import (
	"fmt"
	"github.com/sample_golang_module/example/stringsx"
)
func main() {
	fmt.Println(stringsx.Hello("cnych"))
}
```

在程序中使用`github.com/sample_golang_module/example/stringsx`这个包，在导入之前，使用`go get`命令将这个包拉到`GOPATH`或者`vendor`目录下即可，将这个包当成module来使用。

```bash
# 在当前项目下初始化module
go mod init    # 不写模块名，则与package名相同

go run main.go

go: finding github.com/sample_golang_module/example v1.0.0
go: downloading github.com/sample_golang_module/example v1.0.0
Hello, cnych

```

执行完成后，上面的命令会自动下载程序中导入的包，下载完成后查看当前项目的`go.mod`文件：

```vim
module <当前项目模块名>
require github.com/cnych/stardust v1.0.0
```

并且还在当前目录下生成一个名为`go.sum`的新文件，里面包含了依赖包的一些hash信息，用来确保文件和版本的正确性：

```vim
github.com/sample_golang_module/example v1.0.0 h1:8EcmmpIoIxq2VrzXdkwUYTD4OcMnYlZuLgNntZ+DxUE=
github.com/sample_golang_module/example v1.0.0/go.mod h1:Qgo0xT9MhtGo0zz48gnmbT9XjO/9kuuWKIOIKVqAv28=
```

模块会被下载到`$GOPATH/pkg/mod`目录下：

```bash
$ ls $GOPATH/pkg/mod/github.com/sample_golang_module/example
example@v1.0.0
```

这样就成功使用了编写的模块。

## 发布一个bugfix版本

发现模块中的Hello函数有bug，需要修复并发布一个新版本：

```go
func Hello(name string) string{
	return fmt.Sprintf("Hello, %s!!!", name)
}
```

在v1这个分支上进行fix，完成之后merge到master分支上去，然后发布一个新的版本，遵从语义化版本骨子额，修正一个bug之后需要添加修正版本号，即v1.0.1：

```bash
git add .
git commit -m "fix Hello function #123"
git tag v1.0.1
git push --tags origin v1
```

## 更新modules

默认情况下，Golang不会自动更新模块，如果自动更新会造成版本管理混乱，所以需要明确告知Golang需要更新的模块，通过以下几种方式：

1. 运行`go get -u xxx`命令来获取最新版的模块
2. 运行`go get package@version`命令来更新指定版本的模块
3. 直接更新`go.mod`文件中的模块依赖版本，然后执行`go mod tidy`命令来更新

更新之后，`go.mod`文件中的依赖模块的版本会变化，`￥GOPATH/pkg/mod`文件夹中会增加新版本的模块。

> 在Go Modules中，每个版本都是独立的文件夹，这样就不会出现版本冲突。

## 主版本升级

根据语义化版本规则，主版本升级的不向后兼容的，从`Go Modules`的角度来看，主版本是一个完全不同的模块了，因为两个大版本之间是互相不兼容的。

修改模块中的Hello函数：

```go
func Hello(name, lang string) (string, error) {
	switch lang {
	case "en":
		return fmt.Sprintf("Hi, %s!", name), nil
	case "zh":
		return fmt.Sprintf("你好, %s!", name), nil
	case "fr":
		return fmt.Sprintf("Bonjour, %s!", name), nil
	default:
		return "", fmt.Errorf("unknow language")
	}
}
```

> 这里需要切换到master分支进行修改，因为v1分支和现在修改的内容是完全不同的版本。

函数有两个参数，返回值也有两个，与v1不兼容。需要更新版本到v2.0.0，**通过更改v2版本的模块路径来区分两个大版本**。比如`github.com/sample_golang_module/example/v2`，这样v2版本的模块和v1版本的模块就是两个完全不同的模块了，在使用新版模块时在模块名称后面加上v2即可。

```vim
module github.com/cnych/stardust/v2
```

接下来的操作给当前版本添加一个v2.0.0的`git tag`或者创建一个名为v2的分支，这样可以将版本之间的影响降到最低：

```bash
git add .
git commit -m "change Hello function to support lang"
git checkout -b v2
git tag v2.0.0
git push origin v2 --tags
```

v2 版本的模块就发布成功，之前程序也不会有任何的影响，还是继续使用现有的 v1.0.1 版本，而且使用 `go get -u` 命令也不会拉取最新的 v2.0.0 版本代码。

用户要试用v2.0.0版本的模块，只需要单独引入v2版本的模块即可：

```go
package main
import (
	"fmt"
	"github.com/sample_golang_module/example/stringsx"
	stringsV2 "github.com/sample_golang_module/example/v2/stringsx"
)
func main() {
	fmt.Println(stringsx.Hello("cnych"))
	if greet, err := stringsV2.Hello("cnych", "zh"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(greet)
	}
}
```

go run main.go命令会自动拉取v2模块的代码：

```
$ go run main.go
go: finding github.com/sample_golang_module/example/v2 v2.0.0
go: downloading github.com/sample_golang_module/example/v2 v1.0.0
Hi, cnych！！！
你好, cnych!
```

在同一个 go 文件中就使用两个不兼容版本的模块。同样这个时候再次查看下 go.mod 文件的变化：

```vim
module <模块名>
require github.com/sample_golang_module/example v1.0.1
require github.com/sample_golang_module/example/v2 v2.0.0
```

默认情况下，Golang 是不会从 go.mod 文件中删除依赖项的，如果我们有不使用的一些依赖项需要清理，可以使用 tidy 命令：

```bash
go mod tidy
```

该命令会清除没有使用的模块，也会更新模块到指定的最新版本。

## Vendor

Go Modules 默认会忽略 `vendor/` 这个目录，但是如果还想将依赖放入 vendor 目录的话，可以执行下面的命令：

```
go mod vendor
```

该命令会在项目根目录下面创建一个 `vendor/` 的文件夹，里面会包含所有的依赖模块代码，并且会在该目录下面添加一个名为 modules.txt 的文件，用来记录依赖包的一些信息，比较类似于 govendor 中的 vendor.json 文件。

> 不过建议还是不要使用该命令，尽量去忘掉 vendor 的存在。

## 镜像仓库

如果有一些依赖包下载不下来的，我们可以使用 GOPROXY 这个参数来设置模块代理，比如：

```bash
$ export GOPROXY="https://goproxy.io"
```

阿里云也提供了 Go Modules 代理仓库服务：http://mirrors.aliyun.com/goproxy/，使用很简单就两步：

1. 使用 go1.11 以上版本并开启 go module 机制：`export GO111MODULE=on`
2. 导出 GOPROXY 环境变量：`export GOPROXY=https://mirrors.aliyun.com/goproxy/`


如果你想上面的配置始终生效，可以将这两条命令添加到. bashrc 中去。

## 搭建私有仓库

除了使用公有的 Go Modules 代理仓库服务之外，很多时候我们在公司内部需要搭建私有的代理服务，特别是在使用 CI/CD 的时候，如果有一个私有代理仓库服务，会大大的提供应用的构建效率。

可以使用 Athens 来搭建私有的代理仓库服务，搭建非常简单，直接用 docker 镜像运行一个服务即可：

```bash
export ATHENS_STORAGE=~/athens-storage
mkdir -p $ATHENS_STORAGE
docker run -d -v $ATHENS_STORAGE:/var/lib/athens \
 -e ATHENS_DISK_STORAGE_ROOT=/var/lib/athens \
 -e ATHENS_STORAGE_TYPE=disk \
 --name goproxy \
 --restart always \
 -p 3000:3000 \
 gomods/athens:latest
```

其中 ATHENS_STORAGE 是用来存放我们下载下来的模块的本地路径，另外 ATHENS 还支持其他类型的存储，比如 内存, AWS S3 或 Minio，都是 OK 的。

然后修改 GOPROXY 配置：
```bash
export GOPROXY=http://127.0.0.1:3000
```

## 总结

一句话：Go Modules 真的用起来非常爽，特别是消除了 GOPATH，这个东西对于 Golang 初学者来说是非常烦人的，很难理解为什么需要进入到特定目录下面才可以编写 Go 代码，现在不用担心了，直接使用 Go Modules 就行。
