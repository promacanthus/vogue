# 00-DDD in Golang

域驱动设计（DDD）是一种软件开发方法，通过将实现与不断发展的模型相连接，简化了开发人员面临的复杂性。

使用DDD的原因：

1. 提供解决困难问题的原则和模式
2. 将复杂的设计基于领域模型
3. 在技​​术专家和领域专家之间发起创造性的协作，以迭代地完善解决领域问题的概念模型

DDD包含4层：

1. 领域层：这里定义领域模型和应用程序逻辑。
2. 基础层：这里包括独立于应用程序的所有内容：外部库，数据库引擎等。
3. 应用层：这里作为领域层和接口层之间的通道。将请求从接口层发送到领域层，由领域层处理请求并返回响应。
4. 接口层：这里包含与其他系统交互的所有内容，例如Web服务，RMI接口或Web应用程序以及批处理前端。

![DDD](../../images/DDD-in-golang.jpg)

## 项目初始化

构建一个食物推荐API。

首先是初始化依赖管理，在这个项目中将使用`go.mod`，在项目的根目录`food-app/`目录下执行如下命令进行初始化操作：

```go
go mod init food-app
```

如下所示是整个目录的组织结构：

```bash
food-app$ tree -a
.
├── application
├── .circleci
├── domain
├── .env
├── .gitignore
├── go.mod
├── infrastructure
├── interfaces
├── main.go
├── README.md
└── utils

6 directories, 5 files
```

在这个项目中将使用postgres和redis数据库来持久化数据。在`.env`文件中定义连接信息。文件内容如下：

```bash
cat .env

#Postgres
APP_ENV=local
...

#Redis
REDIS_HOST=127.0.0.1
...
```

该文件位于项目的根目录。

## 领域层

首先考虑领域层，在领域中有几个模式，分别是：实体、值对象、聚合、服务等。这个项目只是一个简单的例子，因此只考虑领域中的两种模式：实体和聚合。

```bash
food-app/domain$ tree -a
.
├── repository
│   ├── food_repository.go
│   └── user_repository.go
└── entity
    ├── food.go
    └── user.go

```

### 实体

这是我们定义`Schema`的地方。例如，我们可以定义一个`user`结构体，将该实体视为领域的蓝图。

如上面的`user.go`文件中定义了`user`结构体其中包含用户信息，同时也添加了辅助工具函数，将用于验证和审查输入值。用一个叫`Hash`的方法用于哈希密码。它被定义在`utils`目录下。

使用`[Gorm](http://gorm.io/)`来实现ORM。在定义food实体是也用了同样的方法。

### 聚合

在代码中，聚合使用仓储模式实现，该存储库定义了基础结构实现的方法的集合。这描述了与给定数据库或第三方API交互的方法的数量。

方法都定义在一个接口中，并将在基础层实现。

## 基础层

这一层实现在repository中定义的方法。这些方法与数据库或者第三方API交互。在本文中只有和数据库的交互。

在这里使用`UserRepo`结构体来实现`UserRepository`接口。

然后创建`db.go`文件来配置数据库。在这个文件中定义`Repositories`结构体，用于保存应用中全部的repositories。在这里我们有user和food两个repository。

repository有一个数据库实例，用于传递给user和food的构造函数，即`NewUserRepository`和`NewFoodRepository`。

