// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 一个简单的Go语言热编译工具。
//
// gobuild会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，
// 一旦触发，就会调用`go build`编译Go源文件并执行。
//  // 监视当前目录下的文件，若发生变化，则触发go build -main="*.go"
//  gobuild
//
//  // 监视当前目录和term目录下的文件，若发生变化，则触发go build -main="main.go"
//  gobuild -main=main.go ~/Go/src/github.com/issue9/term
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/issue9/term/colors"
)

// 当前程序的版本号
const version = "0.2.4.150607"

const usage = `gobuild 用于热编译Go程序。
 
命令行语法:
 gobuild [options] [dependents]
 
 options:
  -h    显示当前帮助信息；
  -v    显示gobuild和go程序的版本信息；
  -o    执行编译后的可执行文件名；
  -ext  需要监视的扩展名，默认值为"go"，区分大小写。
        若需要监视所有类型文件，请使用*，传递空值代表不监视任何文件；
  -main 指定需要编译的文件，默认为""。
 
 dependents:
  指定其它依赖的目录，只能出现在命令的尾部。
 
 
常见用法:
 
 gobuild
   监视当前目录，若有变动，则重新编译当前目录下的*.go文件；
 
 gobuild -main=main.go
   监视当前目录，若有变动，则重新编译当前目录下的main.go文件；
 
 gobuild -main="main.go" dir1 dir2
   监视当前目录及dir1和dir2，若有变动，则重新编译当前目录下的main.go文件；
`

func main() {
	// 检测基本环境是否满足
	if gopath := os.Getenv("GOPATH"); len(gopath) == 0 {
		log(erro, "未设置环境变量GOPATH")
		return
	}

	// 初始化flag
	showHelp := false
	showVersion := false
	var mainFiles, outputName, extString string

	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&outputName, "o", "", "指定输出名称")
	flag.StringVar(&extString, "ext", "go", "指定监视的文件扩展名")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")
	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		colors.Print(colors.Stdout, colors.Green, colors.Default, "gobuild: ")
		colors.Println(colors.Stdout, colors.Default, colors.Default, version)
		colors.Print(colors.Stdout, colors.Green, colors.Default, "Go: ")
		goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
		colors.Println(colors.Stdout, colors.Default, colors.Default, goVersion)
		return
	}

	if len(extString) == 0 { // 允许不监视任意文件，便输出一信息来警告
		log(warn, "将ext设置为空值，意味着不监视任何文件的改变，这将没有任何意义！")
	}

	b := newBuilder(mainFiles, outputName, strings.Split(extString, ","), flag.Args())
	b.build()

	done := make(chan bool)
	<-done
}
