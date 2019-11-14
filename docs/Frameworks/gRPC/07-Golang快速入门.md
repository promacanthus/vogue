# 07-Golang快速入门

本指南以简单的工作示例为您介绍Go中的gRPC。

## 先决条件

### Go 版本

gRPC要求Go版本大于等于1.6：

```bash
sugoi@Sugoi-PC:~$ go version 
go version go1.12.7 linux/amd64
```

有关安装说明，请遵循以下指南：[入门-Go编程语言](https://golang.org/doc/install)。

### 安装gRPC

使用下面的指令安装gRPC:

```bash
go get -u google.golang.org/grpc
```

### 安装protocol buffers v3

安装用于生成gRPC服务代码的protoc编译器。最简单的方法是从[此处](https://github.com/google/protobuf/releases)下载适用于您的平台的预编译二进制文件（`protoc-<version>-<platform>.zip`）：

1. 解压缩此文件
2. 更新环境变量PATH以包含protoc二进制文件的路径

```bash
# protoc environment
export PATH=$PATH:/opt/protoc/bin
```

接下来，为Go安装protoc插件。

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

protoc编译器插件`protoc-gen-go`将安装在`$GOBIN`中，默认为`$GOPATH/bin`,它必须在你的`$PATH`中，这样protocol的编译器`protoc`才能找到它。

```bash
# Golang environment
export GOPATH=/home/sugoi/go
export PATH=$PATH:/opt/go/bin:$GOPATH/bin
```

## 下载示例

使用`go get google.golang.org/grpc`获取的grpc代码也包含它示例。它们可以在示例`dir:$GOPATH/src/google.golang.org/grpc/examples`下找到。

## 编译示例

进入示例所在的目录下：

```bash
cd $GOPATH/src/google.golang.org/grpc/examples/helloworld
```

gRPC服务在`.proto`文件中定义，该文件用于生成相应的`.pb.go`文件。`.pb.go`文件是通过使用protocol编译器（`protoc`）编译`.proto`文件生成的。

此示例已生成`helloworld.pb.go`文件（通过编译`helloworld.proto`生成），并且可以在此目录中找到：`$GOPATH/src/google.golang.org/grpc/examples/helloworld/helloworld`

`helloworld.pb.go`文件包含：

- 生成客户端和服务端的代码
- 用于填充，序列化和检索`HelloRequest`和`HelloReply`消息类型的代码

## 尝试一下

使用`go run`代码来编译并执行服务端和客户端代码，在示例目录下执行如下命令：

```bash
# 在第一个命令行终端
go run greeter_server/main.go

# 在第二个命令行终端
go run greeter_client/main.go
```

如果一切顺利的话，将会在运行客户端的命令行终端中看到输出信息：`Greeting: Hello World`。

这样就成功的使用gRPC运行客户端-服务端应用程序。

## 更新gRPC服务端

现在让我们看看如何使用其他方法更新应用程序中服务端的代码以供客户端调用。我们的gRPC服务是使用protocol buffers定义的；可以在[gRPC简介](../gRPC/01-gRPC简介.md)和[gRPC基础:Go](https://grpc.io/docs/tutorials/basic/go/)中找到更多关于如何在`.proto`文件中定义服务的知识。现在需要知道的是，服务端和客户端的stub都有一个`SayHello` RPC方法，该方法从客户端获取`HelloRequest`参数并从服务器返回`HelloReply`，此方法定义如下：

```go
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) ;
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

让我们更新一下代码，使`Greeter`服务有两种方法。确保在示例目录下（`$GOPATH/src/google.golang.org/grpc/examples/helloworld`）。

编辑`helloworld/helloworld.proto`并使用新的`SayHelloAgain`方法更新它，新方法具有相同的请求和响应类型：

```go
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) ;
  // Sends another greeting
  rpc SayHelloAgain (HelloRequest) returns (HelloReply) ;
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

## 生成gRPC代码

接下来，更新应用程序使用的gRPC代码以使用新的服务定义，还在上面的例子目录（`$GOPATH/src/google.golang.org/grpc/examples/helloworld`）

```bash
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```

这会使我们的更改重新生成`helloworld.pb.go`文件。

## 更新并运行应用

现在有了新生成的服务端和客户端代码，还需要在示例应用程序的手动编写部分中实现并调用新方法。

### 更新服务端

编辑`greeter_server/main.go`并向其添加以下函数：

```go
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}
```

### 更新客户端

编辑`greeter_client/main.go`并向其添加以下函数：

```go
r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
if err != nil {
        log.Fatalf("could not greet: %v", err)
}
log.Printf("Greeting: %s", r.Message)
```

### 运行

```bash
# 第一个命令行终端，运行服务端
go run greeter_server/main.go

# 第二个命令行终端，运行客户端
go run greeter_client/main.go

# 得到如下输出
Greeting: Hello world
Greeting: Hello again world
```
