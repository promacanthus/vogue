---
title: 13-Charts调试与测试
date: 2020-04-14T10:09:14.130627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- 云原生
- Helm
summary: 13-Charts调试与测试
showInMenu: false

---

## 调试模板

调试模板比较麻烦，因为模板在Tiller服务器而不是Helm客户端上渲染。然后渲染的模板被发送到KubernetesAPI服务器，可能由于格式以外的原因，服务器可能会拒绝接收这些YAML文件。

有几个命令可以帮助您进行调试：

1. `helm lint`是验证 chart 是否遵循最佳实践的首选工具
2. `helm install --dry-run --debug`：让服务器渲染模板，然后返回结果清单文件的方法
3. `helm get manifest`：查看服务器上安装的模板的方法

## 测试chart

一个chart包含许多一起工作的Kubernetes资源和组件。在开发charts时，需要编写一些测试来验证charts在安装时是否按预期工作。这些测试也有助于使用者了解charts应该做什么。

测试在Helm chart中的`templates/tests/test-connection.yaml`中是一个pod定义，指定一个给定的命令来运行容器，如下所示。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: {{ include "mysql.fullname" . }}-test-connection
  labels:
{{ include "mysql.labels" . | indent 4 }}
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args:  ['{{ include "mysql.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
```

在 Helm 中，有两个测试 hook：`helm.sh/hook:test-success`和`helm.sh/hook:test-failure`。

- test-success 表示测试pod应该成功完成。Pod中的容器应该`exit 0`.
- test-failure 是一种断言测试容器不能成功完成的方式。Pod中的容器未`exit 0`，则表示成功。

主要测试如下内容：

1. 验证来自`values.yaml`文件的配置是否正确注入
2. 确保用户名和密码正常工作
3. 确保不正确的用户名和密码不起作用
4. 断言服务已启动并正确进行负载平衡

可以使用该`helm test`命令在release中运行Helm中的预定义测试。对于chart使用者来说，这是一种很好的方式来检查发布的chart（或应用程序）是否按预期工作。

> 可以在单个yaml文件中定义尽可能多的测试，也可以在`templates/tests`目录中的多个yaml文件中进行分布测试。
