---
title: 04-Git多人单分支集成协作
date: 2020-04-14T10:09:14.298627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 常用工具
- Git
summary: 04-Git多人单分支集成协作
showInMenu: false

---

## 不同人修改不用文件

```bash
# 用户一的操作

git branch -av
# 查看仓库中全部的分支（包括本地和远程）

git checkout -b feature/add_git_commands origin/feature/add_git_commands
#  基于远程仓库的分支origin创建本地分支，并切换到新建的分支上

git branch -v
# 只查看本地分支

git push <远程仓库名称，默认是origin>/<分支的名称>
# 将当前分支推送到远程仓库中对应的分支
# 因为创建时有本地分支和远程仓库中的分支有关联关系，上述命令可以缺省远程仓库中分支的选择


# 用户二的操作

git fetch <远程仓库的名字，默认为origin>
# 将远程仓库中分支的变更都同步到本地仓库

git merge <远程仓库名>/<远程仓库分支名>
# 将本地当前分支与远程仓库中的指定分支合并
# 填写merge信息，然后 :wq!

git push
```

## 不同人修改同文件不同区域

开发之前先同步远程仓库

```bash
# 操作方式一

git pull
# 同步远程仓库中的代码，并且与本地分支进行合并



#  操作方式二

git fetch
# 将本地仓库中远程分支与远程仓库中对应的分支进行同步

git merge <远程仓库名>/<远程仓库分支名>
# 将本地当前分支与本地仓库的远程分支合并
# 填写merge信息，然后 :wq!
```

## 不同人修改同文件同区域

```bash
<<<<<
# 这里是本地的文件内容
====
# 这里是远程仓库的文件内容
>>>>>

# 直接在文件上进行修改，将需要的内容保留，不需要的内容删除

git commit
# 提交解决的冲突

git push
# 将解决后的commit推送到远程仓库中

git merge --abort
# 放弃解决的冲突
```

## 同时变更文件名和文件内容

git存放blob文件时是以文件内容来区分的，并不以文件名来区分；此处的变更文件名操作和变更文件内容的操作能够自动被git处理，原因就在于blob文件并没有发生修改的冲突。

如果既变更了文件名又修改了文件，同时另一个人也修改了该文件的同一位置的内容，就会被git识别为冲突，而不能自动进行处理了。

```bash
git  mv file1 file2

git pull
# git能够智能感知到文件名的变化，因为是git是基于文件内容区分的
```

## 同时修改同一文件名

git会报告冲突，自行判断要使用哪个文件名，然后在提交。

```bash
git add  <需要的filename>
git rm <不需要的filename>
```

## 禁止向集成分支执行push -f操作

```bash
git push -f
# 强制更新
```

## 禁止向集成分支执行变更历史操作

公共分支禁止rebase。
