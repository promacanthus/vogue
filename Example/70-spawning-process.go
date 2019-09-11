// 有时候Go程序需要生成其他非Go进程
// 例如一些网站的语法高亮，是通过从Go程序生成pygmentize进程来实现

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	// 一个简单的命令，该命令不带参数或输入，只是将内容输出到stdout
	// exec.Command帮助程序创建一个对象来表示此外部进程
	dateCmd := exec.Command("date")

	// .Output是另一个帮助程序，它处理运行命令，等待命令完成和收集输出的常见情况
	// 如果没有错误，dateOut将保存带有日期信息的字节
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(">date")
	fmt.Println(string(dateOut))

	// 稍微复杂一点的命令，将数据传输到stdin上的外部进程并从stdout收集结果
	grepCmd := exec.Command("grep", "hello")

	// 在这里，显式地获取输入/输出管道，
	// 启动进程，向其写入一些输入，读取输出结果，最后等待进程退出
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()
	grepIn.Write([]byte("hello grep\ngoodbye grep"))
	grepIn.Close()
	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()

	// 在上面的例子中省略了错误检查，可以使用if err != nil 模式来表示错误
	// 使用StdoutPipe收集结果，使用StderrPipe收集错误
	fmt.Println(">grep hello")
	fmt.Println(string(grepBytes))

	// 在生成命令时，需要提供一个显式描述的命令和参数数组，而不是只能传入一个命令行字符串
	// 如果要使用字符串生成完整命令，可以使用bash的-c选项
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(">ls -a -l -h")
	fmt.Println(string(lsOut))
}

// 生成的程序返回的输出与我们直接从命令行运行它们的输出相同。
