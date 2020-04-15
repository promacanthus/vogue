---
title: 12-Charts注意事项
date: 2020-04-14T10:09:14.130627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- 云原生
- Helm
summary: 12-Charts注意事项
showInMenu: false

---

## 引用字符串，不要引用整数值

```go
// 使用字符串时，引用字符串（用双引号引起来）比把它们留为空白字符更安全
name: {{.Values.MyName | quote}}

// 使用整数时 不要引用整数值。否则，可能会导致Kubernetes内部的解析错误
port: {{.Values.Port}}
```

以上两种情况，在设置env变量值时不适用

```go
//  env变量值都需要引起来
env:
  - name: HOST
    value: "http://host"
  - name: PORT
    value: "1234"
```

## 管理空格

当模板引擎运行时，它将删除`{{   }}`中的空白内容，但是按原样保留剩余的空白。

> **换行符也是空格。**

**YAML中的缩进空格是严格的，因此管理空格变得非常重要**。Helm模板有几个工具可以使用。

**第一种：**

使用特殊字符修饰模板声明的大括号语法，以告诉模板引擎填充空格。

   1. `{{-` 添加破折号和空格，表示应该将格左移，
   2. `-}}` 添加空格和破折号，表示应该删除右空格

> 确保`-`和其他指令之间有空格

   1. `- 3` 意思是删除左空格并打印3
   2. `-3` 意思是打印-3

**第二种：**

告诉模板系统如何缩进比掌握模板指令的间距更容易，可以使用indent函数：`{{ indent 2 value: "true" }}`。

## 小心随机值生成

Helm中有一些函数允许生成随机数据、加密秘钥等，在升级过程中模板会被重新渲染，当模板运行产生与上次运行不同的数据时，将触发该资源的更新。

## 避免使用块

Go 模板语言提供了一个 block 关键字，允许开发人员提供一个默认的实现，后续将被覆盖。在 Helm chart 中，块不是重写的最佳工具，因为如果提供了同一个块的多个实现，那么所选哪个是不可预知的。

建议是改为使用include。

## 一般约定

### Chart名称

Chart 名称应该是小写字母和数字组成，字母开头：

```yaml
drupal
cert-manager
oauth2-proxy

# Chart 名称中不能使用大写字母和下划线。Chart 名称不应使用点。
```

**包含 chart 的目录必须与 chart 具有相同的名称**。因此，chart cert-manager 必须在名为 cert-manager/ 的目录中创建。这不仅仅是一种风格的细节，而是 Helm Chart 格式的要求。

### 版本号

Helm 使用 [SemVer](https://semver.org/)2 来表示版本号。当 SemVer 版本存储在 Kubernetes 标签中时，通常会将 + 字符更改为一个 _ 字符，因为标签不允许 + 标志作为值。

### 格式化YAML

YAML 文件应该使用两个空格缩进（而不是制表符）。

### 通过版本来限制Tiller

一个 Chart.yaml 文件可以指定一个 tillerVersion SemVer 约束：

```yaml
name: mychart
version: 0.2.0
tillerVersion: ">=2.4.0"

```

当模板使用 Helm 旧版本不支持的新功能时，应该设置此限制。虽然此参数将接受复杂的 SemVer 规则，但最佳做法是默认为格式 >=2.4.0，其中 2.4.0 引入了 chart 中使用的新功能的版本。

此功能是在Helm 2.4.0中引入的，因此任何2.4.0版本以下的Tiller都会忽略此字段。

## Values

变量名称应该以小写字母开头，单词应该用 camelcase 分隔：

```yaml
# 正确写法：

chicken: true
chickenNoodleSoup: true

# 不正确写法：

Chicken: true  # 可能与内置变量冲突
chicken-noodle-soup: true # 不要在变量名称中使用短横
```

### 展平或嵌套值

YAML 是一种灵活的格式，并且值可以嵌套或扁平化。

```yaml
# 嵌套：
server:
  name: nginx
  port: 80

# 展平：
serverName: nginx
serverPort: 80

# 在大多数情况下，展平应该比嵌套更受青睐,因为对模板开发人员和用户来说更简单。
```

为了获得最佳安全性，必须在每个级别检查嵌套值：

```yaml
{{if .Values.server}}
  {{default "none" .Values.server.name}}
{{end}}
# 对于每一层嵌套，都必须进行存在检查。但对于展平配置，可以跳过这些检查，使模板更易于阅读和使用。

{{default "none" .Values.serverName}}
# 当有大量相关变量时，且至少有一个是非可选的，可以使用嵌套值来提高可读性。
```

### 使类型清晰

YAML 的类型强制规则有时是违反直觉的。
例如：

- foo: false
- foo: "false"
- foo: 12345678 ：在某些情况下，大整数将被转换为科学记数法

避免类型转换错误的最简单方法是明确地表示字符串（引用所有字符串），并隐含其他所有内容。

> 通常，为了避免整型转换问题，最好将整型存储为字符串，并在模板中使用 `{{int $value}}` 将字符串转换为整数。

在大多数情况下，显式类型标签受到重视，所以 `foo: !!string 1234` 应该将 1234 视为一个字符串。但是，**YAML 解析器消费标签，因此类型数据在解析后会丢失**。

### 考虑用户如何使用你的 values

有几种潜在的 values 来源：

1. chart 的 values.yaml 文件
2. 由 helm install -f 或 helm upgrade -f 提供的 value 文件
3. 传递给 --set 或的 --set-string 标志 helm install 或 helm upgrade 命令
4. 通过 --set-file 将 文件内容传递给 helm install or helm upgrade

在设计 value 的结构时，请记住 chart 的用户可能希望通过 -f 标志或 --set 选项覆盖它们。

由于 --set 在表现力方面比较有限，编写 values.yaml 文件的第一个指导原则可以轻松使用 --set 覆盖。

出于这个原因，使用 map 来构建 value 文件通常会更好。

```yaml
# 难以配合 --set 使用
servers:
  - name: foo
    port: 80
  - name: bar
    port: 81
```

Helm <=2.4 时，以上不能用 --set 来表示。在 Helm 2.5 中，访问 foo 上的端口是 --set servers[0].port=80。用户不仅难以弄清楚，而且如果稍后 servers 改变顺序，则容易出错。

```yaml
# 使用方便

servers:
  foo:
    port: 80
  bar:
    port: 81
# 访问 foo 的端口更为方便：--set servers.foo.port=80
```

### 文档'values.yaml'

应该记录'values.yaml'中的每个定义的属性。文档字符串应该以它描述的属性的名称开始，然后至少给出一个单句描述。

```yaml
# 不正确
# the host name for the webserver
serverHost = example
serverPort = 9191

# 正确
# serverHost is the host name for the webserver
serverHost = example
# serverPort is the HTTP listener port for the webserver
serverPort = 9191
```

使用参数名称开始每个注释，它使文档易于grep，并使文档工具能够可靠地将文档字符串与其描述的参数关联起来。

## templates 目录结构

templates 目录的结构应如下所示：

- 如果产生 YAML 输出，模板文件应该有扩展名 `.yaml`。扩展名`.tpl`可用于产生不需要格式化内容的模板文件。
- 模板文件名应该使用横线符号（my-example-configmap.yaml），而不是 camelcase。
- 每个资源定义应该在它自己的模板文件中。
- 模板文件名应该反映名称中的资源种类。例如 foo-pod.yaml， bar-svc.yaml

### 定义模板的名称

定义的模板（在 `{{define}}` 指令内创建的模板）可以全局访问。这意味着 chart 及其所有子 chart 都可以访问所有使用 `{{ define }}`创建的模板。出于这个原因，所有定义的模板名称应该是带有某个 namespace。

```yaml
# 正确
{{- define "nginx.fullname"}}
{{/* ... */}}
{{end -}}

# 不正确
{{- define "fullname" -}}
{{/* ... */}}
{{end -}}

# 强烈建议通过 helm create 命令创建新 chart
```

### 格式化模板

**模板应该使用两个空格缩进（不是制表符）**。

模板指令在大括号之后和大括号之前应该有空格：

```yaml
# 正确
{{.foo}}
{{print "foo"}}
{{- print "bar" -}}

# 不正确
{{.foo}}
{{print "foo"}}
{{-print "bar"-}}

```

模板应尽可能地填充空格：

```yaml
foo:
  {{- range .Values.items}}
  {{.}}
  {{end -}}
  
```

块（如控制结构）可以缩进以指示模板代码的流向

```yaml
{{if $foo -}}
  {{- with .Bar}}Hello{{ end -}}
{{- end -}}
```

但是，由于 YAML 是一种面向空格的语言，因此代码缩进有时经常不能遵循该约定。

### 生成模板中的空格

最好将生成的模板中的空格保持最小。特别是，许多空行不应该彼此相邻。但偶尔空行（特别是逻辑段之间）很好。

```yaml
# 最佳实践
apiVersion: batch/v1
kind: Job
metadata:
  name: example
  labels:
    first: first
    second: second
```

## YAML注释与模板注释

```yaml
# YAML 注释 用#
# This is a comment
type: sprocket

# 模板注释 用/**/
{{- /*
This is a comment.
*/ -}}
type: frobnitz
```

记录模板功能时应使用模板注释，如解释定义的模板：

```yaml
{{- /*
mychart.shortname provides a 6 char truncated version of the release name.
*/ -}}
{{define "mychart.shortname" -}}
{{.Release.Name | trunc 6}}
{{- end -}}
```

在模板内部，当 Helm 用户可能（有可能）在调试过程中看到注释时，可以使用 YAML 注释。

```yaml
# This may cause problems if the value is more than 100Gi
memory: {{.Values.maxMem | quote}}
```

上面的注释在用户运行 helm install --debug 时可见，而在 `{{- /* */ -}}` 部分中指定的注释不是。

## 标签OR注释

在下列条件下，元数据项应该是标签：

- Kubernetes 使用它来识别此资源
- 为了查询系统目的，向操作员暴露是非常有用的

例如，使用 `helm.sh/chart: NAME-VERSION` 作为标签，以便操作员可以方便地查找要使用的特定 chart 的所有实例。

如果元数据项不用于查询，则应将其设置为注释。

**Helm hook 总是注释**。

### 标准标签

下表定义Helm chart 使用的通用标签。Helm 本身从不要求特定的标签。

- 标记为 REC 的标签是表示推荐的，应放置在 chart 上以保持全局一致性。
- 标记 OPT 是表示可选的。这些都是惯用的或通常使用的，但不是经常用于运维目的。

名称|状态|描述
---|---|---
app.kubernetes.io/name|REC|这应该是应用程序名称，反映整个应用程序。 通常使用`{{template“name” .}}`来实现此目的。 许多Kubernetes清单都使用它，而不是Helm特有的。
helm.sh/chart|REC|这应该是chart名字和版本: `{{.Chart.Name}}-{{ .Chart.Version \replace "+" "_" }}`.
app.kubernetes.io/managed-by|REC|这里总是被设置为 `{{.Release.Service}}`. 这是为了找到由Tiller管理的所有东西。
app.kubernetes.io/instance|REC|这里应该是 `{{.Release.Name}}`. 它有助于区分同一应用程序的不同实例。
app.kubernetes.io/version|OPT|应用程序的版本可以被设置为 `{{.Chart.AppVersion}}`.
app.kubernetes.io/component|OPT|这是用于标记应用可能在应用程序中扮演的不同角色的通用标签。 例如 `app.kubernetes.io/component: frontend`
app.kubernetes.io/part-of|OPT|当多个chartsor软件一起构成一个应用时. 例如，应用软件和数据库来构成网站。 这可以设置为支持的顶级应用程序。
