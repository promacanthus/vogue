---
title: "05 HPA"
date: 2020-06-28T09:32:50+08:00
draft: true
---

## 简介

Pod 水平自动伸缩（Horizontal Pod Autoscaler）特性， 可以基于：

1. CPU利用率自动伸缩如下API对象中的Pod数量：
   1. `replication controller`
   2. `deployment`
   3. `replica set`
2. 其他应程序提供的度量指标`custom metrics`

HPA无法缩放的对象，比如`DaemonSets`。

HPA由Kubernetes API资源和控制器实现，控制器会周期性的获取平均CPU利用率，并与目标值相比来调整`replication controller`或者`deployment`中副本的数量。

`controller manager`的`--horizontal-pod-autoscaler-sync-period`参数指定周期（默认值15秒），从`resource metrics API`和`custom metrics API`获取指标。

通常情况下，控制器将从一系列的聚合 API 中获取指标数据：

- `metrics.k8s.io`：由 `metrics-server`（需要额外启动）提供
- `custom.metrics.k8s.io`
- `external.metrics.k8s.io`

控制器也可以直接从 Heapster 获取指标。

> 注意：FEATURE STATE: Kubernetes 1.11 [deprecated] 自 Kubernetes 1.11起，从 Heapster 获取指标特性已废弃。

HPA的扩容算法：`期望副本数 = ceil[当前副本数 * ( 当前指标 / 期望指标 )]`

例如，CPU利用率：

- 当前指标为200m，目标设定值为100m,那么由于`200.0/100.0=2.0`， 副本数量将会翻倍。
- 当前指标为50m，副本数量将会减半，因为`50.0/100.0=0.5`。
- 如果计算出的缩放比例接近1.0（跟据`--horizontal-pod-autoscaler-tolerance`参数全局配置的容忍值，默认为`0.1`）， 将会放弃本次缩放。

### 自定义指标API

2018年1月22日，新的指标监控目标是暴露一个API可以让HPA用来获取任意的指标，就像Master的指标API，这个新的API围绕获取的指标构建，这些指标通过引用Kubernetes对象（或组）和指标名称，因此，这些API对于其他的想要消费自定义指标的消费者（尤其是控制器）将很有用。

API的根路径类似这样：`/apis/custom-metrics/v1alpha1`，方便起见，下面省略相同的根路径：

- 检索给定名称的全局对象的指标（如，`Node`、`PersistentVolune`）：`/{object-type}/{object-name}/{metric-name...}`
- 检索给定类型的全局对象的指标：`/{object-type}/*/{metric-name...}`
- 检索给定标签匹配的给定类型的全局对象的指标：`/{object-type}/*/{metric-name...}?labelSelector=foo`
- 检索给定命名空间中的对象的指标：`/namespaces/{namespace-name}/{object-type}/{object-name}/{metric-name...}`
- 检索给定类型的所有命名空间对象的指标：`/namespaces/{namespace-name}/{object-type}/*/{metric-name...}`
- 检索与给定标签匹配的给定类型的所有命名空间对象的指标：`/namespaces/{namespace-name}/{object-type}/*/{metric-name...}?labelSelector=foo`
- 检索描述给定名称空间的指标：`/namespaces/{namespace-name}/metrics/{metric-name}`

### metrics-server

[Metrics Server](https://github.com/kubernetes-sigs/metrics-server)是集群范围的资源使用情况数据聚合器。

- 如果使用`kube-up.sh`脚本部署集群，则会默认部署为`Deployment`
- 如果使用其他方式部署，可以使用`[deployment components.yaml](https://github.com/kubernetes-sigs/metrics-server/releases)`文件来部署

Metrics Server从`Summary API`中收集指标数据，这些指标数据是每个节点上的`kubelet`暴露出来的。
