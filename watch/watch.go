// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package watch 监视文件变化并编译
package watch

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
)

// Watch 执行热编译服务
//
// 如果初始化参数有误，则反错误信息，如果是编译过程中出错，将直接将错误内容输出到 logs。
func Watch(ctx context.Context, p *message.Printer, l Logger, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}

	b := opt.newBuilder(p, l)

	env, err := goVersion()
	if err != nil {
		return err
	}
	b.systemLog(Info, localeutil.Phrase("当前环境参数如下：%s", env))

	b.systemLog(Info, localeutil.Phrase("给程序传递了以下参数：%s", b.appArgs)) // 输出提示信息

	switch { // 提示扩展名
	case len(b.exts) == 0: // 允许不监视任意文件，但输出警告信息
		b.systemLog(Warn, localeutil.StringPhrase("将 ext 设置为空值，意味着不监视任何文件的改变！"))
	case len(b.exts) > 0:
		b.systemLog(Info, localeutil.Phrase("系统将监视以下类型的文件：%s", b.exts))
	}

	b.systemLog(Info, localeutil.Phrase("输出文件为：%s", b.appName)) // 提示 appName

	return b.watch(ctx, opt.paths)
}

// 开始监视 paths 中指定的目录或文件
func (b *builder) watch(ctx context.Context, paths []string) error {
	go b.build() // 第一次主动编译程序，后续的才是监视变化。

	watcher, err := b.initWatcher(paths)
	if err != nil {
		return err
	}
	defer watcher.Close()

	var buildTime time.Time
	for {
		select {
		case <-ctx.Done():
			b.killApp()
			b.killGo()
			b.systemLog(Info, localeutil.StringPhrase("用户取消"))
			return nil
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				b.systemLog(Ignore, localeutil.Phrase("watcher.Events:忽略 %s 事件", event.String()))
				continue
			}

			if b.isIgnore(event.Name) { // 不需要监视的扩展名
				b.systemLog(Ignore, localeutil.Phrase("watcher.Events:忽略不被监视的文件：%s", event.String()))
				continue
			}

			if time.Since(buildTime) <= b.watcherFreq {
				b.systemLog(Ignore, localeutil.Phrase("watcher.Events:忽略短期内频繁修改的文件：%s", event.Name))
				continue
			}
			buildTime = time.Now()

			b.systemLog(Info, localeutil.Phrase("watcher.Events:%s 事件触发了编译", event.String()))
			go b.build()
		case err := <-watcher.Errors:
			b.systemLog(Warn, localeutil.Phrase("watcher.Errors：%s", err.Error()))
			return nil
		} // end select
	}
}

func (b *builder) initWatcher(paths []string) (*fsnotify.Watcher, error) {
	b.systemLog(Info, localeutil.StringPhrase("初始化监视器..."))

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	paths = b.filterPaths(paths)

	ps := strings.Join(paths, ",\n")
	b.systemLog(Ignore, localeutil.Phrase("以下路径或是文件将被监视：%s", ps))

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

func goVersion() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimPrefix(buf.String(), "go version ")), nil
}
