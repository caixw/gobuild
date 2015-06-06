gobuild [![Build Status](https://travis-ci.org/caixw/gobuild.svg?branch=master)](https://travis-ci.org/caixw/gobuild)
======

一个简单的Go代码热编译工具。

gobuild会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，
一旦触发，就会调用`go build`编译Go源文件并执行。
```go
 // 监视当前目录下的文件，若发生变化，则触发go build -main="*.go"
 gobuild

 // 监视当前目录和term目录下的文件，若发生变化，则触发go build -main="main.go"
 gobuild -main=main.go ~/Go/src/github.com/issue9/term
```


### 安装

```shell
go get github.com/caixw/gobuild
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/caixw/gobuild)
[![GoDoc](https://godoc.org/github.com/caixw/gobuild?status.svg)](https://godoc.org/github.com/caixw/gobuild)


### 版权

本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
