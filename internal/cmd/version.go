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

func initVersion(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("version", p.Sprintf("显示版本信息"), func(w io.Writer) error {
		fmt.Fprintf(w, "gobuild %s build with %s\n", version, runtime.Version())
		if v, err := local.GoVersion(); err != nil {
			fmt.Fprintln(w, p.Sprintf("获取本地环境出错：%s", err.Error()))
		} else {
			fmt.Fprintln(w, p.Sprintf("本地环境 %s", v))
		}
		return nil
	})
}
