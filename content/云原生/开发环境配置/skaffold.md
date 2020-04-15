---
title: skaffold
date: 2020-04-14T10:09:14.226627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- 云原生
- 开发环境配置
summary: skaffold
showInMenu: false

---

## 安装

下载最新稳定版，并放在PATH中。

```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
chmod +x skaffold
sudo mv skaffold /usr/local/bin
```

为了使Skaffold保持最新状态，会对Google服务器进行更新检查，以查看是否有新版本的Skaffold。默认情况下，此行为是启用的。作为副作用，该请求被记录。要禁用更新检查，有两个选择： 

1. 将`SKAFFOLD_UPDATE_CHECK`环境变量设置为false
2. 使用以下命令在skaffold的全局配置中将其关闭`skaffold config set -g update-check false`

> Skaffold将使用minikube中托管的Docker守护程序构建应用程序。如果想针对其他Kubernetes集群进行部署，例如此类，GKE集群，将必须安装Docker才能构建此应用程序。

## 常用指令

## `skaffold init`

Skaffold需要skaffold.yaml，对于受支持的项目Skaffold可以生成一个简单的配置。在项目的根目录下运行`skaffold init`为应用程序配置Skaffold。

### `skaffold dev`

`skaffold dev`可以在应用程序上进行连续的本地开发。在开发模式下，Skaffold将监视应用程序的源文件，并在检测到更改时将重建镜像（或将文件同步到正在运行的容器中），推送任何新镜像并将应用程序重新部署到群集中。`skaffold dev`被认为是Skaffold的主要操作模式，因为它允许在迭代应用程序时连续地利用Skaffold的所有功能。

监视本地源代码，并在每次检测到更改时执行Skaffold管道。 `skaffold.yaml`提供了工作流程的规范，例如管道为：

1. 使用Dockerfile从源代码构建Docker镜像
2. 使用其内容的sha256哈希值标记Docker镜像
3. 更新Kubernetes清单k8s-pod.yaml来运行上一步构建的镜像
4. 使用`kubectl apply -f`部署Kubernetes清单
5. 从已部署的应用流回日志

修改go文件源代码后，保存文件时，Skaffold将看到此更改，并重复skaffold.yaml中描述的工作流程，重新构建并重新部署应用程序。

#### Dev loop

运行`skaffold dev`时，Skafold首先将对skaffold.yaml中指定的所有模块进行完整构建和部署，类似于`skaffold run`。 成功构建和部署后，Skaffold将开始监视项目中指定的所有模块的所有源文件依赖性。在对这些源文件进行更改时，Skaffold将重建关联的模块，并将新更改重新部署到集群。

dev循环将一直运行，直到用户使用`Ctrl + C`取消Skaffold进程为止。收到此信号后，Skaffold将清除活动集群上所有已部署的模块，这意味着Skaffold将不会放弃在运行的整个生命周期中创建的任何Kubernetes资源。 可以选择使用`--no-prune`标志禁用此功能。

#### 执行顺序

Skaffold在dev循环期间执行的操作有优先顺序，因此行为始终是可预测的。操作顺序为：

1. 文件同步
2. 构建
3. 部署

#### 文件监视程序和监视模式

Skaffold根据使用的构建器和模块的根目录来计算每个模块的依赖性。一旦计算出所有源文件依赖关系，在dev模式下，Skaffold将连续监视这些文件的后台更改，并在检测到更改时有条件地重新运行循环。

默认情况下，Skaffold使用fsnotify监视本地文件系统上的事件。

Skaffold支持：

- 轮询（polling）模式：在此模式下以可配置的间隔检查文件系统的更改
- 手动（manual）模式：在此模式下，Skaffold等待用户输入以检查文件更改

可以通过`--trigger`标志配置这些监视模式。

#### Control API

默认情况下，每次在本地更改文件时，dev循环将执行所有操作（根据需要），手动触发模式下操作除外。但是，可以通过Skaffold API的用户输入来关闭各个动作。

使用此API，即使被监视的文件在文件系统上已更改，用户也可以告诉Skaffold在执行任何这些操作之前等待用户输入。这样，用户可以在本地进行迭代时“排队”更改，然后仅在被要求时重新构建和重新部署Skaffold。 当构建的发生频率超出预期时，构建或部署需要很长时间或成本很高时，或者用户希望将其他工具与skaffold dev集成时，这将非常有用。

### `skaffold run`

如果希望只构建和部署一次，那么运行`skaffold run`。 Skaffold将只执行一次skaffold.yaml中描述的工作流程。
