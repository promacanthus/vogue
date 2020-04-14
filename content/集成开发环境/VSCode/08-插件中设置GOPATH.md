---
title: 08-插件中设置GOPATH
date: 2020-04-14T10:09:14.278627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 集成开发环境
- VSCode
summary: 08-插件中设置GOPATH
showInMenu: false

---

在任何使用都可以使用`Go: Current GOPATH`命令来查看扩展程序使用的`GOPATH`。

## 从环境变量获得`GOPATH`

开箱即用，扩展程序使用环境变量`GOPATH`中设置的值。从Go1.8版本开始，如果没有设置这个环境变量，那么会使用`go env`中描述的`GOPATH`值。

## 通过`go.gopath`设置`GOPATH`

在**用户设置**中设置`go.gopath`会覆盖从上面的逻辑中派生的`GOPATH`。

在**工作空间**中设置`go.gopath`会覆盖用户设置中的`go.gopath`。可以在此文件中将多个文件夹目录设置为`GOPATH`。

> **注意，多个目录之间使用`:`分隔（在window系统中使用`;`分隔）**。

## 通过`go.inferGopath`设置`GOPATH`

设置`go.interGopath`会覆盖在`go.gopath`中设置的值。

如果`go.inferGopath`设置为`true`,那么扩展程序将会尝试从工作空间（例如，VSCode中打开的目录）的路径中推断`GOPATH`。扩展程序会从`src`目录的路径开始向上搜索，并将`GOPATH`设置为高于该目录的一个级别，其中也包括全局的`GOPATH`。

运行命令`go env GOPATH`来查找全局的`GOPATH`。

> 例如，如果项目是这样的`/aaa/bbb/ccc/src/...`,那么打开目录`/aaa/bbb/ccc/src`或者其中的任何子内容，都将会导致扩展程序自动向上搜索，在路径中找到`src`目录，然后将`GOPATH`设置为高于它一个级别的值，即`GOPATH=/aaa/bbb/ccc`

当处理具有不同`GOPATH`的Go项目时，这个设置非常有用。

不需要在每个工作空间中设置`GOPATH`，也不需要在工作空间中设置全部的目录然后用`:`或`;`来分隔。

只需要将`go.inferGopath`设置为`true`，扩展程序就会自动搜索并使用正确的`GOPATH`。

## 通过`go.toolsGopath`设置Go工具的`GOPATH`

使用`go get`命令在`GOPATH`中安装Go工具，要防止Go工具扰乱`GOPATH`,可以使用`go.toolsGopath`设置来提供单独的`GOPATH`仅用于Go工具。

第一次设置`go.toolsGopath`时，需要运行`GO: Install/Update Tools`命令，这样Go工具就能够安装在指定的位置。

如果没有设置`go.toolsGopath`或者在它指定的位置没有找到Go工具，那么会使用上面几节中描述的逻辑来找`GOPATH`中的Go工具。如果在那里也没有找到Go工具，那么就会在路径找中查找，这些路径时`PATH`环境变量的一部分。
