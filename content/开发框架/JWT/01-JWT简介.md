---
title: "01 JWT简介"
date: 2020-05-08T14:42:52+08:00
draft: true
---

新闻：免费获取[JWT手册](https://auth0.com/resources/ebooks/jwt-handbook?_ga=2.201301742.175960213.1588919939-1075582363.1588919939)并深入学习JWT！

## 什么是JWT

JSON Web Token（JWT）是一个开放标准（[RFC 7519](https://tools.ietf.org/html/rfc7519)），它定义了一种紧凑且自包含的方式，用于在各方之间以JSON对象安全地传输信息。由于此信息是经过数字签名的，因此可以被验证和信任。可以使用secret（使用`HMAC`算法）或使用`RSA`/`ECDSA`的公私密钥对对JWT进行签名。

尽管可以对JWT进行加密以在双方之间提供保密性，但我们将重点放在已签名的令牌。签名的令牌可以验证其中包含的声明的完整性，而加密的令牌则将这些声明信息隐藏了起来。当使用公/私密钥对对令牌进行签名时，签名可以证明只有持有私钥的一方才是对其进行签名的一方。

## 何时使用JWT

以下是JSON Web令牌常用的场景：

- 授权：这是使用JWT的最常见方案。一旦用户登录，每个后续请求将包括JWT，从而允许用户访问该令牌允许的路由，服务和资源。**单点登录**是当今广泛使用JWT的一项功能，因为它的开销很小并且可以在**不同的域**中轻松使用。（解决跨域认证的另一种方式，session持久化）。
- 信息交换：JWT是在各方之间安全地传输信息的一种好方法。因为可以对JWT进行签名（例如，使用公/私密钥对），所以可以确定信息发送者的身份。另外，由于签名是使用`header`和`payload`计算的，因此还可以验证内容是否被篡改。

## JWT的结构

JWT以紧凑的形式由三部分组成，这些部分由点（`.`）分隔，分别是：

- Header
- Payload
- Signature

因此，JWT通常如下所示，`xxxxx.yyyyy.zzzzz`。

让我们分解不同的部分。

### Header

Header通常由两部分组成：令牌的类型（在这里只能是`JWT`）和所使用的签名算法，例如`HMAC`、`SHA256`或`RSA`。例如：

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

然后，此JSON被`Base64Url`编码以形成JWT的第一部分。

上面代码中：

- alg属性表示签名的算法（algorithm），默认是 HMAC SHA256（写成 HS256）；
- typ属性表示这个令牌（token）的类型（type），JWT 令牌统一写为JWT。

### Payload

令牌的第二部分是payload，其中包含声明。声明是有关实体（通常是用户）和其他数据的语句。共有三种类型的声明：

- registered
- public
- private

#### registered

这是一组预定义的字段，不是强制性的，而是建议使用的，以提供一组有用的可互操作的声明。官方定义了7个字段：

- iss（issuer）：发起者
- exp（expiration time）：过期时间
- sub（subject）：主题
- aud（audience）：受众
- nbf（Not Before）：生效时间
- iat（Issued At）：签发时间
- jti（JWT ID）：编号

> 请注意，声明名称时仅使用三个字符代表，因为JWT是紧凑的。

#### public

这些可以由使用JWT的人员随意定义。但是为避免冲突，应在IANA JSON Web Token注册表中定义它们，或将其定义为包含抗冲突命名空间的URI。

#### private

这些是自定义声明，目的是在同意使用它们的各方之间共享信息，既不是注册声明也不是公共声明。

一个`payload`示例如下：

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

然后，对payload进行`Base64Url`编码，以形成JWT的第二部分。

> 请注意，对于已签名的令牌，此信息尽管可以**防止篡改**，但任何人都可以读取。除非将其**加密**，否则请勿将机密信息放入JWT的有效负载或报头元素中。

### Signature

要创建签名部分，必须获取编码后的`Header`、编码后的`payload`，一个`secret`，`Header`中指定的算法，然后就可以对JWT进行签名。

例如，如果要使用`HMAC SHA256`算法，则将通过以下方式创建签名：

```bash
# 生产的内容就是JWT的第三部分
HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)
```

签名用于验证消息在此过程中是否被篡改，并且对于使用私钥进行签名的令牌，它还可以验证JWT的发送者的真实身份。

### 结果

最后生成的结果是由点分隔的被Base64-URLi编码的三个字符串，可以在HTML和HTTP环境中轻松传递这些字符串，与基于XML的标准（例如SAML）相比，它更紧凑。

下图显示了一个JWT，它已对`header`和`payload`进行了编码，并用一个`secret`进行了签名。

如果想使用JWT，可以使用[jwt.io Debugger](https://jwt.io/#debugger-io)解码，验证和生成JWT。

## JWT工作原理

### 授权

在认证登录中，当用户使用其凭据成功登录时，将返回JWT。由于令牌是凭据，因此必须格外小心以防止安全问题。通常，令牌的保留时间不应超过要求的时间。

[由于缺乏安全性，也不应该将敏感的会话数据存储在浏览器存储中](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#local-storage)。

每当用户想要访问受保护的路由或资源时，用户代理（一般是浏览器）应该发送JWT，通常的做法是在`Authorization`请求头中使用**Bearer schema**。请求头的内容应如下所示：`Authorization: Bearer <token>`。

在某些情况下，这是一种无状态的授权机制。服务器的受保护路由将在`Authorization`请求头中检查有效的JWT，如果存在，则将允许用户访问受保护的资源。如果JWT包含必要的数据，则可以减少查询数据库中某些操作的需求，尽管这种情况并非总是如此。

**如果在`Authirization`请求头中发送令牌，则跨域资源共享（CORS）不会成为问题，因为它不使用cookie**。

下图显示了如何获取JWT并将其用于访问API或资源：

![image](/images/client-credentials-grant.png)

1. 应用程序或客户端向授权服务器请求授权。这是一种不同的授权流程执行。例如，典型的符合[OpenID Connect](https://openid.net/connect/)规范的Web应用程序将使用[授权代码流程](https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth)访问`/oauth/authorize`端点。
2. 授予授权后，授权服务器将访问令牌返回给应用程序。
3. 应用程序使用访问令牌来访问受保护的资源（例如API）。

请注意，在使用签名令牌时，令牌或令牌中包含的所有信息都会暴露给用户或其他方，虽然他们无法更改它，因此不应将机密信息放入令牌中。

## 为何使用JWT

让我们谈谈与Simple Web Tokens（SWT）和Security Assertion Markup Language Tokens（SAML）相比，JSON Web Tokens（JWT）的好处。

1. JSON没有XML冗长，因此在编码时JSON会更小一些，从而使JWT比SAML更紧凑。这使得JWT是在HTML和HTTP环境中传递的不错的选择。
2. 在安全方面，SWT只能通过`HMAC`算法使用共享密码进行对称签名。但是，JWT和SAML令牌可以使用`X.509`证书形式的公用/专用密钥对进行签名。与签名JSON的简单性相比，使用XML Digital Signature签名XML而不引入模糊的安全漏洞是非常困难的。
3. JSON解析器在大多数编程语言中都很常见，因为它们直接映射到对象。相反，XML中没有很自然的文档-对象映射。因此使用JWT比SMAL断言更容易。
4. 在用法方面，JWT是在互联网上广泛使用的，这说明在多平台（有其实移动平台）场景下，客户端侧对JWT进行处理是很容易的。

如果想了解有关JWT的更多信息，或者在应用程序中使用它进行身份验证，请查看Auth0上的[JWT登录页面](http://auth0.com/learn/json-web-tokens?_ga=2.222409243.1097272832.1589187420-1075582363.1588919939)。
