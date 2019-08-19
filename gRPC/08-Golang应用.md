# gRPC Basic - Go

本教程提供了一个基本对Go程序员的指导关于如何使用gRPC。

通过练习这个例子，可以学习到如下知识：

1. 在`.proto`文件中定义服务。
2. 使用`protocol buffers`编译器生成服务端和客户端代码。
3. 使用Go gRPC API为你的服务编写一个简单的客户端和服务端。

假定已阅读[概述](https://grpc.io/docs/)并熟悉[`protocol buffers`](https://developers.google.com/protocol-buffers/docs/overview)。请注意，本教程中的示例使用`protocol buffers`的`proto3`版本：可以在[proto3语言指南](https://developers.google.com/protocol-buffers/docs/proto3)和[Go代码生成指南](https://developers.google.com/protocol-buffers/docs/reference/go-generated)中找到更多信息。

## 为何使用gRPC

我们的示例是一个简单的路径映射应用程序，允许客户端获取有关其路径`Feature`的信息，并创建其`RouteSummary`后与其他服务端或客户端交换路径信息，如，进行流量更新。

使用gRPC，只需要在`.proto`文件中定义一次服务，然后就可以使用gRPC支持的任何语言实现客户端和服务端，而这些语言又可以在谷歌内部的服务器或个人平板电脑等各种环境中运行，不同的语言和环境之间通信的复杂性都由gRPC处理。我们还获得了使用`protocol buffers`的所有优势，包括**高效的序列化**，**简单的IDL**和**简单的接口更新**。

## 样例代码与设置

本例子的样例代码在GitHub的[grcp-go仓库](https://github.com/grpc/grpc-go/tree/master/examples/route_guide)中。执行如下指令来克隆仓库中的代码：

```bash
go get google.golang.org/grpc
```

然后，切换到样例目录`$GOPATH/src/grpc-go/examples/route_guide`中：

```bash
cd $GOPATH/src/google.golang.org/grpc/examples/route_guide
```

还应该安装相关工具来生成服务端和客户端接口代码，如果还没有准备好，请按照[Go快速入门指南](/gRPC/07-Golang快速入门.md)中的设置说明进行操作。

## 定义服务

第一步，使用`protocol buffers`定义：

1. gRPC服务
2. 方法请求类型
3. 方法响应的类型

完整的[`.proto`文件](/gRPC/route_guide.proto)在这里。

定义一个服务需要在`.proto`文件中指定服务的名字：

```go
service RouteGuide {
   ...
}
```

然后在定义的服务中定义`rpc`方法，指定方法请求和响应的类型。gRPC允许定义四种类型的服务方法，在示例服务`RouterGuide`中都用到了，具体如下。

- **简单RPC**：客户端使用`stub`向服务端发送请求然后等待响应返回，就像普通的方法调用：

  ```go
  // 获得给定位置的特征。
  rpc GetFeature(Point) returns (Feature) {}
  ```

- **服务端侧流数据RPC**：客户端向服务端发送单个请求并获取返回流以读取消息序列。客户端从返回的流中读取，直到没有更多消息。正如在示例中所写的那样，将`stream`关键字放在**响应类型**之前来指定服务器端流方法。

  ```go
  // 获得给定Rectangle中可用的特征。
  // 得到的结果是流式传输而不是一次返回（例如，在响应消息中有重复字段），
  // 因为rectangle可能覆盖很大面积并且包含大量的特征。
  rpc ListFeatures(Rectangle) returns (stream Feature) {}
  ```

- **客户端侧流数据RPC**：客户端多次使用提供的流写入一系列消息并将它们发送到服务端。一旦客户端写完消息，它就等待服务端全部读取并返回响应的响应。通过在**请求类型**之前放置`stream`关键字来指定客户端流方法。

  ```go 
  // 接受正在遍历的路径上的Points消息类型的流，在遍历完成时返回RouteSummary。
  rpc RecordRoute(stream Point) returns (RouteSummary) {}
  ```

- **双向流数据RPC**：双方使用读写流发送一系列消息，这两个流独立运行，因此客户端和服务端可以按照自己喜欢的顺序进行读写。（例如，服务端可以在写入响应之前等待接收所有客户端消息，或者它可以交替地读取消息然后写入响应消息，或者其他一些读写组合）。每个流中的消息顺序都能得到保证。可以通过在请求和响应之前放置`stream`关键字来指定此类方法。

  ```go
  // 接受在遍历路径时发送的RouteNotes消息类型的流，同时接收其他RouteNotes消息类型的消息（例如，来自其他用户）。
  rpc RouteChat(stream RouteNote) returns (stream RouteNote) {}
  ```

`.proto`文件还包含服务方法中使用的所有请求和响应类型的`protocol buffers`消息类型定义，例如，这里是Point消息类型：

```go
message Point {
  int32 latitude = 1;
  int32 longitude = 2;
}
```

## 生成客户端和服务端代码

第二步，使用`protoc`(`protocol buffers`的编译器)和特定的`gRPC-GO`插件(`protoc-gen-go`)，根据`.proto`文件中定义的服务生成gRPC客户端和服务端接口。这与[快速入门](/gRPC/07-Golang快速入门)中的操作一样。

在`route_guide`的示例目录中运行如下命令：

```bash
protoc -I routeguide/ routeguide/route_guide.proto --go_out=plugins=grpc:routeguide

#  在目录下生成如下文件

route_guide.pb.go
```

route_guide.pb.go文件包含：

- 填充，序列化和检索请求和响应消息类型的所有`protocol buffers`代码
- 实现`RouteGuide`服务中定义的客户端接口类型或stub（存根）和方法
- 实现`RouteGuide`服务中定义的服务端接口类型和方法

## 创建服务端

首先创建`RouteGuide`服务端。如果只对创建gRPC客户端感兴趣，可以跳过本节直接阅读[创建客户端](##创建客户端)。

要使`RouteGuide`服务能够正常提供它的服务，有两个部分需要完成：

1. 实现从服务定义中生成的服务接口：它执行服务的实际“工作”
2. 运行gRPC服务以监听来自客户端的请求并将其分派给正确的服务端实现

在`grpc-go/examples/route_guide/server/server.go`中可以看到`RouteGuide`服务端样例。

### 实现 RouteGuide

如下所示，服务端有一个叫`routeGuideServer`的结构体类型它实现了自动生成的`RouteGuideServer`接口。

```go
type routeGuideServer struct {
        ...
}
...

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
        ...
}
...

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
        ...
}
...

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
        ...
}
...

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
        ...
}
...
```

#### 简单RPC

`routeGuideServer`实现了所有的服务方法，先看看最简单的类型`GetFeature()`，它只是从客户端获取一个`Point`并从某个`Feature`自己的数据库中返回相应的特征信息。

```go
func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
  for _, feature := range s.savedFeatures {
    if proto.Equal(feature.Location, point) {
      return feature, nil
    }
  }
  // No feature was found, return an unnamed feature
  return &pb.Feature{"", point}, nil
}
```

该方法传递RPC的上下文对象和客户端`Point`类型的`protocol buffers`请求。它返回一个`Feature`类型的`protocol buffers`对象，其中包含响应信息和错误。在此方法中，我们使用适当的信息填充`Feature`，然后将其与`nil`错误一起返回来告诉gRPC已经完成了RPC的处理，然后`Feature`就可以返回给客户端。

#### 服务端侧流数据RPC

现在看一下流数据RPC。`ListFeatures()`方法是服务端侧流数据RPC，因此需要将多个`Feature`返回给客户端。

```go
func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
  for _, feature := range s.savedFeatures {
    if inRange(feature.Location, rect) {
      if err := stream.Send(feature); err != nil {
        return err
      }
    }
  }
  return nil
}
```

如上所示，此方法不是在方法参数中获取简单的请求和响应对象，而是得到一个请求对象（在`Feature`中客户端想要查找的`Rectangle`）和一个特定的`RouteGuide_ListFeaturesServer`对象来编写我们的响应。

在该方法中，填充尽可能多的需要被返回的`Feature`对象，并使用其`Send()`方法将`Feature`都写入`RouteGuide_ListFeaturesServer`。最后，就像在简单RPC方法中那样，返回一个`nil`错误告诉gRPC已经完成了写响应信息的操作。如果在此调用中发生任何错误，那么将返回非零错误; gRPC层会将其转换为适当的RPC状态，以便在线路上发送。

#### 客户端侧流数据RPC

现在来看一些更复杂的东西：客户端侧流数据方法`RecordRoute()`，从客户端获取`Points`类型的流并返回单个包含有关传输链路信息的`RouteSummary`。如下所示，该方法根本没有请求参数。相反，它获取`RouteGuide_RecordRouteServer`流，服务端可以使用该流来读取和写入消息，服务端可以使用它的`Recv()`方法接收客户端消息，并使用它的`SendAndClose()`方法返回其单个响应。

```go
func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
  var pointCount, featureCount, distance int32
  var lastPoint *pb.Point
  startTime := time.Now()
  
  for {
    point, err := stream.Recv()
    if err == io.EOF {
      endTime := time.Now()
      return stream.SendAndClose(&pb.RouteSummary{
        PointCount:   pointCount,
        FeatureCount: featureCount,
        Distance:     distance,
        ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
      })
    }

    if err != nil {
      return err
    }

    pointCount++
    for _, feature := range s.savedFeatures {
      if proto.Equal(feature.Location, point) {
        featureCount++
      }
    }

    if lastPoint != nil {
      distance += calcDistance(lastPoint, point)
    }

    lastPoint = point
  }
}
```

在方法体中，使用`RouteGuide_RecordRouteServer`的`Recv()`方法重复读取客户端请求一个请求对象（在本例中为一个`Point`）的请求，直到没有更多消息，服务端需要检查每次调用`Recv()`方法后从中返回的错误`err`。

- 如果错误是`nil`，那么说明流仍然处于正常状态并且可以继续进行读取操作
- 如果错误是`io.EOF`，那么说明消息流已经结束，服务端可以返回它的`RouteSummary`
- 如果错误是任何其他值，那么将“按原样”返回错误，以便gRPC层将其转换为RPC状态

#### 双向流数据RPC

最后，看一下双向流数据PRC `RouteChat()`方法。

```go
func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
  for {
    in, err := stream.Recv()
    if err == io.EOF {
      return nil
    }

    if err != nil {
      return err
    }

    key := serialize(in.Location)

                ... // 此处代码  寻找要发送给客户端的note

    for _, note := range s.routeNotes[key] {
      if err := stream.Send(note); err != nil {
        return err
      }
    }

  }
}
```

得到一个`RouteGuide_RouteChatServer`流，就像客户端侧流数据示例那样，可用于读取和写入消息。但是，在这里通过方法的流来返回值，哪怕此时客户端仍在向其消息流写入消息。

这里的读写语法与客户端流方法中的语法非常相似，只是此处服务端使用流的`Send()`方法而不是`SendAndClose()`方法，因为它正在写入多个响应。尽管每一方都会按照写入的顺序获取对方的消息，但客户端和服务端都可以按任意顺序进行读写，因为，这些流完全独立运行。

### 启动服务端

实现了所有方法后，还需要启动一个gRPC服务端，以便客户端可以实际使用我们的服务。以下代码段显示了如何为`RouteGuide`服务执行此操作：

```go
flag.Parse()
lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

if err != nil {
        log.Fatalf("failed to listen: %v", err)
}

grpcServer := grpc.NewServer()
pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})

... // 此处代码  确定是否启动TLS

grpcServer.Serve(lis)
```

要构建和启动一个服务端，我们需要：

1. 使用`lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port)).`指定需要监听客户端请求的端口
2. 使用`grpc.NewServer()`创建一个gRPC服务端实例
3. 使用gRPC服务端注册我们实现的服务
4. 使用端口详细信息调用服务端的`Serve()`方法进行阻塞等待，直到进程被终止或调用`Stop()`。

## 创建客户端

在本节中，将介绍如何为`RouteGuide`服务创建Go客户端，可以在`grpc-go/examples/route_guide/client/client.go`中查看完整的示例客户端代码。

### 创建一个stub（存根）

要调用服务端的方法，首先需要创建一个gRPC通道来与服务端通信。通过将服务器地址和端口号传递给`grpc.Dial()`来创建它，如下所示：

```go
conn, err := grpc.Dial(*serverAddr)
if err != nil {
    ...
}
defer conn.Close()
```

如果请求的服务需要身份认证，可以使用`DialOptions`在`grpc.Dial()`方法中设置身份验证凭据（例如，TLS，GCE凭据，JWT凭证，但是，在这里的`RouteGuide`服务中不需要执行此操作。

一旦gRPC通道建立起来，就需要有一个客户端stub(存根)来执行RPC调用。我们使用从`.proto`文件中生成的在`pb`包中提供的`NewRouteGuideClient()`方法创建一个客户端stub。

```go
client := pb.NewRouteGuideClient(conn)
```

### 调用服务端方法

现在来看看如何调用服务方法。请注意，在`gRPC-Go`中，RPC以**阻塞/同步**模式运行，这意味着RPC调用等待服务端响应，并将返回响应或错误。

#### 简单RPC

调用简单RPC的`GetFeature()`方法几乎与调用本地方法一样简单。

```go
feature, err := client.GetFeature(context.Background(), &pb.Point{409146138, -746188906})
if err != nil {
        ...
}
```

如上所示：在之前获得的`stub`上调用该方法，在方法参数中，创建并填充请求`protocol buffers`对象（在例子中为`Point`）。 还传递一个`context.Context`对象，它允许我们在必要时更改RPC的行为，例如超时/取消执行中的RPC请求。如果调用没有返回错误，那么可以从第一个返回值中读取服务端的响应信息。

```go
log.Println(feature)
```
#### 服务端侧流数据RPC

这是调用服务端流方法`ListFeatures()`的地方，该方法返回地理`Feature`流。如果已经阅读过[创建服务端](##创建服务端)，其中一些部分可能看起来非常熟悉：流数据RPC在服务端和客户端之间都以类似的方式实现。

```go
rect := &pb.Rectangle{ ... }  // initialize a pb.Rectangle
stream, err := client.ListFeatures(context.Background(), rect)
if err != nil {
    ...
}
for {
    feature, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
    }
    log.Println(feature)
}
```

与简单的RPC一样，传递上下文和请求给方法，但是在这里不会返回响应对象，而是返回`RouteGuide_ListFeaturesClient`的实例。这个客户端实例可以使用`RouteGuide_ListFeaturesClient`流来读取服务端的响应。

使用`RouteGuide_ListFeaturesClient`的`Recv()`方法重复的在服务器的响应中读取protocol buffers对象（在本例中为Feature），直到没有更多消息：这个客户端实例需要检查每次调用`Recv()`返回的错误`err`。

- 如果是`nil`，那么说明这个流仍然是正常的并且可以继续读取
- 如果是`io.EOF`，那么说明消息流已经结束
- 否则必须有一个RPC错误，它通过`err`参数传递。

#### 客户端侧流数据RPC

客户端流方法`RecordRoute`类似于服务器端方法，不过只传递给方法一个上下文，然后会获取一个`RouteGuide_RecordRouteClient`流的返回，这个流可以用来写和读消息。

```go
// Create a random number of random points
r := rand.New(rand.NewSource(time.Now().UnixNano()))
pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
var points []*pb.Point

for i := 0; i < pointCount; i++ {
	points = append(points, randomPoint(r))
}

log.Printf("Traversing %d points.", len(points))

stream, err := client.RecordRoute(context.Background())
if err != nil {
	log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
}

for _, point := range points {
	if err := stream.Send(point); err != nil {
		if err == io.EOF {
			break
		}
		log.Fatalf("%v.Send(%v) = %v", stream, point, err)
	}
}

reply, err := stream.CloseAndRecv()
if err != nil {
	log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
}

log.Printf("Route summary: %v", reply)
```

`RouteGuide_RecordRouteClient`有一个`Send()`方法，可以用它向服务端发送请求。一旦使用`Send()`将客户端的请求向流中写入完成，就需要在流上调用`CloseAndRecv()`方法让gRPC知道我们已完成写入并期望收到响应。可以从`CloseAndRecv()`返回的错误`err`中获取RPC状态，如果状态为`nil`，那么`CloseAndRecv()`的第一个返回值将是有效的服务端响应。

#### 双向流数据RPC

最后是双向流数据RPC`RouteChat()`。与`RecordRoute`的情况一样，只给方法传递一个上下文对象，然后获得返回一个可用于写入和读取消息的流。但是，这次通过方法的流返回值，而服务端仍在向其消息流写入消息。

```go
stream, err := client.RouteChat(context.Background())
waitc := make(chan struct{})
go func() {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
	}
}()
for _, note := range notes {
	if err := stream.Send(note); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}
stream.CloseSend()
<-waitc
```

这里的读写语法与客户端流方法的语法非常相似，只是在完成调用后使用流的`CloseSend()`方法。尽管每一方都会按照写入顺序获取对方的消息，但客户端和服务端都可以按任意顺序进行读写，这些流完全独立运行。

## 搞起来

要编译和运行服务端，假设位于`$GOPATH/src/google.golang.org/grpc/examples/route_guide`文件夹中，只需：

```go
go run server/server.go
```

同样，要运行客户端：

```go
go run client/client.go
```