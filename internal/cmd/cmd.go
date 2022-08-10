// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild"
	"github.com/caixw/gobuild/locales"
	"github.com/caixw/gobuild/log"
)

// 与版号相关的变量
var (
	buildDate  string // 由链接器提供此值
	commitHash string // 由链接器提供此值
	version    = "0.10.0"
)

func init() {
	if len(buildDate) > 0 {
		version += ("+" + buildDate)
	}

	if commitHash != "" {
		version += ("." + commitHash)
	}
}

func Exec() {
	p := getPrinter()

	var showHelp, showVersion, showIgnore bool
	var exts string
	var freq int
	opt := &gobuild.WatchOptions{Printer: p}

	flag.BoolVar(&showHelp, "h", false, p.Sprintf("显示帮助信息"))
	flag.BoolVar(&showVersion, "v", false, p.Sprintf("显示版本号"))
	flag.BoolVar(&opt.Recursive, "r", true, p.Sprintf("是否查找子目录"))
	flag.BoolVar(&showIgnore, "i", false, p.Sprintf("是否显示被标记为 IGNORE 的日志内容"))
	flag.StringVar(&opt.OutputName, "o", "", p.Sprintf("指定输出名称，程序的工作目录随之改变"))
	flag.StringVar(&opt.AppArgs, "x", "", p.Sprintf("传递给编译程序的参数"))
	flag.IntVar(&freq, "freq", 1, p.Sprintf("监视器的更新频率"))
	flag.StringVar(&exts, "ext", "go", p.Sprintf("指定监视的文件扩展，区分大小写"))
	flag.StringVar(&opt.MainFiles, "main", "", p.Sprintf("指定需要编译的文件"))
	flag.Usage = func() {
		bs := &bytes.Buffer{}
		flag.CommandLine.SetOutput(bs)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stdout, p.Sprintf(usage, bs.String()))
	}
	flag.Parse()
	opt.Exts = strings.Split(exts, ",")
	opt.WatcherFrequency = time.Duration(freq) * time.Second

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

	logs := log.NewConsole(showIgnore)
	defer logs.Stop()

	if err := gobuild.Watch(context.Background(), logs.Logs, opt); err != nil {
		panic(err)
	}
}

func getPrinter() *message.Printer {
	tag, _ := localeutil.DetectUserLanguageTag()
	c := catalog.NewBuilder(catalog.Fallback(tag))
	if err := localeutil.LoadMessageFromFSGlob(c, &localeFS, "*.yaml", yaml.Unmarshal); err != nil {
		panic(err)
	}
	if err := localeutil.LoadMessageFromFSGlob(c, locales.Locales, "*.yaml", yaml.Unmarshal); err != nil {
		panic(err)
	}
	return message.NewPrinter(tag, message.Catalog(c))
}
