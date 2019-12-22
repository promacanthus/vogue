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
