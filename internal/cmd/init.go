// SPDX-License-Identifier: MIT

package cmd

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	i "github.com/caixw/gobuild/internal/init"
)

var initFS *flag.FlagSet

func initInit(o *cmdopt.CmdOpt, p *message.Printer) {
	initFS = o.New("init", p.Sprintf("初始化项目"), doInit(p))
}

func doInit(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if initFS.NArg() == 0 {
			return errors.New(p.Sprintf("未指定模块名称"))
		}
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		return i.Init(wd, initFS.Arg(0))
	}
}
