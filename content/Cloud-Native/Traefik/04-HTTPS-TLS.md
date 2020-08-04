---
title: "04 HTTPS & TLS"
date: 2020-08-04T10:26:07+08:00
draft: true
---

- [0.1. 总览](#01-总览)
- [0.2. TLS](#02-tls)
  - [0.2.1. 证书定义](#021-证书定义)
    - [0.2.1.1. Traefik自定执行](#0211-traefik自定执行)
    - [0.2.1.2. 用户自定义](#0212-用户自定义)
  - [0.2.2. 证书存储](#022-证书存储)
  - [0.2.3. 默认证书](#023-默认证书)
  - [0.2.4. TLS选项](#024-tls选项)
    - [0.2.4.1. TLS最低版本](#0241-tls最低版本)
    - [0.2.4.2. TLS最高版本](#0242-tls最高版本)
    - [0.2.4.3. 密码套件](#0243-密码套件)
    - [0.2.4.4. 曲线首选项](#0244-曲线首选项)
    - [0.2.4.5. 严格的SNI检查](#0245-严格的sni检查)
    - [0.2.4.6. 首选服务器密码套件](#0246-首选服务器密码套件)
    - [0.2.4.7. 客户端身份验证（mTLS）](#0247-客户端身份验证mtls)
- [0.3. 实战操作](#03-实战操作)
  - [0.3.1. 定义IngressRoute](#031-定义ingressroute)
  - [0.3.2. 定义Services](#032-定义services)
  - [0.3.3. 定义Deployments](#033-定义deployments)
  - [0.3.4. 转发端口](#034-转发端口)
  - [0.3.5. Traefik Routers¶](#035-traefik-routers)

## 0.1. 总览

Traefik支持`HTTPS`和`TLS`，这涉及配置的两个部分：路由器和`TLS`连接（及其基础证书）。

当路由器必须处理`HTTPS`流量时，应该在路由器定义中使用`tls`字段指定它。下面几部分将说明如何配置`TLS`连接，包括：

1. 如何获取TLS证书：
   1. 通过动态配置中的定义
   2. 通过`Let's Encrypt（ACME）`
2. 如何配置`TLS`选项和证书存储

## 0.2. TLS

Transport Layer Security

### 0.2.1. 证书定义

#### 0.2.1.1. Traefik自定执行

使用`Let's Encrypt（ACME）`。

#### 0.2.1.2. 用户自定义

即使Traefik已经处于运行状态，也可以在动态配置的`tls.certificates`字段进行添加或删除TLS证书的操作。

```yaml
# Dynamic configuration

tls:
  certificates:
    - certFile: /path/to/domain.cert
      keyFile: /path/to/domain.key
    - certFile: /path/to/other-domain.cert
      keyFile: /path/to/other-domain.key
```

> 注意：在上面的例子中，使用`file provider`来处理这些定义。这是配置证书（以及选项和存储）的唯一可用方法。但是，在Kubernetes中者当属必须是通过`secret`对象提供的。

### 0.2.2. 证书存储

在Traefik中，按照证书存储区进行分组存储，例如下面的定义：

```yaml
# Dynamic configuration

tls:
  stores:
    default: {}
```

> 注意：除了默认定义（命名为`default`）以外的任何存储定义都将被忽略，因此，只有一个全局可用的TLS存储。

在`tls.certificates`字段中，可以指定存储列表以指示证书的存储位置：

```yaml
# Dynamic configuration

tls:
  certificates:
    - certFile: /path/to/domain.cert
      keyFile: /path/to/domain.key
      stores:
        - default
    # Note that since no store is defined,
    # the certificate below will be stored in the `default` store.
    - certFile: /path/to/other-domain.cert
      keyFile: /path/to/other-domain.key
```

> 注意：`stores`字段的列表中的内容会被忽略，然后被自动设置为`default`

### 0.2.3. 默认证书

对于没有SNI（Server Name Indication，服务器名称指示）或没有匹配域名的连接，Traefik可以使用默认证书。此默认证书应在TLS存储中定义：

```yaml
# Dynamic configuration

tls:
  stores:
    default:
      defaultCertificate:
        certFile: path/to/cert.crt
        keyFile: path/to/cert.key
```

如果未提供默认证书，则Traefik会生成并使用自签名证书。

> SNI是TLS的一个扩展协议，在该协议下，在握手过开始时客户端告诉它症状连接的服务器要连接的主机名称。这允许服务器在相同的IP地址和TCP端口号上呈现多个证书，并且因此允许在相同的IP地址上提供多个安全（HTTPS）网站（或其他任何基于TLS的服务），而不需要所有这些站点使用相同的证书。它与`HTTP/1.1`基于名称的虚拟主机的概念相同，但是用于HTTPS。所需的主机名未加密，因此窃听者可以查看请求的网站。
>
> 为了使SNI协议起作用，绝大多数访问者必须使用实现它的Web浏览器。使用未实现SNI浏览器的用户将被提供默认证书，因此很可能会收到证书警告。

### 0.2.4. TLS选项

TLS选项允许配置TLS连接的某些参数。

#### 0.2.4.1. TLS最低版本

```yaml
# Dynamic configuration

tls:
  options:
    default:
      minVersion: VersionTLS12

    mintls13:
      minVersion: VersionTLS13


---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  minVersion: VersionTLS12

---
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: mintls13
  namespace: default

spec:
  minVersion: VersionTLS13
```

#### 0.2.4.2. TLS最高版本

不推荐使用此设置来禁用TLS1.3。正确的方法是更新客户端以支持TLS1.3。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      maxVersion: VersionTLS13

    maxtls12:
      maxVersion: VersionTLS12

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  maxVersion: VersionTLS13

---
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: maxtls12
  namespace: default

spec:
  maxVersion: VersionTLS12
```

#### 0.2.4.3. 密码套件

Traefik使用Golang编写，更多支持的密码套件可以查看`crypto/tls`库，点[这里](https://pkg.go.dev/crypto/tls?tab=doc#pkg-constants)。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      cipherSuites:
        - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  cipherSuites:
    - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
```

> 注意：
>
> - 为TLS 1.2及以下版本定义的密码套件不能在TLS 1.3中使用，反之亦然。
> - 使用TLS 1.3，密码套件不可配置（在这种情况下，所有受支持的密码套件都是安全的）。

#### 0.2.4.4. 曲线首选项

此选项允许以特定顺序设置首选椭圆曲线。

可以使用由`crypto`库定义的曲线名称（例如`CurveP521`）和RFC定义的名称（例如`secp521r1`）。 有关更多信息，请参见[CurveID](https://pkg.go.dev/crypto/tls?tab=doc#CurveID)类型定义。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      curvePreferences:
        - CurveP521
        - CurveP384

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  curvePreferences:
    - CurveP521
    - CurveP384
```

#### 0.2.4.5. 严格的SNI检查

通过严格的SNI检查，Traefik不拒绝来自未指定`server_name`扩展名或与`tlsOption`上配置的任何证书都不匹配的客户端发起的连接。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      sniStrict: true

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  sniStrict: true
```

#### 0.2.4.6. 首选服务器密码套件

此选项允许服务器选择它最喜欢的密​​码套件，而不是客户端的密码套件。请注意，设置`minVersion`或`maxVersion`后，此功能会自动启用。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      preferServerCipherSuites: true

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  preferServerCipherSuites: true
```

#### 0.2.4.7. 客户端身份验证（mTLS）

Traefik通过`clientAuth`字段支持客户端和服务器的认证。

对于需要验证客户端证书的身份验证策略，应在`clientAuth.caFiles`中设置证书的证书颁发机构。

`clientAuth.clientAuthType`选项控制行为，如下所示：

- `NoClientCert`：忽略任何客户端证书。
- `RequestClientCert`：要求提供证书，但是如果没有提供证书，则继续进行。
- `RequireAnyClientCert`：需要证书，但不验证它是否由`clientAuth.caFiles`中列出的CA签名。
- `VerifyClientCertIfGiven`：如果提供了证书，则验证它是否由`clientAuth.caFiles`中列出的CA签名。否则继续进行，无需任何证书。
- `RequireAndVerifyClientCert`：需要证书，该证书必须由`clientAuth.caFiles`中列出的CA签名。

```yaml
# Dynamic configuration

tls:
  options:
    default:
      clientAuth:
        # in PEM format. each file can contain multiple CAs.
        caFiles:
          - tests/clientca1.crt
          - tests/clientca2.crt
        clientAuthType: RequireAndVerifyClientCert

---
# kubernetes
apiVersion: traefik.containo.us/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: default

spec:
  clientAuth:
    # the CA certificate is extracted from key `tls.ca` of the given secrets.
    secretNames:
      - secretCA
    clientAuthType: RequireAndVerifyClientCert
```

## 0.3. 实战操作

### 0.3.1. 定义IngressRoute

首先，定义`IngressRoute`和`Middleware`。另请注意RBAC授权资源；稍后将通过部署的`serviceAccountName`引用它们。

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressroutes.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRoute
    plural: ingressroutes
    singular: ingressroute
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: middlewares.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: Middleware
    plural: middlewares
    singular: middleware
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressroutetcps.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRouteTCP
    plural: ingressroutetcps
    singular: ingressroutetcp
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: ingressrouteudps.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: IngressRouteUDP
    plural: ingressrouteudps
    singular: ingressrouteudp
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: tlsoptions.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TLSOption
    plural: tlsoptions
    singular: tlsoption
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: tlsstores.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TLSStore
    plural: tlsstores
    singular: tlsstore
  scope: Namespaced

---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: traefikservices.traefik.containo.us

spec:
  group: traefik.containo.us
  version: v1alpha1
  names:
    kind: TraefikService
    plural: traefikservices
    singular: traefikservice
  scope: Namespaced

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: traefik-ingress-controller

rules:
  - apiGroups:
      - ""
    resources:
      - services
      - endpoints
      - secrets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - extensions
    resources:
      - ingresses/status
    verbs:
      - update
  - apiGroups:
      - traefik.containo.us
    resources:
      - middlewares
      - ingressroutes
      - traefikservices
      - ingressroutetcps
      - ingressrouteudps
      - tlsoptions
      - tlsstores
    verbs:
      - get
      - list
      - watch

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: traefik-ingress-controller

roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: traefik-ingress-controller
subjects:
  - kind: ServiceAccount
    name: traefik-ingress-controller
    namespace: default
```

### 0.3.2. 定义Services

然后，定义Services：

- 一个用于Traefik本身
- 另一个用于其路由的应用程序，即本例中的演示HTTP服务器：[whoami](https://github.com/containous/whoami)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: traefik

spec:
  ports:
    - protocol: TCP
      name: web
      port: 8000
    - protocol: TCP
      name: admin
      port: 8080
    - protocol: TCP
      name: websecure
      port: 4443
  selector:
    app: traefik

---
apiVersion: v1
kind: Service
metadata:
  name: whoami

spec:
  ports:
    - protocol: TCP
      name: web
      port: 80
  selector:
    app: whoami
```

### 0.3.3. 定义Deployments

接下来，定义Deployment，即服务背后实际运行的Pod。同样：

- 一个是Traefik的pod
- 另一个是whoami应用程序的pod

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: default
  name: traefik-ingress-controller

---
kind: Deployment
apiVersion: apps/v1
metadata:
  namespace: default
  name: traefik
  labels:
    app: traefik

spec:
  replicas: 1
  selector:
    matchLabels:
      app: traefik
  template:
    metadata:
      labels:
        app: traefik
    spec:
      serviceAccountName: traefik-ingress-controller
      containers:
        - name: traefik
          image: traefik:v2.2
          args:
            - --api.insecure
            - --accesslog
            - --entrypoints.web.Address=:8000
            - --entrypoints.websecure.Address=:4443
            - --providers.kubernetescrd
            - --certificatesresolvers.myresolver.acme.tlschallenge
            - --certificatesresolvers.myresolver.acme.email=foo@you.com
            - --certificatesresolvers.myresolver.acme.storage=acme.json
            # Please note that this is the staging Let's Encrypt server.
            # Once you get things working, you should remove that whole line altogether.
            - --certificatesresolvers.myresolver.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory
          ports:
            - name: web
              containerPort: 8000
            - name: websecure
              containerPort: 4443
            - name: admin
              containerPort: 8080

---
kind: Deployment
apiVersion: apps/v1
metadata:
  namespace: default
  name: whoami
  labels:
    app: whoami

spec:
  replicas: 2
  selector:
    matchLabels:
      app: whoami
  template:
    metadata:
      labels:
        app: whoami
    spec:
      containers:
        - name: whoami
          image: containous/whoami
          ports:
            - name: web
              containerPort: 80
```

### 0.3.4. 转发端口

请注意，不应让下面的ingressRoute资源自动应用于集群。 因为，一旦Traefik的ACME提供者检测到拥有TLS路由器，它将尝试为相应域生成证书。

但是的Traefik pod无法从外部到达，这将使ACME TLS challenge失败。 因此，为了使整个工作正常进行，必须延迟应用ingressRoute资源，直到正确设置了端口转发为止。

```bash
kubectl port-forward --address 0.0.0.0 service/traefik 8000:8000 8080:8080 443:4443 -n default
```

请注意，由于Linux上特权端口的限制，上述命令可能无法在端口443上进行侦听。在这种情况下，您以使用一些技巧，例如使用以下命令提高kubectl的上限 setcaps，使用authbind或在主机和WAN之间设置NAT。

### 0.3.5. Traefik Routers¶

现在，可以应用实际的ingressRoutes了：

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: simpleingressroute
  namespace: default
spec:
  entryPoints:
    - web
  routes:
  - match: Host(`your.example.com`) && PathPrefix(`/notls`)
    kind: Rule
    services:
    - name: whoami
      port: 80

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: ingressroutetls
  namespace: default
spec:
  entryPoints:
    - websecure
  routes:
  - match: Host(`your.example.com`) && PathPrefix(`/tls`)
    kind: Rule
    services:
    - name: whoami
      port: 80
  tls:
    certResolver: myresolver
```

给它几秒钟的时间来完成ACME TLS challenge，然后您可以从外部访问whoami pod（通过Traefik路由）。

```bash
curl [-k] https://your.example.com/tls
curl http://your.example.com:8000/notls
```

请注意，只要使用的是Let's Encrypt的服务器，就必须使用`-k`参数，因为它不是一个被认证的CA只是在系统上手动添加的自认证机构。
