# gRPC

## RPC框架原理

RPC框架的目标就是让**远程服务调用更加简单、透明**，RPC框架负责屏蔽底层的传输方式（**TCP**或者**UDP**）、序列化方式（**XML**/**Json**/**二进制**）和通信细节。服务调用者可以像调用本地接口一样调用远程的服务提供者，而不需要关心底层通信细节和调用过程。

RPC框架的调用原理图如下所示：

![image](/images/RPC.png)

## 主流的RPC框架

业界主流的RPC框架整体上分为三类：

1. 支持**多语言**的RPC框架:
   1. Google的**gRPC**
   2. Apache（Facebook）的**Thrift**
2. 只支持特定语言的RPC框架:新浪微博的**Motan**
3. 支持**服务治理**等服务化特性的分布式服务框架，其底层内核仍然是RPC框架:阿里的
**Dubbo**

随着微服务的发展，基于语言中立性原则构建微服务，逐渐成为一种主流模式，例如:

- 对于后端并发处理要求高的微服务，比较适合采用**Go语言**构建
- 对于前端的Web界面，则更适合**Java**和**JavaScript**

因此，基于多语言的RPC框架来构建微服务，是一种比较好的技术选择。

>例如Netflix的API服务编排层和后端的微服务之间就采用gRPC进行通信。

## gRPC简介
本文档向介绍gRPC和protocol buffers。gRPC可以使用protocol buffers作为其接口定义语言（Interface Definition Language,IDL）和其基础消息交换格式。

在gRPC中，客户端应用程序可以直接调用不同计算机上的应用程序中的方法，就像它是本地对象一样，可以更轻松地创建分布式应用程序和服务。与许多RPC系统一样，gRPC**基于定义服务的思想**，指定可以使用其参数和返回类型远程调用的方法。

- 在服务端，服务器实现此接口并运行gRPC服务器来处理客户端调用
- 在客户端，客户端有一个存根（在某些语言中称为客户端），它提供与服务器相同的方法

gRPC的调用示例如下所示：

![image](/images/landing-2.svg)

gRPC客户端和服务端可以在各种环境中相互运行和通信（从Google内部的服务器到桌面应用），并且可以使用任何gRPC支持的语言编写。因此，可以使用Go，Python或Ruby轻松创建gRPC客户端与使用Java编写的gRPC服务端通信。此外，最新的Google API将具有gRPC版本的接口，可以轻松地在编写的应用程序中构建Google提供的功能和服务。

## gRPC特点

1. 语言中立，支持多种语言；
2. 基于IDL文件定义服务，通过proto3工具生成指定语言的数据结构、服务端接口以及客户端Stub；
3. 通信协议基于标准的HTTP/2设计，支持双向流、消息头压缩、单TCP的多路复用、服务端推送等特性，这些特性使得gRPC在移动端设备上更加省电和节省网络流量；
4. 序列化支持PB（ProtocolBuffer）和JSON，PB是一种语言无关的高性能序列化框架，基于HTTP/2 + PB，保障了RPC调用的高性能。

## 使用Protocol Buffers

默认情况下，gRPC使用[protocol buffers](https://developers.google.com/protocol-buffers/docs/overview)，这是Google成熟的开源机制，用于序列化结构化数据（尽管它可以与其他数据格式（如JSON）一起使用）。

使用protocol buffers的第一步是：定义要在proto文件中序列化的数据的结构（这是一个扩展名为`.proto`的普通文本文件）。

protocol buffers数据被构造为消息，其中每个消息是包含一系列称为字段的键值对的信息小的逻辑记录。这是一个简单的例子：

```go
message Person {
  string name = 1;
  int32 id = 2;
  bool has_ponycopter = 3;
}
```

一旦指定了数据结构，就可以使用protocol buffers编译器`protoc`从原型定义生成**首选语言**的**数据访问类**:

1. 为每个字段提供了简单的访问器，如`name()`和`set_name()`
2. 将整个结构序列化/解析为原始字节的方法

    > 例如，如果选择的语言是C++，则运行编译器上面的例子将生成一个名为`Person`的类。然后，可以在应用程序中使用此类来填充，序列化和检索Person protocol buffers消息。

正如将在示例中更详细地看到的那样，可以在普通的proto文件中定义gRPC服务，并将RPC方法参数和返回类型指定为protocol buffers消息：

```go
// The greeter service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
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

gRPC还可以使用带有特殊gRPC插件的`protoc`来生成proto文件中的代码。如果，使用gRPC插件，将获得生成的gRPC客户端和服务端代码，以及用于填充，序列化和检索消息类型的常规protocol buffers代码。
我们将在下面更详细地看一下这个例子。

可以在[Protocol Buffers文档](https://developers.google.com/protocol-buffers/docs/overview)中找到有关protocol buffers的更多信息，并了解如何使用所选语言的快速入门来获取和安装带有gRPC插件的`protoc`。

## protocol buffers的版本

虽然在很早之前，protocol buffers已经可供开源用户使用，但在示例中使用了一种新的protocol buffers，称为proto3，它具有略微简化的语法，一些有用的新功能，并支持更多语言。

目前提供Java，C++，Python，Objective-C，C＃，lite-runtime（Android Java），Ruby和JavaScript，都来自[protocol buffers GitHub repo](https://github.com/google/protobuf/releases)，以及来自[golang/protobuf GitHub](https://github.com/golang/protobuf)的Go语言生成器repo，还有更多语言在开发中。

可以在[proto3语言指南](https://developers.google.com/protocol-buffers/docs/proto3)和每种语言的[参考文档](https://developers.google.com/protocol-buffers/docs/reference/overview)中找到更多信息，[参考文档](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec)还包括`.proto`文件格式的正式规范。

通常，虽然可以使用proto2（当前的默认protocol buffers版本），但建议将proto3与gRPC一起使用，因为它允许使用全系列的gRPC支持的语言，以及避免与proto2客户端与proto3服务端通信时的兼容性问题，反之亦然。
