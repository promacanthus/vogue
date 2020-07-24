---
title: "04 LDAP"
date: 2020-07-24T10:26:53+08:00
draft: true
---

- [0.1. 简介](#01-简介)
- [0.2. 协议概观](#02-协议概观)
- [0.3. 目录结构](#03-目录结构)
- [0.4. 操作命令](#04-操作命令)
  - [0.4.1. Add](#041-add)
  - [0.4.2. Bind（认证）](#042-bind认证)
  - [0.4.3. Delete](#043-delete)
  - [0.4.4. Search and Compare](#044-search-and-compare)
    - [0.4.4.1. baseObject](#0441-baseobject)
    - [0.4.4.2. scope](#0442-scope)
    - [0.4.4.3. filter](#0443-filter)
    - [0.4.4.4. derefAliases](#0444-derefaliases)
    - [0.4.4.5. attributes](#0445-attributes)
    - [0.4.4.6. sizeLimit, timeLimit](#0446-sizelimit-timelimit)
    - [0.4.4.7. typesOnly](#0447-typesonly)
  - [0.4.5. Modify](#045-modify)
  - [0.4.6. Modify DN](#046-modify-dn)
  - [0.4.7. Extended operations](#047-extended-operations)
    - [0.4.7.1. StartTLS](#0471-starttls)
  - [0.4.8. Abandon](#048-abandon)
  - [0.4.9. Unbind](#049-unbind)
- [0.5. URI scheme](#05-uri-scheme)
- [0.6. Schema](#06-schema)

## 0.1. 简介

Lightweight Directory Access Protocol (LDAP) 是一个开源的与供应商无关的工业级标准应用层协议，用于通过IP协议访问和维护分布式目录信息服务。

不论是开发内网还是公网应用程序，域目录服务都发挥着重要作用，它允许在整个网络中共享有关用户，系统，网络，服务和应用程序的信息。

因此，目录服务可以提供任何组织的记录集合，通常具有分层结构，例如公司的电子邮件目录。类似的电话目录就是一个带有客户地址和电话号码的列表。

> LDAP的最新版本是V3，发布在[RFC 4511](https://tools.ietf.org/html/rfc4511)。[Request for Comments 简称RFCs]。

LDAP的常见用法是提供一个中心位置来存储用户名和密码。这允许许多不同的应用程序和服务连接到LDAP服务器以验证用户。

LDAP是基于[X.500](https://en.wikipedia.org/wiki/X.500)标准中包含的标准的更简单子集。因此，LDAP也称为`X.500-lite`。

## 0.2. 协议概观

客户端启动一个LDAP会话来连接到一台LDAP服务器（称为目录系统代理DSA），默认连接到`TCP`和`UDP`的`389`端口上，或者在`LDAPS`（基于SSL的LDAP）的`636`端口上。

然后客户端发送一个操作请求到服务器，服务器返回响应信息。除了某些例外，客户端无需在发送下一个请求之前等待上一个响应返回，服务器可以按任意顺序返回响应。所有请求和响应信息都使用基本编码规则（Basic Encoding Rules,BER）编码后再进行网络传输。

客户端可以执行如下请求操作：

- StartTLS：使用LDAPv3 TLS扩展进行安全连接
- Bind：认证并指定LDAP协议版本
- Search：搜索或检索目录条目
- Compare：测试一个给定的条目是否包含给定的属性值
- Add：添加一个新条目
- Delete：删除一个条目
- Modify：修改一个条目
- Modify Distinguished Name (DN)：移动或重命名一个条目
- Abandon：终止上一个请求
- Extended Operation：用于定义其他操作的通用操作
- Unbind：关闭连接（不是bind的逆操作）

另外，服务器可以发送“未经请求的通知（Unsolicited Notifications）”，这不是对任何请求的响应信息。例如，连接超时之前。

保护LDAP通信的一种常见可选方案是使用SSL隧道。LDAP over SSL的默认端口是`636`。在LDAPv2中普遍使用基于LDAP over SSL，但从未在任何正式规范中对其进行标准化。LDAPv2已于2003年正式停用。

## 0.3. 目录结构

LDAP协议提供了一个目录，该目录遵循`X.500`模型的1993版：

- 一个条目由一组属性组成
- 一个属性具有名称（表示属性类型或属性描述）和一个或多个值，这些属性在scheme中定义
- 每个条目都有一个唯一的标识符：它的专有名称（Distinguished Name，DN），由它的相对专有名称（Relative Distinguished Name，RDN）组成，和条目中的某些属性组成，后面紧跟着父条目的DN。

> 将DN视为完整文件路径，将RDN视为其父文件夹中的相对文件名（例如，如果`/foo/bar/myfile.txt`是DN，则`myfile.txt`是RDN）。

当条目在路径树中移动时，DN可能会在条目的生存期内发生变化。为了可靠且明确地标识条目，可以在条目的操作属性集中提供`UUID`。

当以LDAP数据交换格式（LDAP Data Interchange Format，LDIF）表示时，条目看起来像这样（LDAP本身是二进制协议）：

```go
 dn: cn=John Doe,dc=example,dc=com  
 cn: John Doe
 givenName: John
 sn: Doe
 telephoneNumber: +1 888 555 6789
 telephoneNumber: +1 888 555 1232
 mail: john@example.com
 manager: cn=Barbara Doe,dc=example,dc=com
 objectClass: inetOrgPerson
 objectClass: organizationalPerson
 objectClass: person
 objectClass: top
```

- `dn`是条目的专有名称，它既不是属性也不是条目的一部分，`dc=example,dc=com`是DN的父条目，其中`dc`表示域组件(Domain Component)
- `cn`是条目的相对专有名称

剩下的行表示的是这个条目的属性。属性名称通常是助记符字符串，例如：

- `cn`：（common name）代表通用名称
- `dc`：（domain component）代表域组件
- `mail`：代表电子邮件地址
- `sn`：（surname）代表姓氏

服务器持有一个从特定条目开始的子路径树，例如：`dc=example,dc=com`以及它的子级。

服务器可能还会保留对其他服务器的引用，因此尝试访问`ou = department，dc = example，dc = com`可能会返回对拥有目录树部分的服务器的引用或延续引用。然后，客户端可以联系另一台服务器。

某些服务器还支持链式（chaining），这意味着该服务器将会与另一台服务器通信并将结果返回给客户端。

LDAP很少定义顺序：服务器可以按任意顺序返回属性的值，条目中的属性以及通过搜索操作找到的条目。

> 正式定义描述：条目定义为一组属性，而属性是一组值，并且这些组不需要排序。

## 0.4. 操作命令

### 0.4.1. Add

ADD操作将新条目插入目录服务器数据库。如果目录中已经存在Add请求中的DN，则服务器将不会添加重复条目，但会将添加结果中的结果代码设置为十进制的68“`entryAlreadyExists`”。

- 在尝试查找条目时，符合LDAP的服务器将永远不会取消在Add请求中传输过来的专有名称的引用，即专有名称不会被解除别名。
- 符合LDAP的服务器将确保专有名称和所有属性符合命名标准。
- 要添加的条目必须不存在，并且直接上级必须存在。

```go
dn: uid=user,ou=people,dc=example,dc=com
changetype: add
objectClass:top
objectClass:person
uid: user
sn: last-name
cn: common-name
userPassword: password
```

在上面的例子中：

- `uid=user,ou=people,dc=example,dc=com`必须不存在
- `ou=people,dc=example,dc=com`必须已存在

### 0.4.2. Bind（认证）

创建LDAP会话时，即LDAP客户端连接到服务器时，该会话的身份验证状态设置为匿名。 BIND操作建立会话的身份验证状态。

Simple BIND和SASL PLAIN可以以纯文本形式发送用户的DN和密码，因此，应使用TLS对使用Simple BIND或SASL PLAIN的连接进行加密。服务器通常查看给定条目中的`userPassword`属性来检查密码。Anonymous BIND（具有空DN和密码）会将连接重置为匿名状态。

SASL（简单身份验证和安全层，Simple Authentication and Security Layer）BIND提供多种认证机制（例如， Kerberos或通过TLS发送的客户端证书。

BIND还通过以整数形式发送版本号来设置LDAP协议版本。如果客户端请求了服务器不支持的版本，则服务器必须在BIND响应中将结果代码设置为协议错误代码。通常，客户端应使用`LDAPv3`，这是协议中的默认设置，但并非总是LDAP库中的默认设置。

在`LDAPv2`中，BIND必须是会话中的第一个操作，但是从`LDAPv3`开始，它不是必需的。在`LDAPv3`中，每个成功的BIND请求都会更改会话的身份验证状态，而每个失败的BIND请求都会重置会话的身份验证状态。

### 0.4.3. Delete

要删除一个条目，LDAP客户端应将格式正确的删除请求发送到服务器。

- 删除请求必须包含要删除条目的DN
- Request controls也可以附加到删除请求中
- 服务器在处理删除请求时不会取消引用别名
- 删除请求只能删除叶条目（没有下属的条目）
  - 一些服务器支持操作属性`hasSubordinates`，该属性指示条目是否具有任何从属条目
  - 一些服务器支持操作属性`numSubordinates`，该属性指示从属于包含`numSubordinates`属性的条目的条目数
- 某些服务器支持删除子树的`request control`，该请求允许删除DN以及从属于DN的所有对象。删除请求受访问控制的约束，即是否允许具有给定身份验证状态的连接删除给定条目由服务器特定的访问控制机制决定

### 0.4.4. Search and Compare

Search操作用于搜索和读取条目，其参数如下。

#### 0.4.4.1. baseObject

相对于要执行搜索的基础对象条目（可能的根）的名称。

#### 0.4.4.2. scope

在`baseObject`下面要搜索哪些元素。可以是：

- `BaseObject`：仅搜索命名的条目，通常用于读取一个条目
- `singleLevel`：在base DN下方的条目
- `WholeSubtree`：从base DN开始的整个子树

#### 0.4.4.3. filter

过滤范围内元素。

例如：

```go
(＆(objectClass = person)(|(givenName = John)(mail = john *)))
```

选择匹配`givenName`和`mail`的`objectClass`属性中的`person`元素

请注意，常见的误区是LDAP数据**区分大小写**，而实际上匹配规则和排序规则来匹配或比较与相对值的关系。

如果要filter来匹配属性值的大小写，则必须使用可扩展的匹配过滤器，例如：

```go
(＆(objectClass = person)(|(givenName:caseExactMatch:=John)(mail:caseExactSubstringsMatch:=john*)))
```

#### 0.4.4.4. derefAliases

是否以及如何遵循别名条目（引用其他条目的条目）

#### 0.4.4.5. attributes

在结果条目中返回哪些属性。

#### 0.4.4.6. sizeLimit, timeLimit

返回的最大条目数，以及允许搜索运行的最长时间。

> 注意，这些值不能覆盖服务器对大小限制和时间限制的设定值。

#### 0.4.4.7. typesOnly

仅返回属性类型，而不返回属性值。

服务器返回匹配的条目和可能的延续引用。这些可以以任意顺序返回。最终结果将包括结果代码。

Compare操作采用DN、属性名称和属性值，并检查命名条目是否包含具有该值的属性。

### 0.4.5. Modify

LDAP客户端使用MODIFY操作来请求LDAP服务器对现有条目进行更改。尝试修改不存在的条目将失败。修改请求受服务器使用的访问控制的约束。

MODIFY操作要求指定条目的DN，并进行一系列更改。序列中的每个更改必须是以下之一：

- add：添加一个新值，该值必须不存在于属性中
- delete：删除现有值
- replace：用新值替换现有值

如下示例，向属性添加值的LDIF：

```go
dn: dc=example,dc=com
changetype: modify
add: cn
cn: the-new-cn-value-to-be-added
-
```

要替换现有属性的值，请使用`replace`关键字。如果属性是多值的，则客户端必须指定要更新的属性的值。

要从条目中删除属性，请使用`delete`关键字和`changetype`指示符`Modify`。如果属性是多值的，则客户端必须指定要删除的属性的值。

还有一个`Modify-Increment`扩展，它允许将可递增的属性值增加指定的数量。

如下示例，使用LDIF将`employeeNumber`递增5：

```go
dn: uid=user.0,ou=people,dc=example,dc=com
changetype: modify
increment: employeeNumber
employeeNumber: 5
-
```

当LDAP服务器处于复制拓扑中时，LDAP客户端应考虑使用`post-read control`来验证更新，而不是在更新后进行搜索。

`post-read control`的设计使应用程序无需在更新后发出搜索请求，这是一种糟糕的形式（因为最终一致性模型），仅仅为了检查一个新的条目在更新后是否生效。

LDAP客户端不应假定每个请求都连接到同一目录服务器，因为可能在LDAP客户端和服务器之间存在负载平衡器或LDAP代理。

### 0.4.6. Modify DN

Modify DN（移动/重命名条目）采用新的RDN，还可以选择新的父级DN，以及一个标志（该标志指示是否删除条目中与旧RDN匹配的值）。 服务器可能支持整个目录子树的重命名。

更新操作是原子性的，而其他的操作将看到新条目或旧条目。

另一方面，LDAP没有定义多个操作的事务：如果读取一个条目然后对其进行修改，则另一个客户端可能同时已更新了该条目。

服务器可以实现支持此功能的扩展。

### 0.4.7. Extended operations

扩展操作是一种通用的LDAP操作，可以定义不属于原始协议规范的新操作。

`StartTLS`是最重要的扩展之一。

其他示例包括`Cancel`和`Password Modify`。

#### 0.4.7.1. StartTLS

`StartTLS`操作在连接上建立TLS，可以提供：

- 数据机密性（以防止第三方查看数据）
- 数据完整性保护（以防止数据被篡改）

在TLS协商期间，服务器发送其`X.509`证书以证明其身份。客户也可以发送证书以证明其身份。

然后，客户端可以使用`SASL/EXTERNAL`。 通过使用`SASL/EXTERNAL`，客户端请求服务器从较低级别提供的凭据（例如TLS）中派生其身份。

尽管从技术上讲，服务器可以使用在任何较低级别建立的任何身份信息，但是通常服务器将使用TLS建立的身份信息。

服务器还通常在单独的端口上支持非标准的`LDAPS`(Secure LDAP或LDAP over SSL)协议，默认情况下为`636`。

LDAPS与LDAP有两种不同：

1. 连接时，客户端和服务器建立TLS，然后再传输任何LDAP消息（不执行StartTLS操作）
2. TLS关闭后必须关闭LDAPS连接

某些`LDAPS`客户端库仅对通信进行加密，并不会根据提供的证书中的名称检查主机名。

### 0.4.8. Abandon

Abandon操作请求服务器中止由消息ID命名的操作。服务器不需要执行该请求。Abandon或成功的Abandon操作都不会返回响应。

类似的`Cancel`扩展操作会返回响应，但并非所有实现都支持此操作。

### 0.4.9. Unbind

Unbind操作将放弃所有未完成的操作并关闭连接，不返回响应。该操作是历史遗留的，并不是Bind操作的逆操作。

客户端可以通过简单地关闭连接来中止会话，但更合适的操作是应该使用Unbind操作。

Unbind允许服务器正常关闭连接并释放资源，否则该资源将保留一段时间，直到发现客户端放弃连接为止。它还指示服务器取消可以取消的操作，并且不发送对不能取消的操作的响应。

## 0.5. URI scheme

LADP的统一资源标识符方案已存在，客户端在不同程度上都支持该方案，服务器返回引用或延续引用，参考[RFC 4516](https://tools.ietf.org/html/rfc4516)。

```go
ldap://host:port/DN?attributes?scope?filter?extensions
```

以下描述的大多数部分都是可选的。

- host：是要搜索的LDAP服务器的FQDN或IP地址
- port：是LDAP服务器的网络端口（默认端口`389`）
- DN：是用作搜索基础的专有名称
- attributes：是要用逗号分隔的属性列表
- scope：指定搜索范围，可以是`base`（默认），`one`或`sub`
- filter：是一个搜索过滤器。例如`(objectClass=*)`
- extensions：是LDAP URL格式的扩展名。

例如：

```go
"ldap://ldap.example.com/cn=John%20Doe,dc=example,dc=com"
```

指向`ldap.example.com`中`John Doe`条目中的所有用户属性

```go
"ldap:///dc=example,dc=com??sub?(givenName=John)"
```

搜索默认服务器中的条目。

> 注意：三元斜杠表示省略主机，而双问号表示省略属性

**与其他URL一样，特殊字符必须进行百分比编码**。

对于基于SSL的LDAP，存在类似的非标准ldaps URI方案。这不应与带有TLS的LDAP混淆，后者是通过使用标准ldap方案的StartTLS操作来实现的。

## 0.6. Schema

子树中条目的内容由directory schema控制，一组与目录信息树（DIT）的结构有关的定义和约束所控制。

Directory Server的模式定义了一组规则，用于管理服务器可以保存的信息种类。它包含许多元素，包括：

- Attribute Syntaxes：提供有关可以存储在属性中的信息种类的信息
- Matching Rules：提供有关如何与属性值进行比较的信息
- Matching Rule Uses：指出哪些属性类型可以与特定匹配规则结合使用
- Attribute Types：定义对象标识符（OID）和一组可能引用给定属性的名称，并将该属性与语法和一组匹配规则相关联
- Object Classes：定义命名的属性集合，并将它们分类为必需和可选属性的集合
- Name Forms：为应包含在条目的RDN中的属性集定义规则
- Content Rules：定义有关可以与条目结合使用的对象类和属性的其他约束
- Structure Rule：定义规则以管理给定条目可能具有的从属条目的类型

属性是负责将信息存储在目录中的元素，Schema定义了可以在条目中使用属性的规则，这些属性可能具有的值的种类以及客户端如何与这些值进行交互。

客户端可以通过检索适当的子schema或子条目来了解服务器支持的schema元素。

该schema定义对象类。每个条目必须具有一个`objectClass`属性，其中包含在schema中定义的命名类。条目的类别的schema定义了该条目可以代表哪种对象-例如 个人，组织或领域。对象类定义还定义了必须包含值的属性列表和可能包含值的属性列表。
