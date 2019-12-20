# 03-Helm命令

Kubernetes 包管理工具。

```bash
# 要开始使用Helm，请运行`helm init`命令
# 这将把Tiller安装到正在运行的Kubernetes集群中，它还将设置任何必要的本地配置
helm init
```

共用的命令操作包括：

```bash
- helm search:    搜索charts
- helm fetch:     下载一个chart到本地来查看
- helm install:   安装一个chart到kubernetes集群中
- helm list:      列出安装的Release
```

环境变量:

```bash
- $HELM_HOME:           为Helm文件设置备用位置。 默认情况下，它们存储在"~/.helm"中
- $HELM_HOST:           设置另一个Tiller主机，格式是 host:port
- $HELM_NO_PLUGINS:     设置 HELM_NO_PLUGINS=1 来取消插件
- $TILLER_NAMESPACE:    设置另一个Tiller的namespace(默认是 "kube-system")
- $KUBECONFIG:          设置另一个kubernetes配置文件 (默认是 "~/.kube/config")
- $HELM_TLS_CA_CERT:    用于验证Helm客户端和Tiller服务器证书的TLS CA证书的路径 (默认是 "$HELM_HOME/ca.pem")
- $HELM_TLS_CERT:       用于向Tiller进行身份验证的TLS客户端证书文件的路径 (默认是 "$HELM_HOME/cert.pem")
- $HELM_TLS_KEY:        用于向Tiller进行身份验证的TLS客户端密钥文件的路径 (默认是 "$HELM_HOME/key.pem")
- $HELM_TLS_ENABLE:     启用Helm和Tiller之间的TLS连接 (默认是 "false")
- $HELM_TLS_VERIFY:     启用Helm和Tiller之间的TLS连接并验证Tiller服务器证书 (默认是 "false")
- $HELM_TLS_HOSTNAME:   用于验证Tiller服务器证书的主机名或IP地址 (默认是 "127.0.0.1")
- $HELM_KEY_PASSPHRASE: 设置为PGP私钥的密码。 如果设置，则在签署helm charts时不会提示输入密码
```

用法:

```bash
  helm [command]

可选的 Commands:
  completion  为指定的shell（bash或zsh）生成自动补全脚本
  create      创建具有给定名称的新chart
  delete      删除Kubernetes集群中指定的release
  dependency  管理chart的依赖
  fetch       从仓库下载指定chart并可选择是否在本地路径下解压缩
  get         下载指定的release
  help        帮助
  history     获取release的历史记录
  home        输出 HELM_HOME 的路径
  init        在客户端和服务端同时初始化Helm
  inspect     检查chart
  install     安装chart压缩包
  lint        检查依赖和模板配置是否正确
  list        列出所有的release
  package     将chart目录打包成chart压缩包
  plugin      增加、列出、或删除Helm的插件
  repo        增加、列出、删除、更新和索引chart 仓库
  reset       从集群中卸载Tiller
  rollback    回滚release到前一个版本
  search      在charts中根据关键字搜索
  serve       启动一个本地的http web 服务器
  status      显示指定release的状态
  template    本地模板渲染
  test        测试release
  upgrade     升级release
  verify      验证给定路径上的chart是否已签名且有效的
  version     输出服务端和客户端版本信息

Flags:
      --debug                           启用详细输出
  -h, --help                            输出Helm的帮助信息
      --home string                     指定Helm配置文件的位置，用于覆盖$HELM_HOME (默认是 "/root/.helm")
      --host string                     指定Tiller地址，用于覆盖$HELM_HOST
      --kube-context string             指定要使用的kubeconfig上下文的名称
      --kubeconfig string               指定要使用的kubeconfig文件的绝对路径
      --tiller-connection-timeout int   指定Helm将等待建立与tiller的连接持续时间（秒）（默认为300）
      --tiller-namespace string         指定Tiller的namespace (默认是 "kube-system")

使用 "helm [command] --help" 来输出更多关于该指令的信息
```
