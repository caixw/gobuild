// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild"
)

func initInit(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("init", p.Sprintf("初始化项目"), p.Sprintf("初始化项目"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			var name string
			if fs.NArg() > 0 {
				name = fs.Arg(0)
			} else {
				name = filepath.Base(wd) // 顺序不能乱，要先拿 name！
				wd = filepath.Dir(wd)
			}

			return gobuild.Init(wd, name)
		}
	})
}
