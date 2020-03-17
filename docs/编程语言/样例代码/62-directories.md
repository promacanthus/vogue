# 62-directories

```go
// Go有几个有用的功能来处理文件系统中的目录。

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 在当前工作目录下创建一个新的子目录
	err := os.Mkdir("subdir", 0755)
	check(err)

	// 在创建临时目录时，最好将其删除
	// os.RemoveAll将会删除整个目录树，类似 rm -rf
	defer os.RemoveAll("subdir")

	// 辅助函数用于创建一个新的空文件
	createEmptyFile := func(name string) {
		d := []byte("")
		check(ioutil.WriteFile(name, d, 0644))
	}

	createEmptyFile("subdir/file1")
	err = os.MkdirAll("subdir/parent/child", 0755)
	check(err)

	createEmptyFile("subdir/parent/file2")
	createEmptyFile("subdir/parent/file3")
	createEmptyFile("subdir/parent/child/file4")

	// ReadDir函数会列出目录的内容，返回一个os.FileInfo对象组成的切片
	c, err := ioutil.ReadDir("subdir/parent")
	check(err)

	fmt.Println("Listing subdir/parent")
	for _, entry := range c {
		fmt.Println("", entry.Name(), entry.IsDir())
	}

	// Chdir用于切换当前工作目录，类型与cd命令
	err = os.Chdir("subdir/parent/child")
	check(err)

	// 在列出当前目录时看到子目录 subdir/parent/child 的内容
	c, err = ioutil.ReadDir(".")
	check(err)

	fmt.Println("Listing subdir/parent/child")
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// 回到开始的路径
	err = os.Chdir("../../..")
	check(err)

	// 可以递归的访问目录，包括其中所以的子目录
	// Walk函数接收回调函数来处理每个文件和目录的访问
	fmt.Println("Visiting subdir")
	err = filepath.Walk("subdir", visit)
}

func visit(p string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fmt.Println("", p, info.IsDir())
	return nil
}

```