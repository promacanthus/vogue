# minikube

Minikube是一种可以轻松在本地运行Kubernetes的工具。 Minikube在笔记本电脑上的虚拟机（VM）内运行一个单节点Kubernetes集群，以供希望试用Kubernetes或每天进行开发的用户使用。

Minikube支持以下Kubernetes功能：

- DNS
- NodePorts
- ConfigMaps and Secrets
- Dashboards
- Container Runtime: Docker, CRI-O, and containerd
- Enabling CNI (Container Network Interface)
- Ingress

## 安装前检查

要检查Linux是否支持虚拟化，运行以下命令并验证输出是否为非空：

```bash
egrep -q 'vmx|svm' /proc/cpuinfo && echo yes || echo no
```

## 安装kubectl

参考这里安装[kubectl](../开发环境/kubectl.md)。

## 安装Hypervisor

安装Hypervisor，可选[KVM](https://www.linux-kvm.org/page/Main_Page)或者[VirtualBox](https://www.virtualbox.org/wiki/Downloads)。

Minikube还支持`--vm-driver = none`选项，该选项在主机而不是VM中运行Kubernetes组件。使用此驱动程序需要Docker和Linux环境，但不需要管理程序。

> 如果在Debian或衍生产品中使用`none`驱动程序，安装Docker使用`.deb`软件包，而不要使用对Minikube不起作用的snap软件包。可以从[Docker](https://www.docker.com/products/docker-desktop)下载`.deb`软件包。

**警告**：`none`驱动程序会导致安全和数据丢失问题。使用`--vm-driver = none`之前，请查阅本文档以获取更多信息。

## 安装minikube

```bash
# 在github获取对应版本的release软件包
https://github.com/kubernetes/minikube/releases

# 命令行下载
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube
sudo mv minikube /usr/local/bin

# 使用install快速导入到PATH中
sudo mkdir -p /usr/local/bin/
sudo install minikube /usr/local/bin/
```

## 确认安装

要确认虚拟机管理程序和Minikube均已成功安装，可以运行以下命令来启动本地Kubernetes集群：

> 注意：要使用`minikube start`来设置`--vm-driver`，请输入安装的虚拟机监控程序的名称，并使用小写字母，替换下面命令中的<driver_name>。

```bash
# 启动
minikube start --vm-driver = <driver_name>
sudo minikube start --vm-driver=none

# 将none设置为默认的驱动
sudo minikube config set vm-driver none

# 增加内存分配，默认2G
minikube config set memory 4096

# 输出
Starting local Kubernetes cluster...
Running pre-create checks...
Creating machine...
Starting local Kubernetes cluster...、

# 检查
minikube status

# 正常运行显示如下输出
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured

# 停止
minikube stop

# 清除本地集群
minikube delete

# 国内使用镜像
minikube start --image-mirror-country='cn' --image-repository='registry.cn-hangzhou.aliyuncs.com/google_containers'
```

目前支持的`--vm-driver`值的完整列表：

- virtualbox
- vmwarefusion
- kvm2 (driver installation)
- hyperkit (driver installation)
- hyperv (driver installation) Note that the IP below is dynamic and can change. It can be retrieved with minikube ip.
- vmware (driver installation) (VMware unified driver)
- none (Runs the Kubernetes components on the host and not in a virtual machine. You need to be running Linux and to have Docker installed.)

设置普通用户运行

```bash
sudo mv /home/sugoi/.kube /home/sugoi/.minikube $HOME
sudo chown -R $USER $HOME/.kube $HOME/.minikube
```
