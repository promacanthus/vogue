---
title: "01 tcpdump"
date: 2020-06-15T11:00:21+08:00
draft: true
---

- [0.1. tcpdump](#01-tcpdump)
  - [0.1.1. 常用参数](#011-常用参数)
  - [0.1.2. 常用过滤表达式](#012-常用过滤表达式)
- [0.2. 案例](#02-案例)

## 0.1. tcpdump

`tcpdump` 仅支持命令行格式使用，常用在 Linux 服务器中抓取和分析网络包。

tcpdump工作的位置：

- 进来的顺序 Wire -> NIC -> `tcpdump` -> netfilter/iptables
- 出去的顺序 iptables -> `tcpdump` -> NIC -> Wire

因此，`iptables`链中的规则会影响到`tcpdump`抓到的包。

```bash
sugoi@sugoi:~$ tcpdump -h
tcpdump version 4.9.3
libpcap version 1.9.1 (with TPACKET_V3)
OpenSSL 1.1.1f  31 Mar 2020
Usage: tcpdump [-aAbdDefhHIJKlLnNOpqStuUvxX#] [ -B size ] [ -c count ]
                [ -C file_size ] [ -E algo:secret ] [ -F file ] [ -G seconds ]
                [ -i interface ] [ -j tstamptype ] [ -M secret ] [ --number ]
                [ -Q in|out|inout ]
                [ -r file ] [ -s snaplen ] [ --time-stamp-precision precision ]
                [ --immediate-mode ] [ -T type ] [ --version ] [ -V file ]
                [ -w file ] [ -W filecount ] [ -y datalinktype ] [ -z postrotate-command ]
                [ -Z user ] [ expression ]
```

### 0.1.1. 常用参数

- `-i`参数指定抓取的网络接口
- `-c`参数指定抓取的网络数据包的数量
- `-w`参数指定抓取的数据包要保存的文件（后缀通常为`.pcap`）
- `-nn`参数不解析IP地址和端口号的名称

### 0.1.2. 常用过滤表达式

|过滤表达式|选项|例子|
|---|---|---|
|主机过滤|host、src host、dst host|tcpdump -nn host 192.168.0.2|
|端口过滤|port、src port、dst port|tcpdump -nn port 80|
|协议过滤|ip、ip6、arp、tcp、udp、icmp|tcpdump -nn tcp|
|逻辑表达式|and、or、not|tcpdump -nn host 192.168.0.2 and port 80|
|特定状态的tcp包|`tcp[tcpflags]`|tcpdump -nn "`tcp[tcpflags] & tcp-syn !=0`"|

## 0.2. 案例

使用`ping`命令来学习`tcpdump`工具的使用。

```bash
sugoi@sugoi:~$ ping -h

Usage
  ping [options] <destination>

Options:
  <destination>      dns name or ip address
  -a                 use audible ping
  -A                 use adaptive ping
  -B                 sticky source address
  -c <count>         stop after <count> replies
  -D                 print timestamps
  -d                 use SO_DEBUG socket option
  -f                 flood ping
  -h                 print help and exit
  -I <interface>     either interface name or address
  -i <interval>      seconds between sending each packet
  -L                 suppress loopback of multicast packets
  -l <preload>       send <preload> number of packages while waiting replies
  -m <mark>          tag the packets going out
  -M <pmtud opt>     define mtu discovery, can be one of <do|dont|want>
  -n                 no dns name resolution
  -O                 report outstanding replies
  -p <pattern>       contents of padding byte
  -q                 quiet output
  -Q <tclass>        use quality of service <tclass> bits
  -s <size>          use <size> as number of data bytes to be sent
  -S <size>          use <size> as SO_SNDBUF socket option value
  -t <ttl>           define time to live
  -U                 print user-to-user latency
  -v                 verbose output
  -V                 print version and exit
  -w <deadline>      reply wait <deadline> in seconds
  -W <timeout>       time to wait for response

IPv4 options:
  -4                 use IPv4
  -b                 allow pinging broadcast
  -R                 record route
  -T <timestamp>     define timestamp, can be one of <tsonly|tsandaddr|tsprespec>

IPv6 options:
  -6                 use IPv6
  -F <flowlabel>     define flow label, default is random
  -N <nodeinfo opt>  use icmp6 node info query, try <help> as argument

For more details see ping(8).
```

- `-I`参数指定发出`ping`命令的网络接口
- `-c`参数指定发出请求的次数

查看网络接口：

```bash
sugoi@sugoi:~$ ifconfig
...

enp0s31f6: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500

...
```

启动`tcpdump`工具开始抓取指定网络接口上对指定IP地址的通信数据包：

```bash
sugoi@sugoi:~$ sudo tcpdump -i enp0s31f6 icmp and host 180.101.49.12 -nn
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on enp0s31f6, link-type EN10MB (Ethernet), capture size 262144 bytes
10:18:23.033206 IP sugoi > 180.101.49.12: ICMP echo request, id 12, seq 1, length 64
10:18:23.040568 IP 180.101.49.12 > sugoi: ICMP echo reply, id 12, seq 1, length 64
10:18:24.035305 IP sugoi > 180.101.49.12: ICMP echo request, id 12, seq 2, length 64
10:18:24.042935 IP 180.101.49.12 > sugoi: ICMP echo reply, id 12, seq 2, length 64
10:18:25.036649 IP sugoi > 180.101.49.12: ICMP echo request, id 12, seq 3, length 64
10:18:25.044844 IP 180.101.49.12 > sugoi: ICMP echo reply, id 12, seq 3, length 64
10:18:26.037946 IP sugoi > 180.101.49.12: ICMP echo request, id 12, seq 4, length 64
10:18:26.046127 IP 180.101.49.12 > sugoi: ICMP echo reply, id 12, seq 4, length 64
```

输出的数据格式为：`时间戳 协议 源地址 > 目的地址 网络包详情`

以百度（180.101.49.12）为例，执行`ping`命令结果如下：

```bash
sugoi@sugoi:~$ ping -I enp0s31f6 -c 4 180.101.49.12
PING 180.101.49.12 (180.101.49.12) from 172.26.160.3 enp0s31f6: 56(84) bytes of data.
64 bytes from 180.101.49.12: icmp_seq=1 ttl=52 time=8.23 ms
64 bytes from 180.101.49.12: icmp_seq=2 ttl=52 time=7.70 ms
64 bytes from 180.101.49.12: icmp_seq=3 ttl=52 time=8.27 ms
64 bytes from 180.101.49.12: icmp_seq=4 ttl=52 time=8.26 ms

--- 180.101.49.12 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3006ms
rtt min/avg/max/mdev = 8.230/8.440/8.625/0.175 ms
```
