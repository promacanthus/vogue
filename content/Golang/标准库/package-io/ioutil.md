---
title: package ioutil
date: 2020-04-14T10:09:14.274627+08:00
draft: false
---

```go
import "io/ioutil"
```

`ioutil`包实现了I/O实用工具。

## 变量

```go
var Discard io.Writer = devNull(0)
```

`Discard`是一个`io.Writer`，在该`Writer`上所有`Write`调用都可以在不执行任何操作的情况下成功完成。

## 函数`NopCloser`

```go
func NopCloser(r io.Reader) io.ReadCloser
```

`NopCloser`返回一个`ReadCloser`，其中包含一个无操作（`no-op`）的`Close`方法，用于包装提供的`Reader r`。

## 函数`ReadAll`

```go
func ReadAll(r io.Reader) ([]byte, error)
```

`ReadAll`从`r`读取，直到出现错误或`EOF`并返回它读取到的数据。成功的调用将返回`err==nil`，而不是`err==EOF`。因为`ReadAll`被定义为从`src`读取直到`EOF`，所以它不会将来自`Read`的`EOF`视为要报告的错误。

```go
r := strings.NewReader("Go is a general-purpose language designed with systems programming in mind.")

b, err := ioutil.ReadAll(r)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("%s", b)

// 输出

Go is a general-purpose language designed with systems programming in mind.
```

## 函数`ReadDir`

```go
func ReadDir(dirname string) ([]os.FileInfo, error)
```

`ReadDir`读取由`dirname`命名的目录，并返回按`filename`排序的目录条目列表。

```go
files, err := ioutil.ReadDir(".")
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    fmt.Println(file.Name())
}
```

## 函数`ReadFile`

```go
func ReadFile(filename string) ([]byte, error)
```

`ReadFile`读取由`filename`命名的文件并返回内容。成功的调用将返回`err==nil`，而不是`err==EOF`。因为`ReadFile`读取整个文件，所以它不会将`Read`中的`EOF`视为要报告的错误。

```go
content, err := ioutil.ReadFile("testdata/hello")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("File contents: %s", content)

// 输出

File contents: Hello, Gophers!
```

## 函数`TempDir`

```go
func TempDir(dir, prefix string) (name string, err error)
```

`TempDir`在目录`dir`中创建一个**新的临时目录**，其名称以`prefix`开头，并返回新目录的路径。如果`dir`是空字符串，`TempDir`将使用临时文件的默认目录（请参阅`os.TempDir`）。同时调用`TempDir`的多个程序将不会选择相同的目录。调用者有责任在不再需要时删除目录。

```go
content := []byte("temporary file's content")
dir, err := ioutil.TempDir("", "example")
if err != nil {
    log.Fatal(err)
}

defer os.RemoveAll(dir) // clean up

tmpfn := filepath.Join(dir, "tmpfile")
if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
    log.Fatal(err)
}
```

## 函数`TempFile`

```go
func TempFile(dir, pattern string) (f *os.File, err error)
```

`TempFile`在目录`dir`中创建一个**新的临时文件**，打开文件进行读写，并返回生成的`*os.File`。文件名是通过获取`pattern`并在末尾添加随机字符串生成的。如果`pattern`包含`“*”`，则随机字符串将替换最后一个`“*”`。如果`dir`是空字符串，则`TempFile`使用临时文件的默认目录（请参阅`os.TempDir`）。同时调用`TempFile`的多个程序不会选择相同的文件。调用者可以使用`f.Name()`来查找文件的路径名。当不再需要时，调用者有责任删除该文件。

```go
content := []byte("temporary file's content")
tmpfile, err := ioutil.TempFile("", "example")
if err != nil {
    log.Fatal(err)
}

defer os.Remove(tmpfile.Name()) // clean up

if _, err := tmpfile.Write(content); err != nil {
    log.Fatal(err)
}
if err := tmpfile.Close(); err != nil {
    log.Fatal(err)
}
```

```go
content := []byte("temporary file's content")
tmpfile, err := ioutil.TempFile("", "example.*.txt")
if err != nil {
    log.Fatal(err)
}

defer os.Remove(tmpfile.Name()) // clean up

if _, err := tmpfile.Write(content); err != nil {
    tmpfile.Close()
    log.Fatal(err)
}
if err := tmpfile.Close(); err != nil {
    log.Fatal(err)
}
```

## 函数`WriteFile`

```go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

`WriteFile`将数据写入由`filename`命名的文件。如果该文件不存在，则`WriteFile`使用权限`perm`创建它;否则`WriteFile`会在写入之前截断它。

```go
message := []byte("Hello, Gophers!")
err := ioutil.WriteFile("testdata/hello", message, 0644)
if err != nil {
    log.Fatal(err)
}
```
