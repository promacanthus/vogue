# 02-Module的坑

## 1. 多版本依赖冲突问题

修改项目中的`go.mod`文件，将冲突的版本进行指定替换，如下所示：

```go
module example

go 1.13

require(
    github.com/ugorji/go v1.1.7 // indirect
    ...
)

replace (
    github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
)
```

## 2. 导入本地模块问题

本地项目目录结构如下，在一个项目中有两个modules，根据两个`go.mod`文件的位置即可确定，分别是mod-a和mod-b。

```bash
mod-a/
├── mod-b
│   ├── go.mod
│   ├── go.sum
│   ├── pkg
│   └── main.go
├── go.mod
├── go.sum
└── main.go
```

现在要在mod-a中导入mod-b，那么只需要修改mod-a的`go.mod`文件即可，具体如下：

```go
module mod-a

go 1.13

require (
    mod-b v0.0.0-00010101000000-000000000000
    ...
)

replace (
    mod-b => ./mod-b
)
```
