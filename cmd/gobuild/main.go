// SPDX-License-Identifier: MIT

// 一个简单的 Go 语言热编译工具
//
// 监视指定目录(可同时监视多个目录)下文件的变化，触发`go build`指令，
// 实时编译指定的 Go 代码，并在编译成功时运行该程序。
// 具体命令格式可使用`gobuild -h`来查看。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/caixw/gobuild"
)

// 与版号相关的变量
var (
	buildDate  string // 由链接器提供此值
	commitHash string // 由链接器提供此值
	version    = "0.9.6"
)

func init() {
	if len(buildDate) > 0 {
		version += ("+" + buildDate)
	}

	if commitHash != "" {
		version += ("." + commitHash)
	}
}

func main() {
	var showHelp, showVersion, showIgnore bool
	opt := &gobuild.Options{}

	flag.BoolVar(&showHelp, "h", false, "显示帮助信息；")
	flag.BoolVar(&showVersion, "v", false, "显示版本号；")
	flag.BoolVar(&opt.Recursive, "r", true, "是否查找子目录；")
	flag.BoolVar(&showIgnore, "i", false, "是否显示被标记为 IGNORE 的日志内容；")
	flag.StringVar(&opt.OutputName, "o", "", "指定输出名称，程序的工作目录随之改变；")
	flag.StringVar(&opt.AppArgs, "x", "", "传递给编译程序的参数；")
	flag.StringVar(&opt.Exts, "ext", "go", "指定监视的文件扩展，区分大小写。* 表示监视所有类型文件，空值代表不监视任何文件；")
	flag.StringVar(&opt.MainFiles, "main", "", "指定需要编译的文件；")
	flag.Usage = usage
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		return
	case showVersion:
		fmt.Fprintln(os.Stdout, "gobuild", version)
		fmt.Fprintln(os.Stdout, "build with", runtime.Version())
		return
	}

	if flag.NArg() == 0 {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		opt.Dirs = []string{wd}
	} else {
		opt.Dirs = flag.Args()
	}

	logs := gobuild.NewConsoleLogs(showIgnore)
	defer logs.Stop()

	if err := gobuild.Build(context.Background(), logs.Logs, opt); err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stdout, `gobuild 是 Go 的热编译工具，监视文件变化，并编译和运行程序。

命令行语法：
 gobuild [options] [dependents]

 options:`)

	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()

	fmt.Fprintln(os.Stdout, `
 dependents:
  指定其它依赖的目录，只能出现在命令的尾部。


常见用法:

 gobuild
   监视当前目录，若有变动，则重新编译当前目录下的 *.go 文件；

 gobuild -main=main.go
   监视当前目录，若有变动，则重新编译当前目录下的 main.go 文件；

 gobuild -main="main.go" dir1 dir2
   监视当前目录及 dir1 和 dir2，若有变动，则重新编译当前目录下的 main.go 文件；


NOTE: 不会监视隐藏文件和隐藏目录下的文件。

源代码采用 MIT 开源许可证，并发布于 https://github.com/caixw/gobuild`)
}
