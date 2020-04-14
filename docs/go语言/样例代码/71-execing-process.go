---
title: 71-execing-process.go
date: 2019-11-25T11:15:47.534182+08:00
draft: false
hideLastModified: false
summaryImage: ""
keepImageRatio: true
tags:
- ""
- Go语言
- 样例代码
summary: 71-execing-process.go
showInMenu: false

---

// 上一个例子生成一个外部的进程，当我们需要一个正在运行的Go进程可访问的外部进程时，会这样做
// 有时只想用另一个（也许是非Go进程）替换当前的Go进程。为此，将使用Go的经典exec函数来实现

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// 下面的示例将执行ls命令
	// Go需要一个想要执行的二进制文件的绝对路径
	// 使用exec.LookPath函数来找打它
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	// Exec需要切片形式的参数（与一个大字符串相对应）
	// 第一个参数应该是程序名称
	args := []string{"ls", "-a", "-l", "-h"}

	// Exec还需要一组环境变量才能使用。在这里，只提供当前的环境变量
	env := os.Environ()

	// 这是实际的syscall.Exec调用
	// 如果此调用成功，那么进程的执行将在此处结束，并由/bin/ls -a -l -h进程替换
	// 如果有错误，将获得返回值
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

// 请注意，Go不提供经典的Unix fork函数
// 通常这不是问题，因为启动goroutine，产生进程和exec'ing进程涵盖了fork的大多数用例
