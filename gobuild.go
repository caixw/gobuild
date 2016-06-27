// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 一个简单的 Go 语言热编译工具。
//
// 监视指定目录(可同时监视多个目录)下文件的变化，触发`go build`指令，
// 实时编译指定的 Go 代码，并在编译成功时运行该程序。
// 具体命令格式可使用`gobuild -v`来查看。
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

// 当前程序的版本号
const version = "0.5.17+20160627"

const usage = `gobuild 是 Go 的热编译工具，监视文件变化，并编译和运行程序。

命令行语法:
 gobuild [options] [dependents]

 options:
  -h    显示当前帮助信息；
  -v    显示 gobuild 和 Go 程序的版本信息；
  -o    编译后的可执行文件名，程序的工作目录随之改变；
  -r    是否搜索子目录，默认为 true；
  -i    是否显示被标记为 IGNORE 的日志内容，默认为 false，即不显示；
  -ext  需要监视的扩展名，默认值为"go"，区分大小写，会去掉每个扩展名的首尾空格。
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


NOTE: 不会监视隐藏文件和隐藏目录下的文件。


源代码采用 MIT 开源许可证，并发布于 https://github.com/caixw/gobuild
`

func main() {
	// 检测基本环境是否满足
	if gopath := os.Getenv("GOPATH"); len(gopath) == 0 {
		log(erro, "未设置环境变量GOPATH")
		return
	}

	// 初始化flag
	var showHelp, showVersion, recursive bool
	var mainFiles, outputName, extString string

	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.BoolVar(&recursive, "r", true, "是否查找子目录")
	flag.BoolVar(&showIgnoreLog, "i", false, "是否显示被标记为IGNORE的日志内容")
	flag.StringVar(&outputName, "o", "", "指定输出名称")
	flag.StringVar(&extString, "ext", "go", "指定监视的文件扩展名")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		return
	case showVersion:
		printVersion()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log(erro, "获取当前工作目录时，发生以下错误:", err)
		return
	}

	// 初始化goCmd的参数
	args := []string{"build", "-o", outputName}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}

	b := &builder{
		exts:      getExts(extString),
		appName:   getAppName(outputName, wd),
		goCmdArgs: args,
	}

	b.watch(recursivePaths(recursive, append(flag.Args(), wd)))
	go b.build()

	done := make(chan bool)
	<-done
}

// 输出版本号
func printVersion() {
	colors.Print(colors.Stdout, colors.Green, colors.Default, "gobuild: ")
	colors.Println(colors.Stdout, colors.Default, colors.Default, version)
	colors.Print(colors.Stdout, colors.Green, colors.Default, "Go: ")
	goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	colors.Println(colors.Stdout, colors.Default, colors.Default, goVersion)
}

// 根据 recursive 值确定是否递归查找 paths 每个目录下的子目录。
func recursivePaths(recursive bool, paths []string) []string {
	if !recursive {
		return paths
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log(erro, "在遍历监视目录时，发生以下错误:", err)
		}

		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			log(erro, "在遍历监视目录时，发生以下错误:", err)
		}
	}

	return ret
}

// 将 extString 分解成数组，并清理掉无用的内容，比如空字符串
func getExts(extString string) []string {
	exts := strings.Split(extString, ",")
	ret := make([]string, 0, len(exts))

	for _, ext := range exts {
		ext = strings.TrimSpace(ext)

		if len(ext) == 0 {
			continue
		}
		if ext[0] != '.' {
			ext = "." + ext
		}
		ret = append(ret, ext)
	}

	switch {
	case len(ret) == 0: // 允许不监视任意文件，但输出一信息来警告
		log(warn, "将ext设置为空值，意味着不监视任何文件的改变！")
	case len(ret) > 0:
		log(info, "系统将监视以下类型的文件:", ret)
	}

	return ret
}

func getAppName(outputName, wd string) string {
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}

	// 转成绝对路径
	outputName, err := filepath.Abs(outputName)
	if err != nil {
		log(erro, err)
	}

	log(info, "输出文件为:", outputName)

	return outputName
}
