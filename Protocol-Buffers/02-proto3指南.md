# `proto3`指南

本指南介绍如何使用`protocol buffers`语言构建`protocol buffers`数据，包括：

- `.proto`文件语法
- 如何从`.proto`文件生成数据访问类

它涵盖了`protocol buffers`语言的`proto3`版本：有关较早的`proto2`语法的信息，请参阅[`Proto2`语言指南](https://developers.google.com/protocol-buffers/docs/proto)。

这是一个参考指南 - 对于使用本文档中描述的许多功能的分步示例，请参阅所选语言的[教程](https://developers.google.com/protocol-buffers/docs/tutorials)（目前仅限`proto`2;更多`proto3`文档即将推出）。

## 定义消息类型

首先让看一个非常简单的例子。假设要定义搜索请求消息格式，其中每个搜索请求都有：

- 一个字符串类型的查询
- 所查询的特定页码
- 每页返回的结果数

这是用于定义消息类型的`.proto`文件。

```protobuf
syntax = "proto3";      // 指定正在使用proto3语法，默认认为使用proto2，必须是文件的第一个非空注释行

message SearchRequest {     // 消息格式以名称-值对的形式指定三个字段
  string query = 1;         // 每个字段有一个名称和类型
  int32 page_number = 2;
  int32 result_per_page = 3;
}
```

### 指定字段类型

在上面的例子中，所有的字段都是[标量类型](##标量值类型)：两个整型一个字符串类型。同时也可以给字段指定组合类型（包括[枚举](##枚举)或其他类型）。

### 分配字段编号

如上所示，消息定义中的每个字段都定义一个**唯一的编号**。这些字段的编号用于在消息的[二进制格式](/Protocol-Buffers/04-编码.md)中标识字段，一旦消息类型被使用就不能再更改。请注意：

- `1到15`范围内的字段编号需要一个字节进行编码，包括字段的编号和字段的类型（可以在[`protocol buffers`编码](/Protocol-Buffers/04-编码.md#消息结构)中找到更多相关信息）
- `16到2047`范围内的字段编号占用两个字节。 因此，应该为非常频繁出现的消息元素保留数字1到15。请记住为将来可能添加的频繁出现的元素留出一些空间

可以指定的最小字段数为1，最大字段数为536,870,911（2的29次方-1）。

> 不能使用数字19000到19999（`FieldDescriptor::kFirstReservedNumber`到`FieldDescriptor::kLastReservedNumber`），因为它们是为`protocol buffers`实现而保留的。

如果在`.proto`中使用这些保留数字之一，`protocol buffers`编译器会发出警告。同样，不能使用任何以前[保留](#保留字段)的字段编号。

### 自定字段规则

消息的字段可以是以下之一：

- 单数：格式良好的消息可以包含零个或一个（但不超过一个）这样的字段。这是`proto3`语法的默认字段规则。
- 重复：该字段可以在格式良好的消息中重复任意次数（包括零）。将保留重复值的顺序。

在`proto3`中，标量数字类型的重复字段默认使用**压缩**编码。

在[`Protocol Buffer Encoding`](/Protocol-Buffers/04-编码.md#压缩重复字段)中找到有关压缩编码的更多信息。

### 添加更多消息类型

可以在单个`.proto`文件中定义多种消息类型。如果要定义多个相关消息，这非常有用。例如，如果要定义与`SearchResponse`消息类型对应的回复消息格式，则可以将其添加到相同的`.proto`：

```protobuf
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}

message SearchResponse {
 ...
}
```

### 添加注释

要为`.proto`文件添加注释，请使用`C/C++`样式`//`和`/* ... */`语法。

```protobuf
/* SearchRequest 代表一个查询请求,
 * 带有分页选项以指示要包含在响应中的结果。*/

message SearchRequest {
  string query = 1;
  int32 page_number = 2;  // 我们需要的页码
  int32 result_per_page = 3;  // 每页返回的结果数
}
```

### 保留字段

如果通过完全删除字段或将其注释来[更新](#更新消息类型)消息类型，未来的用户可以在对类型进行更新时再次使用该字段编号。如果以后加载相同`.proto`文件的旧版本，这可能会导致严重问题，包括数据损坏，隐私错误等。确保不会发生这种情况的一种方法是**指定已删除字段或字段的编号为保留的**（并且/或命名可能导致JSON序列化问题）。如果将来的任何用户尝试使用这些字段标识符，`protocol buffers`编译器将会发出警告。

```protobuf
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "foo", "bar";
}
```

请注意，不能在同一保留语句中混合字段名称和字段编号。

### `.proto`文件将生成什么

在`.proto`文件上运行[`protocol buffers`编译器](#生成自定义的类)时，编译器会根据文件中的描述生成所选语言的代码，这些代码是需要使用的消息类型，包括：获取和设置字段值，将消息序列化为输出流，并从输入流中解析您的消息。

- 对于`C++`，编译器会从每个`.proto`生成一个`.h`和`.cc`文件，并为文件中描述的每种消息类型提供一个类。
- 对于`Java`，编译器生成一个`.java`文件，其中包含每种消息类型的类，以及用于创建消息类实例的特殊`Builder`类。
- `Python`有点不同：`Python`编译器生成一个模块，其中包含`.proto`中每种消息类型的静态描述符，然后与元类一起使用，以在运行时创建必要的`Python`数据访问类。
- 对于`Go`，编译器会生成一个`.pb.go`文件，其中包含文件中每种消息类型的类型。
- 对于`Ruby`，编译器生成一个带有包含消息类型的`Ruby`模块的`.rb`文件。
- 对于`Objective-C`，编译器从每个`.proto`生成一个`pbobjc.h`和`pbobjc.m`文件，其中包含文件中描述的每种消息类型的类。
- 对于`C＃`，编译器从每个`.proto`生成一个`.cs`文件，其中包含文件中描述的每种消息类型的类。
- 对于`Dart`，编译器会生成一个`.pb.dart`文件，其中包含文件中每种消息类型的类。

可以按照所选语言的教程（即将推出的proto3版本）了解有关为每种语言使用API​​的更多信息。有关更多API详细信息，请参阅相关[API参考](https://developers.google.com/protocol-buffers/docs/reference/overview)（proto3版本即将推出）。

## 标量值类型

## 默认值

## 枚举

## 使用其他消息类型

## 嵌套类型

## 更新消息类型

## 未知字段

## 生成自定义的类
