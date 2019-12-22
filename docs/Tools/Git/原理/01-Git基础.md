# 01-Git基础

## 版本控制演变过程

### VCS出现之前

1. 用目录拷区别不同版本
2. 公共文件容易被覆盖
3. 成员沟通成本很高，代码集成效率低下

### 集中式VCS/SVN

1. 有集中的版本管理服务器
2. 具备文件版本管理和分支管理能力
3. 集成效率有明显地提高
4. 客户端必须**时刻**和服务器相连

### 分布式VCS/SVN

1. 服务器和客户端都有完整的版本库
2. 脱离服务端，客户端照样可以管理版本
3. 查看历史和版本比较等多数操作，都不需要访问服务器，比集中式VCS更能提高版本管理效率

### Git的特点

1. 最优的存储能力
2. 非凡的性能
3. 开源
4. 很容易做备份
5. 支持离线操作
6. 容易定制工作流程

## Git最小化配置

配置user.name和user.email

``` bash
git config --global user.name 'your_name'
git config --global user.email 'your_email'
```

### Git的三个作用域

优先级： local > global > system

- 缺省等同于local

``` bash
git config --local #local只对仓库有效
git config --global #global对登录用户所有仓库有效
git config --system #system对系统的所有用户有效
```

- 显示config的配置，加--list

``` bash
git config --list --local
git config --list --global
git config --list --system
```

- 设置，缺省等同于local

``` bash
git config --local
git config --global
git config --system
```

- 清除设置

``` bash
git config --unset --local user.name
git config --unset --global user.name
git config --unset --system user.name
```

## 创建Git仓库

### 两种场景

- 把已有的项目代码纳入到Git管理

``` bash
cd <项目代码所在文件夹>
git init
```

- 新建的项目直接使用Git管理

``` bash
cd <某个文件夹>
git init project_name #会在当前路径下创建和项目名称相同的文件夹
cd project_name
```

## Git暂存区

```bash
git add .   # 将当前路径下所有的文件都添加到暂存区中
git add  <file> # 将<file>t添加到暂存区中
git add -u  # 将已经被git追踪的文件进行更新

git reset -head # 重置HEAD、索引和工作区

git  mv a b # 将被git追踪的文件a重命名为b并添加到暂存区中
```

## Git版本历史

```bash
git log oneline # 查看简洁的历史
git log -n<number> #显示最近的n次日志
git log --all   # 查看全部日志
git log --graph # 显示为图片
```

- 查看本地分支

```bash
git branch -v
```

## Git图形界面

```bash
# 安装gitk工具
sudo apt-get install gitk

# 在git仓库目录下打开工具
gitk
```

## .git目录

```bash
drwxr-xr-x  branches/
-rw-r--r--  COMMIT_EDITMSG
-rw-r--r--  config  # 本地仓库(local)相关配置
-rw-r--r--  description
-rw-rw-r--  FETCH_HEAD
-rw-rw-r--  gitk.cache
-rw-rw-r--  HEAD    # 指向仓库当前工作的分支
drwxr-xr-x  hooks/
-rw-rw-r--  index
drwxr-xr-x  info/
drwxr-xr-x  logs/
drwxr-xr-x  objects/    #文件夹中的子文件夹都是以哈希值的前两位字符命名每个object由40位字符组成，前两位字符用来当文件夹，后38位做文件
drwxr-xr-x  refs/headers    # 分支
drwxr-xr-x  refs/tags   # 里程碑
```

```bash
git cat-file    # 命令 显示版本库对象的内容、类型及大小信息
git cat-file -t b44dd71d62a5a8ed3   # 显示版本库对象的类型
git cat-file -s b44dd71d62a5a8ed3   # 显示版本库对象的大小
git cat-file -p b44dd71d62a5a8ed3   # 显示版本库对象的内容
```

git中的对象：

- commit：一次提交生成一个tree
- tag：里程碑
- tree：保存对应commit时间点，仓库中文件与目录的结构以及其中的内容
- blob：表示一个文件，与文件名无关与文件内容有关

> 通常，blob表示一个文件，tree表示一个文件夹。

**在Git中，文件内容相同的文件就是唯一的一个blob**。

> 没有文件也就是没有blob对象的目录是不会被git管理的，因为git要对文件进行版本管理，所以没有必要对空目录生成对象。基于这一点，假设`readme`文件的全路径是这样：`[仓库根目录]/doc/readme`。那么tree的数量与全路径中`“/”`的数量一致。

**即，有几层文件夹，就有几个tree。**

一个`commit`对应一个`tree`，这个是root节点。

## 分离头指针(HEAD)

```bash
git checkout <commit ID>
```

可用于在对应commit下进行实验性尝试，尝试完成后直接切换回原分支即可。

> Git认为没有与分支或tag绑定的commit都应该丢弃。

### HEAD与branch

- HEAD可以指向分支的最后一次commit
- HEAD也可以指向某一个具体的commit

> HEAD永远指向commit。

```bash
git diff    <commit ID>  <commit ID>    # 比较两次commit的差异
git diff  HEAD  HEAD^1      # 比较HEAD与HEAD前的一个commit
# HEAD^^ == HEAD~2
```
