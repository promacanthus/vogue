---
title: "05 Pending的请求"
date: 2020-07-29T18:40:54+08:00
draft: true
---

## 现象

发起请求后，Chrome浏览器显示的状态一直是pending，等待响应，最后超时了。

## 工具

Chrome的日志采集工具：`chrome://net-export/`
Chrome的日志分析工具：`https://netlog-viewer.appspot.com/#import`

正常的是这样的，有send，有read：

![image](/images/2020-07-30_11-18.png)

异常的就是这两个步骤中有一个步骤出现问题了。
