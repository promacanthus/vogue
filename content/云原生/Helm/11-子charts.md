---
title: 11-子charts.md
date: 2020-04-14T10:09:14.126627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 云原生
- Helm
summary: 11-子charts.md
showInMenu: false

---

# 11-子charts

chart有称为子chart的依赖关系，它们也有自己的值和模板。关于子chart的注意点：

1. 子chart被认为是“独立的”，子chart不能明确依赖于其父chart
2. 子chart无法访问其父项的值
3. 父chart可以覆盖子chart的值
4. Helm有全局值的概念，可以被所有charts访问

> 所有子chart都保存在父chart的`charts/`目录中。

**每个子chart都是独立的 chart**。

## 覆盖子chart中的值

在父chart的values.yaml文件中添加如下信息：

```yaml
favorite:
  drink: coffee
  food: pizza
pizzaToppings:
  - mushrooms
  - cheese
  - peppers
  - onions

mysubchart:                   # 子chart的名字，这部分以下是所有内容都会发送给子chart
  dessert: ice cream          # 子chart的values.yaml中的值
```

这里有一个重要的细节需要注意：没有改变子chart模板指向 `.Values.mysubchart.dessert`。从该子chart模板的角度来看，该值仍位于 `.Values.dessert`。随着模板引擎一起传递值，它会设置范围。所以对于 mysubchart 模板，只有指定给 mysubchart 的值才会在 `.Values` 里。

## 全局chart值

全局值是可以从任何chart或子chart用完全相同的名称访问的值。**全局值需要明确声明，不能像使用现有的非全局值一样来使用全局值**。

values 数据类型有一个保留部分，称为 `Values.global`, 可以设置全局值。如下所示：

```yaml
favorite:
  drink: coffee
  food: pizza
pizzaToppings:
  - mushrooms
  - cheese
  - peppers
  - onions

mysubchart:
  dessert: ice cream

global:                 # 全局值
  salad: caesar
```

## 与子chart共享模板

父chart和子chart可以共享命名模板。任何chart中的任何定义块都可用于其他chart，如下所示：

```yaml
{{- define "labels"}}
  from: mychart
{{ end }}
# 命名模板的名字是全局共享的，任何chart都可以用
```
