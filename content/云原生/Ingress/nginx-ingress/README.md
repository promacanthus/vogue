---
title: README.md
date: 2020-04-14T10:09:14.130627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 云原生
- Ingress
- nginx-ingress
summary: README.md
showInMenu: false

---

# 安装指南

## 先决条件和通用部署命令

**注意：**

1. 默认配置从所有命名空间中监视Ingress对象。使用标志`--watch-namespace`将范围限制为特定命名空间。
2. 如果多个Ingress为同一主机定义不同的路径，则`ingress-controller`将合并定义。
3. 如果正在使用GKE，则需要使用以下命令将用户初始化为集群管理员。

```bash
kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole cluster-admin \
  --user $(gcloud config get-value account)
```

**所有部署都需要强制执行以下命令**。

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/mandatory.yaml
```

### 裸机部署

不同的云平台部署过程不同，[可以参考这里](https://kubernetes.github.io/ingress-nginx/deploy/#provider-specific-steps)。

1. 使用NodePort暴露`Ingress-controller`访问端口。

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/baremetal/service-nodeport.yaml
```

2. 验证部署

检查`ingress-controller`是否已启动，请运行以下命令：

```bash
kubectl get pods --all-namespaces -l app.kubernetes.io/name=ingress-nginx --watch
```

3. 查看已安装的版本

要检查处于运行状态的`ingress-controller`是哪个版本的，进入这个pod然后执行`nginx-ingress-controller version`命令。

```bash
POD_NAMESPACE=ingress-nginx
POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app.kubernetes.io/name=ingress-nginx -o jsonpath='{.items[0].metadata.name}')

kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- /nginx-ingress-controller --version
```

### 使用Helm部署

从charts官方的仓库中使用这个[chart](https://github.com/kubernetes/charts/tree/master/stable/nginx-ingress)，来安装`ingress-controller`以`my-nginx`作为实例名字。

```bash
helm install stable/nginx-ingress --name my-nginx

# 如果kubernetes集群开启了RBAC，执行如下命令
helm install stable/nginx-ingress --name my-nginx --set rbac.create=true

# 检查安装的版本
POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=ingress-nginx -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POD_NAME -- /nginx-ingress-controller --version
```
