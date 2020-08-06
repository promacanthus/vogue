---
title: "23 Interface重要性"
date: 2020-08-06T18:36:22+08:00
draft: true
---

## 场景与接口定义

在线商城，需要在Go后台提供存储与查询产品的服务，既实现一个负责保存和检索产品的存储库。

```bash
productrepo/
└── api.go

0 directories, 1 file
```

创建一个`productrepo`包和一个`api.go`文件，该API应该暴露出存储库里所有的产品方法。

```golang
// api.go
package productrepo

type ProductRepository interface {
 StoreProduct(name string, id int)
 FindProductByID(id int)
}

```

在`productrepo`包下，定义了`ProductRepository`接口，它代表的就是存储库。该接口中定义两个方法：

1. `StoreProduct()`方法用于存储产品信息
2. `FindProductByID()`方法通过产品ID查找产品信息

## 接口实现

存储库接口定义完成后，就需要有实体对象去实现该接口。

```bash
productrepo/
├── api.go
└── mock.go

0 directories, 2 files
```

在`productrepo`包下，新建`mock.go`文件，定义`mockProductRepo`对象。

```golang
// mock.go
package productrepo

import "fmt"

// 实现接口的空实体对象
type mockProductRepo struct {}

func (m mockProductRepo) StoreProduct(name string, id int) {
 fmt.Println("mocking the StoreProduct func")
}

func (m mockProductRepo) FindProductByID(id int) {
 fmt.Println("mocking the FindProductByID func")
}
```

在示例代码中mock出`ProductRepository`接口所需的方法。

然后在`api.go`中增加`New()`方法，它返回的一个实现了`ProductRepository`接口的对象。

```golang
// api.go
func New() ProductRepository {
 return mockProductRepo{}
}
```

## 为什么使用接口

对于已经定义的`ProductRepository`接口，可以有多种对象去实现它。

- 对于小型的个人项目来说可以不要接口，直接创建一个实体对象。
- 在复杂的实际应用项目中，通常会有很多种存储对象：
  - 使用本地MySQL存储，
  - 连接到云数据库(例如阿里云、谷歌云和腾讯云等)存储。

不同的数据库存储都需要实现`ProductRepository`接口定义的`StoreProduct()`方法和`FindProductByID()`方法。

以本地MySQL存储库为例，它要管理产品对象，需要实现`ProductRepository`接口。

```bash
productrepo/
├── api.go
├── mock.go
└── mysql.go

0 directories, 3 files
```

在`productrepo`包下，新建`mysql.go`文件，定义了`mysqlProductRepo`对象并实现接口方法。

```golang
// mysql.go
package productrepo

import "fmt"

type mysqlProductRepo struct {
}

func (m mysqlProductRepo) StoreProduct(name string, id int) {
 fmt.Println("mysqlProductRepo: mocking the StoreProduct func")
 // In a real world project you would query a MySQL database here.
}

func (m mysqlProductRepo) FindProductByID(id int) {
 fmt.Println("mysqlProductRepo: mocking the FindProductByID func")
 // In a real world project you would query a MySQL database here.
}
```

相似地，当项目中同时需要把产品信息存储到云端时，以阿里云为例。

```bash
productrepo/
├── aliyun.go
├── api.go
├── mock.go
└── mysql.go

0 directories, 4 files
```

在`productrepo`包下，新建`aliyun.go`文件，定义了`aliCloudProductRepo`对象并实现接口方法。

```golang
// aliyun.go
package productrepo

import "fmt"

type aliCloudProductRepo struct {

}

func (m aliCloudProductRepo) StoreProduct(name string, id int) {
 fmt.Println("aliCloudProductRepo: mocking the StoreProduct func")
 // In a real world project you would query an ali Cloud database here.
}

func (m aliCloudProductRepo) FindProductByID(id int) {
 fmt.Println("aliCloudProductRepo: mocking the FindProductByID func")
 // In a real world project you would query an ali Cloud database here.
}
```

此时，更新`api.go`中定义的`New()`方法。

```golang
// api.go
func New(environment string) ProductRepository {
 switch environment {
 case "aliCloud":
  return aliCloudProductRepo{}
 case "local-mysql":
  return mysqlProductRepo{}
 }
 return mockProductRepo{}
}
```

通过将环境变量`environment`传递给`New()`函数，它将基于该环境值返回`ProductRepository`接口的正确实现对象。

```bash
.
├── go.mod
├── main.go
└── productrepo
    ├── aliyun.go
    ├── api.go
    ├── mock.go
    └── mysql.go

1 directory, 6 files
```

定义程序入口`main.go`文件以及main函数。

```golang
// main.go
package main

import "productrepo"

func main() {
 env := "aliCloud"
 repo := productrepo.New(env)
 repo.StoreProduct("HuaWei mate 40", 105)
}
```

通过使用`productrepo.New()`方法基于环境值来获取`ProductRepository`接口对象。

如果需要切换产品存储库，则只需要使用对应的`env`值调用`productrepo.New()`方法即可。

## 如果不使用接口

### 需要为每个对象增加初始化方法

```golang
// mysql.go
func NewMysqlProductRepo() *mysqlProductRepo {
 return &mysqlProductRepo{}
}

//aliyun.go
func NewAliCloudProductRepo()  *aliCloudProductRepo{
 return &aliCloudProductRepo{}
}

// mock.go
func NewMockProductRepo() *mockProductRepo {
 return &mockProductRepo{}
}
```

### 调用对象处产生大量重复代码

```golang
// main.go
package main

import "productrepo"

func main() {
 env := "aliCloud"
 switch env {
 case "aliCloud":
  repo := productrepo.NewAliCloudProductRepo()
  repo.StoreProduct("HuaWei mate 40", 105)
    // the more function to do, the more code is repeated.
 case "local-mysql":
  repo := productrepo.NewMysqlProductRepo()
  repo.StoreProduct("HuaWei mate 40", 105)
    // the more function to do, the more code is repeated.
 default:
  repo := productrepo.NewMockProductRepo()
  repo.StoreProduct("HuaWei mate 40", 105)
    // the more function to do, the more code is repeated.
 }
}
```

在项目演进过程中，会迭代很多存储库对象，而通过`ProductRepository`接口，可以轻松地实现扩展，而不必反复编写相同逻辑的代码。

## 总结

开发中，常常提到要功能模块化，上面示例：通过接口为载体，一类服务就是一个接口，实现接口即服务。
