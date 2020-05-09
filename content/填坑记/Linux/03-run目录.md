---
title: "03 Run目录"
date: 2020-05-09T11:35:30+08:00
draft: true
---

Linux 系统在运行时数据方面的工作方式有一些小但重大的变化。它重新组织了文件系统中可访问的方式和位置，这在 Linux 文件系统中提供了更多一致性。

查看`/run`路径：

```bash
sugoi@sugoi:/run$ df -k .
文件系统         1K-块      已用    可用        已用%   挂载点
tmpfs          1168928    1908   1167020      1%     /run
```

被识别为`tmpfs`临时文件系统，即`/run`中的文件和目录是直接存储在内存中。

> tmpfs表示保存在内存或基于磁盘的交换分区中的数据，看起来像已挂载的文件系统，这样可以使其更易于访问和管理。

`/run`中保存了各种数据，每个目录都是运行中的进程所使用的文件，如：

- `×.pid`的各种系统进程ID
- 为了与 `/run` 的变化保持一致，一些运行时数据的旧位置现在是符号链接。`/var/run` 现在是指向 `/run` 的指针，`/var/lock` 指向 `/run/lock` 的指针，可以保证旧的引用按预期工作。

## 命令介绍

```bash
df --help
用法：df [选项]... [文件]...
Show information about the file system on which each FILE resides,
or all file systems by default.
...

  -k                    即--block-size=1K
...

```
