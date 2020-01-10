# 66-cmd-line-subcmd

```go
// 一些命令行工具，比如go工具或git有许多子命令，每个子命令都有自己的一组标志。
// 例如： go build 和 go get 是go工具的两个不同子命令
// 使用flag包可以轻松定义具有自己标志的简单子命令

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 使用flag.NewFlagSet函数声明子命令
	// 然后继续定义特定一该子命令的命令行标识
	fooCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	fooEnable := fooCmd.Bool("enable", false, "enable")
	fooName := fooCmd.String("name", "", "name")

	// 对于不同的子命令可以定义不同的命令行标识
	barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	barLevel := barCmd.Int("level", 0, "level")

	// 子命令应该是程序的第一个参数
	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands'")
		os.Exit(1)
	}

	// 检查调用了哪个子命令
	switch os.Args[1] {
	// 对于每个子命令，解析自己的命令行标识并访问尾随的位置参数
	case "foo":
		fooCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'foo'")
		fmt.Println("enable:", *fooEnable)
		fmt.Println("name:", *fooName)
		fmt.Println("tail:", fooCmd.Args())
	case "bar":
		barCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'bar'")
		fmt.Println("level:", *barLevel)
		fmt.Println("tail:", barCmd.Args())
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}
}

```