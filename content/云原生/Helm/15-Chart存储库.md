---
title: 15-Chart存储库
date: 2020-04-14T10:09:14.130627+08:00
draft: false
---

Chart存储库是可以存储和共享charts的地方，Helm官方维护了一个[Chart存储库](https://github.com/helm/charts)。使用Helm可以轻松创建和运行自己的Chart存储库。

## 自建Chart存储库

chart repo是容纳一个或多个封装的chart的HTTP服务器。
> 虽然Helm可用于管理本地chart目录，但在共享chart时，首选机制是chart repo库。Helm附带了用于开发人员测试的内置服务器（helm server）。

作为repo库服务器的条件：

1. 可以提供YAML文件和tar文件
2. 可以回答GET请求的HTTP服务器

repo库的主要特征是存在一个名为`index.yaml`的特殊文件，它具有：

1. repo库提供所有软件包的列表
2. 允许检索和验证这些软件包的元数据

在客户端，repo库使用`helm repo`命令进行管理，Helm不提供将chart上传到远程存储服务器的工具（这样做会增加部署服务器的需求，从而增加配置repo库的难度）。

> Chart库是能提供YAML和tar文件并回答GET请求的HTTP服务器，因此托管Chart存储库时，很多选择。例如，Google云存储桶，AmazonS3存储桶，GithubPages或者创建自己的Web服务器。

### Chart库结构

Chart存储库分为两部分：

- index.yaml：包含Chart库中所有charts的索引
- 若干已打包的charts

例如有一个Chart存储库`https://example.com/charts`，如下所示：

```bash
charts/
  |
  |- index.yaml
  |
  |- alpine-0.1.2.tgz
  |
  |- alpine-0.1.2.tgz.prov

# charts和index.yaml文件可以位于同一台服务器，也可位于不同服务器
```

提供的chart的下载URL：`https://example.com/charts/alpine-0.1.2.tgz`。

### 索引文件

一个有效的Chart存储库必须包含一个索引文件`index.yaml`，索引文件中包含库中每个chart的元数据，如chart的Chart.yaml文件的内容。

`helm repo index`命令，根据已经打包的charts文件给本地目录生成索引文件，索引文件示例如下：

```yaml
apiVersion: v1
entries:
  alpine:
    - created: 2016-10-06T16:23:20.499814565-06:00
      description: Deploy a basic Alpine Linux pod
      digest: 99c76e403d752c84ead610644d4b1c2f2b453a74b921f422b9dcb8a7c8b559cd
      home: https://k8s.io/helm
      name: alpine
      sources:
      - https://github.com/helm
      urls:
      - https://technosophos.github.io/tscharts/alpine-0.2.0.tgz
      version: 0.2.0
    - created: 2016-10-06T16:23:20.499543808-06:00
      description: Deploy a basic Alpine Linux pod
      digest: 515c58e5f79d8b2913a10cb400ebb6fa9c77fe813287afbacf1a0b897cd78727
      home: https://k8s.io/helm
      name: alpine
      sources:
      - https://github.com/helm
      urls:
      - https://technosophos.github.io/tscharts/alpine-0.1.0.tgz
      version: 0.1.0
  nginx:
    - created: 2016-10-06T16:23:20.499543808-06:00
      description: Create a basic nginx HTTP server
      digest: aaff4545f79d8b2913a10cb400ebb6fa9c77fe813287afbacf1a0b897cdffffff
      home: https://k8s.io/helm
      name: nginx
      sources:
      - https://github.com/helm/charts
      urls:
      - https://technosophos.github.io/tscharts/nginx-1.1.0.tgz
      version: 1.1.0
generated: 2016-10-06T16:23:20.499029981-06:00
```

生成的索引和包可以从由网络服务器提供，使用`helm server`可以启动本地服务器，在本地测试所有内容。

```bash
# 启动一个本地web服务器，在./charts目录找到chart提供服务
$ helm serve --repo-path ./charts
Regenerating index. This may take a moment.
Now serving you on 127.0.0.1:8879
# server命令将在启动过程中自动生成一个index.yaml文件
```

## 托管Chart存储库

### 普通web服务器

配置普通Web服务器来提供Chart存储库服务，只需执行以下操作：

1. 将索引和charts置于服务器目录中
2. 确保index.yaml可以在没有认证要求的情况下访问
3. 确保yaml文件的正确内容类型（text/yaml或text/x-yaml）

> 如果要在`$WEBROOT/charts`以外的目录为chart提供服务，请确保Web根目录中有一个`charts/`目录，并将索引文件和chart放入该文件夹内。

## 管理Chart存储库

### 将chart存储到Chart存储库

**Chart存储库中的chart必须被正确打包并有正确的版本号**。具体步骤如下：

```bash
# 打包
helm package docs/examples/alpine/

mkdir fantastic-charts
mv alpine-0.1.0.tgz fantastic-charts/

# 使用本地路径和远程Chart存储库URL，并在指定路径中生成index.yaml
helm repo index fantastic-charts --url https://fantastic-charts.storage.googleapis.com

```

### 添加新的chart到Chart存储库

每次将新chart添加到存储库时，都必须重新生成索引。`helm repo index`命令将`index.yaml`从头开始完全重建该文件，但仅包括它在本地找到的charts。

可以使用`--merge`标志向现有`index.yaml`文件增量添加新chart。

如果同时生成了出处文件`provenance`，也要一起上传。

### 共享charts

只要知道Chart存储库的URL就能够共享charts。

使用`helm repo add [NAME] [URL]`命令将Chart存储库添加到使用者的helm客户端中，并且可以给存储库取一个别名。

```bash
helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com
helm repo list
fantastic-charts    https://fantastic-charts.storage.googleapis.com
```

如果Chart存储库由HTTP基本认证支持，也可以在此处提供用户名和密码：

```bash
helm repo add fantastic-charts https://fantastic-charts.storage.googleapis.com --username my-username --password my-password
helm repo list
fantastic-charts    https://fantastic-charts.storage.googleapis.com
```

**如果Chart存储库不包含index.yaml文件，则添加不成功**。

添加成功后，可以搜索charts，通过`helm repo update`指令来获取最新的charts。

>`helm repo add`和`helm repo update`命令获取`index.yaml`文件并将它们存储在 `$HELM_HOME/repository/cache/`目录中。这是`helm search`找到有关charts信息的地方。
