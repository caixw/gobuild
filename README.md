gobuild [![Build Status](https://travis-ci.org/caixw/gobuild.svg?branch=master)](https://travis-ci.org/caixw/gobuild)
======

gobuild是一个简单的Go代码热编译工具。
会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，
一旦触发，就会调用`go build`编译Go源文件并执行。


命令行语法:
```shell
 gobuild [options] [dependents]
```

```txt
options:
 -h    显示当前帮助信息；
 -v    显示gobuild和go程序的版本信息；
 -o    执行编译后的可执行文件名；
 -r    是否搜索子目录，默认为true；
 -ext  需要监视的扩展名，默认值为"go"，区分大小写，会去掉每个扩展名的首尾空格。
       若需要监视所有类型文件，请使用*，传递空值代表不监视任何文件；
 -main 指定需要编译的文件，默认为""。
```

```txt
dependents:
 指定其它依赖的目录，只能出现在命令的尾部。
```


常见用法:
```go
 // 监视当前目录下的文件，若发生变化，则触发go build -main="*.go"
 gobuild

 // 监视当前目录和~/Go/src/github.com/issue9/term目录下的文件，
 // 若发生变化，则触发go build -main="main.go"
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
