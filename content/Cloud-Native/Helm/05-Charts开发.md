---
title: 05-Charts开发
date: 2020-04-14T10:09:14.126627+08:00
draft: false
---

## 前置条件

在编写chart的环境下安装一个Helm客户端，具体步骤参照[Helm安装](../Helm/02-Helm安装.md)部分。

## 自定义chart

Helm相关的命令查看[Helm命令](../Helm/03-Helm命令.md)部分。

```bash
helm create portal  # portal为chart的名字

tree portal

portal/
├── charts                          # 依赖目录，此chart依赖的任何其他charts
├── Chart.yaml                      # 此chart的YAML文件
├── templates                       # 模板目录，当与值组合时，将生成有效的kubernetes manifest文件
│   ├── deployment.yaml             # kubernetes Deployment
│   ├── _helpers.tpl                # 用于修改要生成的kubernetes对象配置的模板，可被被整个chart复用，模板partials默认位置
│   ├── ingress.yaml                # kubernetes Ingress
│   ├── NOTES.txt                   # 包含使用说明的纯文本文件，也可用模板生成其中的内容
│   ├── service.yaml                # kubernetes Service
│   └── tests                       # 测试chart是否如预期运行
│       └── test-connection.yaml    # 可编写多个测试文件，或在一个文件中编写多个测试pod
└── values.yaml                     # 此chart的默认配置值的YAML文件,声明的变量会被传递到templates中

3 directories, 8 files

# .helmignore文件，构建包时要忽略的模式，每行一个模式，支持shell全局匹配，相对路径匹配和否定（前缀为！）
# 所有需要的kubernetes对象都可以在templates文件夹中创建
```

**.helmignore例子：**

```bash
# comment
.git
*/temp*
*/*/temp*
temp?
```

**NOTES.txt例子：**

在 chart install 或 chart upgrade 结束时，Helm 可以为用户打印出一大堆有用的信息。这些信息是使用模板高度定制的。这个文件是纯文本的，但是它像一个模板一样处理，并且具有所有可用的普通模板函数和对象。

```yaml
Thank you for installing {{ .Chart.Name }}.

Your release is named {{ .Release.Name }}.

To learn more about the release, try:

helm status {{ .Release.Name }}
helm get {{ .Release.Name }}
```

## 模板

templates目录下是yaml文件的模板，遵循[Go template](https://golang.org/pkg/text/template/)语法。

**deployment.yaml** 文件的内容如下：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "portal.fullname" . }}
  labels:
{{ include "portal.labels" . | indent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app.kubernetes.io/name: {{ include "portal.name" . }}
      app.kubernetes.io/instance: {{ .Release.Name }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: {{ include "portal.name" . }}
        app.kubernetes.io/instance: {{ .Release.Name }}
    spec:
    {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
    {{- end }}
      containers:
        - name: {{ .Chart.Name }}
          image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
    {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
    {{- end }}
    {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
    {{- end }}

```

其中，双大括号括起来的部分是Go template，其中的Values是`values.yaml`文件中定义的。

**values.yaml**文件的内容如下所示：

```yaml
# Default values for portal.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1

image:
  repository: nginx
  tag: stable
  pullPolicy: IfNotPresent

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths: []

  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}
```

举个例子，比如deployment的镜像：

```yaml
# deployment.yaml
image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
imagePullPolicy: {{ .Values.image.pullPolicy }}

# values.yaml
image:
  repository: nginx
  tag: stable
  pullPolicy: IfNotPresent
```

都是一一对应的关系，根据实际的应用需要进行修改对应的values.yaml中的值。

其中`.Values.image.repository`表示从顶层命名空间开始，先找到Values，然后在里面找到image对象，在image中找到repository对象。用dot（`.`）来分隔每一个namespace。

Helm中除了有Values.yaml文件，还有[内置的对象](04-Charts简介.md)。

## 验证chart

使用helm部署kubernetes的应用的时候，实际上是将`templates`渲染成最终的kubernetes能够识别的`yaml`格式。

在部署之前可以使用`helm install --dry-run --debug <chart-dir>`命令来验证chart配置，该输出中包含了模板的变量配置和最终渲染的yaml文件。
> 该命令还是会向Tiller服务器请求一个round-trip。

输出内容如下：

```yaml
[debug] Created tunnel using local port: '44431'

[debug] SERVER: "127.0.0.1:44431"

[debug] Original chart version: ""
[debug] CHART PATH: /home/sugoi/文档/Helm/Chart/portal

NAME:   independent-camel
REVISION: 1
RELEASED: Mon Jun  3 18:12:22 2019
CHART: portal-0.1.0
USER-SUPPLIED VALUES:
{}

COMPUTED VALUES:
affinity: {}
fullnameOverride: ""
image:
  pullPolicy: IfNotPresent
  repository: nginx
  tag: stable
imagePullSecrets: []
ingress:
  annotations: {}
  enabled: false
  hosts:
  - host: chart-example.local
    paths: []
  tls: []
nameOverride: ""
nodeSelector: {}
replicaCount: 1
resources: {}
service:
  port: 80
  type: ClusterIP
tolerations: []

HOOKS:
---
# independent-camel-portal-test-connection
apiVersion: v1
kind: Pod
metadata:
  name: "independent-camel-portal-test-connection"
  labels:
    app.kubernetes.io/name: portal
    helm.sh/chart: portal-0.1.0
    app.kubernetes.io/instance: independent-camel
    app.kubernetes.io/version: "1.0"
    app.kubernetes.io/managed-by: Tiller
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args:  ['independent-camel-portal:80']
  restartPolicy: Never
MANIFEST:

---
# Source: portal/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: independent-camel-portal
  labels:
    app.kubernetes.io/name: portal
    helm.sh/chart: portal-0.1.0
    app.kubernetes.io/instance: independent-camel
    app.kubernetes.io/version: "1.0"
    app.kubernetes.io/managed-by: Tiller
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: portal
    app.kubernetes.io/instance: independent-camel
---
# Source: portal/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: independent-camel-portal
  labels:
    app.kubernetes.io/name: portal
    helm.sh/chart: portal-0.1.0
    app.kubernetes.io/instance: independent-camel
    app.kubernetes.io/version: "1.0"
    app.kubernetes.io/managed-by: Tiller
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: portal
      app.kubernetes.io/instance: independent-camel
  template:
    metadata:
      labels:
        app.kubernetes.io/name: portal
        app.kubernetes.io/instance: independent-camel
    spec:
      containers:
        - name: portal
          image: "nginx:stable"
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {}
```

可以看到Deployment和Service的名字的前半部分是两个随机单词，后半部分是values.yaml中配置的值。

## 安装chart

有5种不同的方式来表达要安装chart到kubernetes集群中：

1. 通过chart的引用: `helm install stable/mariadb`
2. 通过chart包的路径: `helm install ./nginx-1.2.3.tgz`
3. 通过一个解压后的chart包的路径: `helm install ./nginx`
4. 通过绝对的URL: `helm install https://example.com/charts/nginx-1.2.3.tgz`
5. 通过chart引用和repo的url: `helm install --repo https://example.com/charts/nginx`

## 打包chart

修改chart.yaml中的helm chart配置信息，然后使用如下命令将chart打成压缩文件：

```bash
helm package .    # 打包出portal-0.1.0.tgz
```

## 依赖

新版已经没有requirement.yaml文件，所有的依赖都在charts文件夹中，使用`helm lint` 命令可以和检查依赖和模板配置是否正确。
