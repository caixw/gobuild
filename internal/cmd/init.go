// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/caixw/gobuild"
	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"
)

var initFS *flag.FlagSet

func initInit(o *cmdopt.CmdOpt, p *message.Printer) {
	initFS = o.New("init", p.Sprintf("初始化项目"), doInit(p))
}

func doInit(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		var name string
		if initFS.NArg() > 0 {
			name = initFS.Arg(0)
		} else {
			name = filepath.Base(wd) // 顺序不能乱，要先拿 name！
			wd = filepath.Dir(wd)
		}

		return gobuild.Init(wd, name)
	}
}
