// SPDX-FileCopyrightText: 2015-2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/caixw/gobuild"
	"github.com/goccy/go-yaml"
	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
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
			logs := gobuild.NewConsoleLogger(watchShowIgnore, os.Stdout, nil, nil)

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			data, err := os.ReadFile(filepath.Join(wd, configFilename))
			if errors.Is(err, os.ErrNotExist) {
				panic("配置文件不存在") // 由调用方 Watch 保证配置文件必定存在
			} else if err != nil {
				return err
			}

			o := &gobuild.WatchOptions{}
			if err := yaml.Unmarshal(data, o); err != nil {
				return err
			}

			return gobuild.Watch(context.Background(), p, logs, o)
		}
	})
}
