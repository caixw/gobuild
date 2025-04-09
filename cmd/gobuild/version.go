// SPDX-FileCopyrightText: 2015-2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
)

const (
	showVersion      = localeutil.StringPhrase("显示版本信息")
	showVersionUsage = localeutil.StringPhrase("显示版本信息 usage")
)

func initVersion(o *cmdopt.CmdOpt, p *message.Printer) {
	o.New("version", showVersion.LocaleString(p), showVersionUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			version := "unknown"
			if info, ok := debug.ReadBuildInfo(); ok {
				version = info.Main.Version
			}
			fmt.Fprintf(w, "gobuild %s build with %s\n", version, runtime.Version())

			if v, err := goVersion(); err != nil {
				fmt.Fprintln(w, localeutil.Phrase("获取本地环境出错：%s", err.Error()).LocaleString(p))
			} else {
				fmt.Fprintln(w, localeutil.Phrase("本地环境 %s", v).LocaleString(p))
			}
			return nil
		}
	})
}

// goVersion 返回本地 Go 的版本信息
func goVersion() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimPrefix(buf.String(), "go version ")), nil
}
