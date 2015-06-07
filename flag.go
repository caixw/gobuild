// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/issue9/term/colors"
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

type args struct {
	outputName string
	mainFiles  string
	exts       []string // 被忽略的扩展名
	paths      []string
}

// 初始化flag相关功能，包括flag.Parse()。
// 分析命令行参数，并将一些简单参数进行处理，比如-v,-h等。
func parseFlag() *args {
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
		return nil
	}

	if showVersion {
		colors.Print(colors.Stdout, colors.Green, colors.Default, "gobuild: ")
		colors.Println(colors.Stdout, colors.Default, colors.Default, version)

		colors.Print(colors.Stdout, colors.Green, colors.Default, "Go: ")
		goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
		colors.Println(colors.Stdout, colors.Default, colors.Default, goVersion)

		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		log(erro, "获取当前工作目录时，发生以下错误:", err)
		os.Exit(2)
	}
	// 确定编译后的文件名
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}

	return &args{
		outputName: outputName,
		mainFiles:  mainFiles,
		exts:       strings.Split(extString, ","),
		paths:      append(flag.Args(), wd),
	}
}
