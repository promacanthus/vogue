---
title: "10 Goquery踩坑"
date: 2020-06-16T10:41:51+08:00
draft: true
---

```go
// package
import "github.com/PuerkitoBio/goquery"

// go mod

require github.com/PuerkitoBio/goquery v1.5.1
```

API文档在`pkg.go.dev`点[这里](https://pkg.go.dev/github.com/PuerkitoBio/goquery@v1.5.1?tab=doc)，Github仓库在[这里](https://github.com/PuerkitoBio/goquery)。

goquery在Github上近9k的star，golang著名的爬虫框架[colly](http://go-colly.org/)也用的它。

## goquery简介

`goquery`实现了与`jQuery`类似的功能，包括使用可链接语法操纵和查询HTML文档。

它为Go语言带来了一种类似于`jQuery`的语法和函数。它基于Go的`net/html`包和CSS Selector库[`cascadia`](https://github.com/andybalholm/cascadia)。由于`net/html`解析器返回的是节点，而不是功能完整的`DOM`树，因此`jQuery`的有状态操作函数（例如`height()`，`css()`，`detach()`）已被省去。

> `net/html`解析器读取`UTF-8`编码的文件（Go默认处理的就是UTF-8编码的文件），所以要确保被操作的源文档是`UTF-8`编码的`HTML`文件。有关如何执行此操作的各种选项，可查看Github仓库的[Wiki](https://github.com/PuerkitoBio/goquery/wiki/Tips-and-tricks)。

在语法上，它尽可能接近`jQuery`，并在可能的情况下使用相同的方法名称，并且具有熟悉而类似的可链接接口。

> `jQuery`是一个广受欢迎的库，因此，参照和遵循它的API来编写类似的`HTML`操作库会更好，这也是Go语言的精神（如`fmt`包的实现），即使`jQuery`有些方法看起来不那么直观（如`index()`等）。

**注意：`goquery`以来`net/html`库，因此需要`Go1.1+`以上版本**。

根据方法的不同类型分类到不同的文件中，三个点（`...`）表示该方法可以重载（overloads）。

- `array.go`（数组类位置选择操作）：`Eq()`，`First()`，`Get()`，`Index...()`，`Last()`，`Slice()`
- `expend.go`（扩充选择的集合）：`Add...()`，`AndSelf()`，`Union()`是`AddSelection()`的别名
- `filter.go`（过滤选择的集合）：`End()`，`Filter...()`，`Has...()`，`Intersection()`是`FilterSelection()`的别名，`Not...()`
- `iteration.go`（遍历节点）：`Each()`，`EachWithBreak()`，`Map()`
- `manipulation.go`（修改HTML）：`After...()`，`Append...()`，`Before...()`，`Clone()`，`Empty()`，`Prepend...()`，`Remove...()`，`ReplaceWith...()`，`Unwrap()`，`Wrap...()`，`WrapAll...()`，`WrapInner...()`
- `property.go`（检查并获取节点属性值）：`Attr*()`, `RemoveAttr()`, `SetAttr()`，`AddClass(), HasClass(), RemoveClass(), ToggleClass()`，`Html()`，`Length()`，`Size()`是`Length()`的别名，`Text()`
- `query.go`（判断节点身份）：`Contains()`，`Is...()`
- `traversal.go`（遍历`HTML`文档树）：，`Children...()`，`Contents()`，`Find...()`，`Next...()`，`Parent[s]...()`，`Prev...()`，`Siblings...()`
- `type.go`(`goquery`公开的类型)：`Document`，`Selection`，`Matcher`
- `utilities.go`（辅助函数，而不是`* Selection`的方法，这不是`jQuery`的一部分）：`NodeName`，`OuterHtml`

了解了具体的功能和提供的函数，具体就是在调用上面函数的时候提供CSS选择器作为参数。

## CSS选择器

分类|类别|描述|语法|例子
---|---|---|---|---
基本选择器|通用选择器|选择所有元素|`*` `ns|*` `*|*`|`*`匹配文档是所有元素
基本选择器|类型选择器|按照给定的节点名称，选择所有匹配的元素|`elementname`|input 匹配任何 `<input>` 元素
基本选择器|类选择器|按照给定的 class 属性的值，选择所有匹配的元素|`.classname`|`.index` 匹配任何 `class` 属性中含有 "index" 的元素
基本选择器|ID选择器|按照 id 属性选择一个与之匹配的元素。需要注意的是，一个文档中，每个 ID 属性都应当是唯一的|`#idname`|`#toc` 匹配 ID 为 "toc" 的元素
基本选择器|属性选择器|按照给定的属性，选择所有匹配的元素|`[attr]` `[attr=value]` `[attr~=value]` `[attr|=value]` `[attr^=value]` `[attr$=value]` `[attr*=value]`|`[autoplay]` 选择所有具有 autoplay 属性的元素
分组选择器|选择器列表|`,`将不同的选择器组合在一起，它选择所有能被列表中的任意一个选择器选中的节点|`A,B`|`div, span` 会同时匹配 `<span>` 元素和 `<div>` 元素
组合选择器|后代组合器|` `空格组合器选择前一个元素的后代节点|`A B`|`div span` 匹配所有位于任意 `<div>` 元素之内的 `<span>` 元素
组合选择器|直接子代组合器|`>`组合器选择前一个元素的直接子代的节点|`A > B`|`ul > li` 匹配直接嵌套在 `<ul>` 元素内的所有 `<li>` 元素
组合选择器|一般兄弟组合器|`~`组合器选择兄弟元素，即后一个节点在前一个节点后面的任意位置，并且共享同一个父节点|`A ~ B`|`p ~ span` 匹配同一父元素下，`<p>` 元素后的所有 `<span>` 元素
组合选择器|紧邻兄弟组合器|`+`组合器选择相邻元素，即后一个元素紧跟在前一个之后，并且共享同一个父节点|`A + B`|`h2 + p` 会匹配所有紧邻在 `<h2>` 元素后的 `<p>` 元素
组合选择器|列组合器|`||`组合器选择属于某个表格行的节点|`A || B`|`col || td` 会匹配所有 `<col>` 作用域内的 `<td>` 元素
伪选择器|伪类|`:`伪选择器支持按照未被包含在文档树中的状态信息来选择元素||`a:visited` 匹配所有曾被访问过的 `<a>` 元素
伪选择器|伪元素|`::`伪选择器用于表示无法用`HTML`语义表达的实体||`p::first-line` 匹配所有 `<p>` 元素的第一行

## goquery使用样例

大部分都和上面的一样，比较特殊的列在下面。

选择器|说明
---|---
Find(“div[lang]")|筛选含有lang属性的div元素
Find(“div[lang=zh]")|筛选lang属性为zh的div元素
Find(“div[lang!=zh]")|筛选lang属性不等于zh的div元素
Find(“div[lang¦=zh]")|筛选lang属性为zh或者zh-开头的div元素
Find(“div[lang*=zh]")|筛选lang属性包含zh这个字符串的div元素
Find(“div[lang~=zh]")|筛选lang属性包含zh这个单词的div元素，单词以空格分开的
Find(“div[lang$=zh]")|筛选lang属性以zh结尾的div元素，区分大小写
Find(“div[lang^=zh]")|筛选lang属性以zh开头的div元素，区分大小写

下面的操作在选择器选出的内容中再进行过滤。

类别|描述|语法|例子
---|---|---|---
内容过滤器|筛选出的元素要包含指定的文本|`Find(":contains(text)")`|<ul><li>`Find("div:contains(DIV2)")`，选择出的`div`元素要包含`DIV2`文本<li>`Find(":empty)`，选出的元素都不能有子元素<li>`Find("span:has(div)")`，选出包含`div`的`span`元素，与`has`类似的`contains`</ul>
`:first-child`/`:last-child`|选出其父元素的第一个子元素|`Find(":first-child")`|`Find("div:first-child")`
`:first-of-type`/`:last-of-type`|选出其父元素的第一个该类型子元素|`Find(":first-of-type")`|`Find("div:first-of-type")`
`:nth-child(n)`/`:nth-last-child(n)`|选出其父元素的第n个子元素|`Find(":nth-child(n)")`|`Find("div:nth-child(3)")`
`:nth-of-type(n)`/`:nth-last-of-type(n)`|选出其父元素的第n个该类型子元素|`Find(":nth-of-type(n)")`|`Find("div:nth-of-type(3)")`
`:only-child`|选出其父元素中只有该元素（数量唯一）的子元素|`Find(":only-child")`|`Find("div:only-child")`
`:only-of-type`|选出其父元素中只有该类型元素（类型唯一）的子元素|`Find(":only-of-type")`|`Find("div:only-of-type")`
