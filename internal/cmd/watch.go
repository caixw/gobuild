// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"errors"
	"flag"
	"io"
	"io/fs"
	"os"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild"
	i "github.com/caixw/gobuild/internal/init"
)

var (
	watchFS         *flag.FlagSet
	watchShowIgnore bool
)

func initWatch(o *cmdopt.CmdOpt, p *message.Printer) {
	watchFS = o.New("watch", p.Sprintf("热编译代码"), doWatch(p))
	watchFS.BoolVar(&watchShowIgnore, "i", false, p.Sprintf("是否显示被标记为 IGNORE 的日志内容"))
}

func doWatch(p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		data, err := os.ReadFile(i.ConfigFilename)
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New(p.Sprintf("未找到配置文件：%s", i.ConfigFilename))
		} else if err != nil {
			return err
		}

		o := &gobuild.WatchOptions{}
		if err := yaml.Unmarshal(data, o); err != nil {
			return err
		}
		o.Printer = p

		if watchFS.NArg() == 0 {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			o.Dirs = []string{wd}
		} else {
			o.Dirs = watchFS.Args()
		}

		logs := NewConsole(watchShowIgnore)
		defer logs.Stop()

		return gobuild.Watch(context.Background(), logs.Logs, o)
	}
}
