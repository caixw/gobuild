// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"errors"
	"flag"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild"
	i "github.com/caixw/gobuild/internal/init"
	"github.com/caixw/gobuild/watch"
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
		logs := newConsoleLogs(watchShowIgnore, os.Stderr, os.Stdout)
		defer logs.Stop()

		cancel, err := watchWithCancel(p, logs)
		if err != nil {
			return err
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		if err := watcher.Add(i.ConfigFilename); err != nil {
			return err
		}
		defer watcher.Close()

		buildTime := time.Now()
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod || time.Since(buildTime) <= watch.MinWatcherFrequency {
					continue
				}

				buildTime = time.Now()
				logs.Logs <- &watch.Log{Type: watch.LogTypeInfo, Message: p.Sprintf("配置文件被修改，重启热编译程序！")}
				cancel()
				if cancel, err = watchWithCancel(p, logs); err != nil {
					return err
				}
			case err := <-watcher.Errors:
				return err
			}
		}
	}
}

func watchWithCancel(p *message.Printer, logs *console) (context.CancelFunc, error) {
	data, err := os.ReadFile(i.ConfigFilename)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New(p.Sprintf("未找到配置文件：%s", i.ConfigFilename))
	} else if err != nil {
		return nil, err
	}

	o := &gobuild.WatchOptions{}
	if err := yaml.Unmarshal(data, o); err != nil {
		return nil, err
	}
	o.Printer = p

	if watchFS.NArg() == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		o.Dirs = []string{wd}
	} else {
		o.Dirs = watchFS.Args()
	}

	ctx, cancel := context.WithCancel(context.Background())
	go gobuild.Watch(ctx, logs.Logs, o)
	return cancel, nil
}
