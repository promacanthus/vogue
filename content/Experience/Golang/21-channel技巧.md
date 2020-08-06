---
title: "21 Channel技巧"
date: 2020-08-06T17:53:42+08:00
draft: true
---

## 赢者为王模式

赢者为王模式：核心思想就是**同时开几个协程做同样的事情**，谁先搞定，就用谁的结果。在Go语言的channel支持下，很容易实现这种并发方式。

### 例子

假设把同一份资源，存储在网络上的5个服务器上（镜像、备份等），现在需要获取这个资源，可以同时开5个协程，访问这5个服务器上的资源，谁先获取到，就用谁的，这样就可以最快速度获取，排除掉网络慢的服务器。

```golang
func main() {
    txtResult := make(chan string, 5)
    go func() {txtResult <- getTxt("res1.flysnow.org")}()
    go func() {txtResult <- getTxt("res2.flysnow.org")}()
    go func() {txtResult <- getTxt("res3.flysnow.org")}()
    go func() {txtResult <- getTxt("res4.flysnow.org")}()
    go func() {txtResult <- getTxt("res5.flysnow.org")}()
    println(<-txtResult)
}

// 模拟函数
func getTxt(host string) string{
    //省略网络访问逻辑，直接返回模拟结果
    //http.Get(host+"/1.txt")
    return host+"：模拟结果"
}
```

这种并发模式适合多个协程做同一件事情，只要有一个协程干成了就OK了。这种模式的优点主要有两个：

1. 最大程度减少耗时
2. 提高成功率

## 最终成功模式

最终成功模式：核心思想是同时并发的从10个文件中成功读取任意5个文件，可以开启5个协程，也可以开启3个，但是**必须成功读取5个才算成功**，否则就是失败。

### 两种常见实现思路

1. 先并发获取，存放起来，然后再一个个判断是否获取成功，如果有的没有成功再重新获取,注意获取的文件不能重复。这种方式是取到结果后进行判断是否成功，然后根据情况再决定是否重新获取,要去重，要判断，业务逻辑比较复杂。

2. 并发的时候就保证成功，里面可能是个for循环，直到成功为止，然后再返回结果。这种思路缺陷也很明显，如果这个文件损坏，那么就会一直死循环下去，要避免死循环，就要加上重试次数。

### 更优的实现思路

使用多个协程，但是发现如果有文件读取不成功，会通过channel的方式标记，换一个文件读取。因为一共10个文件，这个不行，换一个，不能在一个文件上等死，只要成功读取5个就可以了。

实现代码如下：

```golang
// Read reads from readers in parallel. Returns p.dataBlocks number of bufs.
func (p *parallelReader) Read(dst [][]byte) ([][]byte, error) {
    newBuf := dst
    //省略不太相关代码
    var newBufLK sync.RWMutex

    //省略无关
    //channel开始创建，要发挥作用了。这里记住几个数字：
    //readTriggerCh大小是10，p.dataBlocks大小是5
    readTriggerCh := make(chan bool, len(p.readers))
    for i := 0; i < p.dataBlocks; i++ {
        // Setup read triggers for p.dataBlocks number of reads so that it reads in parallel.
        readTriggerCh <- true
    }

    healRequired := int32(0) // Atomic bool flag.
    readerIndex := 0
    var wg sync.WaitGroup
    // readTrigger 为 true, 意味着需要用disk.ReadAt() 读取下一个数据
    // readTrigger 为 false, 意味着读取成功了，不再需要读取
    for readTrigger := range readTriggerCh {
        newBufLK.RLock()
        canDecode := p.canDecode(newBuf)
        newBufLK.RUnlock()
        //判断是否有5个成功的，如果有，退出for循环
        if canDecode {
            break
        }
        //读取次数上限，不能大于10
        if readerIndex == len(p.readers) {
            break
        }
        //成功了，退出本次读取
        if !readTrigger {
            continue
        }
        wg.Add(1)
        //并发读取数据
        go func(i int) {
            defer wg.Done()
            //省略不太相关代码
            _, err := rr.ReadAt(p.buf[bufIdx], p.offset)
            if err != nil {
                //省略不太相关代码
                // 失败了，标记为true，触发下一个读取.
                readTriggerCh <- true
                return
            }
            newBufLK.Lock()
            newBuf[bufIdx] = p.buf[bufIdx]
            newBufLK.Unlock()
            // 成功了，标记为false，不再读取
            readTriggerCh <- false
        }(readerIndex)
        //控制次数，同时用来作为索引获取和存储数据
        readerIndex++
    }
    wg.Wait()

    //最终结果判断，如果OK了就正确返回，如果有失败的，返回error信息。
    if p.canDecode(newBuf) {
        p.offset += p.shardSize
        if healRequired != 0 {
            return newBuf, errHealRequired
        }
        return newBuf, nil
    }

    return nil, errErasureReadQuorum
}
```

前提是从10个数据里读取任意5个。

1. 初始化的chan大小是10，但是通过for循环只存放了5个true
2. 然后对chan循环读取数据，如果是true就开启go协程获取数据，如果是false就终止这次循环
3. 当前在这之前还会判断下是否已经成功获取了5个，如果是的话，直接跳出整个for循环
4. 通过readerIndex每次尝试获取一个数据，如果成功塞一个false到chan中，如果失败则塞个true
5. 这样不成功的readerIndex不再尝试读取，失败了就通过true标记尝试读取下一个readerIndex
6. 通过chan这种巧妙的方式不断循环，直到成功读取5个，或者把10个数据都读一遍为止
7. 最终再基于是否成功读取到5个数据，做最终的判断，是返回成功数据，还是错误

利用channel来做标记和循环取数据，是一种非常好的方式，简化了代码逻辑，整体看起来非常清晰。
