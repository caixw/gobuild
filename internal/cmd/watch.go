// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild"
	"github.com/caixw/gobuild/watch"
)

var (
	watchFS         *flag.FlagSet
	watchShowIgnore bool
)

func initWatch(o *cmdopt.CmdOpt, p *message.Printer) {
	watchFS = o.New("watch", p.Sprintf("热编译代码"), doWatch(p))
	watchFS.BoolVar(&watchShowIgnore, "i", false, p.Sprintf("是否显示被标记为 IGNORE 的日志内容"))
}

func doWatch(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		logs := watch.NewConsoleLogger(watchShowIgnore, os.Stderr, os.Stdout)

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		return gobuild.WatchConfig(wd, p, logs)
	}
}
