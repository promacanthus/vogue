---
title: 00-Endpoints(REST)
date: 2020-04-14T10:09:14.254627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 开发框架
- Go-kit
- Examples
summary: 00-Endpoints(REST)
showInMenu: false

---

查看[原文](https://www.ru-rocker.com/2017/02/17/micro-services-using-go-kit-http-endpoint/)。

在本文中，将使用Go编程语言和go-kit作为标准微服务库创建一个简单的微服务。该服务将通过REST端点公开。

## 微服务

在当今世界，微服务体系结构已得到普及。 我不会特别解释什么是微服务架构，因为互联网上已经讨论了很多。 但是，我将提供两个有关微服务的良好网站。

- 首先是我最喜欢的[martinfowler.com](https://martinfowler.com/)，可以在[这里](https://martinfowler.com/articles/microservices.html)看到他的惊人解释。

- 另一个来自[microservices.io](https://microservices.io/)，里面有很多关于模式和示例的好文章。

### Go-lang

Go是一种开放源代码的编程语言，可轻松构建简单，可靠和高效的软件。该语言由Google设计，旨在解决Google的问题。因此，人们可以希望这种语言能够大规模运行，适用于大型程序和依赖项。

### Go-kit

Go-kit确实有助于简化构建微服务架构的过程。 这是因为它具有许多功能，例如服务连接性，指标和日志记录。 因此，我要特别感谢Peter Bourgon（@[peterbourgon](https://peter.bourgon.org/)）和所有提供此出色库的贡献者。

## 用例

我在本文中的重点是将REST端点公开为服务。服务本身将使用HTTP POST方法公开端点。然后它将以JSON格式返回lorem ipsum消息。

端点URL格式为`/lorem/{type}/{min}/{max}`，其中包含以下说明：

- type将是lorem类型，即word（单词），sentence（句子）和paragraph（段落）
- min和max表示生成器生成的最少和最多字母数

## 一步步操作

在逐步开始之前，此示例需要几个库。那些是：

1. [go-kit libraries](https://godoc.org/github.com/go-kit/kit)
2. [golorem libraries](https://godoc.org/github.com/drhodes/golorem)：用于生成lorem ipsum文本
3. [gorilla mux libraries](https://godoc.org/github.com/gorilla/mux)： 用于http处理程序

### 第一步：创建Service

```shell
go get github.com/go-kit/kit
go get github.com/drhodes/golorem
go get github.com/gorilla/mux
```

无论使用哪种最酷的工具，都需要在第一时间创建业务逻辑。如用例所述，我们的业务逻辑是根据单词，句子或段落创建lorem ipsum文本。因此，让我们在工作区下创建lorem文件夹。就我而言，我的文件夹是`$GOPATH/github.com/ru-rocker/gokit-playground/lorem`。然后在该文件夹下创建文件`service.go`并添加以下代码：

```go
// Define service interface
type Service interface {
    Word(min, max int) string               // generate a word with at least min letters and at most max letters.
    Sentence(min, max int) string         // generate a sentence with at least min words and at most max words.
    Paragraph(min, max int) string       // generate a paragraph with at least min sentences and at most max sentences.
}

// Implement service with empty struct
type LoremService struct {}
```

现在我们已经有了接口，但是，没有方法实现的接口将毫无意义。为service.go，添加以下实现：

```go
// Implement service functions
func (LoremService) Word(min, max int) string {
    return golorem.Word(min, max)
}

func (LoremService) Sentence(min, max int) string {
    return golorem.Sentence(min, max)
}

func (LoremService) Paragraph(min, max int) string {
    return golorem.Paragraph(min, max)
}
```

注意，对每个方法实现都使用了golorem函数。

### 第2步：建模Request和Response

由于此服务是HTTP的一部分，因此下一步是对请求和响应进行建模。如果回头看一下用例，将发现实现需要三个属性，分别是：type、min、max。

对于响应本身，只需要两个字段。

- 这些字段是包含lorem ipsum文本的消息
- 另一个是错误字段，当出现错误时，它将给出错误描述

因此，让我们创建另一个文件，并为其命名为`endpoints.go`，并添加以下代码：

```go
//request
type LoremRequest struct {
    RequestType string
    Min int
    Max int
}

//response
type LoremResponse struct {
    Message string `json:"message"`
    Err     error `json:"err,omitempty"` //omitempty 表示，如果值为nil，则不会显示此字段
}
```

### 第3步：创建端点

端点是go-kit中的特殊功能，可以将它包装到`http.Handler`中。为了使我们的service变成`endpoint.Endpoint`函数，我们将要制作一个函数来处理LoremRequest，在函数内部做一些逻辑，然后返回LoremResponse。对于endpoints.go，添加如下代码：

```go
var (
    ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

// endpoints wrapper
type Endpoints struct {
    LoremEndpoint endpoint.Endpoint
}

// creating Lorem Ipsum Endpoint
func MakeLoremEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(LoremRequest)

        var (
            txt string
            min, max int
        )

        min = req.Min
        max = req.Max

        if strings.EqualFold(req.RequestType, "Word") {
            txt = svc.Word(min, max)
        } else if strings.EqualFold(req.RequestType, "Sentence"){
            txt = svc.Sentence(min, max)
        } else if strings.EqualFold(req.RequestType, "Paragraph") {
            txt = svc.Paragraph(min, max)
        } else {
            return nil, ErrRequestTypeNotFound
        }

        return LoremResponse{Message: txt}, nil
    }
}
```

### 第4步：传输

在处理`http`请求和响应之前，首先需要创建从`struct`到`json`或相反的编码器和解码器。为此，创建一个新文件，命名为`transport.go`并添加以下代码：

```go
// decode url path variables into request
func decodeLoremRequest(_ context.Context, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    requestType, ok := vars["type"]
    if !ok {
        return nil, ErrBadRouting
    }

    vmin, ok := vars["min"]
    if !ok {
        return nil, ErrBadRouting
    }

    vmax, ok := vars["max"]
    if !ok {
        return nil, ErrBadRouting
    }

    min, _ := strconv.Atoi(vmin)
    max, _ := strconv.Atoi(vmax)
    return LoremRequest{
        RequestType: requestType,
        Min: min,
        Max: max,
    }, nil
}

// errorer由可能包含错误的所有具体响应类型实现
// 它使我们能够更改HTTP响应代码，而无需触发端点（传输级）错误
type errorer interface {
    error() error
}

// encodeResponse是对客户端的所有响应类型进行编码的常用方法
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    if e, ok := response.(errorer); ok && e.error() != nil {
        // 不是Go kit传输错误，而是业务逻辑错误
        // 提供这些作为HTTP错误
        encodeError(ctx, e.error(), w)
        return nil
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    return json.NewEncoder(w).Encode(response)
}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
    if err == nil {
        panic("encodeError with nil error")
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": err.Error(),
    })
}
```

一旦声明了编码器和解码器功能，就可以创建`http`处理程序了。在`transport.go`中添加以下代码：

```go
var (
    // 如果缺少预期的路径变量，则返回ErrBadRouting
    ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// Make Http Handler
func MakeHttpHandler(ctx context.Context, endpoint Endpoints, logger log.Logger) http.Handler {
    r := mux.NewRouter()
    options := []httptransport.ServerOption{
        httptransport.ServerErrorLogger(logger),
        httptransport.ServerErrorEncoder(encodeError),
    }

    //POST /lorem/{type}/{min}/{max}
    // 注意：请看这行。这是描述URL路径和HTTP请求方法的方法
    r.Methods("POST").Path("/lorem/{type}/{min}/{max}").Handler(httptransport.NewServer(
        ctx,
        endpoint.LoremEndpoint,
        decodeLoremRequest,
        encodeResponse,
        options...,
    ))
    return r
}
```

### 第5步：main

到目前为止，已经为服务提供了服务/业务层，端点和传输。设置完成后，就该创建main函数了。main函数基本上是构造端点并使其可用于HTTP传输。

因此，在lorem文件夹下，创建另一个文件夹，名为`lorem.d`（点d表示守护程序）。可以随意命名，然后创建文件`main.go`，并添加以下代码：

```go
func main() {
    ctx := context.Background()
    errChan := make(chan error)

    var svc lorem.Service
    svc = lorem.LoremService{}
    endpoint := lorem.Endpoints{
        LoremEndpoint: lorem.MakeLoremEndpoint(svc),
    }

    // Logging domain.
    var logger log.Logger
    {
        logger = log.NewLogfmtLogger(os.Stderr)
        logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
        logger = log.NewContext(logger).With("caller", log.DefaultCaller)
    }

    r := lorem.MakeHttpHandler(ctx, endpoint, logger)

    // HTTP transport
    go func() {
        fmt.Println("Starting server at port 8080")
        handler := r
        errChan <- http.ListenAndServe(":8080", handler)
    }()

    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errChan <- fmt.Errorf("%s", <-c)
    }()
    fmt.Println(<- errChan)
}
```

### 第6步：运行样例

现在该运行示例了。打开shell输入如下命令：

```shell
cd $GOPATH
go run src/github.com/ru-rocker/gokit-playground/lorem/lorem.d/main.go
```

使用curl或postman测试端点。

## 样例代码

每当您对本文感兴趣并愿意了解更多信息时，都可以在我的[github](https://github.com/ru-rocker/gokit-playground)上进行检查。
