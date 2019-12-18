# 14-Helm注意事项

## 服务依赖管理

所有使用helm部署的应用中如果没有指定chart的名字都会生成一个随机的Release name。而真正的资源对象的名字是在YAML文件中定义的名字App name，两者连接才是资源对象的名字：`Release name-App name`。

> 使用helm chart部署的包含依赖关系的应用，都使用同一套Release name，在编写YAML文件时，要注意**服务发现**时需要配置的服务地址。

使用环境变量的方式，如下配置：

```yaml
env:
 - name: SERVICE_NAME
   value: "{{ .Release.Name }}-{{ .Values.image.env.SERVICE_NAME }}"

# 使用Go template的语法
# {{ .Values.image.env.SERVICE_NAME }}的值从values.yaml文件中获得

# valus.yaml 配置如下
image:
  env:
    SERVICE_NAME: k8s-app-monitor-test
```

## 解决本地chart依赖

```bash
# 1. 在本地当期配置的目录下启动helm server，不指定参数，直接使用默认端口
helm server

# 2. 将该repo加入到repo list中
helm repo add local http://localhost:8879

# 3. 在本地浏览器访问http://localhost:8879，查看到本地所有的chart

# 4. 下载依赖到本地
helm dependency update

# 所有的chart都会下载到本地的charts目录
```

## helm命令自动补全

```bash
# zsh
source <(helm completion zsh)

# bash
source <(helm completion bash)

```
