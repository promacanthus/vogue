# 04-Charts简介

Helm使用称为chart的包装格式。chart是描述相关的一组Kubernetes资源的文件集合。

> 单个chart可能用于部署简单pod，或者一些复杂的应用程序堆栈。

chart通过创建为[特定目录树文件](/05-自定义chart.yaml)，将它们打包到版本化的压缩包，然后进行部署。

## Chart.yaml

```yaml
apiVersion: chart API 版本, 总是 "v1" (必选)
name: chart的名字 (必选)
version: SemVer 2 版本 (必选)
kubeVersion: A SemVer range of compatible Kubernetes versions (可选)
description: A single-sentence description of this project (可选)
keywords:
  - A list of keywords about this project (可选)
home: The URL of this project's home page (可选)
sources:
  - A list of URLs to source code for this project (可选)
maintainers: # (可选)
  - name: The maintainer's name (每个维护者必选)
    email: The maintainer's email (每个维护者可选)
    url: A URL for the maintainer (每个维护者可选)
engine: gotpl # The name of the template engine (可选, 默认是gotpl)
icon: A URL to an SVG or PNG image to be used as an icon (可选).
appVersion: The version of the app that this contains (可选). This needn't be SemVer.
deprecated: Whether this chart is deprecated (可选, boolean)
tillerVersion: The version of Tiller that this chart requires. This should be expressed as a SemVer range: ">2.0.0" (可选)
```

## 通过charts/目录手动管理依赖性

通过将依赖的charts复制到`charts/`目录来明确表达这些依赖关系。

> 依赖关系可以是这些chart归档（foo-1.2.3.tgz）或解压缩的chart目录（但它的名字不能从`_`或`.`开始，这些文件会被chart加载器忽略）。

例如，WordPress chart依赖Apache chart和MySQL chart，则在WordPress chart的`charts/`目录中提供如下charts：

```bash
wordpress:
  Chart.yaml
  # ...
  charts/
    apache/
      Chart.yaml检查依赖和模板配置是否正确
      # ...
    mysql/
      Chart.yaml
      # ...

# 将依赖项放入charts目录，使用 helm fetch 命令
```

### 使用依赖关系的操作方面影响

指定chart依赖关系后，如何影响使用`helm intall`和`helm upgrade`的chart安装？

假设，名为“A”的chart创建一下kubernetes对象：

- namespace：A-Namespace
- StatefulSet：A-StatefulSet
- service：A-Service

同时，A依赖于创建对象的chart “B”：

- namespace：B-Namespace
- deployment：B-Deployment
- service：B-Service

安装/升级chart A之后，会创建/修改单个Helm版本。该版本会按照以下顺序创建/更新所有上述kubernetes对象：

- A-Namespace
- B-Namespace
- A-StatefulSet
- B-Deployment
- A-Service
- B-Service

因为，当helm安装/升级chart时，chart中的kubernetes对象及其所有依赖项都是：

- 聚合成一个单一的集合
- 按类型排序
- 按名称排序

**按上述顺序创建/更新**。

单个release是使用chart及其依赖关系创建的所有对象。（kubernetes类型的安装顺序有kind_sorter.go中InstallOrder给出）

## 模板Template和值Values

Helm  chart模板使用Go模板语言[Go template language](https://golang.org/pkg/text/template/)编写，其中添加了来自[Sprig库](https://github.com/Masterminds/sprig)的50个左右的附加模板函数以及一些专用函数。

所有模板文件都存储在chart的`templates/`目录下。**当Helm渲染charts时，它将通过模板引擎渲染传递该目录中的每个文件**。

模板的值有两种提供方式：

1. chart开发人员可能会在chart内部提供一个`values.yaml`文件，该文件可以包含默认值。
2. chart用户可能会提供（使用命令helm install -f）一个包含值的YAML文件。当用户提供自定义值时，这些值将覆盖chart中`values.yaml`文件中的值。

### 预定义值

通过`values.yaml`文件（或通过`--set`标志）提供的值，可以从`.Values`模板中的对象访问。**也可以在模板中访问其他预定义的数据片段**。

以下值是预定义的，可以用于每个模板，并且**不能被覆盖**。与所有值一样，**名称区分大小写**：

|内置对象|描述|
---|---
Release|release本身
Release.Name|release的名称（不是chart的名称）
Release.Time|chart版本上次更新的时候。这将匹配Last Released发布对象上的时间
Release.Namespace|release 发布的 namespace
Release.Service|处理 release 的服务。**始终是Tiller**
Release.Revision|release的修订版本号。它从 1 开始，并随着你每次`helm upgrade`增加
Release.IsUpgrade|如果当前操作是升级或回滚，则设置为 true
Release.IsInstall|如果当前操作是安装，则设置为 true
||
Values|从 values.yaml 文件和用户提供的文件传入模板的值。默认情况下，Values 是空的。
||
Chart|Chart.yaml 的内容。chart 版本可以从 `Chart.Version` 和维护人员 `Chart.Maintainers`一起获得
||
Files|这提供对 chart 中所有非特殊文件的访问。虽然无法使用它来访问模板，但可以使用它来访问 chart 中的其他文件（除非它们被排除使用 .helmignore）
Files.Get/Files.GetString|是按名称获取文件的函数（`.Files.Get ”file name“`）
Files.GetBytes|是将文件内容作为字节数组而不是字符串获取的函数。这对于像图片这样的东西很有用
||
Capabilities|提供关于 Kubernetes 集群支持的功能的信息
Capabilities.APIVersions |是一组版本信息
Capabilities.APIVersions.Has $version |指示是否在群集上启用版本（例如 `batch/v1`）
Capabilities.KubeVersion |提供查找 Kubernetes 版本的方法。具有以下值：`Major，Minor，GitVersion，GitCommit，GitTreeState，BuildDate，GoVersion，Compiler，Platform`
Capabilities.TillerVersion |提供查找 Tiller 版本的方法。具有以下值：`SemVer，GitCommit，GitTreeState`
||
Template|包含有关正在执行的当前模板的信息
Name|到当前模板的 namespace 文件路径（例如 `mychart/templates/mytemplate.yaml`）
BasePath|当前 chart 模板目录的 namespace 路径（例如 `mychart/templates`）

**这些值可用于任何顶级模板。这并不意味着它们将在任何地方都要有**。

> 注意，任何位置的chart.yaml字段将被删除，不会在chart对象内部被访问。因此，chart.yaml不能用于将任意结构化的数据传递到模板中，values.yaml文件可以用于传递。

### 值 values.yaml文件

Values的内容来源有四个地方：

1. chart中的values.yaml文件
2. 子chart来自父chart的values.yaml文件
3. value文件通过`helm install`或`helm upgrade`的`-f`标志运行用户提供的YAML值
4. 通过`--set`标志运行用户提供的YAML值

**上述四个来源优先级逐渐升高，高优先级会覆盖低优先级**。

```bash
helm install --values=myvals.yaml wordpress
# 以这种方式传递值时，将被合并到默认values文件中

# values.yaml文件的内容
imageRegistry: "quay.io/deis"
dockerTag: "latest"
pullPolicy: "Always"
storage: "s3"

# myvals.yaml文件的内容
storage: "gcs"

# 执行上述命令合并后的values.yaml文件
imageRegistry: "quay.io/deis"
dockerTag: "latest"
pullPolicy: "Always"
storage: "gcs"

# 注意，最后一个字段被覆盖了，其他的不变
```

注意：

1. 包含在charts中的默认values文件的名称必须为values.yaml，在命令行上指定的文件可以被指明为任何名称
2. 如果在`helm  install`或`helm upgrade` 使用`--set`，这些值仅在客户端转换为YAML

#### 删除默认key

如果您需要从默认值中删除一个键，可以覆盖该键的值为 null，在这种情况下，Helm 将从覆盖值合并中删除该键。

如，stable 版本的 Drupal chart 允许配置 liveness 探测器，如果配置自定义的 image。以下是默认值：

```yaml
livenessProbe:
  httpGet:
    path: /user/login
    port: http
  initialDelaySeconds: 120
```

如果尝试覆盖 liveness Probe 处理程序 exec 而不是 httpGet，使用 `--set livenessProbe.exec.command=[cat,docroot/CHANGELOG.txt]`，Helm 会将默认和重写的键合并在一起，从而产生以下 YAML：

```yaml
livenessProbe:
  httpGet:
    path: /user/login
    port: http
  exec:
    command:
    - cat
    - docroot/CHANGELOG.txt
  initialDelaySeconds: 120
```

但是，Kubernetes 会报错，因为无法声明多个 liveness Probe 处理程序。为了克服这个问题，可以指示 Helm 将 `livenessProbe.httpGet` 设置为空来删除它：

```bash
helm install stable/drupal --set image=my-registry/drupal:0.1.0 --set livenessProbe.exec.command=[cat,docroot/CHANGELOG.txt] --set livenessProbe.httpGet=null
```

### 范围scope、依赖dependencies、值values

`values.yaml`文件可以声明顶级chart的值，也可以为chart的`charts/`目录中包含的任何chart声明值，即values可以为chart及其任何依赖项提供值。

例如，WordPress依赖Apache和MySQL，values文件可以为所有这些组件提供值：

```yaml
title: "My WordPress Site" # Sent to the WordPress template

mysql:
  max_connections: 100 # Sent to MySQL
  password: "secret"

apache:
  port: 8080 # Passed to Apache
```

- 高级别的chart可以访问下面定义的所有变量：所以WordPress chart可以访问MySQL密码`.Values.mysql.password`
- 低级别的chart无法访问父级别chart中的内容：所有MySQL无法访问`title`属性，也无法访问`.Values.apache.port`

#### 全局值

从` 2.0.0-Alpha.2` 开始，Helm 支持特殊的 “全局” 值：

```yaml
title: "My WordPress Site" # Sent to the WordPress template

global:                   # 新增加的字段，为全局值
  app: MyWordPress        # 此值可供所有chart使用 .Values.global.app

mysql:
  max_connections: 100 # Sent to MySQL
  password: "secret"

apache:
  port: 8080 # Passed to Apache
```

这提供了一种与所有子chart共享一个顶级变量的方法，这对设置metadata中像标签这样的属性很有用。

> 如果子chart声明了一个全局变量，则该全局将向下传递（到子chart的子chart），但不向上传递到父chart（子chart无法影响到父chart的值）。

**父chart的全局变量优先与子chart中的全局变量**。

当涉及到编写模板和values文件时，可以参考一下几个标准：

- [Go templates](https://godoc.org/text/template)
- [The YAML format](https://yaml.org/spec/)

## chart起始包

`helm create`命令采用可选`--starter`选项，可以指定起始chart。

起始chart只是普通的chart，位于`$HELM_HOME/starters`。作为chart开发人员，可以创作专门设计用作起始chart。记住这些chart时应考虑以下因素：

- chart.yaml将被生成器覆盖
- 用户将期望修改这样的chart内容，因此文档应该指出用户如何做到这一点
- 所有templates目录下的匹配项`<CHARTNAME>`将被替换为指定的chart名称，以便起始chart可用作模板，另外，values.yaml的`<CHARTNAME>`也会被替换
- 目前添加chart的唯一方法是手动将其复制到`$HELM_HOME/starters`,在chart文档中，需要解释该过程
