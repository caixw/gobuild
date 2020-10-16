// SPDX-License-Identifier: MIT

// Package gobuild 提供了对 Go 语言热编译的支持
package gobuild

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type builder struct {
	exts        []string  // 需要监视的文件扩展名
	appName     string    // 输出的程序文件
	wd          string    // 工作目录
	appCmd      *exec.Cmd // appName 的命令行包装引用，方便结束其进程。
	appArgs     []string  // 传递给 appCmd 的参数
	goCmdArgs   []string  // 传递给 go build 的参数
	logs        chan *Log
	watcherFreq time.Duration

	// 退出标记
	exit chan struct{}

	env string // 当前系统的环境信息
}

// Build 执行热编译服务
func Build(ctx context.Context, logs chan *Log, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}

	b, err := opt.newBuilder(logs)
	if err != nil {
		return err
	}

	b.log(LogTypeInfo, fmt.Sprintf("当前环境参数如下：%s", b.env))

	b.log(LogTypeInfo, fmt.Sprint("给程序传递了以下参数：", b.appArgs)) // 输出提示信息

	switch { // 提示扩展名
	case len(b.exts) == 0: // 允许不监视任意文件，但输出一信息来警告
		b.log(LogTypeWarn, "将 ext 设置为空值，意味着不监视任何文件的改变！")
	case len(b.exts) > 0:
		b.log(LogTypeInfo, fmt.Sprint("系统将监视以下类型的文件:", b.exts))
	}

	b.log(LogTypeInfo, fmt.Sprint("输出文件为:", b.appName)) // 提示 appName

	w, err := b.initWatcher(opt.paths)
	if err != nil {
		return err
	}
	defer w.Close()

	go b.watch(ctx, w)

	<-b.exit
	return context.Canceled
}

func (opt *Options) newBuilder(logs chan *Log) (*builder, error) {
	b := &builder{
		exts:        opt.exts,
		appName:     opt.appName,
		wd:          filepath.Dir(opt.appName),
		appArgs:     opt.appArgs,
		goCmdArgs:   opt.goCmdArgs,
		logs:        logs,
		watcherFreq: opt.WatcherFrequency,
		exit:        make(chan struct{}, 1),
	}

	var buf bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	b.env = strings.TrimSpace(strings.TrimPrefix(buf.String(), "go version "))
	return b, nil
}

func (b *builder) log(typ int8, msg ...interface{}) {
	b.logs <- &Log{
		Type:    typ,
		Message: fmt.Sprint(msg...),
	}
}

// 确定文件 path 是否属于被忽略的格式
func (b *builder) isIgnore(path string) bool {
	if b.appCmd != nil && b.appCmd.Path == path { // 忽略程序本身的监视
		return true
	}

	for _, ext := range b.exts {
		if ext == "*" {
			return false
		}
		if strings.HasSuffix(path, ext) {
			return false
		}
	}

	return true
}

// 开始编译代码
func (b *builder) build() {
	b.log(LogTypeInfo, "编译代码...")

	goCmd := exec.Command("go", b.goCmdArgs...)
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		b.log(LogTypeError, "编译失败:", err)
		return
	}

	b.log(LogTypeSuccess, "编译成功!")

	b.restartApp()
}

// 重启被编译的程序
func (b *builder) restartApp() {
	defer func() {
		if err := recover(); err != nil {
			b.log(LogTypeError, "restart.defer:", err)
		}
	}()

	// kill process
	if b.appCmd != nil && b.appCmd.Process != nil {
		b.log(LogTypeInfo, "中止旧进程:", b.appName)
		if err := b.appCmd.Process.Kill(); err != nil {
			b.log(LogTypeError, "kill:", err)
		}
		b.log(LogTypeSuccess, "旧进程被终止!")
	}

	b.log(LogTypeInfo, "启动新进程:", b.appName)
	b.appCmd = exec.Command(b.appName, b.appArgs...)
	b.appCmd.Dir = b.wd
	b.appCmd.Stderr = os.Stderr
	b.appCmd.Stdout = os.Stdout
	if err := b.appCmd.Start(); err != nil {
		b.log(LogTypeError, "启动进程时出错:", err)
	}
}

// 过滤掉不需要监视的目录。以下目录会被过滤掉：
// 整个目录下都没需要监视的文件；
func (b *builder) filterPaths(paths []string) []string {
	ret := make([]string, 0, len(paths))

	for _, dir := range paths {
		fs, err := ioutil.ReadDir(dir)
		if err != nil {
			b.log(LogTypeError, err)
			continue
		}

		ignore := true
		for _, f := range fs {
			if f.IsDir() {
				continue
			}
			if !b.isIgnore(f.Name()) {
				ignore = false
				break
			}
		}
		if !ignore {
			ret = append(ret, dir)
		}
	}

	return ret
}

func (b *builder) initWatcher(paths []string) (*fsnotify.Watcher, error) {
	b.log(LogTypeInfo, "初始化监视器...")

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	paths = b.filterPaths(paths)

	b.log(LogTypeIgnore, "以下路径或是文件将被监视:")
	for _, path := range paths {
		b.log(LogTypeIgnore, path)
	}

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

// 开始监视 paths 中指定的目录或文件。
func (b *builder) watch(ctx context.Context, watcher *fsnotify.Watcher) {
	var buildTime time.Time
	for {
		select {
		case <-ctx.Done():
			b.log(LogTypeInfo, context.Canceled)

			b.exit <- struct{}{}
			return
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				b.log(LogTypeIgnore, "watcher.Events:忽略 CHMOD 事件:", event)
				continue
			}

			if b.isIgnore(event.Name) { // 不需要监视的扩展名
				b.log(LogTypeIgnore, "watcher.Events:忽略不被监视的文件:", event)
				continue
			}

			if time.Now().Sub(buildTime) <= b.watcherFreq { // 已经记录
				b.log(LogTypeIgnore, "watcher.Events:该监控事件被忽略:", event)
				continue
			}

			buildTime = time.Now()
			b.log(LogTypeInfo, "watcher.Events:触发编译事件:", event)

			go b.build()
		case err := <-watcher.Errors:
			watcher.Close()
			b.log(LogTypeWarn, "watcher.Errors", err)
		} // end select
	}
}
