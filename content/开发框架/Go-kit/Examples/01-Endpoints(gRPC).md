---
title: 01-Endpoints(gRPC)
date: 2020-04-14T10:09:14.254627+08:00
draft: false
---

查看[原文](https://www.ru-rocker.com/2017/02/24/micro-services-using-go-kit-grpc-endpoint/)。

本文将展示如何使用Golang编程语言和Go-kit作为其框架来创建微服务。该服务将通过gRPC协议公开。

## 概览

在上一篇文章中，我谈到了使用go编程语言创建微服务。另外，我介绍了go-kit作为微服务框架。如您所见，两者在构建服务方面都非常出色。就像以前一样，在本文中，将创建另一个端点作为服务，但是，将使用gRPC协议作为通信接口，而不是通过REST暴露服务。

### gRPC

我将简要介绍gRPC，因为[grpc.io](http://www.grpc.io/)已经对此进行了很好的解释。gRPC是Google提供的框架，用于支持远程过程调用。通过使用gRPC，客户端可以直接调用运行在其他计算机上的服务的方法。此外，默认情况下，gRPC使用`protocol buffers`作为序列化结构数据的机制。与XML或JSON相比，此序列化更小且编码和解码更快。

> 如果想要了解`protocol buffers`性能与JSON相比有多出色，可以访问此[链接](https://auth0.com/blog/beating-json-performance-with-protobuf/)。

## 用例

用例仍然与上一篇文章相同。将创建一个服务来生成“ lorem ipsum”文本。但是这次会有所不同，不会在一项服务中创建三个功能（单词，句子和段落），而是创建一个称为Lorem的功能。然后从该函数中，无论是生成单词，句子还是段落，都将分派请求类型。

## 一步步操作

假设我们已经从上一篇文章中获得了所需的库。但这还不够，我们需要安装其他库和`protocol buffers`。

1. 从[这里](https://github.com/google/protobuf/releases)下载`protocol buffers`的编译器，提取bin文件夹并将其导出到`$PATH`中。

```bash
#PROTOC
export PROTOC_HOME=~/opt/protoc-3.2.0-osx-x86_64
export PATH=${PATH}:$PROTOC_HOME/bin/
```

2. 在命令行执行如下操作：

```bash
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

3. 将GOBIN也要添加到PATH中：

```bash
# GOPATH
export GOPATH=~/workspace-golang/
export GOBIN=$GOPATH/bin/
export PATH=${PATH}:$GOBIN
```

然后创建新文件夹，为其命名：`lorem-grpc`。就我而言，我是在`$GOPATH/github.com/ru-rocker/gokit-playground/lorem-grpc`下完成的。我将其称为WORKDIR。

### 第1步：Proto

使用`protocol buffers`第三版的语言（即`proto3`）来定义数据结构，首先在WORKDIR中创建名叫pb的文件夹，然后在其中创建一个proto文件。在这个proto文件中指定需要使用的语法是`proto3`，此外，还将定义服务和请求响应对。

在本文中，服务名称是Lorem，以LoremRequest作为输入参数，并返回LoremResponse。

```protobuf
syntax = "proto3";
package pb;

service Lorem {
    rpc Lorem(LoremRequest) returns (LoremResponse) ;
}

message LoremRequest {
    string requestType = 1;
    int32 min = 2;
    int32 max = 3;
}

message LoremResponse {
    string message = 1;
    string err = 2;
}
```

消息内的字段具有以下格式：

```protobuf
type name = index
```

- type：描述数据类型属性，它可以是string，bool，double，float，int32，int64等。
- index：是数据流的整数值，用于指示字段位置

基于proto文件使用`protocol buffers`编译器和用于`protobuf`的Go插件生成Go文件。在pb文件夹下，执行如下命令：

```bash
protoc lorem.proto --go_out=plugins=grpc:.
```

这将会创建名为`lorem.pb.go`的文件。

### 第2步：定义服务

在WORKDIR下创建文件service.go，然后在其中创建Lorem函数。

```go
package lorem_grpc

import (
    gl "github.com/drhodes/golorem"
    "strings"
    "errors"
    "context"
)

var (
    ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

// Define service interface
type Service interface {
    // generate a word with at least min letters and at most max letters.
    Lorem(ctx context.Context, requestType string, min, max int) (string, error)
}

// Implement service with empty struct
type LoremService struct {}

// Implement service functions
func (LoremService) Lorem(_ context.Context, requestType string, min, max int) (string, error) {
    var result string
    var err error
    if strings.EqualFold(requestType, "Word") {
        result = gl.Word(min, max)
    } else if strings.EqualFold(requestType, "Sentence") {
        result = gl.Sentence(min, max)
    } else if strings.EqualFold(requestType, "Paragraph") {
        result = gl.Paragraph(min, max)
    } else {
        err = ErrRequestTypeNotFound
    }
    return result, err
}
```

### 第3步：创建端点

需要注意的一件事是，使用Endpoints结构实现Service接口，因为，创建gRPC客户端连接时需要此机制。

```go
package lorem_grpc

import (
    "github.com/go-kit/kit/endpoint"
    "context"
    "errors"
)

//request
type LoremRequest struct {
    RequestType string
    Min         int32
    Max         int32
}

//response
type LoremResponse struct {
    Message string `json:"message"`
    Err     string `json:"err,omitempty"`
}

// endpoints wrapper
type Endpoints struct {
    LoremEndpoint endpoint.Endpoint
}

// creating Lorem Ipsum Endpoint
func MakeLoremEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(LoremRequest)

        var (
            min, max int
        )

        min = int(req.Min)
        max = int(req.Max)
        txt, err := svc.Lorem(ctx, req.RequestType, min, max)

        if err != nil {
            return nil, err
        }

        return LoremResponse{Message: txt}, nil
    }
}

// Wrapping Endpoints as a Service implementation.
// Will be used in gRPC client
func (e Endpoints) Lorem(ctx context.Context, requestType string, min, max int) (string, error) {
    req := LoremRequest{
        RequestType: requestType,
        Min: int32(min),
        Max: int32(max),
    }
    resp, err := e.LoremEndpoint(ctx, req)
    if err != nil {
        return "", err
    }
    loremResp := resp.(LoremResponse)
    if loremResp.Err != "" {
        return "", errors.New(loremResp.Err)
    }
    return loremResp.Message, nil
}
```

### 第4步：对请求和响应建模

需要对请求和响应进行编码/解码，为此，创建model.go并定义编码/解码功能。

```go
package lorem_grpc

import (
    "context"
    "github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
)

//Encode and Decode Lorem Request
func EncodeGRPCLoremRequest(_ context.Context, r interface{}) (interface{}, error) {
    req := r.(LoremRequest)
    return &pb.LoremRequest{
        RequestType: req.RequestType,
        Max: req.Max,
        Min: req.Min,
    } , nil
}

func DecodeGRPCLoremRequest(ctx context.Context, r interface{}) (interface{}, error) {
    req := r.(*pb.LoremRequest)
    return LoremRequest{
        RequestType: req.RequestType,
        Max: req.Max,
        Min: req.Min,
    }, nil
}

// Encode and Decode Lorem Response
func EncodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
    resp := r.(LoremResponse)
    return &pb.LoremResponse{
        Message: resp.Message,
        Err: resp.Err,
    }, nil
}

func DecodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
    resp := r.(*pb.LoremResponse)
    return LoremResponse{
        Message: resp.Message,
        Err: resp.Err,
    }, nil
}
```

### 第5步：传输

在此步骤中，需要为grpcServer类型实现LoremServer接口。然后创建函数以返回grpcServer。

```go
package lorem_grpc

import (
    "golang.org/x/net/context"
    grpctransport "github.com/go-kit/kit/transport/grpc"
    "github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
)

type grpcServer struct {
    lorem grpctransport.Handler
}

// implement LoremServer Interface in lorem.pb.go
func (s *grpcServer) Lorem(ctx context.Context, r *pb.LoremRequest) (*pb.LoremResponse, error) {
    _, resp, err := s.lorem.ServeGRPC(ctx, r)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.LoremResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint Endpoints) pb.LoremServer {
    return &grpcServer{
        lorem: grpctransport.NewServer(
            ctx,
            endpoint.LoremEndpoint,
            DecodeGRPCLoremRequest,
            EncodeGRPCLoremResponse,
        ),
    }
}
```

### 第6步：Server

在Go lang中创建gRPC服务器几乎与创建HTTP服务器一样容易。不同的是，我们使用的是tcp协议而不是http。在WORKDIR下，创建文件夹server并创建文件server_grpc_main.go。

```go
package main

import (
    "net"
    "flag"
    "github.com/ru-rocker/gokit-playground/lorem-grpc"
    context "golang.org/x/net/context"
    "google.golang.org/grpc"
    "github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
    "os"
    "os/signal"
    "syscall"
    "fmt"
)

func main() {

    var gRPCAddr = flag.String("grpc", ":8081","gRPC listen address")
    flag.Parse()
    ctx := context.Background()

    // init lorem service
    var svc lorem_grpc.Service
    svc = lorem_grpc.LoremService{}
    errChan := make(chan error)

    // creating Endpoints struct
    endpoints := lorem_grpc.Endpoints{
        LoremEndpoint: lorem_grpc.MakeLoremEndpoint(svc),
    }

    //execute grpc server
    go func() {
        listener, err := net.Listen("tcp", *gRPCAddr)
        if err != nil {
            errChan <- err
            return
        }
        handler := lorem_grpc.NewGRPCServer(ctx, endpoints)
        gRPCServer := grpc.NewServer()
        pb.RegisterLoremServer(gRPCServer, handler)
        errChan <- gRPCServer.Serve(listener)
    }()

    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errChan <- fmt.Errorf("%s", <-c)
    }()

    fmt.Println(<- errChan)
}
```

### 第7步：Client

现在该到客户端了，需要创建返回Service的NewClient函数，因为NewClient函数将转换为Endpoint。

```go
package client

import (
    "github.com/ru-rocker/gokit-playground/lorem-grpc"
    "github.com/ru-rocker/gokit-playground/lorem-grpc/pb"
    grpctransport "github.com/go-kit/kit/transport/grpc"
    "google.golang.org/grpc"
)

// Return new lorem_grpc service
func NewGRPCClient(conn *grpc.ClientConn) lorem_grpc.Service {
    var loremEndpoint = grpctransport.NewClient(
        conn, "Lorem", "Lorem",
        lorem_grpc.EncodeGRPCLoremRequest,
        lorem_grpc.DecodeGRPCLoremResponse,
        pb.LoremResponse{},
    ).Endpoint()

    return lorem_grpc.Endpoints{
        LoremEndpoint:     loremEndpoint,
    }
}
```

请注意，我为此功能返回了Endpoint（请参阅有关实现Service接口的步骤3）。然后在cmd目录下创建可执行客户端。

```go
package main

import (
    "flag"
    "time"
    "log"
    grpcClient "github.com/ru-rocker/gokit-playground/lorem-grpc/client"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "github.com/ru-rocker/gokit-playground/lorem-grpc"
    "fmt"
    "strconv"
)

func main() {
    var (
        grpcAddr = flag.String("addr", ":8081",
            "gRPC address")
    )
    flag.Parse()
    ctx := context.Background()
    conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
        grpc.WithTimeout(1*time.Second))

    if err != nil {
        log.Fatalln("gRPC dial:", err)
    }
    defer conn.Close()

    loremService := grpcClient.New(conn)
    args := flag.Args()
    var cmd string
    cmd, args = pop(args)

    switch cmd {
    case "lorem":
        var requestType, minStr, maxStr string

        requestType, args = pop(args)
        minStr, args = pop(args)
        maxStr, args = pop(args)

        min, _ := strconv.Atoi(minStr)
        max, _ := strconv.Atoi(maxStr)
        lorem(ctx, loremService, requestType, min, max)
    default:
        log.Fatalln("unknown command", cmd)
    }
}

// parse command line argument one by one
func pop(s []string) (string, []string) {
    if len(s) == 0 {
        return "", s
    }
    return s[0], s[1:]
}

// call lorem service
func lorem(ctx context.Context, service lorem_grpc.Service, requestType string, min int, max int) {
    mesg, err := service.Lorem(ctx, requestType, min, max)
    if err != nil {
        log.Fatalln(err.Error())
    }
    fmt.Println(mesg)
}
```

## 运行

为了测试它的工作方式，首先我们需要运行gRPC服务器。

```bash
cd $GOPATH

#Running grpc server
go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/server/src/github.com/ru-rocker/gokit-playground/lorem-grpc/server/server_grpc_main.go
```

接下来，执行客户端：

```bash
cd $GOPATH

#Running client
go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/client/cmd/main.go lorem sentence 10 20
```

输入如下：

```bash
# sentence
go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/client/cmd/client_grpc_main.go lorem sentence 10 20
Concurrunt nota re dicam fias, sim aut pecco, die appetitum.


# word
go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/client/cmd/client_grpc_main.go lorem word 10 20
difficultates

# paragraph
go run src/github.com/ru-rocker/gokit-playground/lorem-grpc/client/cmd/client_grpc_main.go lorem paragraph 10 20
En igitur aequo tibi ita recedimus an aut eum tenacius quae mortalitatis eram aut rapit montium inaequaliter dulcedo aditum. Rerum tempus mala anima volebant dura quae o modis, fama vanescit fit. Nuntii comprehendet ponamus redducet cura sero prout, nonne respondi ipsa angelos comes, da ea saepe didici. Crebro te via hos adsit. Psalmi agam me mea pro da. Audi pati sim da ita praeire nescio faciant. Deserens da contexo e suaveolentiam qualibus subtrahatur excogitanda pusillus grex, e o recorder cor re libidine. Ore siderum ago mei, cura hi deo. Dicens ore curiosarum, filiorum eruuntur munerum displicens ita me repente formaeque agam nosti. Deo fama propterea ab persentiscere nam acceptam sed e a corruptione. Rogo ea nascendo qui, fuit ceterarumque. Drachmam ore operatores exemplo vivunt. Recolo hi fac cor secreta fama, domi, rogo somnis. Sapores fidei maneas saepe corporis re oris quantulum doleam te potu ita lux da facie aut. Benedicendo e tertium nosse agam ne amo, mole invenio dicturus me cognoscere ita aer se memor consulerem ab te rei. Miles ita amaritudo rogo hi flendae quietem invoco quae odor desuper tu. Temptatione dicturus ita mediator ita mundum lux partes miseros percepta seu dicant avaritiam nares contra deseri securus. Ea sobrios tale, rogo sanctis. Ita ne manu uspiam hierusalem, transeam dicite subduntur responsa cor socialiter fit deseri album praeditum.
```

## 总结

创建gRPC端点非常简单，与REST端点只有一点点不同。但是，gRPC比REST运行得更快。我的感觉是将来gRPC将取代JSON。

同时，如果您对此示例感兴趣，可以在我的[github](https://github.com/ru-rocker/gokit-playground)上检查源代码。
