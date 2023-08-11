// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild/internal/local"
)

var (
	mainVersion = "1.1.0"
	metadata    string
	fullVersion = mainVersion

	versionFull bool
)

const (
	showVersion      = localeutil.StringPhrase("显示版本信息")
	showVersionUsage = localeutil.StringPhrase("显示版本信息 usage")
	fullVersionUsage = localeutil.StringPhrase("显示完整的版本号")
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

func initVersion(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("version", showVersion.LocaleString(p), showVersionUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		fs.BoolVar(&versionFull, "f", false, fullVersionUsage.LocaleString(p))

		return func(w io.Writer) error {
			version := mainVersion
			if versionFull {
				version = fullVersion
			}
			fmt.Fprintf(w, "gobuild %s build with %s\n", version, runtime.Version())

			if v, err := local.GoVersion(); err != nil {
				fmt.Fprintln(w, localeutil.Phrase("获取本地环境出错：%s", err.Error()).LocaleString(p))
			} else {
				fmt.Fprintln(w, localeutil.Phrase("本地环境 %s", v).LocaleString(p))
			}
			return nil
		}
	})
}
