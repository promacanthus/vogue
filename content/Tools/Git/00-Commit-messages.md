---
title: "00 Commit规范和Messages"
date: 2020-05-11T13:49:42+08:00
draft: true
---
- [0.1. 规范](#01-规范)
  - [0.1.1. Header](#011-header)
  - [0.1.2. Body](#012-body)
  - [0.1.3. Footer](#013-footer)
  - [0.1.4. commit.template](#014-committemplate)
- [0.2. Messages](#02-messages)
  - [0.2.1. 使用祈使句](#021-使用祈使句)
  - [0.2.2. 首字母大写](#022-首字母大写)
  - [0.2.3. 尽量做到只看注释便可明白而无需查看变更内容](#023-尽量做到只看注释便可明白而无需查看变更内容)
  - [0.2.4. 使用信息本身来解释“原因”、“目的”、“手段”和其他的细节](#024-使用信息本身来解释原因目的手段和其他的细节)
  - [0.2.5. 避免使用无上下文的信息](#025-避免使用无上下文的信息)
  - [0.2.6. 限制每行字数](#026-限制每行字数)
  - [0.2.7. 保持语言的一致性](#027-保持语言的一致性)

## 0.1. 规范

参考 Angular 团队的 commit 规范。

```console
<type>(<scope>): <subject>
// 空一行
<body>
// 空一行
<footer>
```

分别对应 Commit message 的三个部分： Header ， Body 和 Footer 。

### 0.1.1. Header

Header 部分只有一行，包括三个字段： `type` （必需）、 `scope` （可选）和 `subject` （必需）。

- type : 用于说明 commit 的类型。一般有以下几种:
  - feat: 新增feature
  - fix: 修复bug
  - docs: 仅仅修改了文档，如readme.md
  - style: 仅仅是对格式进行修改，如逗号、缩进、空格等。不改变代码逻辑
  - refactor: 代码重构，没有新增功能或修复bug
  - perf: 优化相关，如提升性能、用户体验等
  - test: 测试用例，包括单元测试、集成测试
  - chore: 改变构建流程、或者增加依赖库、工具等
  - revert: 版本回滚
- scope : 用于说明 commit 影响的范围，比如: views, component, utils, test...
- subject : commit 目的的简短描述

### 0.1.2. Body

对本次 commit 修改内容的具体描述, 可以分为多行。如下所示:

```console
# body: 72-character wrapped. This should answer:
# * Why was this change necessary?
# * How does it address the problem?
# * Are there any side effects?
# initial commit
```

### 0.1.3. Footer

一些备注, 通常是 BREAKING CHANGE (当前代码与上一个版本不兼容) 或修复的 bug(关闭 Issue) 的链接。

### 0.1.4. commit.template

```bash
#  这个命令只能设置当前分支的提交模板
git config commit.template   [模板文件名]  

# 这个命令能设置全局的提交模板，注意global前面是两杠
git config  --global commit.template   [模板文件名]
```

新建 `.gitmessage.txt` (模板文件) 内容可以如下:

```txt
# headr: <type>(<scope>): <subject>
# - type: feat, fix, docs, style, refactor, test, chore
# - scope: can be empty
# - subject: start with verb (such as 'change'), 50-character line
#
# body: 72-character wrapped. This should answer:
# * Why was this change necessary?
# * How does it address the problem?
# * Are there any side effects?
#
# footer:
# - Include a link to the issue.
# - BREAKING CHANGE
#
```

## 0.2. Messages

原仓库[地址](https://github.com/RomuloOliveira/commit-messages-guide/blob/master/README_zh-CN.md)。

### 0.2.1. 使用祈使句

```bash
# good
Use InventoryBackendPool to retrieve inventory backend  // 用 InventoryBackendPool 获取库存

# bad
Used InventoryBackendPool to retrieve inventory backend // InventoryBackendPool 被用于获取库存
```

> commit 信息描述的是引用的变更部分实际上**做了什么**，它的效果，而不是因此被做了什么。

### 0.2.2. 首字母大写

```bash
# Good
Add `use` method to Credit model

# Bad
add `use` method to Credit model
```

> 首字母大写的原因是遵守英文句子开头使用大写字母的语法规则。

### 0.2.3. 尽量做到只看注释便可明白而无需查看变更内容

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

### 0.2.4. 使用信息本身来解释“原因”、“目的”、“手段”和其他的细节

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

### 0.2.5. 避免使用无上下文的信息

```bash
# Bad
Fix this

Fix stuff

It should work now

Change stuff

Adjust css
```

### 0.2.6. 限制每行字数

[Pro Git Book](https://git-scm.com/book/zh/v2)建议：

- 主题最多使用50个字符，
- 正文最多使用72个字符。

### 0.2.7. 保持语言的一致性

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
