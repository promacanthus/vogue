---
title: "Swagger Guide"
date: 2020-04-17T18:06:14+08:00
draft: true
---

## 前言

下面是需要用到的工具和库。

|名称|描述|
|---|---|
|[swag](https://github.com/swaggo/swag)|使用Swagger 2.0为Go语言自动生成RESTful API文档|
|[gin-swagger](https://github.com/swaggo/gin-swagger)|gin中间件以使用Swagger 2.0自动生成RESTful API文档|

swag是Golang的工具，将代码注释转换为Swagger2.0文档。除了swag，还需一个web框架的中间件包装器库，当前，swag支持的web框架包括[gin](https://github.com/swaggo/gin-swagger)、[echo](https://github.com/swaggo/echo-swagger)、[buffalo](https://github.com/swaggo/buffalo-swagger)、[net/http](https://github.com/swaggo/http-swagger)，对于不同的框架来说，注释都是一样的。

在这里使用的是gin-swagger，另一个Golang Swagger库[go-swagger](https://goswagger.io)，它似乎更受欢迎并且功能也更强大。

- gin-swagger：简单易用
- go-swagger：对生成的内容进行更多的控制

swag工具的[README](https://github.com/swaggo/swag/blob/master/README.md)很详细。

下面把它翻译成中文，便于日后查看。

## Swag

Swag将Go的注释转换为Swagger2.0文档。我们为流行的 [Go Web Framework](#支持的Web框架) 创建了各种插件，这样可以与现有Go项目快速集成（使用Swagger UI）。

### 目录


## 快速开始

1. 将注释添加到API源代码中，请参阅声明性注释格式。
2. 使用如下命令下载swag：

```bash
go get -u github.com/swaggo/swag/cmd/swag
```

从源码开始构建的话，需要有Go环境（1.9及以上版本）。

或者从github的release页面下载预编译好的二进制文件。

3. 在包含`main.go`文件的项目根目录运行`swag init`。这将会解析注释并生成需要的文件（`docs`文件夹和`docs/docs.go`）。

```bash
swag init
```

确保导入了生成的`docs/docs.go`文件，这样特定的配置文件才会被初始化。如果通用API指数没有写在`main.go`中，可以使用`-g`标识符来告知swag。

```bash
swag init -g http/api.go
```

### swag cli

```bash
swag init -h
NAME:
   swag init - Create docs.go

USAGE:
   swag init [command options] [arguments...]

OPTIONS:
   --generalInfo value, -g value       Go file path in which 'swagger general API Info' is written (default: "main.go")
   --dir value, -d value               Directory you want to parse (default: "./")
   --propertyStrategy value, -p value  Property Naming Strategy like snakecase,camelcase,pascalcase (default: "camelcase")
   --output value, -o value            Output directory for all the generated files(swagger.json, swagger.yaml and doc.go) (default: "./docs")
   --parseVendor                       Parse go files in 'vendor' folder, disabled by default
   --parseDependency                   Parse go files in outside dependency folder, disabled by default
   --markdownFiles value, --md value   Parse folder containing markdown files to use as description, disabled by default
   --generatedTime                     Generate timestamp at the top of docs.go, true by default
```

## 支持的Web框架

- [gin](http://github.com/swaggo/gin-swagger)
- [echo](http://github.com/swaggo/echo-swagger)
- [buffalo](https://github.com/swaggo/buffalo-swagger)
- [net/http](https://github.com/swaggo/http-swagger)

## 如何与Gin集成

[点击此处](https://github.com/swaggo/swag/tree/master/example/celler)查看示例源代码。

1. 使用`swag init`生成Swagger2.0文档后，导入如下代码包：

```go
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```

2. 在`main.go`源代码中添加通用的API注释：

```bash
// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func main() {
    r := gin.Default()

    c := controller.NewController()

    v1 := r.Group("/api/v1")
    {
        accounts := v1.Group("/accounts")
        {
            accounts.GET(":id", c.ShowAccount)
            accounts.GET("", c.ListAccounts)
            accounts.POST("", c.AddAccount)
            accounts.DELETE(":id", c.DeleteAccount)
            accounts.PATCH(":id", c.UpdateAccount)
            accounts.POST(":id/images", c.UploadAccountImage)
        }
    //...
    }
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    r.Run(":8080")
}
//...
```

此外，可以动态设置一些通用的API信息。生成的代码包`docs`到处`SwaggerInfo`变量，使用该变量可以通过编码的方式设置标题、描述、版本、主机和基础路径。使用Gin的示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"

    "./docs" // docs is generated by Swag CLI, you have to import it.
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/

func main() {

    // programatically set swagger info
    docs.SwaggerInfo.Title = "Swagger Example API"
    docs.SwaggerInfo.Description = "This is a sample server Petstore server."
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "petstore.swagger.io"
    docs.SwaggerInfo.BasePath = "/v2"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

    r := gin.New()

    // use ginSwagger middleware to serve the API docs
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run()
}
```

3. 在`controller`代码中添加API操作注释：

```go
package controller

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/swag/example/celler/httputil"
    "github.com/swaggo/swag/example/celler/model"
)

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
    id := ctx.Param("id")
    aid, err := strconv.Atoi(id)
    if err != nil {
        httputil.NewError(ctx, http.StatusBadRequest, err)
        return
    }
    account, err := model.AccountOne(aid)
    if err != nil {
        httputil.NewError(ctx, http.StatusNotFound, err)
        return
    }
    ctx.JSON(http.StatusOK, account)
}

// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts [get]
func (c *Controller) ListAccounts(ctx *gin.Context) {
    q := ctx.Request.URL.Query().Get("q")
    accounts, err := model.AccountsAll(q)
    if err != nil {
        httputil.NewError(ctx, http.StatusNotFound, err)
        return
    }
    ctx.JSON(http.StatusOK, accounts)
}

//...
```

```bash
swag init
```

4. 运行程序，然后在浏览器中访问 http://localhost:8080/swagger/index.html。将看到Swagger 2.0 Api文档，如下所示：

![swagger_index.html](https://raw.githubusercontent.com/swaggo/swag/master/assets/swagger-image.png)

## 开发现状

[Swagger 2.0 文档](https://swagger.io/docs/specification/2-0/basic-structure/)

- [x] Basic Structure
- [x] API Host and Base Path
- [x] Paths and Operations
- [x] Describing Parameters
- [x] Describing Request Body
- [x] Describing Responses
- [x] MIME Types
- [x] Authentication
  - [x] Basic Authentication
  - [x] API Keys
- [x] Adding Examples
- [x] File Upload
- [x] Enums
- [x] Grouping Operations With Tags
- [ ] Swagger Extensions

## 声明式注释格式

### 通用API信息

**示例** [`celler/main.go`](https://github.com/swaggo/swag/blob/master/example/celler/main.go)

| 注释  | 说明                                | 示例                         |
|-------------|--------------------------------------------|---------------------------------|
| title       | **必填** 应用程序的名称。| // @title Swagger Example API   |
| version     | **必填** 提供应用程序API的版本。| // @version 1.0  |
| description | 应用程序的简短描述。|// @description This is a sample server celler server. |
| tag.name    | 标签的名称。| // @tag.name This is the name of the tag                     |
| tag.description   | 标签的描述。| // @tag.description Cool Description         |
| tag.docs.url      | 标签的外部文档的URL。| // @tag.docs.url https://example.com|
| tag.docs.description  | 标签的外部文档说明。| // @tag.docs.description Best example documentation |
| termsOfService | API的服务条款。| // @termsOfService http://swagger.io/terms/                     |
| contact.name | 公开的API的联系信息。| // @contact.name API Support  |
| contact.url  | 联系信息的URL。 必须采用网址格式。| // @contact.url http://www.swagger.io/support|
| contact.email| 联系人/组织的电子邮件地址。 必须采用电子邮件地址的格式。| // @contact.email support@swagger.io                                   |
| license.name | **必填** 用于API的许可证名称。|// @license.name Apache 2.0|
| license.url  | 用于API的许可证的URL。 必须采用网址格式。| // @license.url http://www.apache.org/licenses/LICENSE-2.0.html |
| host        | 运行API的主机（主机名或IP地址）。     | // @host localhost:8080         |
| BasePath    | 运行API的基本路径。 | // @BasePath /api/v1             |
| query.collection.format | 查询或枚举中的默认集合（数组）参数格式：csv，multi，pipes，tsv，ssv。 如果未设置，则默认为csv。| // @query.collection.format multi
| schemes     | 用空格分隔的请求的传输协议。 | // @schemes http https |
| x-name      | 扩展的键必须以x-开头，并且只能使用json值 | // @x-example-key {"key": "value"} |

### 使用Markdown描述

如果文档中的短字符串不足以完整表达，或者需要展示图片，代码示例等类似的内容，则可能需要使用Markdown描述。要使用Markdown描述，请使用一下注释。

| 注释  | 说明                                | 示例                         |
|-------------|--------------------------------------------|---------------------------------|
| title       | **必填** 应用程序的名称。| // @title Swagger Example API   |
| version     | **必填** 提供应用程序API的版本。| // @version 1.0  |
| description.markdown  | 应用程序的简短描述。 从`api.md`文件中解析。 这是`@description`的替代用法。|// @description.markdown No value needed, this parses the description from api.md|
| tag.name    | 标签的名称。| // @tag.name This is the name of the tag                     |
| tag.description.markdown   | 标签说明，这是`tag.description`的替代用法。 该描述将从名为`tagname.md的`文件中读取。  | // @tag.description.markdown         |

