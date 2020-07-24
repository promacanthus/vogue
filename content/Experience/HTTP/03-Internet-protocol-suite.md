---
title: "03 Internet Protocol Suite"
date: 2020-07-24T10:00:26+08:00
draft: true
---

## 网络拓扑和数据流

![image](/images/350px-IP_stack_connections.svg.png)

数据流的定义是在一个简单的网络拓扑中的两台主机（A和B）通过各自路由之间的链路相连。每个主机上的应用程序都执行读取和写入操作，就好像这些进程通过某种数据管道直接相互连接。建立此管道后，由于底层的通信原理是在较低的协议层中实现的，因此每个进程都将隐藏大多数通信细节。

以此类推：

- 在传输层，通信表现为主机到主机，而无需了解应用程序数据结构和连接的路由器
- 在网络层，则在每个路由器处遍历各个网络边界

![image](/images/1920px-UDP_encapsulation.svg.png)

通过RFC 1122中对各层的描述来封装应用程序数据并向下传输。

## 各层协议

### 应用层

- 9P, Plan 9 from Bell Labs distributed file system protocol
- AFP, Apple Filing Protocol
- APPC, Advanced Program-to-Program Communication
- AMQP, Advanced Message Queuing Protocol
- Atom Publishing Protocol
- BEEP, Block Extensible Exchange Protocol
- Bitcoin
- BitTorrent
- CFDP, Coherent File Distribution Protocol
- CoAP, Constrained Application Protocol
- DDS, Data Distribution Service
- DeviceNet
- eDonkey
- ENRP, Endpoint Handlespace Redundancy Protocol
- FastTrack (KaZaa, Grokster, iMesh)
- Finger, User Information Protocol
- Freenet
- FTAM, File Transfer Access and Management
- Gopher, Gopher protocol
- HL7, Health Level Seven
- HTTP, Hypertext Transfer Protocol
- H.323, Packet-Based Multimedia Communications System
- IMAP, Internet Message Access Protocol
- IRC, Internet Relay Chat
- IPFS, InterPlanetary File System
- Kademlia
- LDAP, Lightweight Directory Access Protocol
- LPD, Line Printer Daemon Protocol
- MIME (S-MIME), Multipurpose Internet Mail Extensions and Secure MIME
- Modbus
- MQTT Protocol
- Netconf
- NFS, Network File System
- NIS, Network Information Service
- NNTP, Network News Transfer Protocol
- NTCIP, National Transportation Communications for Intelligent Transportation System Protocol
- NTP, Network Time Protocol
- OSCAR, AOL Instant Messenger Protocol
- POP, Post Office Protocol
- PNRP, Peer Name Resolution Protocol
- RDP, Remote Desktop Protocol
- RELP, Reliable Event Logging Protocol
- RFP, Remote Framebuffer Protocol
- Rlogin, Remote Login in UNIX Systems
- RPC, Remote Procedure Call
- RTMP, Real Time Messaging Protocol
- RTP, Real-time Transport Protocol
- RTPS, Real Time Publish Subscribe
- RTSP, Real Time Streaming Protocol
- SAP, Session Announcement Protocol
- SDP, Session Description Protocol
- SIP, Session Initiation Protocol
- SLP, Service Location Protocol
- SMB, Server Message Block
- SMTP, Simple Mail Transfer Protocol
- SNTP, Simple Network Time Protocol
- SSH, Secure Shell
- SSMS, Secure SMS Messaging Protocol
- TCAP, Transaction Capabilities Application Part
- TDS, Tabular Data Stream
- Tor (anonymity network)
- Tox
- TSP, Time Stamp Protocol
- VTP, Virtual Terminal Protocol
- Whois (and RWhois), Remote Directory Access Protocol
- WebDAV
- X.400, Message Handling Service Protocol
- X.500, Directory Access Protocol (DAP)
- XMPP, Extensible Messaging and Presence Protocol
- Z39.50
- DNS, Domain Name Services

### 传输层

- ATP, AppleTalk Transaction Protocol
- CUDP, Cyclic UDP
- DCCP, Datagram Congestion Control Protocol
- FCP, Fibre Channel Protocol
- IL, IL Protocol
- MPTCP, Multipath TCP
- RDP, Reliable Data Protocol
- RUDP, Reliable User Datagram Protocol
- SCTP, Stream Control Transmission Protocol
- SPX, Sequenced Packet Exchange
- SST, Structured Stream Transport
- TCP, Transmission Control Protocol
- UDP, User Datagram Protocol
- UDP-Lite
- µTP, Micro Transport Protocol
- RSVP

### 网络层

- Anti-replay
- Gateway-to-Gateway Protocol
- Internet Control Message Protocol
- Internet Control Message Protocol for IPv6
- Internet Group Management Protocol
- Internet Group Management Protocol with Access Control
- Internet Protocol
- IPv4
- IPv6
- Locator/Identifier Separation Protocol
- Seamoby
- SwIPe (protocol)
- ECN
- IPsec

### 链路层

- Address Resolution Protocol (ARP)
- Reverse Address Resolution Protocol (RARP)
- Neighbor Discovery Protocol (NDP) as ARP for IPv6
- IS-IS (RFC 1142) is another link-state routing protocol
- Open Shortest Path First (OSPF)
- Tunnels（L2TP）
- PPP
- MAC （Ethernet Wi-Fi DSL ISDN FDDI）
