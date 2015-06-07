// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
)

const usage = `gobuild 用于热编译Go程序。
 
用法:
 gobuild [options] [dependents]
 
 options:
  -h    显示当前帮助信息；
  -v    显示gobuild和go程序的版本信息；
  -o    执行编译后的可执行文件名；
  -ext  监视的扩展名，默认值为"go"，区分大小写，若需要监视所有类型文件，请使用*；
  -main 指定需要编译的文件，默认为""。
 
 dependents:
  指定其它依赖的目录，只能出现在命令的尾部。
 
 
常见用法:
 
 gobuild
   监视当前目录，若有变动，则重新编译当前目录下的*.go文件；
 
 gobuild -main=main.go
   监视当前目录，若有变动，则重新编译当前目录下的main.go文件；
 
 gobuild dir1 dir2
   监视当前目录及dir1和dir2，若有变动，则重新编译当前目录下的*.go文件；
`

var (
	showHelp    = false
	showVersion = false
	mainFiles   = ""
	outputName  = ""
	extString   = "go"
)

// 初始化flag相关功能，包括flag.Parse()。
func initFlag() {
	// 初始化flag相关设置
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&outputName, "o", "", "指定输出名称")
	flag.StringVar(&extString, "ext", "go", "指定监视的文件扩展名")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()
}
