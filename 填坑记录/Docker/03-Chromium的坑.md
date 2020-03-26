# 03-Chromium的坑

制作镜像的时候，在里面安装了一个chromium浏览器，但是容易出现崩溃的情况。

## 安装

以Ubuntu基础镜像为例。

```dockerfile
FROM ubuntu:18.04
RUN apt-get update \
    && apt-get install -y --allow-unauthenticated chromium-browser chromium-browser-l10n chromium-codecs-ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && ln -s /usr/bin/chromium-browser /usr/bin/google-chrome-stable \
    && mkdir $HOME \
    && echo "CHROMIUM_FLAGS='--no-first-run --no-sandbox --start-maximized --user-data-dir --disable-software-rasterizer --disable-dev-shm-usage'" > $HOME/.chromium-browser.init
```

## 报错

通过Terminal运行发现报错信息如下。

```bash
[2576:2605:0310/115342.554370:ERROR:zygote_host_impl_linux.cc(259)] Failed to adjust OOM score of renderer with pid 2644: Permission denied (13)
[2765:2771:0310/115428.402995:FATAL:memory.cc(22)] Out of memory. size=262144
...
```

需要增加`CHROMIUM_FLAGS='--no-first-run --no-sandbox --start-maximized --user-data-dir --disable-software-rasterizer --disable-dev-shm-usage'`。

## 参数详解

全部的参数在[这里](https://peter.sh/experiments/chromium-command-line-switches/)。

- `--no-first-run`：跳过“首次运行”任务，无论它实际上是否是“首次运行”。 会被`kForceFirstRun`参数覆盖。这不会删除“首次运行”步骤，因此也不能防止在没有此标志的情况下下次启动chrome时发生首次运行。
- `--no-sandbox`：对通常为沙盒的所有进程类型禁用沙盒。
- `--start-maximized`：无论以前的任何设置如何，都以最大化（全屏）的方式启动浏览器。
- `--user-data-dir`：浏览器存储用户配置文件的目录。
- `--disable-software-rasterizer`：禁止使用3D软件光栅化器。
- `--disable-dev-shm-usage`：在某些VM环境中，`/dev/shm`分区太小，导致Chrome发生故障或崩溃（请[参阅](http://crbug.com/715363)）。 使用此标志解决此问题（临时目录将始终用于创建匿名共享内存文件）。

### `/dev/shm`

`/dev/shm`是Linux下一个非常有用的目录，这个目录不在硬盘上，而是在内存里。因此在Linux下，就不需要大费周折去建 `ramdisk`，直接使用`/dev/shm/`就可达到很好的优化效果。

`/dev/shm/`需要注意的一个是容量问题，在Linux下，它默认最大为内存的一半大小，使用`df -h`命令可以看到。

但它并不会真正的占用这块内存:

- 如果`/dev/shm/`下没有任何文件，它占用的内存实际上就是0字节；
- 如果它最大为1G，里头放有100M文件，那剩余的900M仍然可为其它应用程序所使用，但它所占用的100M内存，是绝不会被系统回收重新划分的。

默认系统就会加载`/dev/shm` ，它就是所谓的`tmpfs`，这与ramdisk（虚拟磁盘）不一样。`tmpfs` 可以使用 RAM，也可以使用交换分区来存储。

> 传统的虚拟磁盘是个块设备，并需要一个 `mkfs` 之类的命令才能真正地使用它，`tmpfs` 是一个文件系统，而不是块设备，只是安装它，就可以使用。

默认的最大一半内存大小在某些场合可能不够用，并且默认的inode数量很低一般都要调高些，这时可以用mount命令来修改`/dev/shm`。

```bash
mount -o size=1500M -o nr_inodes=1000000 -o noatime,nodiratime -o remount /dev/shm

# 在/etc/fstab中增加配置
tmpfs /dev/shm tmpfs defaults,size=1.5G 0 0
```
