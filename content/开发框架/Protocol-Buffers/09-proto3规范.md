---
title: 09-proto3规范
date: 2020-04-14T10:09:14.258627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- Protocol-Buffers
summary: 09-proto3规范
showInMenu: false

---

这是`Protocol Buffers`语言（`proto3`）第3版的语言规范参考。使用Extended Backus-Naur Form（EBNF）指定语法：

```bash
|   alternation
()  grouping
[]  option (zero or one time)
{}  repetition (any number of times)
```

有关使用proto3的更多信息，请参阅[语言指南](../Protocol-Buffers/02-proto3指南/)。

## 词汇元素

### 字母和数字

```bash
letter = "A" … "Z" | "a" … "z"
decimalDigit = "0" … "9"
octalDigit   = "0" … "7"
hexDigit     = "0" … "9" | "A" … "F" | "a" … "f"
```

### 身份标识

```bash
ident = letter { letter | decimalDigit | "_" }
fullIdent = ident { "." ident }
messageName = ident
enumName = ident
fieldName = ident
oneofName = ident
mapName = ident
serviceName = ident
rpcName = ident
messageType = [ "." ] { ident "." } messageName
enumType = [ "." ] { ident "." } enumName
```

### 整数

```bash
intLit     = decimalLit | octalLit | hexLit
decimalLit = ( "1" … "9" ) { decimalDigit }
octalLit   = "0" { octalDigit }
hexLit     = "0" ( "x" | "X" ) hexDigit { hexDigit }
```

### 浮点数

```bash
floatLit = ( decimals "." [ decimals ] [ exponent ] | decimals exponent | "."decimals [ exponent ] ) | "inf" | "nan"
decimals  = decimalDigit { decimalDigit }
exponent  = ( "e" | "E" ) [ "+" | "-" ] decimals
```

### 布尔

```bash
boolLit = "true" | "false"
```

### 字符串

```bash
strLit = ( "'" { charValue } "'" ) |  ( '"' { charValue } '"' )
charValue = hexEscape | octEscape | charEscape | /[^\0\n\\]/
hexEscape = '\' ( "x" | "X" ) hexDigit hexDigit
octEscape = '\' octalDigit octalDigit octalDigit
charEscape = '\' ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | '\' | "'" | '"' )
quote = "'" | '"'
```

### 空声明

```bash
emptyStatement = ";"
```

### 常量

```bash
constant = fullIdent | ( [ "-" | "+" ] intLit ) | ( [ "-" | "+" ] floatLit ) | strLit | boolLit
```

## 句法

语法(`syntax`)语句用于定义`protobuf`版本。

```bash
syntax = "syntax" "=" quote "proto3" quote ";"

# 例如

syntax = "proto3";
```

## 导入语句

`import`语句用于导入另一个`.proto`的定义。

```bash
import = "import" [ "weak" | "public" ] strLit ";"

# 例如

import public "other.proto";
```

## 包

包说明符可用于防止协议消息类型之间的名称冲突。

```bash
package = "package" fullIdent ";"

# 例如

package foo.bar;
```

## 可用选项

选项可用于`proto`文件，消息，枚举和服务。可用选项可以是`protobuf`定义的选项或自定义选项。有关更多信息，请参阅[语言指南中的选项](../Protocol-Buffers/02-proto3指南/)。

```protobuf
option = "option" optionName  "=" constant ";"
optionName = ( ident | "(" fullIdent ")" ) { "." ident }

// 例如
option java_package = "com.example.foo";
```

## 字段

字段是`protocol buffers`消息的基本元素。字段可以是普通字段，`oneof`字段或`map`字段。字段具有类型和字段编号。

```protobuf
type = "double" | "float" | "int32" | "int64" | "uint32" | "uint64"
      | "sint32" | "sint64" | "fixed32" | "fixed64" | "sfixed32" | "sfixed64"
      | "bool" | "string" | "bytes" | messageType | enumType
fieldNumber = intLit;
```

### 普通字段

每个字段都有类型，名称和字段编号。它可能有字段选项。

```protobuf
field = [ "repeated" ] type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
fieldOptions = fieldOption { ","  fieldOption }
fieldOption = optionName "=" constant

// 例如

foo.bar nested_message = 2;
repeated int32 samples = 4 [packed=true];
```

### `Oneof`集合`oneof`字段

`oneof`由`oneof`字段和`oneof`名称组成。

```protobuf
oneof = "oneof" oneofName "{" { oneofField | emptyStatement } "}"
oneofField = type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"

// 例如

oneof foo {
    string name = 4;
    SubMessage sub_message = 9;
}
```

### `Map`字段

`map`字段具有键类型，值类型，名称和字段编号。键类型可以是任何整数或字符串类型。

```protobuf
mapField = "map" "<" keyType "," type ">" mapName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
keyType = "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" |
          "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string"

// 例如

map<string, Project> projects = 3;
```

## 保留的(Reserved)

保留语句声明了一系列不能在此消息中使用的字段编号或字段名称。

```protobuf
reserved = "reserved" ( ranges | fieldNames ) ";"
ranges = range { "," range }
range =  intLit [ "to" ( intLit | "max" ) ]
fieldNames = fieldName { "," fieldName }

// 例如

reserved 2, 15, 9 to 11;
reserved "foo", "bar";
```

## 顶级定义

### 枚举定义

枚举定义由名称和枚举主体组成。枚举主体可以有选项和枚举字段。枚举定义必须以枚举值零开始。

```protobuf
enum = "enum" enumName enumBody
enumBody = "{" { option | enumField | emptyStatement } "}"
enumField = ident "=" intLit [ "[" enumValueOption { ","  enumValueOption } "]" ]";"
enumValueOption = optionName "=" constant

// 例如

enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
```

### 消息定义

消息由消息名称和消息正文组成。消息正文可以包含字段，嵌套枚举定义，嵌套消息定义，可用选项，`oneof`字段，`map`字段和保留语句。

```protobuf
message = "message" messageName messageBody
messageBody = "{" { field | enum | message | option | oneof | mapField |
reserved | emptyStatement } "}"

// 例如

message Outer {
  option (my_option).a = true;
  message Inner {   // Level 2
    int64 ival = 1;
  }
  map<int32, string> my_map = 2;
}
```

### 服务定义

```protobuf
service = "service" serviceName "{" { option | rpc | emptyStatement } "}"
rpc = "rpc" rpcName "(" [ "stream" ] messageType ")" "returns" "(" [ "stream" ]
messageType ")" (( "{" {option | emptyStatement } "}" ) | ";")

// 例如

service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
```

## `Proto`文件

```proto
proto = syntax { import | package | option | topLevelDef | emptyStatement }
topLevelDef = message | enum | service
```

示例`.proto`文件

```protobuf
syntax = "proto3";
import public "other.proto";
option java_package = "com.example.foo";
enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
message outer {
  option (my_option).a = true;
  message inner {   // Level 2
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
```
