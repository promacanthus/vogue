---
title: CORS和Cookie.md
date: 2020-04-14T10:09:14.238627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 填坑记
- HTTP
summary: CORS和Cookie.md
showInMenu: false

---

# CORS和Cookie

## 问题

Chrome 80 版本开始，跨站访问时直接不携带cookie进行请求的发送，其他的浏览器可以正常访问。

因此，定位到和浏览器版本有关，新增了一个cookie的属性，samesite。

## CORS

跨域资源共享(CORS) 是一种机制，它使用额外的 HTTP 头来告诉浏览器，让运行在一个 origin (domain) 上的Web应用被准许访问来自不同源服务器上的指定的资源。

当一个资源从与该资源本身所在的服务器不同的域、协议或端口请求一个资源时，资源会发起一个跨域 HTTP 请求。

> 比如，站点 `http://domain-a.com` 的某 HTML 页面通过 `<img>` 的 src 请求 `http://domain-b.com/image.jpg`。

网络上的许多页面都会加载来自不同域的CSS样式表，图像和脚本等资源。出于安全原因，浏览器限制从脚本内发起的跨源HTTP请求，也可能是跨站请求可以正常发起，但是返回结果被浏览器拦截了。

> 例如，`XMLHttpRequest`和`Fetch API`遵循同源策略。 这意味着使用这些API的Web应用程序只能从加载应用程序的同一个域请求HTTP资源，除非响应报文包含了正确CORS响应头。

跨域资源共享（ CORS ）机制允许 Web 应用服务器进行跨域访问控制，从而使跨域数据传输得以安全进行。现代浏览器支持在 API 容器中（例如 XMLHttpRequest 或 Fetch ）使用 CORS，以降低跨域 HTTP 请求所带来的风险。

> CORS请求失败会产生错误，但是为了安全，在JavaScript代码层面是无法获知到底具体是哪里出了问题。只能查看浏览器的控制台以得知具体是哪里出现了错误。

## Cookie

cookie 属性:“path, domain, expire, HttpOnly, Secure”，现在增加了一个新属性 SameSite，一种新的防止跨站点请求伪造（cross site request forgery）的 http 安全特性。

Chrome 80 默认将没有设置SameSite设置为SameSite=Lax。

|值|描述|
|---|---|
|Strict|最为严格，完全禁止第三方Cookie，跨站点时，任何情况下都不会发送Cookie|
|Lax|稍稍放宽，大多数情况也是不发送第三方Cookie，但是导航到目标网址的 `Get` 请求除外|
|None|网站可以选择显式关闭SameSite属性，将其设为None。前提是必须同时设置Secure属性（Cookie 只能通过 HTTPS 协议发送），否则无效|

## HTTP 响应首部字段

更多内容点击[这里](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Access_control_CORS)。

响应首部中可以携带一个 `Access-Control-Allow-Origin` 字段，其语法如下:

```http
Access-Control-Allow-Origin: <origin> | *
```

其中，origin 参数的值指定了允许访问该资源的外域 URI。对于不需要携带身份凭证的请求，服务器可以指定该字段的值为通配符，表示允许来自所有域的请求。

例如，下面的字段值将允许来自 `http://mozilla.com` 的请求：

```http
Access-Control-Allow-Origin: http://mozilla.com
```

如果服务端指定了具体的域名而非“`*`”，那么响应首部中的 `Vary` 字段的值必须包含 `Origin`。这将告诉客户端：服务器对不同的源站返回不同的内容。

### 导航到目标网站的Get请求

导航到目标网址的 GET 请求，只包括三种情况：链接，预加载请求，GET 表单。详见下表。

|请求类型|示例|正常情况|Lax|
|---|---|---|---|
|链接|`<a href="..."></a>`|发送 Cookie|发送 Cookie|
|预加载|`<link rel="prerender" href="..."/>`|发送 Cookie|发送 Cookie|
|GET 表单|`<form method="GET" action="...">`|发送 Cookie|发送 Cookie|
|POST 表单|`<form method="POST" action="...">`|发送 Cookie|不发送|
|iframe|`<iframe src="..."></iframe>`|发送 Cookie|不发送|
|AJAX|`$.get("...")`|发送 Cookie|不发送|
|Image|`<img src="...">`|发送 Cookie|不发送|

设置了Strict或Lax以后，基本就杜绝了 [CSRF](https://baike.baidu.com/item/%E8%B7%A8%E7%AB%99%E8%AF%B7%E6%B1%82%E4%BC%AA%E9%80%A0/13777878?fromtitle=CSRF&fromid=2735433) 攻击。当然，前提是用户浏览器支持 SameSite 属性。

### 修改Chrome设置

1. 谷歌浏览器地址栏输入：`chrome://flags/`
2. 找到：`SameSite by default cookies`、`Cookies without SameSite must be secure`
3. 设置上面这两项设置成 Disable
