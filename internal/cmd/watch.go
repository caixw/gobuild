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

func initWatch(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("watch", p.Sprintf("热编译代码"), p.Sprintf("热编译代码 usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var watchShowIgnore bool
		fs.BoolVar(&watchShowIgnore, "i", false, p.Sprintf("是否显示被标记为 IGNORE 的日志内容"))

		return func(w io.Writer) error {
			logs := watch.NewConsoleLogger(watchShowIgnore, os.Stderr, os.Stdout)

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			return gobuild.WatchConfig(wd, p, logs)
		}
	})
}
