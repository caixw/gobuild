// SPDX-License-Identifier: MIT

package watch

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/message"

	"github.com/caixw/gobuild/log"
)

type builder struct {
	exts        []string  // 需要监视的文件扩展名
	appName     string    // 输出的程序文件
	wd          string    // 工作目录
	appCmd      *exec.Cmd // appName 的命令行包装引用，方便结束其进程。
	appArgs     []string  // 传递给 appCmd 的参数
	goCmdArgs   []string  // 传递给 go build 的参数
	logs        chan<- *log.Log
	watcherFreq time.Duration
	p           *message.Printer

	env string // 当前系统的环境信息
}

func (opt *Options) newBuilder(logs chan<- *log.Log) (*builder, error) {
	b := &builder{
		exts:        opt.Exts,
		appName:     opt.appName,
		wd:          filepath.Dir(opt.appName),
		appArgs:     opt.appArgs,
		goCmdArgs:   opt.goCmdArgs,
		logs:        logs,
		watcherFreq: opt.WatcherFrequency,
		p:           opt.Printer,
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

// 输出不翻译的内容
func (b *builder) log(typ int8, msg ...interface{}) {
	b.logs <- &log.Log{
		Type:    typ,
		Message: b.p.Sprint(msg...),
	}
}

// 输出翻译后的内容
func (b *builder) logf(typ int8, key message.Reference, msg ...interface{}) {
	b.logs <- &log.Log{
		Type:    typ,
		Message: b.p.Sprintf(key, msg...),
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
	b.logf(log.Info, "编译代码...")

	goCmd := exec.Command("go", b.goCmdArgs...)
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		b.logf(log.Error, "编译失败：%s", err.Error())
		return
	}

	b.logf(log.Success, "编译成功!")

	b.restartApp()
}

// 重启被编译的程序
func (b *builder) restartApp() {
	defer func() {
		if err := recover(); err != nil {
			b.logf(log.Error, "重启失败：%v", err)
		}
	}()

	// kill process
	if b.appCmd != nil && b.appCmd.Process != nil {
		b.logf(log.Info, "中止旧进程：%s", b.appName)
		if err := b.appCmd.Process.Kill(); err != nil {
			b.logf(log.Error, "中止旧进程失败：%s", err.Error())
		}
		b.logf(log.Success, "旧进程被终止!")
	}

	b.logf(log.Info, "启动新进程：%s", b.appName)
	b.appCmd = exec.Command(b.appName, b.appArgs...)
	b.appCmd.Dir = b.wd
	b.appCmd.Stderr = os.Stderr
	b.appCmd.Stdout = os.Stdout
	if err := b.appCmd.Start(); err != nil {
		b.logf(log.Error, "启动进程时出错：%s", err)
	}
}

// 过滤掉不需要监视的目录。以下目录会被过滤掉：
// 整个目录下都没需要监视的文件；
func (b *builder) filterPaths(paths []string) []string {
	ret := make([]string, 0, len(paths))

	for _, dir := range paths {
		fs, err := os.ReadDir(dir)
		if err != nil {
			b.logf(log.Error, err)
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
	b.logf(log.Info, "初始化监视器...")

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	paths = b.filterPaths(paths)

	b.logf(log.Ignore, "以下路径或是文件将被监视：")
	for _, path := range paths {
		b.log(log.Ignore, path) // 路径不翻译
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
			b.logf(log.Info, context.Canceled.Error())
			return
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				b.logf(log.Ignore, "watcher.Events:忽略 %s 事件", event.String())
				continue
			}

			if b.isIgnore(event.Name) { // 不需要监视的扩展名
				b.logf(log.Ignore, "watcher.Events:忽略不被监视的文件：%s", event.Name)
				continue
			}

			if time.Since(buildTime) <= b.watcherFreq { // 已经记录
				b.logf(log.Ignore, "watcher.Events:忽略短期内频繁修改的文件：%s", event.Name)
				continue
			}

			buildTime = time.Now()
			b.logf(log.Info, "watcher.Events:%s 事件触发了编译", event.String())

			go b.build()
		case err := <-watcher.Errors:
			watcher.Close()
			b.logf(log.Warn, "watcher.Errors：%s", err.Error())
		} // end select
	}
}
