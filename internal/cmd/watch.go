// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"flag"
	"io"
	"os"
	"strings"
	"time"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild"
	"github.com/caixw/gobuild/log"
)

var (
	watchFS *flag.FlagSet

	showIgnore bool
	exts       string
	freq       int
	opt        = &gobuild.WatchOptions{}
)

func initWatch(o *cmdopt.CmdOpt, p *message.Printer) {
	watchFS = o.New("watch", p.Sprintf("热编译代码"), doWatch(p))

	watchFS.BoolVar(&opt.Recursive, "r", true, p.Sprintf("是否查找子目录"))
	watchFS.BoolVar(&showIgnore, "i", false, p.Sprintf("是否显示被标记为 IGNORE 的日志内容"))
	watchFS.StringVar(&opt.OutputName, "o", "", p.Sprintf("指定输出名称，程序的工作目录随之改变"))
	watchFS.StringVar(&opt.AppArgs, "x", "", p.Sprintf("传递给编译程序的参数"))
	watchFS.IntVar(&freq, "freq", 1, p.Sprintf("监视器的更新频率"))
	watchFS.StringVar(&exts, "ext", "go", p.Sprintf("指定监视的文件扩展，区分大小写"))
	watchFS.StringVar(&opt.MainFiles, "main", "", p.Sprintf("指定需要编译的文件"))
	opt.Exts = strings.Split(exts, ",")
	opt.WatcherFrequency = time.Duration(freq) * time.Second
	opt.Printer = p

}

func doWatch(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if flag.NArg() == 0 {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			opt.Dirs = []string{wd}
		} else {
			opt.Dirs = flag.Args()
		}

		logs := log.NewConsole(showIgnore)
		defer logs.Stop()

		return gobuild.Watch(context.Background(), logs.Logs, opt)
	}
}
