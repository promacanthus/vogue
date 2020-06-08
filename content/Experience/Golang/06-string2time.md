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

看一看`time`包中定义的常量，就能理解啦，这些是预定义的布局，用于`Time.Format`和`time.Parse`。布局中使用的参考时间是特定时间：`Mon Jan 2 15:04:05 MST 2006`，换成Unix时间是`1136239445`

```go
const (
    ANSIC       = "Mon Jan _2 15:04:05 2006"
    UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
    RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
    RFC822      = "02 Jan 06 15:04 MST"
    RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
    RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
    RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
    RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
    RFC3339     = "2006-01-02T15:04:05Z07:00"
    RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
    Kitchen     = "3:04PM"
    // Handy time stamps.
    Stamp      = "Jan _2 15:04:05"
    StampMilli = "Jan _2 15:04:05.000"
    StampMicro = "Jan _2 15:04:05.000000"
    StampNano  = "Jan _2 15:04:05.000000000"
)
```
