---
title: 02-proto3指南.md
date: 2020-04-14T10:09:14.258627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 开发框架
- Protocol-Buffers
summary: 02-proto3指南.md
showInMenu: false

---

# 02-proto3指南

本指南介绍如何使用`protocol buffers`语言构建`protocol buffers`数据，包括：

- `.proto`文件语法
- 如何从`.proto`文件生成数据访问类

它涵盖了`protocol buffers`语言的`proto3`版本：有关较早的`proto2`语法的信息，请参阅[`Proto2`语言指南](https://developers.google.com/protocol-buffers/docs/proto)。

这是一个参考指南，对于使用本文档中描述的许多功能的分步示例，请参阅所选语言的[教程](https://developers.google.com/protocol-buffers/docs/tutorials)（目前仅限`proto`2，更多`proto3`文档即将推出）。

## 定义消息类型

首先看一个非常简单的例子。假设要定义搜索请求消息格式，其中每个搜索请求都有：

- 一个字符串类型的查询
- 所查询的特定页码
- 每页返回的结果数

这是用于定义消息类型的`.proto`文件。

```protobuf
 // 指定正在使用proto3语法，默认使用proto2，必须是文件的第一个非空注释行
syntax = "proto3";

message SearchRequest {     // 消息格式以名称-值对的形式指定三个字段
  string query = 1;         // 每个字段有一个名称和类型
  int32 page_number = 2;
  int32 result_per_page = 3;
}
```

### 指定字段类型

在上面的例子中，所有的字段都是标量类型：两个整型一个字符串类型。同时也可以给字段指定组合类型（包括枚举或其他类型）。

### 分配字段编号

如上所示，消息定义中的每个字段都定义一个**唯一的编号**。这些字段的编号用于在消息的[二进制格式](../Protocol-Buffers/04-编码.md)中标识字段，一旦消息类型被使用就不能再更改。请注意：

- `1到15`范围内的字段编号需要一个字节进行编码，包括字段的编号和字段的类型（可以在[`protocol buffers`编码](../Protocol-Buffers/04-编码.md)中找到更多相关信息）
- `16到2047`范围内的字段编号占用两个字节。 因此，应该为非常频繁出现的消息元素保留数字1到15，请记住为将来可能添加的频繁出现的元素留出一些空间

可以指定的最小字段数为1，最大字段数为536,870,911（2的29次方-1）。

> 不能使用数字19000到19999（`FieldDescriptor::kFirstReservedNumber`到`FieldDescriptor::kLastReservedNumber`），因为它们是为`protocol buffers`实现而保留的。

如果在`.proto`中使用这些保留数字之一，`protocol buffers`编译器会发出警告。同样，不能使用任何以前保留的字段编号。

### 自定字段规则

消息的字段可以是以下之一：

- 单数：格式良好的消息可以包含零个或一个（但不超过一个）这样的字段。这是`proto3`语法的默认字段规则。
- 重复：该字段可以在格式良好的消息中重复任意次数（包括零）。将保留重复值的顺序。

在`proto3`中，标量数字类型的重复字段默认使用**压缩**编码。

在[`Protocol Buffer Encoding`](../Protocol-Buffers/04-编码.md)中找到有关压缩编码的更多信息。

### 添加更多消息类型

可以在单个`.proto`文件中定义多种消息类型。如果要定义多个相关消息，这非常有用。例如，如果要定义与`SearchResponse`消息类型对应的回复消息格式，则可以将其添加到相同的`.proto`文件中：

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

如果通过完全删除字段或将其注释来更新消息类型，未来的用户可以在对类型进行更新时再次使用该字段编号。如果以后加载相同`.proto`文件的旧版本，这可能会导致严重问题，包括数据损坏，隐私错误等。确保不会发生这种情况的一种方法是**指定已删除字段或字段的编号为保留的**（否则可能导致JSON序列化问题）。如果将来的任何用户尝试使用这些字段标识符，`protocol buffers`编译器将会发出警告。

```protobuf
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "foo", "bar";
}
```

请注意，不能在同一保留语句中混合字段名称和字段编号。

### `.proto`文件将生成什么

在`.proto`文件上运行`protocol buffers`编译器时，编译器会根据文件中的描述生成所选语言的代码，这些代码是需要使用的消息类型，包括：获取和设置字段值，将消息序列化为输出流，并从输入流中解析消息。

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

标量消息字段可以具有以下类型之一：该表显示`.proto`文件中指定的类型，以及自动生成的类中的相应类型：

| .proto 类型 | 注释                                                                 | Go类型  |
| ----------- | -------------------------------------------------------------------- | ------- |
| double      |                                                                      | float64 |
| float       |                                                                      | float32 |
| int32       | 使用可变长度编码，编码负数的效率低（如果字段可能有负值，改用sint32） | int32   |
| int64       | 使用可变长度编码，编码负数的效率低（如果字段可能有负值，改用sint64） | int64   |
| uint32      | 使用可变长度编码                                                     | uint32  |
| uint64      | 使用可变长度编码                                                     | uint64  |
| sint32      | 使用可变长度编码，有符号int值（这比常规int32更有效地编码负数）       | int32   |
| sint64      | 使用可变长度编码，有符号int值（这比常规int64更有效地编码负数）       | int64   |
| fixed32     | 总是四个字节，如果值大于2的28次方则比uint32更有效                    | uint32  |
| fixed64     | 总是八个字节，如果值大于2的56次方则比uint32更有效                    | uint64  |
| sfixed32    | 总是四个字节                                                         | int32   |
| sfixed64    | 总是八个字节                                                         | int64   |
| bool        |                                                                      | bool    |
| string      | 字符串必须始终包含UTF-8编码或7位ASCII文本，且不能超过2的32次方       | string  |  |
| bytes       | 可以包含不超过2的32次方的任意字节序列                                | []bytes |

在[`protocol buffers`编码](../Protocol-Buffers/04-编码.md)中可以找到更多关于在序列化消息时这些类型是如何被编码的信息。

## 默认值

在解析消息时，如果编码消息不包含某个特定的单数元素，则解析对象中相应的字段将被设置为该字段的默认值。这些默认值根据类型而不同：

- 对于字符串，默认值为空字符串。
- 对于字节，默认值为空字节。
- 对于布尔型，默认值为false。
- 对于数字类型，默认值为零。
- 对于枚举，默认值是**第一个定义的枚举值**，该值必须为0。
- 对于消息字段，未设置该字段。它的确切值取决于编程语言。有关详细信息，请参阅[生成代码指南](https://developers.google.com/protocol-buffers/docs/reference/overview)。

重复（repeated）字段的默认值为空（通常是相应编程语言的空列表）。

> 请注意，对于**标量消息字段**，一旦解析了消息，就无法确定字段是否显式设置为默认值（例如，布尔值是否设置为`false`）或根本没有设置值，因此，在定义消息类型时要注意。例如，如果不希望默认情况下也发生这种行为，那么当设置为`false`时，没有一个布尔值可以打开某些行为。另外请注意，如果标量消息字段设置为其默认值，则该值不会在传输时序列化。

有关默认值如何在生成的代码中工作的更多详细信息，请参阅所选语言的[生成代码指南](https://developers.google.com/protocol-buffers/docs/reference/overview)。

## 枚举

在定义消息类型时，可能希望其中一个字段只有一个预定义的值列表。例如，假设要为每个`SearchRequest`添加语料库（`corups`）字段，其中语料库可以是`UNIVERSAL`，`WEB`，`IMAGES`，`LOCAL`，`NEWS`，`PRODUCTS`或`VIDEO`。可以非常简单地通过向消息定义添加枚举（`enum`），并为每个可能的值添加常量。

在下面的例子中，添加了一个名为`Corpus`的枚举，其中包含所有可能的值，以及一个类型为Corpus的字段：

```protobuf
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
  enum Corpus {
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 2;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
  }
  Corpus corpus = 4;
}
```

如上所示，`Corpus`枚举的第一个常量映射为零：每个枚举定义**必须**包含一个映射到零的常量作为其第一个元素。这是因为：

- 必须有一个零值，以便可以使用0作为数字默认值。
- 零值必须是第一个元素，以便与`proto2`语义兼容，其中第一个枚举值始终是默认值。

可以通过为不同的枚举常量指定相同的值来定义别名。为此，需要将`allow_alias`选项设置为`true`，否则`proto`编译器将在找到别名时生成错误消息。

```protobuf
enum EnumAllowingAlias {
  option allow_alias = true;    // 开启配置
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 1;
}
enum EnumNotAllowingAlias {
  UNKNOWN = 0;
  STARTED = 1;
  // RUNNING = 1;  // 取消注释此行将导致Google内部的编译错误和外部的警告消息。
}
```

枚举常量必须在32位整数范围内。由于枚举值在传输时使用[`varint`编码](/Protocol-Buffers/04-编码.md)，因此，负值效率低，不建议使用。可以在消息定义中定义枚举，如上例所示，也可以在外部定义枚举，这些枚举可以在`.proto`文件中的任何消息定义中重用。还可以使用语法`MessageType.EnumType`将一个消息中声明的枚举类型用作不同消息中字段的类型。

在使用枚举的`.proto`文件上运行`protocol buffers`编译器时，生成的代码将具有相应的`Java`或`C++`枚举类型，在Python中使用一个特殊的`EnumDescriptor`类，用于在运行时生成类中创建一组带有整数值的符号常量。

在反序列化期间，无法识别的枚举值将保留在消息中，如何表示这种值取决于具体的编程语言。

- 在支持具有超出指定符号范围的值的开放枚举类型的编程语言中（例如`C++`和`Go`），未知的枚举值仅作为其基础整数表示存储。
- 在支持封闭枚举类型（如`Java`）的编程语言中，枚举中的大小写用于表示无法识别的值，并且可以使用特殊访问器访问基础整数。

在任何一种情况下，如果消息被序列化了，那么无法识别的值仍然会和消息一起被序列化。

有关如何在应用程序中使用消息枚举的详细信息，请参阅所选语言的[生成代码指南](https://developers.google.com/protocol-buffers/docs/reference/overview)。

### 保留值

如果通过完全删除枚举条目或将其注释掉来更新枚举类型，那么未来的用户可以在对类型进行更新时重用该数值。如果后面又加载相同`.proto`文件的旧版本，这可能会导致严重问题（包括数据损坏，隐私错误等）。确保不会发生这种情况的一种方法是设置已删除条目的数值（和`/`或名称，也可能导致JSON序列化问题））为`reserved`。如果将来的任何用户尝试使用这些标识符，`protocol buffers`编译器将会发出警告。可以使用`max`关键字指定保留（`reserved`）的数值范围达到最大可能值。

```protobuf
enum Foo {
  reserved 2, 15, 9 to 11, 40 to max;
  reserved "FOO", "BAR";
}
```

请注意，不能在同一保留（`reserved`）语句中混合字段名称和数值。

## 使用其他消息类型

可以使用其他消息类型作为字段类型。例如，假设在每个`SearchResponse`消息中包含`Result`消息，为此，可以在同一`.proto`中定义`Result`消息类型，然后在`SearchResponse`中指定`Result`类型的字段：

```protobuf
message SearchResponse {
  repeated Result results = 1;
}

message Result {
  string url = 1;
  string title = 2;
  repeated string snippets = 3;
}
```

### 导入定义

在上面的示例中，`Result`消息类型在与`SearchResponse`相同的文件中定义，如果要用作字段类型的消息类型在另一个`.proto`文件中定义，可以通过导入来使用其他`.proto`文件中的定义。

要导入另一个`.proto`的定义，请在文件顶部添加一个`import`语句：

```protobuf
import "myproject/other_protos.proto";
```

默认情况下，只能使用直接导入的`.proto`文件中的定义。但是，有时可能需要将`.proto`文件移动到新位置。那么，可以在旧位置放置一个虚拟`.proto`文件，以使用`import public`将所有导入转发到新位置，而不是直接移动`.proto`文件并在一次更改中更新所有导入它的文件。`import public`的依赖可以被任何包含`import public`语句的`proto`文件传递。 例如：

```protobuf
// new.proto
// 所有定义都移动到了这里。

// old.proto
// 这是proto文件，所有的clients都导入了它。
import public "new.proto";
import "other.proto";

// client.proto
import "old.proto";
// 使用old.proto和new.proto中的定义，但不使用other.proto
```

`proto`编译器使用`-I/--proto_path`标志在编译器命令行中指定的一组目录中搜索导入的文件。如果没有给出标志，它将查找调用编译器的目录。通常，应将`--proto_path`标志设置为项目的根目录，并对所有导入使用完全限定名称。

### 使用`proto2`消息类型

可以导入[`proto2`](https://developers.google.com/protocol-buffers/docs/proto)消息类型并在`proto3`消息中使用它们，反之亦然。但是，`proto2`枚举不能直接用于`proto3`语法（如果已导入的`proto2`消息使用它们就没关系）。

## 嵌套类型

可以在其他消息类型中定义和使用消息类型，如下例所示：此处`Result`消息在`SearchResponse`消息中定义：

```protobuf
message SearchResponse {
  message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
  }
  repeated Result results = 1;
}
```

如果要在其父消息类型之外重用此消息类型，请将其称为`Parent.Type`：

```protobuf
message SomeOtherMessage {
  SearchResponse.Result result = 1;
}
```

可以根据需要深入嵌套消息：

```protobuf
message Outer {       // Level 0
  message MiddleAA {  // Level 1
    message Inner {   // Level 2
      int64 ival = 1;
      bool  booly = 2;
    }
  }
  message MiddleBB {  // Level 1
    message Inner {   // Level 2
      int32 ival = 1;
      bool  booly = 2;
    }
  }
}
```

## 更新消息类型

如果现有的消息类型不再满足需求，例如，希望消息格式具有额外的字段，但仍然希望使用旧格式创建的代码。**在不破坏任何现有代码的情况下更新消息类型非常简单**。请记住以下规则：

- 请勿更改任何现有字段的字段编号。
- 如果添加新字段，则使用“旧”消息格式序列化的任何消息仍可由新生成的代码进行解析。应该记住这些元素的默认值，以便新代码可以正确地与旧代码生成的消息进行交互。同样的新代码创建的消息可以由旧代码解析，旧的二进制文件在解析时只是忽略新字段。有关详细信息，请参阅“未知字段”部分。
- 在更新的消息类型中不再使用的字段编号就可以删除。 有时想要重命名该字段，可能添加前缀“`OBSOLETE_`”，或者将字段编号设置为保留（`reserved`），以便`.proto`的未来用户不会意外地重复使用该编号。
- `int32`，`uint32`，`int64`，`uint64`和`bool`都是兼容的，这意味着可以将字段从这些类型之一更改为另一种类型，而不会破坏向前或向后兼容性。如果在传输中解析出一个不适合相应类型的数字，将获得与在`C++`中将该数字转换为该类型相同的效果（例如，如果将64位数字作为int32读取，它将被截断为32位）。
- `sint32`和`sint64`彼此兼容，但与其他整数类型不兼容。
- 只要`byte`是有效的`UTF-8`，`string`和`byte`是兼容的。
- 如果`byte`包含消息的编码版本，则嵌入消息与`byte`兼容。
- `fixed32`与`sfixed32`兼容，`fixed64`与`sfixed64`兼容。
- `enum`在传输格式中与`int32`，`uint32`，`int64`和`uint64`兼容（请注意，如果值不合适，将截断值）。但请注意，在反序列化消息时，客户端代码可能会以不同方式对待它们：例如，无法识别的`proto3`枚举类型将保留在消息中，但在反序列化消息时如何表示它是依赖于编程语言的。`Int`字段总是保留它们的值。
- 将单个值更改为新`oneof`的成员是安全且二进制兼容的。如果确保没有代码一次设置多个字段，那么将多个字段移动到新的`oneof`可能是安全的。将任何字段移动到某个现有的`oneof`中都是不安全的。

## 未知字段

未知字段是格式良好的`protocol buffers`序列化数据，它表示解析器无法识别的字段。例如，当旧二进制文件解析具有新字段的新二进制文件发送的数据时，这些新字段将成为旧二进制文件中的未知字段。

最初，`proto3`消息在解析期间总是丢弃未知字段，但在3.5版本中，重新引入了未知字段的保存以匹配`proto2`行为。在版本3.5及更高版本中，未知字段在解析期间保留并包含在序列化输出中。

## `Any`

`Any`消息类型允许将消息用作嵌入类型，而无需使用这些消息的`.proto`定义。`Any`包含任意序列化消息（如`byte`），并带有一个`URL`作为该消息类型的全局唯一标识符用于表示和解析它。要使用`Any`类型，需要导入`google/protobuf/any.proto`。

```protobuf
import "google/protobuf/any.proto";

message ErrorStatus {
  string message = 1;
  repeated google.protobuf.Any details = 2;
}
```

`type.googleapis.com/packagename.messagename`是给定消息类型的默认类型`URL`。

不同的编程语言实现将支持运行时库来帮助程序以类型安全的方式打包和解压缩`Any`值。例如，在`Java`中，`Any`类型将具有特殊的`pack()`和`unpack()`访问器，而在`C++`中则有`PackFrom()`和`UnpackTo () `方法：

```c++
// 在Any中存储任意消息类型。
NetworkErrorDetails details = ...;
ErrorStatus status;
status.add_details()->PackFrom(details);

// 从Any读取任意消息。
ErrorStatus status = ...;
for (const Any& detail : status.details()) {
  if (detail.Is<NetworkErrorDetails>()) {
    NetworkErrorDetails network_error;
    detail.UnpackTo(&network_error);
    ... processing network_error ...
  }
}
```

**目前，正在开发用于处理`Any`类型的运行时库。**

如果已熟悉[`proto2`语法](https://developers.google.com/protocol-buffers/docs/proto)，则`Any`类型将替换[扩展](https://developers.google.com/protocol-buffers/docs/proto#extensions)。

## `Oneof`

如果有一个包含许多字段的消息，并且最多只能同时设置一个字段，则可以使用`oneof`来强制执行此操作同时还能节省内存。

`oneof`字段与正常的字段一样，只是在同一个`oneof`中的所有字段共享内存，并且最多可以同时设置一个字段。设置`oneof`中的任何一个成员时都会自动清除所有其他成员。可以使用`case()`或`WhichOneof()`方法检查`oneof`中的哪个值（如果有）被设置了，具体使用哪个取决于选择的编程语言。

### 使用`Oneof`

要在`.proto`中定义`oneof`，请使用`oneof`关键字，并在后面跟着`oneof`名称，在本例中为`test_oneof`：

```protobuf
message SampleMessage {
  oneof test_oneof {
    string name = 4;
    SubMessage sub_message = 9;
  }
}
```

然后，将`oneof`字段添加到`oneof`定义中。可以添加任何类型的字段，但不能使用`repeated`字段。

在生成的代码中，`oneof`字段与常规字段具有相同的`getter`和`setter`。还可以获得一种特殊方法来检查`oneof`中设置了哪个值（如果有）。可以在相关[API参考](https://developers.google.com/protocol-buffers/docs/reference/overview)中找到有关所选语言的`oneof API`的更多信息。

### `Oneof`的功能

- 设置`oneof`字段将自动清除`oneof`的所有其他成员。因此，如果设置多个字段，则只有设置的最后一个字段仍然具有值。

    ```c++
    SampleMessage message;
    message.set_name("name");
    CHECK(message.has_name());
    message.mutable_sub_message();   // Will clear name field.
    CHECK(!message.has_name());
    ```

- 如果解析器在传输中遇到同一个`oneof`的多个成员，则在解析的消息中仅使用看到的最后一个成员。
- `oneof`不能是`repeated`。
- 如果将`oneof`字段设置为默认值（例如将`int32` `oneof`字段设置为0），那么该`oneof`字段的`“case”`将会被设置，并且这些字段的值将在传输时序列化。
- 如果使用的是`C++`，请确保代码不会导致内存崩溃。以下示例代码将崩溃，因为通过调用`set_name()`方法删除了`sub_message`。

    ```c++
    SampleMessage message;
    SubMessage* sub_message = message.mutable_sub_message();
    message.set_name("name");      // Will delete sub_message
    sub_message->set_...            // Crashes here
    ```

- 同样在`C++`中，如果`Swap()`两条`oneof`的消息，则每条消息都将以另一条消息的`“case”`作为结束：在下面的示例中，`msg1`将会有`sub_message`同时`msg2`将会有`name`。

    ```c++
    SampleMessage msg1;
    msg1.set_name("name");
    SampleMessage msg2;
    msg2.mutable_sub_message();
    msg1.swap(&msg2);
    CHECK(msg1.has_sub_message());
    CHECK(msg2.has_name());
    ```

### 向后兼容性问题

添加或删除`oneof`字段时要小心。如果在检查`oneof`的值返回`None/NOT_SET`，这可能意味着`oneof`尚未设置或已设置为`oneof`的另一个不同版本。没有办法区分，因为没有办法知道传输中的未知字段是否是`oneof`的成员。

#### 标签重用问题

- **将字段移入或移出`oneof`**：在序列化或解析消息后，可能会丢失一些信息（某些字段将被清除）。但是，可以安全地将单个字段移动到新的`oneof`字段中，并且如果已知只有一个字段被设置，则可以移动多个字段。
- **删除`oneof`字段并将其添加回**：在序列化或解析消息后，这可能会清除当前设置的`oneof`字段。
- **拆分或合并`oneof`**：这与移动常规字段有类似的问题。

## Maps

如果要在数据定义中创建关联映射，`protocol buffers`提供了一种方便的快捷方式语法：

```protobuf
map<key_type, value_type> map_field = N;
```

其中`key_type`可以是任何**整数**或**字符串**类型（除了浮点类型和字节之外的任何标量类型）。请注意，枚举不是有效的`key_type`。 `value_type`可以是`map`之外的任何类型。

如果要创建一个`map`，其中每个`Project`消息与字符串键相关联，则可以像下面这样定义它：

```protobuf
map<string, Project> projects = 3;
```

- `map`的字段不能是`repeated`。
- 传输格式的顺序和`map`值的迭代顺序是未定义的，因此不能依赖`map`中的项目按特定顺序排序。
- 从`.proto`文件中生成文本格式时，`map`按键排序，数字键按数字排序。
- 在传输或合并时进行解析，如果有重复的`map`键，那么就使用最后的那个键。从文本格式解析`map`时，如果存在重复键，则解析可能会失败。
- 如果给`map`提供了键却没有提供值，那么字段序列化的具体行为就取决于具体的编程语言。在`c++`,`Java`，`Python`中使用类型的默认值进行序列化，在其他编程语言中，没有值被序列化。

目前，所有`proto3`支持的编程语言都能生成`map`API，更多关于所选语言的`map`API的参考查看[API参考文档](https://developers.google.com/protocol-buffers/docs/reference/overview)。

### 向后兼容性

在传输时，`map`的语法等效于下面的示例，因此不支持`map`的`protocol buffers`实现仍然可以处理传输的数据：

```protobuf
message MapFieldEntry {
  key_type key = 1;
  value_type value = 2;
}

repeated MapFieldEntry map_field = N;
```

支持`map`的任何`protocol buffers`实现都必须生成和接受上述定义所表示的可接受的数据。

## Packages

可以将可选的`package`说明符添加到`.proto`文件，以防止`protocol buffers`消息类型之间的名称冲突。

```protobuf
package foo.bar;
message Open { ... }
```

然后，可以在定义消息类型的字段时使用`package`说明符：

```protobuf
message Foo {
  ...
  foo.bar.Open open = 1;
  ...
}
```

`package`说明符影响生成代码的方式取决于选择的编程语言：

- 在`C++`中，生成的类包含在`C++`命名空间中。例如，`Open`将位于命名空间`foo::bar`中。
- 在`Java`中，除非在`.proto`文件中明确提供选项`java_package`，否则该包将用作`Java`包。
- 在`Python`中，将忽略`package`指令，因为`Python`模块是根据它们在文件系统中的位置进行组织的。
- 在`Go`中，除非在`.proto`文件中明确提供选项`go_package`，否则该包将用作`Go`包名称。
- 在`Ruby`中，生成的类包含在嵌套的`Ruby`命名空间中，转换为所需的`Ruby`大写形式（首字母大写;如果第一个字符不是字母，则`PB_`前置）。例如，`Open`将位于名称空间`Foo::Bar`中。
- 在`C＃`中，转换为`PascalCase`后，包将用作命名空间，除非在`.proto`文件中明确提供选项`csharp_namespace`。例如，`Open`将位于名称空间`Foo.Bar`中。

### 包和名称解析

`protocol buffers`语言中的类型名称解析与`C++`类似：首先搜索最里面的范围，然后搜索下一个范围，依此类推，每个包被认为是其父包的“内部”。一个`'.'` （例如`.foo.bar.Baz`）意味着从最外层的范围开始。

`protocol buffers`编译器通过解析导入的`.proto`文件来解析所有类型名称。每种编程语言的代码生成器都知道如何引用该语言中的每种类型，即使它具有不同的范围规则。

## 定义服务

如果要在RPC（远程过程调用）系统中使用自定义消息类型，可以在`.proto`文件中定义RPC服务接口`protocol buffers`编译器将以选择的编程语言生成服务接口代码和`stub`。

> 例如，要定义一个RPC服务，该服务获取`SearchRequest`请求并返回`SearchResponse`响应消息，可以在`.proto`文件中定义它，如下所示：

```protobuf
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
```

与`protocol buffers`一起使用的最简单的RPC系统是[gRPC](https://grpc.io/)：一种由Google开发的语言平台中立的开源RPC系统。gRPC特别适用于`protocol buffers`，并允许使用特定的`protocol buffers`编译器插件直接从`.proto`文件生成相关的RPC代码。

如果不想使用gRPC，也可以将`protocol buffers`与自定义的RPC实现一起使用。可以在[`Proto2`语言指南](https://developers.google.com/protocol-buffers/docs/proto#services)中找到更多相关信息。

还有一些正在进行的第三方项目为`Protocol Buffers`开发RPC实现。有关我们了解的项目的链接列表，请参阅[第三方加载项wiki页面](https://github.com/protocolbuffers/protobuf/blob/master/docs/third_party.md)。

## JSON映射

`Proto3`支持`JSON`中的规范编码，使得在系统之间共享数据变得更加容易。在下表中逐个类型地描述编码。

如果`JSON`编码数据中缺少某个值，或者其值为`null`，则在解析为`protocol buffers`时，它将被解释为相应的默认值。如果某个字段在`protocol buffers`中具有默认值，则默认情况下将在`JSON`编码的数据中省略该字段以节省空间。一种实现方式是提供可选选项在`JSON`编码的输出中输出字段和它的默认值。

| proto3                 | JSON          | JSON example                              | Notes                                                                                                                                                                                                                                                                              |
| ---------------------- | ------------- | ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| message                | object        | `{"fooBar": v, "g": null, …}`             | 生成JSON对象。`message`字段名称映射到小驼峰命名并成为JSON对象的`key`。如果指定了`json_name`字段这个选项，则将指定的值作为`key`。解析器接受小驼峰命名的名称（或`json_name`选项指定的名称）和原始的`proto`字段名称。`null`是所有字段类型都可接受的值，并被视为相应字段类型的默认值。 |
| enum                   | string        | `"FOO_BAR"`                               | 使用`proto`中指定的枚举值的名称。解析器接受枚举名称和整数值。                                                                                                                                                                                                                      |
| map<K,V>               | object        | `{"k": v, …}`                             | 所有键都转换为字符串。                                                                                                                                                                                                                                                             |
| repeated V             | array         | `[v, …]`                                  | `null`被接受为空的`list[]`。                                                                                                                                                                                                                                                       |
| bool                   | true, false   | `true, false`                             |
| string                 | string        | `"Hello World!"`                          |
| bytes                  | base64 string | `"YWJjMTIzIT8kKiYoKSctPUB+"`              | JSON值将是使用带填充的标准base64编码方式编码的字符串数据。带有/不带填充的标准或URL安全的base64编码方式也是可接受的。                                                                                                                                                               |
| int32, fixed32, uint32 | number        | `1, -10, 0`                               | 十进制数形式的JSON值，接受数字或字符串。                                                                                                                                                                                                                                           |
| int64, fixed64, uint64 | string        | `"1", "-10"`                              | 十进制字符串形式的JSON值，接受数字或字符串。                                                                                                                                                                                                                                       |
| float, double          | number        | `1.1, -10.0, 0, "NaN", "Infinity"`        | 一个或多个特殊字符串`“NaN”`，`“Infinity”`和`“-Infinity”`形式的JSON值，接受数字或字符串，指数表示法也被接受。                                                                                                                                                                       |
| Any                    | object        | `{"@type": "url", "f": v, … }`            | 如果`Any`包含具有特殊JSON映射的值，则它将按如下方式转换：`{“@ type”：xxx，“value”：yyy}`。 否则，该值将转换为JSON对象，并将插入`“@type”`字段以指示实际数据类型。                                                                                                                   |
| Timestamp              | string        | `"1972-01-01T10:00:20.021Z"`              | 使用`RFC 3339`，其中生成的输出将始终被**Z-标准化**并使用0,3,6或9个小数位。也接受“Z”以外的偏移。                                                                                                                                                                                    |
| Duration               | string        | `"1.000340012s", "1s"`                    | 生成的输出始终包含0,3,6或9个小数位，具体取决于所需的精度，后跟后缀`“s”`。接受任何小数位（也可以没有小数位），只要它们符合纳秒精度并且需要后缀`“s”`。                                                                                                                               |
| Struct                 | object        | `{ … }`                                   | 任意JSON对象，查看`struct.proto`文件                                                                                                                                                                                                                                               |
| Wrapper types          | various types | `2, "2", "foo", true, "true", null, 0, …` | `Wrappers`在JSON中使用与包装基元类型相同的表示形式，除了在数据转换和传输期间允许和保留`null`。                                                                                                                                                                                     |
| FieldMask              | string        | `"f.fooBar,h"`                            | 查看`field_mask.proto`文件                                                                                                                                                                                                                                                         |
| ListValue              | array         | `[foo, bar, …]`                           |
| Value                  | value         |                                           | 任意JSON值                                                                                                                                                                                                                                                                         |
| NullValue              | null          |                                           | JSON中的`null`                                                                                                                                                                                                                                                                     |
| Empty                  | object        | `{}`                                      | 空的JSON对象                                                                                                                                                                                                                                                                       |

### JSON选项

`proto3 JSON`实现可以提供以下可用选项：

- **输出字段的默认值**：在`proto3 JSON`的输出中默认省略字段的默认值。有一个选项可以提供覆盖此行为并输出字段的默认值。
- **忽略未知字段**：默认情况下，`proto3 JSON`解析器会拒绝未知字段，可以提供一个选项来忽略解析未知字段。
- **使用`proto`字段名称而不是小驼峰命名的名称**：默认情况下，`proto3 JSON`会将字段名称转换为小驼峰命名并将其用作JSON的名称。有一个选项可以提供使用`proto`字段名称作为JSON的名称。`proto3 JSON`解析器需要接受转换后的小驼峰命名的名称和`proto`字段名称。
- **将枚举值作为整数而不是字符串输出**：默认情况下，在JSON输出中使用枚举值的名称。可以提供一个选项以使用枚举值的数值。

## 可用选项

`.proto`文件中的各个声明可以使用许多选项进行注释。选项不会更改声明的整体含义，但可能会影响该声明在特定上下文中被处理的方式。可用选项的完整列表在`google/protobuf/descriptor.proto`中定义。

- 一些选项是**文件级选项**，这意味着它们应该在顶级范围内编写，而不是在任何`message`、`enum`或`service`的定义中。
- 一些选项是**消息级选项**，这意味着它们应该写在`message`定义中。
- 一些选项是**字段级选项**，这意味着它们应该写在字段定义中。

> 可用选项也可以写在枚举类型，枚举值，服务类型和服务方法上，但是，目前没有任何支持这些的可用选项。

以下是一些最常用的选项：

- `java_package`（文件级选项）：用于生成的Java类的包。如果`.proto`文件中没有给出显式的`java_package`选项，那么默认情况下将使用`proto`包（.`proto`文件中的`“package”`关键字指定的包）。但是，`proto`包通常不能生成好的`Java`包，因为`proto`包不会以反向域名开头。如果不生成Java代码，则此选项无效。

    ```protobuf
    option java_package = "com.example.foo";
    ```

- `java_multiple_files`（文件级选项）：生成在包级别中定义的顶级`message`、`enum`和`service`，而不是在`.proto`文件之后命名的外部类中。

    ```protobuf
    option java_multiple_files = true;
    ```

- `java_outer_classname`（文件级选项）：生成最外层的Java类的类名（以及文件名）。如果`.proto`文件中没有指定显式的`java_outer_classname`，则通过将`.proto`文件名转换为驼峰命名来构造类名（因此`foo_bar.proto`变为`FooBar.java`）。如果不生成Java代码，则此选项无效。

    ```protobuf
    option java_outer_classname = "Ponycopter";
    ```

- optimize_for（文件级选项）：可以设置为`SPEED`，`CODE_SIZE`或`LITE_RUNTIME`。这会影响`C++`和`Java`代码生成器（可能还有第三方生成器），以下列方式：

    ```protobuf
    option optimize_for = CODE_SIZE;
    ```

    - `SPEED`（default）：`protocol buffers`编译器将生成用于对消息类型进行序列化、解析和执行其他常见操作的代码。此代码经过高度优化。
    - `CODE_SIZE`：`protocol buffers`编译器将生成最小的类，并依赖于基于反射的共享代码来实现序列化、解析和各种其他操作。因此生成的代码将比使用`SPEED`小得多，但操作会更慢。生成的类仍将实现与`SPEED`模式完全相同的公共API。此模式在包含大量`.proto`文件的应用程序中最有用，并且不需要所有这些文件都非常快速。
    - `LITE_RUNTIME`：`protocol buffers`编译器将生成仅依赖于`“lite”`的运行时库的类（即依赖于`libprotobuf-lite`而不是`libprotobuf`）。`lite`运行时库比完整库小得多（大约小一个数量级），其中省略了描述符和反射等功能。这对于在移动电话等受限平台上运行的应用程序尤其有用。编译器仍将生成所有方法的快速实现，就像在`SPEED`模式下那样。生成的类将仅实现每种语言的`MessageLite`接口，该接口仅提供完整`Message`接口的方法的子集。
`
- `cc_enable_arenas`（文件级选项）：为`C++`生成的代码启用[竞技场分配](https://developers.google.com/protocol-buffers/docs/reference/arenas)。
- `objc_class_prefix`（文件级选项）：设置`Objective-C`类前缀，该前缀由此`.proto`文件提供给所有`Objective-C`生成的类和枚举，没有默认值。应该使用[Apple建议](https://developer.apple.com/library/ios/documentation/Cocoa/Conceptual/ProgrammingWithObjectiveC/Conventions/Conventions.html#//apple_ref/doc/uid/TP40011210-CH10-SW4)的3-5个大写字符之间的前缀。请注意，Apple保留所有2个字母的前缀。
- `deprecated`（文件级选项）：如果设置为`true`，则表示该字段已弃用，新代码不应使用该字段。在大多数语言中，这没有实际效果。在`Java`中，这将成为`@Deprecated`注释。将来，其他特定语言的代码生成器可能会在字段的访问器上生成弃用注释，这将导致在编译尝试使用该字段的代码时发出警告。如果任何人都没有使用该字段，并且想要阻止新用户使用该字段，请考虑使用保留语句替换字段声明。

    ```protobuf
    int32 old_field = 6 [deprecated=true];
    ```

### 自定义选项

`Protocol Buffers`还允许定义和使用自定义的选项。这是大多数人不需要的**高级功能**。如果确实认为需要创建自定义的选项，请参阅[`proto2`语言指南](https://developers.google.com/protocol-buffers/docs/proto.html#customoptions)以获取详细信息。请注意，创建自定义的选项使用的[扩展](https://developers.google.com/protocol-buffers/docs/proto.html#extensions)仅允许用于`proto3`中的自定义的选项。

## 生成自定义的类

需要使`.proto`文件中定义的消息类型来生成`Java`，`Python`，`C++`，`Go`，`Ruby`，`Objective-C`或`C＃`代码，这需要在`.proto`上运行`protocol buffers`编译器`protoc`。如果尚未安装编译器，请下载该软件包并按照自述文件中的说明进行操作。对于`Go`，还需要为编译器安装一个特殊的代码生成器插件，可以在`GitHub`上的`golang/protobuf`存储库中找到这个插件和安装说明。

协议编译器的调用如下：

```bash
protoc --proto_path=IMPORT_PATH \
       --cpp_out=DST_DIR \
       --java_out=DST_DIR \
       --python_out=DST_DIR \
       --go_out=DST_DIR \
       --ruby_out=DST_DIR \
       --objc_out=DST_DIR \
       --csharp_out=DST_DIR \
       path/to/file.proto
```

- `IMPORT_PATH`指定解析导入指令时查找`.proto`文件的目录，如果省略，则使用当前目录。可以通过多次传递`--proto_path`选项来指定多个导入目录,将按顺序搜索。`-I=IMPORT_PATH`可以用作`--proto_path`的缩写形式。
- 可以提供一个或多个输出指令：
    - `--cpp_out`在`DST_DIR`中生成`C++`代码。有关更多信息，请参阅[`C++`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/cpp-generated)。
    - `--java_out`在`DST_DIR`中生成`Java`代码。有关更多信息，请参阅[`Java`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/java-generated)。
    - `--python_out`在`DST_DIR`中生成`Python`代码。有关更多信息，请参阅[`Python`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/python-generated)。
    - `--go_out`在`DST_DIR`中生成`Go`代码。有关更多信息，请参阅[`Go`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/go-generated)。
    - `--ruby_out`在`DST_DIR`中生成`Ruby`代码。`Ruby`生成的代码参考即将推出！
    - `--objc_out`在`DST_DIR`中生成`Objective-C`代码。有关更多信息，请参阅[`Objective-C`生成的代码参考](https://developers.google.com/protocol-buffers/docs/reference/objective-c-generated)。
    - `--csharp_out`在`DST_DIR`中生成`C＃`代码。有关更多信息，请参阅[`C＃`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/csharp-generated)。
    - `--php_out`在`DST_DIR`中生成`PHP`代码。有关更多信息，请参阅[`PHP`生成代码参考](https://developers.google.com/protocol-buffers/docs/reference/php-generated)。

    为方便起见，如果`DST_DIR`以`.zip`或.`jar`结尾，编译器会将输出写入具有给定名称的单个`ZIP`格式存档文件。`.jar`输出还将根据`Java JAR`规范的要求提供清单文件。请注意，如果输出存档已存在，则会被覆盖; 编译器不够智能，无法将文件添加到现有存档中。

- 必须提供一个或多个`.proto`文件作为输入。可以一次指定多个`.proto`文件。虽然文件是相对于当前目录命名的，但每个文件必须驻留在其中一个`IMPORT_PATH`中，以便编译器可以确定其规范名称。
