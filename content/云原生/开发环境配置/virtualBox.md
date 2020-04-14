---
title: virtualBox
date: 2020-04-14T10:09:14.226627+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- 云原生
- 开发环境配置
summary: virtualBox
showInMenu: false

---

## 安装

将以下行添加到您的`/etc/apt/sources.list`中，将`<mydist>`替换为对应的发行版代号，18.04为bionic

```bash
echo "deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian bionic contrib" >> /etc/apt/sources.list

wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -

sudo apt-get update
sudo apt-get install virtualbox
```

遇到以下签名无效时该怎么办：BADSIG ...从存储库刷新软件包时？

```bash
sudo -s -H
apt-get clean
rm /var/lib/apt/lists/*
rm /var/lib/apt/lists/partial/*
apt-get clean
apt-get update
```
