---
title: "02 jwt-go"
date: 2020-05-11T18:49:17+08:00
draft: true
---

## 社区的库

https://jwt.io/ 上列出了10个支持JWT的Golang库。

|名称|地址|star|支持的registered字段|支持的签名算法|
|---|---|---|---|---|
|jwt-go|github.com/dgrijalva/jwt-go|7571|exp、nbf|大部分|
|jose2go|github.com/dvsekhvalnov/jose2go|135|无|除了EdDSA|
|jose|github.com/SermoDigital/jose|838|全部|除了EdDSA|
|jwt|github.com/robbert229/jwt|79|exp、nbf、iat|HS256、HS512|
|go-jose|github.com/square/go-jose|1530|无|除了EdDSA|
|jwx|github.com/lestrrat-go/jwx|275|全部|除了EdDSA|
|jwt-auth|github.com/adam-hanna/jwt-auth|179|exp、nbf、jti|大部分|
|gojwt|github.com/nickvellios/gojwt|5|exp|HS256|
|jwt/v3|github.com/gbrlsnchs/jwt/v3|332|全部|除了EdDSA|
|jwt|github.com/pascaldekloe/jwt|179|全部|全部|
|jwt|gitlab.com/rucuriousyet/jwt|2|无|HS256、HS512|
|sjwt|github.com/brianvoe/sjwt|65|全部|HS256|
|jwt|github.com/cristalhq/jwt|57|全部|全部|

`payload`中`registered`预定义的字段介绍：

- iss（issuer）：发起者
- exp（expiration time）：过期时间
- sub（subject）：主题
- aud（audience）：受众
- nbf（Not Before）：生效时间
- iat（Issued At）：签发时间
- jti（JWT ID）：编号

## jwt-go使用指南

jwt-go在Go官方新的包管理站点的[地址](https://pkg.go.dev/github.com/dgrijalva/jwt-go?tab=doc)。

[Medium原文链接](https://medium.com/@Raulgzm/securing-golang-api-using-json-web-token-jwt-2dc363792a48)。

> 原文使用`mux`本文使用`gin`，同时将jwt升级为v4版本，使用经典的`MVC`架构来实现，Github仓库点击[这里](https://github.com/Promacanthus/jwt-go-example)。

JSON Web令牌（JWT）是一种更现代的身份验证方法。随着Web向客户端和服务器之间的更大距离转移，JWT提供了一种很好的替代传统基于cookie的身份验证模型的方法。

JWT为客户端提供了一种对每个请求进行身份验证的方法，而无需维护会话或将登录凭据重复传递给服务器。

![image](/images/1_iR3aP6QPoByZbypgQst6Yw.png)

## 使用基于令牌的方法的好处

- 跨域/CORS：Cookie + CORS在不同的域中不能很好地发挥作用。基于令牌的方法允许对任何域上的任何服务器进行`AJAX`调用，因为使用 `HTTP Header` 来传输用户信息。
- 无状态的：无需保留会话存储，令牌是一个独立的实体，可以传达所有用户信息。
- CDN：可以通过CDN提供应用的所有资产（例如`javascript`，`HTML`，图像等），而服务器端只是API。
- 解耦：不受特定认证方案的束缚。令牌可以在任何地方生成，因此，只需要使用这一种验证方式，就可以从任何地方调用服务端的API。
- 面向移动设备：当开始在原声平台（iOS，Android，Windows 等）上工作时，cookie并不是使用安全API（必须处理cookie容器）的理想选择。采用基于令牌的方法可以大大简化这一过程。
- CRSF（Cross-site request forgery）：由于不依赖Cookie，因此无需防御跨站点请求伪造。

## 样例代码

主要步骤：

1. 创建Token对象
2. 填充Token中`Header`和`Payload`
3. 生成Token中的`Signature`
4. 从HTTP请求中获取并验证Token

```go
import jwtv4 "github.com/dgrijalva/jwt-go/v4"

// example-1
token := jwtv4.New(jwtv4.SigningMethodRS512)

    token.Claims = jwtv4.MapClaims{
        "exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
        "iat": time.Now().Unix(),
        "sub": userUUID,
    }

// example-2
token := jwtv4.NewWithClaims(jwtv4.SigningMethodRS512,jwtv4.MapClaims{
        "exp": time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
        "iat": time.Now().Unix(),
        "sub": userUUID,
    })

    tokenString, err := token.SignedString(j.privateKey)
```

解析：

`Token`是一个结构体，如下：

```go
type Token struct {
    Raw       string                 // 原始令牌，解析令牌时填充
    Method    SigningMethod          // 使用或将要使用的签名方法
    Header    map[string]interface{} // 令牌的第一段
    Claims    Claims                 // 令牌的第二段
    Signature string                 // 令牌的第三段，解析令牌时填充
    Valid     bool                   // 判断令牌是否有效，解析/验证令牌时填充
}
```

创建`Token`对象使用`New`函数：

```go
func New(method SigningMethod) *Token
// 其实New函数是对NewWithClaims函数的封装
func New(method SigningMethod) *Token {
    return NewWithClaims(method, MapClaims{})
}
```

`New`函数需要一个`SigingMethod`类型的形参。

`SigingMethod`是一个接口类型：

```go
type SigningMethod interface {
    Verify(signingString, signature string, key interface{}) error // 如果签名有效，则返回nil
    Sign(signingString string, key interface{}) (string, error)    // 返回编码后的签名或错误
    Alg() string                                                   // 返回此方法的alg标识符（例如：“HS256”）
}
```

根据Go语言的隐式接口规则，只要实现了上述三个方法的类型都可以作为`New`函数的形参。

`jwt-go`包中实现了`ECDSA`、`HMAC`、`RSA`三大类算法。

```go
// ECDSA
var (
    SigningMethodES256 *SigningMethodECDSA
    SigningMethodES384 *SigningMethodECDSA
    SigningMethodES512 *SigningMethodECDSA
)

// HMAC
var (
    SigningMethodHS256  *SigningMethodHMAC
    SigningMethodHS384  *SigningMethodHMAC
    SigningMethodHS512  *SigningMethodHMAC
)

// RSA
var (
    SigningMethodRS256 *SigningMethodRSA
    SigningMethodRS384 *SigningMethodRSA
    SigningMethodRS512 *SigningMethodRSA
)

```

`payload`部分可以使用预定义的，也可以使用自定义的。

JWT官方定义中，生成签名的算法是`HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), secret)`这里需要一个`secret`字段。我们使用`SignedString`方法，这是整个过程中最消耗资源的地方。

```go
func (t *Token) SignedString(key interface{}) (string, error)
// 这里的key，可以用非对称加密算法生成的私钥

// 使用openssl工具生成公私钥对
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
```

到这里JWT制作完成，下面开始从HTTP请求中验证Token。

```go
import (
    jwtv4 "github.com/dgrijalva/jwt-go/v4"
    "github.com/dgrijalva/jwt-go/v4/request"
)

token, err := request.ParseFromRequest(ctx.Request, request.OAuth2Extractor, func(token *jwtv4.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwtv4.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        } else {
            return authBackend.PublicKey, nil
        }
    })
```

jwt的request子模块中有相应的函数`ParseFromRequest`和`ParseFromRequestWithClaims`。

```go
func ParseFromRequest(req *http.Request, extractor Extractor, keyFunc jwt.Keyfunc, options ...ParseFromRequestOption) (token *jwt.Token, err error)
func ParseFromRequestWithClaims(req *http.Request, extractor Extractor, claims jwt.Claims, keyFunc jwt.Keyfunc) (token *jwt.Token, err error)
```

从HTTP请求中提取并解析JWT令牌。它的行为与Parse相同，但是接受`Request`和`Extractor`而不是令牌字符串。`Extractor`接口允许自定义提取令牌的逻辑。库中提供了几种有用的实现，如下所示，同时可以通过`ParseFromRequestOption`来修改解析行为。

```go
// OAuth2访问令牌的提取器。在“Authorization” Header 的“access_token”参数中查找令牌
var OAuth2Extractor = &MultiExtractor{
    AuthorizationHeaderExtractor,
    ArgumentExtractor{"access_token"},
}
```

其中`keyFunc`是`func(*Token) (interface{}, error)`函数类型，解析方法使用此回调函数提供验证密钥。该函数接收已解析但未验证的令牌。这使你可以使用令牌头中的属性（例如“kid”）来标识要使用的密钥。
