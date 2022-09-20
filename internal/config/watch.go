// SPDX-License-Identifier: MIT

package config

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild/watch"
)

// Watch 监视配置文件
//
// 如果配置文件发生变化，那么重启热编译程序；
// 如果 wd 中不存在配置文件，则向上一级查找。
func Watch(wd string, p *message.Printer, logs watch.Logger) error {
	wd, err := getRoot(wd)
	if err != nil {
		return err
	}
	if err := os.Chdir(wd); err != nil { // 切换工作目录为项目根目录
		return err
	}

	cancel, err := newWatcher(p, logs)
	if err != nil {
		return err
	}

	// gobuild.yaml 监视

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add(Filename); err != nil {
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
			logs.Output(watch.LogTypeInfo, p.Sprintf("配置文件被修改，重启热编译程序！"))
			cancel()
			if cancel, err = newWatcher(p, logs); err != nil {
				return err
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}

func newWatcher(p *message.Printer, l watch.Logger) (context.CancelFunc, error) {
	data, err := os.ReadFile(Filename)
	if errors.Is(err, fs.ErrNotExist) {
		panic("配置文件不存在") // 由调用方 Watch 保证配置文件必定存在
	} else if err != nil {
		return nil, err
	}

	o := &watch.Options{}
	if err := yaml.Unmarshal(data, o); err != nil {
		return nil, err
	}
	o.Printer = p
	o.Logger = l

	ctx, cancel := context.WithCancel(context.Background())
	go watch.Watch(ctx, o)
	return cancel, nil
}

// 由 wd 目录向上查找，直到找到包含配置文件的目录。
func getRoot(wd string) (string, error) {
	abs, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return root(abs)
}

func root(parent string) (string, error) {
	wd := filepath.Dir(parent)
	if wd == "" || wd == parent {
		return "", fs.ErrNotExist
	}

	path := wd + string(filepath.Separator) + Filename
	_, err := os.Stat(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return root(wd)
	case err != nil:
		return "", err
	default:
		return wd, nil
	}
}
