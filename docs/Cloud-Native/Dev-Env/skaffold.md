# skaffold

## 安装

下载最新稳定版，并放在PATH中。

```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
chmod +x skaffold
sudo mv skaffold /usr/local/bin
```

为了使Skaffold保持最新状态，会对Google服务器进行更新检查，以查看是否有新版本的Skaffold。默认情况下，此行为是启用的。作为副作用，该请求被记录。要禁用更新检查，有两个选择： 

1. 将`SKAFFOLD_UPDATE_CHECK`环境变量设置为false
2. 使用以下命令在skaffold的全局配置中将其关闭`skaffold config set -g update-check false`

