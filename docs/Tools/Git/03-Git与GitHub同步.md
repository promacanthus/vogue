# 03-Git与GitHub同步

## 配置公私钥

### 创建SSH Key

在用户主目录下，看看有没有`.ssh`目录，如果有，再看看这个目录下有没有`id_rsa`和`id_rsa.pub`这两个文件，如果已经有了，可直接跳到下一步。如果没有，打开Shell（Windows下打开Git Bash），创建SSH Key：

```bash
ssh-keygen -t rsa -b 4096 -C "youremail@example.com"
```

在用户主目录里找到`.ssh`目录，里面有`id_rsa`和`id_rsa.pub`两个文件，这两个就是SSH Key的秘钥对：

- `id_rsa`是私钥，不能泄露出去
- `id_rsa.pub`是公钥，可以放心地告诉任何人

### 登陆GitHub

1. 打开“Account settings”，“SSH Keys”页面
2. 点“Add SSH Key”，填上任意Title，在Key文本框里粘贴id_rsa.pub文件的内容
3. 点“Add Key”，就应该看到已经添加的Key

为什么GitHub需要SSH Key呢？因为GitHub需要识别出你推送的提交确实是你推送的，而不是别人冒充的，而Git支持SSH协议，所以，GitHub只要知道了你的公钥，就可以确认只有你自己才能推送。

当然，GitHub允许你添加多个Key。假定你有若干电脑，你一会儿在公司提交，一会儿在家里提交，只要把每台电脑的Key都添加到GitHub，就可以在每台电脑上往GitHub推送了。

## 把本地仓库同步到GitHub

在本地创建了一个Git仓库后，想在GitHub创建一个Git仓库，并且让这两个仓库进行远程同步，这样，GitHub上的仓库既可以作为备份，又可以让其他人通过该仓库来协作。

1. 登陆GitHub
2. 在右上角找到“Create a new repo”按钮，创建一个新的仓库
3. 在Repository name填入仓库的名字(如 my_repo)，其他保持默认设置，点击“Create repository”按钮，就成功地创建了一个新的Git仓库

    > 在GitHub上的这个my_repo仓库还是空的，GitHub告诉我们，可以从这个仓库克隆出新的仓库，也可以把一个已有的本地仓库与之关联，然后，把本地仓库的内容推送到GitHub仓库。

4. 根据GitHub的提示，在本地的my_repo仓库下运行命令：

    ```bash
    git remote add origin git@github.com:<your Github name>/my_repo.git     # 根据github页面给的提示输入命令即可
    ```

    添加后，远程库的名字就是origin（这是Git默认的叫法，也可以改成别的)，但是origin这个名字一看就知道是远程库。

5. 最后，就可以把本地库的所有内容推送到远程库上：

    ```bash
    git push -u origin master
    ```

    把本地库的内容推送到远程，用`git push`命令，实际上是把当前分支master推送到远程。

    由于远程库是空的，第一次推送master分支时，加上了-u参数，Git不但会把本地的master分支内容推送到远程新的master分支，还会把本地的master分支和远程的master分支关联起来，在以后的推送或者拉取时就可以简化命令。

6. 推送成功后，可以立刻在GitHub页面中看到远程库的内容已经和本地一模一样。

从现在起，只要本地作了提交，就可以通过命令：

```bash
git push origin master
```

把本地master分支的最新修改推送至GitHub，现在，你就拥有了真正的分布式版本库！

## 将本地和远程库集成

1. 将远程库的变更拉取到本地

    ```bash
    git fetch <remote name> <branch name>
    ```

2. 将拉取的远程库和本地库合并

```bash
git merge <remote name>/<branch name>   # 将历史相关（父子关系的，即两个分支非fast forward）的分支

git merge --allow-unrelated-histories <remote name>/<branch name>   # 将历史不相关的分支合并
# 修改merge message
# :wq!
```

注意：

```bash
git pull = git fetch + git merge
```

## 从远程库克隆

假设我们从零开发，那么最好的方式是先创建远程库，然后，从远程库克隆。

1. 登陆GitHub，创建一个新的仓库，名字叫zero
2. 勾选Initialize this repository with a README，这样GitHub会自动为我们创建一个README.md文件。创建完毕后，可以看到README.md文件
3. 远程库已经准备好了，下一步是用命令git clone克隆一个本地库：

    ```bash
    git clone git@github.com:<your github name>/zero.git        #在页面的右上角可以直接复制该链接
    ```

4. 进入本地的目录就可以看到初始化创建的README.md文件

**如果有多个人协作开发，那么每个人各自从远程克隆一份就可以了**。

GitHub给出的地址不止一个，还可以用 `https://github.com/yourgithubname/zero.git` 这样的地址。

实际上，Git支持多种协议，**默认的git://使用ssh**，但也可以使用https等其他协议。

使用https除了速度慢以外，还有个最大的麻烦是每次推送都必须输入口令，但是在某些只开放http端口的公司内部就无法使用ssh协议而只能用https。
