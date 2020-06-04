---
title: "04 Session"
date: 2020-06-04T13:58:06+08:00
draft: true
---

当运行对个Pod后，可以通过Service进行负载均衡，默认的方式是`RoundRobin`。Service官方文档点[这里](https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies)。

## 问题

当Pod中运行的服务与会话有关的时候，也就是要确保每次都将来自特定客户端的连接传递到同一Pod。

## 解决方案

### 第一种情况

如果对集群外暴露的服务只使用了Service，那么使用Kubernetes的Service提供的`service.spec.sessionAffinity`进行控制：

- None：默认值
- ClientIP：基于客户端的IP地址选择会话关联

设置了上面的参数之后，通过适当设置 `service.spec.sessionAffinityConfig.clientIP.timeoutSeconds` 来设置最大会话停留时间，（默认值为 10800 秒，即 3 小时）。

这种情况一般使用：

- NodePort
- HostNetwork

### 第二种情况

使用ingress等，将Service暴露到集群外，下面以Traefik为例。

其他边缘路由不确定，Traefik会影响Service的参数配置。也就是说在Traefik不设置sticky session设置的情况下，Service设置了sessionAffinity也不会生效。

Traefik的设置，这里以Traefik 2.2.1为例

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: my-service-ingressroute
  namespace: othernamespace
spec:
  entryPoints:
    - websecure
  routes:
  - match: Host(`my-service.mydomain.com`)
    kind: Rule
    services:
    - name: my-service
      port: 80
      sticky:
        cookie:
          httpOnly: true
```

Traefik官方文档中关于sticky session的说明点[这里](https://docs.traefik.io/routing/services/#sticky-sessions)。
