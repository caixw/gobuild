// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild/internal/local"
)

var (
	mainVersion = "0.11.0"
	metadata    string
	fullVersion = mainVersion

	versionFull bool
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

func initVersion(o *cmdopt.CmdOpt, p *message.Printer) {
	fs := o.New("version", p.Sprintf("显示版本信息"), doVersion(p))
	fs.BoolVar(&versionFull, "f", false, p.Sprintf("显示完整的版本号"))
}

func doVersion(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		version := mainVersion
		if versionFull {
			version = fullVersion
		}
		fmt.Fprintf(w, "gobuild %s build with %s\n", version, runtime.Version())

		if v, err := local.GoVersion(); err != nil {
			fmt.Fprintln(w, p.Sprintf("获取本地环境出错：%s", err.Error()))
		} else {
			fmt.Fprintln(w, p.Sprintf("本地环境 %s", v))
		}
		return nil
	}
}
