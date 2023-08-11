// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"io"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild"
	"github.com/caixw/gobuild/watch"
)

const (
	watchTitle      = localeutil.StringPhrase("热编译代码")
	watchUsage      = localeutil.StringPhrase("热编译代码 usage")
	showIgnoreUsage = localeutil.StringPhrase("是否显示被标记为 IGNORE 的日志内容")
)

func initWatch(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("watch", watchTitle.LocaleString(p), watchUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var watchShowIgnore bool
		fs.BoolVar(&watchShowIgnore, "i", false, showIgnoreUsage.LocaleString(p))

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
