# kubectl

## 安装

```bash
# 直接下载最新版二进制文件
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

# 下载指版本
curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl

# kubectl赋予执行权限
chmod +x ./kubectl

# 将二进制文件移动的PATH中
sudo mv ./kubectl /usr/local/bin/kubectl

# 测试一下
kubectl version --client

---------

# 通过包管理工具下载
sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

---------

# 通过snap商店下载
snap install kubectl --classic
```

## 配置

为了让kubectl查找和访问Kubernetes集群，它需要一个kubeconfig文件，该文件在您使用[kube-up.sh](https://github.com/kubernetes/kubernetes/blob/master/cluster/kube-up.sh)创建集群或成功部署Minikube集群时会自动创建。默认情况下，kubectl配置位于`~/.kube/config`。

```bash
# 通过获取集群状态来检查kubectl是否已正确配置
kubectl cluster-info
# 如果看到URL响应，则表明kubectl已正确配置为访问集群

# 如果您看到类似以下的消息，则说明kubectl配置不正确或无法连接到Kubernetes集群
The connection to the server <server-name:port> was refused - did you specify the right host or port?

# 如果kubectl cluster-info返回了url响应，但无法访问集群，要检查其配置是否正确，请使用
kubectl cluster-info dump
```

Bash的kubectl完成脚本可以使用命令`kubectl completion bash`生成。在shell中输入完成脚本可启用kubectl自动补全功能。但是，完成脚本取决于[bash-completion](https://github.com/scop/bash-completion)，这意味着必须先安装此软件（可以通过运行`type _init_completion`来测试是否已安装bash-completion）。

使用`apt-get install bash-completion`或者`yum install bash-completion`来安装bash-completion。上面的命令创建`/usr/share/bash-completion/bash_completion`，这是bash-completion的主要脚本。根据软件包管理器的不同，需要手动`source ~/.bashrc`。

重新加载shell并运行`type _init_completion`，如果命令成功，则表示已经安装成功，否则将`source /usr/share/bash-completion/bash_completion`内容添加到`~/.bashrc`中。

```bash
# 启动命令补全

# 将补全脚本添加到.bashrc
echo 'source <(kubectl completion bash)' >>~/.bashrc

# 将补全脚本添加到/etc/bash_completion.d目录
kubectl completion bash >/etc/bash_completion.d/kubectl

# 如果对kubectl有别名，则可以扩展shell补全功能以使用该别名
echo 'alias k=kubectl' >>~/.bashrc
echo 'complete -F __start_kubectl k' >>~/.bashrc
```

## 手动创建 kubeconfig
kubeconfig 由 kube-up 生成，但是，也可以使用下面命令生成自己想要的配置（可以使用任何想要的子集）

```bash
# create kubeconfig entry
kubectl config set-cluster $CLUSTER_NICK \
    --server=https://1.1.1.1 \
    --certificate-authority=/path/to/apiserver/ca_file \
    --embed-certs=true \
# Or if tls not needed, replace --certificate-authority and --embed-certs with
    --insecure-skip-tls-verify=true \
    --kubeconfig=/path/to/standalone/.kube/config

# create user entry
kubectl config set-credentials $USER_NICK \
# bearer token credentials, generated on kube master
    --token=$token \
# use either username|password or token, not both
    --username=$username \
    --password=$password \
    --client-certificate=/path/to/crt_file \
    --client-key=/path/to/key_file \
    --embed-certs=true \
    --kubeconfig=/path/to/standalone/.kube/config

# create context entry
kubectl config set-context $CONTEXT_NAME \
    --cluster=$CLUSTER_NICK \
    --user=$USER_NICK \
    --kubeconfig=/path/to/standalone/.kube/config
```
注：
生成独立的 kubeconfig 时，标识 `--embed-certs` 是必选的，这样才能远程访问主机上的集群。

`--kubeconfig`既是加载配置的首选文件，也是保存配置的文件。如果您是第一次运行上面命令，那么 `--kubeconfig` 文件的内容将会被忽略。

```bash
export KUBECONFIG=/path/to/standalone/.kube/config
```

上面提到的 `ca_file`，`key_file` 和 `cert_file` 都是集群创建时在 master 上产生的文件，可以在文件夹 `/srv/kubernetes` 下面找到。持有的 token 或者 基本认证也在 master 上产生。
如果您想了解更多关于 kubeconfig 的详细信息，运行帮助命令 `kubectl config -h`。
