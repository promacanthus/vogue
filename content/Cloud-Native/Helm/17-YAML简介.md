# 17-YAML简介

YAML具有一些有用的特性，可以使我们的模板更少出错并更易于阅读。

## 标量和集合

根据YAML规范YAML spec，有两种类型的集合，以及许多标量类型。

这两种类型的集合是字典和数组：

```yaml
# 字典
map:
  one: 1
  two: 2
  three: 3

# 数组
sequence:
  - one
  - two
  - three
```

标量值是单个值（与集合相对）。

## YAML中的标量类型

在Helm的YAML语言中，值的标量数据类型由一组复杂的规则确定，包括用于资源定义的Kubernetes schema。在推断类型时，以下规则成立。

```yaml
count: 1      # int
size: 2.34    # float

count: "1" # string, not int
size: '2.34' # string, not float

isGood: true   # bool
answer: "true" # string

# 空值的词是null（not nil）
```

请注意：

- `port: "80"`是有效的，并且能通过模板引擎和YAML分析器
- 但如果Kubernetes预期port为整数，则会失败

在某些情况下，可以使用YAML节点标签强制进行特定的类型推断:

```yaml
coffee: "yes, please"
age: !!str 21       # !!str告诉解析器age是一个字符串
port: !!int "80"    # !!int告诉解析器port是一个整数
```

## YAML中的字符串

在YAML文档中的大部分数据都是字符串，YAML有多种表示字符串的方式。有三种内置方式声明字符串：

1. `way1: bare words`：单词未被引号引起来
2. `way2: "double-quoted strings"`：双引号的字符串可以使用特定的字符\进行转义
3. `way3: 'single-quoted strings'`：单引号字符串是“文字”字符串，并且不使用\转义字符。唯一的转义序列是`''`，它被解码为一个单独的`'`

声明**多行**字符串，所有内置样式必须位于同一行上：

```yaml
coffee: |
  Latte
  Cappuccino
  Espresso

# coffee的值等价于：Latte\nCappuccino\nEspresso\n
```

请注意`|`后的第一行必须正确缩进。

### 控制多行字符串中的空格

用来`|`表示一个多行字符串。但请注意字符串的内容后跟着`\n`。

**`|-`，去掉行尾`\n`**

```yaml
coffee: |-
  Latte
  Cappuccino
  Espresso

# coffee值等价于：Latte\nCappuccino\nEspresso
```

**`|+`，保留所以`\n`**

```yaml
coffee: |+
  Latte
  Cappuccino
  Espresso


another: value
# coffee值等价于：Latte\nCappuccino\nEspresso\n\n\n
```

文本块内部的缩进被保留，并保留换行符：

```yaml
coffee: |-
  Latte
    12 oz
    16 oz
  Cappuccino
  Espresso
# coffee将是Latte\n 12 oz\n 16 oz\nCappuccino\nEspresso
```

## 缩进和模板

在编写模板时，希望将文件内容注入模板。有两种方法可以做到这一点：

1. 使用`{{ .Files.Get "FILENAME" }}`得到chart中的文件的内容
2. 使用`{{ include "TEMPLATE" . }}`渲染模板，然后其内容放入chart

将文件插入YAML时，最好理解上面的多行规则。通常情况下，插入静态文件的最简单方法是做这样的事情：

```yaml
myfile: |
{{ .Files.Get "myfile.txt" | indent 2 }}
```

使用indent 2告诉模板引擎使用两个空格缩进“myfile.txt”中的**每一行**。

请注意，不缩进该模板行。那是因为如果缩进了，第一行的文件内容会缩进两次。

## 折叠多行字符串

有时候想在YAML中用多行代表一个字符串，但是当它被解释时，要把它当作一个长行。这被称为“折叠”。

要声明一个折叠块，使用`>`代替`|`，除最后一个换行符之外的所有内容都将转换为空格。

```yaml
coffee: >
  Latte
  Cappuccino
  Espresso

# coffee的值等价于：Latte Cappuccino Espresso\n

# 请注意，在折叠语法中，缩进文本将导致行被保留

coffee: >-
  Latte
    12 oz
    16 oz
  Cappuccino
  Espresso
# coffee的值等价于：Latte\n 12 oz\n 16 oz\nCappuccino Espresso

# 注意区别，此时除了最后一个\n，其他换行符和空格都保留
```

## 将多个文档嵌入到一个文件中

可以将多个YAML文档放入单个文件中。这是通过在一个新文档前加---,在文档结束加...来完成的

```yaml
---
document:1
...
---
document: 2
...

# 在许多情况下，无论是---或...可被省略
```

> Helm中的某些文件不能包含多个文档。例如，如果文件内部提供了多个values.yaml文档，则只会使用第一个文档。但是，模板文件可能有多个文档。发生这种情况时，文件（及其所有文档）在模板渲染期间被视为一个对象。但是，最终的YAML在被送到Kubernetes之前被分成多个文件。

**建议每个文件在绝对必要时才使用多个文档。在一个文件中有多个文件可能很难调试。**

## YAML是JSON的Superset

因为YAML是JSON的超集，所以任何有效的JSON文档都应该是有效的YAML。

```yaml
{
  "coffee": "yes, please",
  "coffees": [
    "Latte", "Cappuccino", "Espresso"
  ]
}

# 以上是下面另一种表达方式：
coffee: yes, please
coffees:
- Latte
- Cappuccino
- Espresso

# 这两者可以混合使用（小心使用）：

coffee: "yes, please"
coffees: [ "Latte", "Cappuccino", "Espresso"]
# 所有这三个都应该解析为相同的内部表示
```

虽然这意味着诸如values.yaml可能包含JSON数据的文件，但Helm不会将文件扩展名.json视为有效的后缀。

## YAML锚

YAML规范提供了一种方法来存储对某个值的引用，并稍后通过引用来引用该值。YAML将此称为“锚定”：

```yaml
coffee: "yes, please"
favorite: &favoriteCoffee "Cappucino"
coffees:
  - Latte
  - *favoriteCoffee
  - Espresso

```

1. 在上面，&favoriteCoffee设置一个引用到Cappuccino。
2. 之后，该引用被用作*favoriteCoffee。
3. 所以coffees变成了 Latte, Cappuccino, Espresso。

> 虽然在少数情况下锚点是有用的，但它们的一个方面可能导致细微的错误：第一次使用YAML时，引用被扩展，然后被丢弃。

所以如果我们要解码然后重新编码上面的例子，那么产生的YAML将是：

```yaml
coffee: yes, please
favorite: Cappucino
coffees:
- Latte
- Cappucino
- Espresso

```

因为Helm和Kubernetes经常读取，修改并重写YAML文件，锚将会丢失。
