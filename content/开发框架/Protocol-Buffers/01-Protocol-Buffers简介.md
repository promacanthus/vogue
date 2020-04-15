---
title: 01-Protocol-Buffers简介
date: 2020-04-14T10:09:14.254627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Protocol-Buffers
summary: 01-Protocol-Buffers简介
showInMenu: false

---

`Protocol buffers`是一种灵活，高效，自动化的机制，用于序列化结构化的数据（如`XML`），但是它更小，更快，更简单。可以定义数据如何被结构化，然后使用特定的生成的源代码轻松地将结构化数据在各种数据流中写入和读取，这支持各种编程语言。甚至可以更新数据结构，而不会破坏根据“旧”格式编译的已部署程序。

## Protocol buffers如何工作

通过在`.proto`文件中定义`protocol buffers`消息类型来指定希望如何构建序列化信息。每个`protocol buffers`消息都是一个小的逻辑信息记录，包含一系列**名称-值**对。以下是`.proto`文件的一个非常基本的示例，该文件定义了包含有关人员信息的消息：

```go
message Person {
  required string name = 1;
  required int32 id = 2;
  optional string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    required string number = 1;
    optional PhoneType type = 2 [default = HOME];
  }

  repeated PhoneNumber phone = 4;
}
```

如上所示，消息格式很简单：

- 每种消息类型都有一个或多个**唯一编号**的字段，
  - 每个字段都有一个**名称**和一个**值类型**，
    - 其中值类型可以是：
      - 数字（整数或浮点数）
      - 布尔值
      - 字符串
      - 原始字节
      - 甚至（如上例所示）其他`protocol buffers`消息类型，允许分层次地构建数据

可以指定:

- 可选字段（optional）
- 必填字段（required）
- 重复字段（repeated）

可以在[`proto3指南`](../Protocol-Buffers/02-proto3指南.md)中找到有关编写`.proto`文件的更多信息。

一旦定义了消息，就可以在`.proto`文件上运行相应编程语言的`protocol buffers`编译器来生成数据访问类，他们为每个字段提供了简单的访问器，如`name()`和`set_name()`，以及将整个结构序列化或解析为原始字节的方法。

> 例如，如果选择的语言是C++，在上面的例子上运行编译器将生成一个名为`Person`的类。然后，可以在应用程序中使用此类来填充，序列化和检索`Person`的`protocol buffers`消息。如下所示的代码。

```c++
Person person;
person.set_name("John Doe");
person.set_id(1234);
person.set_email("jdoe@example.com");
fstream output("myfile", ios::out | ios::binary);
person.SerializeToOstream(&output);
```

可以在以下位置阅读消息：

```c++
fstream input("myfile", ios::in | ios::binary);
Person person;
person.ParseFromIstream(&input);
cout << "Name: " << person.name() << endl;
cout << "E-mail: " << person.email() << endl;
```

可以在不破坏向后兼容性的情况下为邮件格式添加新字段，**旧的二进制文件在解析时只是忽略新字段**。因此，如果通信协议使用`protocol buffers`作为数据格式，则可以扩展协议，而无需担心破坏现有代码。

将在[API参考](https://developers.google.com/protocol-buffers/docs/reference/overview)部分找到有关使用生成的`protocol buffers`代码的完整参考，在[`protocol buffers`编码](https://developers.google.com/protocol-buffers/docs/encoding)中找到有关`protocol buffers`消息如何编码的更多信息。

## 为什么不直接用XML

对于序列化结构化数据，`protocol buffers`比`XML`具有许多优点：

1. 更简单
2. 缩小3~10倍
3. 快20~100倍
4. 更少歧义
5. 生成更易于编程的数据访问类型

例如，假设要为具有`name`和`email`的`Person`建模。在XML中，需要：

```xml
  <person>
    <name>John Doe</name>
    <email>jdoe@example.com</email>
  </person>
```

而相应的`protocol buffers`消息（`protocol buffers`文本格式）是：

```go
# protocol buffer的文本表示
# 这不是在实际传输中的二进制格式
person {
  name: "John Doe"
  email: "jdoe@example.com"
}
```

当此消息被编码为`protocol buffers`二进制格式（上面的文本格式只是方便人类可读的表示形式，用于调试和编辑）时，它可能是28字节长并且需要大约100-200纳秒来解析。如果删除空格，`XML`版本至少为69个字节，并且需要大约5,000-10,000纳秒才能解析。

此外，操作`protocol buffers`要容易得多：

```c++
cout << "Name: " << person.name() << endl;
cout << "E-mail: " << person.email() << endl;
```

而使用`XML`，必须执行以下操作：

```c++
  cout << "Name: "
       << person.getElementsByTagName("name")->item(0)->innerText()
       << endl;
  cout << "E-mail: "
       << person.getElementsByTagName("email")->item(0)->innerText()
       << endl;
```

但是，`protocol buffers`并不总是比`XML`更好的解决方案。例如:

1. `protocol buffers`不是使用标记对基于文本的文档（例如`HTML`）建模的好方法，**因为无法轻松地将结构与文本交错**。
2. `XML`是人类可读和可编辑的; `protocol buffers`在它原生的格式中不是人类可读和可编辑的。
3. `XML`在某种程度上也是自我描述的。只有拥有消息定义（如`.proto文件`）时，`protocol buffers`才有意义。

## 如何开始使用

[下载这个包](https://developers.google.com/protocol-buffers/docs/downloads)，它包含`Java`，`Python`和`C++`版本的`protocol buffers`编译器的完整源代码，以及`I/O`和测试所需的类。要构建和安装编译器，请按照自述文件中的说明进行操作。

完成所有设置后，请尝试按照所选语言的[教程](https://developers.google.com/protocol-buffers/docs/tutorials)进行操作，这将指导你创建一个使用`protocol buffers`的简单应用程序。

## `proto3`简介

最新的版本3[发布](https://github.com/protocolbuffers/protobuf/releases)了一个新的语言版本：`Protocol Buffers语言版本3`（简称`proto3`），以及现有语言版本（简称`proto2`）中的一些新功能。`Proto3`简化了`protocol buffers`语言，既易于使用，又可以在更广泛的编程语言中使用：当前的版本允许使用`Java`，`C++`，`Python`，`Java Lite`，`Ruby`，`JavaScript`，`Objective`和`C＃`来生成`protocol buffers`代码。此外，可以使用最新的`Go protoc`插件为Go生成`proto3`代码，该插件可从[`golang/protobuf`](https://github.com/golang/protobuf) Github存储库获得。更多的语言正在筹备中。

请注意，两种语言版本的API不完全兼容。为避免给现有用户带来不便，将继续在新`protocol buffers`版本中支持以前的语言版本。

可以在[发行说明](https://github.com/protocolbuffers/protobuf/releases)中看到与当前默认版本的主要差异，并了解`Proto3`语言指南中的`proto3`语法。proto3的完整文档即将推出！

（如果名称`proto2`和`proto3`看起来有点令人困惑，那是因为最初开源`protocol buffers`时，它实际上是Google的第二版语言，也称为`proto2`，这也是开源版本号从v2开始的原因。

## 一点小历史

`protocol buffers`最初是在Google开发的，用于处理索引服务器请求/响应协议。在`protocol buffers`之前，有一种请求和响应的格式，它使用请求和响应的手动编组/解组，并支持许多版本的协议。 这导致一些非常丑陋的代码，如：

```go
if (version == 3) {
   ...
 } else if (version > 4) {
   if (version == 5) {
     ...
   }
   ...
 }
```

明确格式化的协议也使新协议版本的推出变得复杂，因为开发人员必须确保请求的发起者和处理请求的实际服务器之间的所有服务器都能理解新协议，然后才能切换以开始使用新协议。

`protocol buffers`开发用于解决这些问题：

1. 可以轻松引入新字段，而中间服务器不需要检查数据就可以简单地解析它并传递数据而无需了解所有字段。
2. 格式更具自我描述性，可以用各种语言（C++，Java等）处理。

但是，用户仍然需要自己手写解析代码。

随着系统的发展，它获得了许多其他功能和用途：

1. 自动生成序列化和反序列化代码避免了手动解析的需要。
2. 除了用于短生命周期的RPC（远程过程调用）请求之外，大家已经开始使用`protocol buffers`作为一种方便的自描述格式用于持久存储数据（例如，在`Bigtable`中）。
3. 首先服务的RPC接口被声明为`protocol`文件的一部分，然后使用`protocol`编译器生成`stub`类，用户可以通过服务接口的实际实现来覆盖这些`stub`类。

现在`protocol buffers`是Google的数据通用语言，在撰写本文时，Google代码树中定义了306,747种不同的消息类型，跨348,952个`.proto`文件。它们既可用于RPC系统，也可用于各种存储系统中的数据持久存储。
