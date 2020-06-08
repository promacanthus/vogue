---
title: "00 Commit Messages"
date: 2020-05-11T13:49:42+08:00
draft: true
---

原仓库[地址](https://github.com/RomuloOliveira/commit-messages-guide/blob/master/README_zh-CN.md)。

## 使用祈使句

```bash
# good
Use InventoryBackendPool to retrieve inventory backend  // 用 InventoryBackendPool 获取库存

# bad
Used InventoryBackendPool to retrieve inventory backend // InventoryBackendPool 被用于获取库存
```

> commit 信息描述的是引用的变更部分实际上**做了什么**，它的效果，而不是因此被做了什么。

## 首字母大写

```bash
# Good
Add `use` method to Credit model

# Bad
add `use` method to Credit model
```

> 首字母大写的原因是遵守英文句子开头使用大写字母的语法规则。

## 尽量做到只看注释便可明白而无需查看变更内容

```bash
# Good
Add `use` method to Credit model // 为 Credit 模块添加 `use` 方法

# Bad
Add `use` method // 添加 `use` 方法

---

# Good
Increase left padding between textbox and layout frame // 在 textbox 和 layout frame 之间添加向左对齐

# Bad
Adjust css // 就改了下 css
```

它在许多场景中（例如多次 commit、多个更改和重构）非常有用，可以帮助审查人员理解提交者的想法。

## 使用信息本身来解释“原因”、“目的”、“手段”和其他的细节

```bash
# Good
Fix method name of InventoryBackend child classes

Classes derived from InventoryBackend were not
respecting the base class interface.

It worked because the cart was calling the backend implementation
incorrectly.

# Good
Serialize and deserialize credits to json in Cart

Convert the Credit instances to dict for two main reasons:

  - Pickle relies on file path for classes and we do not want to break up
    everything if a refactor is needed
  - Dict and built-in types are pickleable by default

# Good
Add `use` method to Credit

Change from namedtuple to class because we need to
setup a new attribute (in_use_amount) with a new value
```

信息的主题和正文之间用空行隔开。其他空行被视为信息正文的一部分。

像“`-`”、“`*`”和“`\`”这样的字符可以提高可读性。

## 避免使用无上下文的信息

```bash
# Bad
Fix this

Fix stuff

It should work now

Change stuff

Adjust css
```

## 限制每行字数

[Pro Git Book](https://git-scm.com/book/zh/v2)建议：

- 主题最多使用50个字符，
- 正文最多使用72个字符。

## 保持语言的一致性

对于项目所有者而言：选择一种语言并使用该语言编写所有的 commit 信息。理想情况下，它应与代码注释、默认翻译区域（用于本地化项目）等相匹配。

对于贡献者而言：使用与现有 commit 历史相同的语言编写 commit 信息。

```bash
# Good
ababab Add `use` method to Credit model
efefef Use InventoryBackendPool to retrieve inventory backend
bebebe Fix method name of InventoryBackend child classes

# Good (Portuguese example)
ababab Adiciona o método `use` ao model Credit
efefef Usa o InventoryBackendPool para recuperar o backend de estoque
bebebe Corrige nome de método na classe InventoryBackend

# Bad (mixes English and Portuguese)
ababab Usa o InventoryBackendPool para recuperar o backend de estoque
efefef Add `use` method to Credit model
cdcdcd Agora vai
```
