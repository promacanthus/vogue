---
title: "06 string2time"
date: 2020-05-11T13:04:47+08:00
draft: true
---

```go
    const shortForm = "2006-01-02"
    parse, err := time.Parse(shortForm, tmp)
```

这里的`shortForm`的内容不能任意定义，Golang中这些数字都是有特殊函义的，见下面列表：

|单位|可选值|
|---|---|
|月份 |1,01,Jan,January
|日　 |2,02,_2
|时　 |3,03,15,PM,pm,AM,am
|分　 |4,04
|秒　 |5,05
|年　 |06,2006
|周几 |Mon,Monday
|时区时差表示 |-07,-0700,Z0700,Z07:00,-07:00,MST
|时区字母缩写 |MST
