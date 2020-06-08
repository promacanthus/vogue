---
title: 02-Helm安装
date: 2020-04-14T10:09:14.126627+08:00
draft: false
---

- 客户端：helm
- 服务端：Tiller

## 安装

### 1. 下载helm客户端

点击[这里](https://github.com/helm/helm/releases)下载。

选择对应版本的二进制安装包，并执行如下命令：

```bash
tar -zxvf helm-v2.0.0-linux-amd64.tgz
mv linux-amd64/helm /usr/local/bin/helm
```

### 2. 创建ServiceAccount和ClusterRoleBinding

使用集群预创建的ClusterRole `<cluster-admin>`

```bash
kubectl apply -f rabc-config.yaml
```

### 3. 安装Tiller

```bash
# 初始化
helm init --service-account tiller --history-max 200 --tiller-image localhost:5000/tiller:v2.14.0 --stable-repo-url https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
# 更换国内的源

# 修改镜像
kubectl set image deployments/tiller-deploy tiller=localhost:5000/tiller:v2.14.0 -n kube-system

# 卸载
kubectl delete deployment tiller-deploy -n kube-system
# 或者
helm reset
```
