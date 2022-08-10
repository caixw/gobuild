// SPDX-License-Identifier: MIT

package watch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/message"

	"github.com/caixw/gobuild/log"
)

type builder struct {
	exts        []string // 需要监视的文件扩展名
	appName     string   // 输出的程序文件
	logs        chan<- *log.Log
	watcherFreq time.Duration
	p           *message.Printer

	wd      string    // 工作目录
	appCmd  *exec.Cmd // appName 的命令行包装引用，方便结束其进程。
	appArgs []string  // 传递给 appCmd 的参数

	goCmd     *exec.Cmd // go build 进程
	goCmdArgs []string  // 传递给 go build 的参数
}

func (opt *Options) newBuilder(logs chan<- *log.Log) *builder {
	return &builder{
		exts:        opt.Exts,
		appName:     opt.appName,
		logs:        logs,
		watcherFreq: opt.WatcherFrequency,
		p:           opt.Printer,

		wd:      filepath.Dir(opt.appName),
		appArgs: opt.appArgs,

		goCmdArgs: opt.goCmdArgs,
	}
}

// 输出不翻译的内容
func (b *builder) log(typ int8, msg ...interface{}) {
	b.logs <- &log.Log{
		Type:    typ,
		Message: fmt.Sprint(msg...),
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
	if b.goCmd != nil && b.goCmd.Process != nil {
		b.logf(log.Info, "中止旧的编译进程：%s", b.appName)
		if err := b.appCmd.Process.Kill(); err != nil {
			b.logf(log.Error, "中止旧的编译进程失败：%s", err.Error())
		}
		if err := b.appCmd.Wait(); err != nil {
			println("wait:", err.Error())
		}
		b.logf(log.Success, "旧的编译进程被终止!")
		b.appCmd = nil
	}

	b.logf(log.Info, "编译代码...")

	b.goCmd = exec.Command("go", b.goCmdArgs...)
	b.goCmd.Stderr = os.Stderr
	b.goCmd.Stdout = os.Stdout
	if err := b.goCmd.Run(); err != nil {
		b.logf(log.Error, "编译失败：%s", err.Error())
		return
	}

	b.logf(log.Success, "编译成功!")
	b.goCmd = nil

	b.restartApp()
}

// 重启被编译的程序
func (b *builder) restartApp() {
	defer func() {
		if err := recover(); err != nil {
			b.logf(log.Error, "重启失败：%v", err)
		}
	}()

	b.kill()

	b.logf(log.Info, "启动新进程：%s", b.appName)
	b.appCmd = exec.Command(b.appName, b.appArgs...)
	b.appCmd.Dir = b.wd
	b.appCmd.Stderr = os.Stderr
	b.appCmd.Stdout = os.Stdout
	if err := b.appCmd.Start(); err != nil {
		b.logf(log.Error, "启动进程时出错：%s", err)
	}
}

func (b *builder) kill() {
	if b.appCmd != nil && b.appCmd.Process != nil {
		b.logf(log.Info, "中止旧进程：%s", b.appName)
		if err := b.appCmd.Process.Kill(); err != nil {
			b.logf(log.Error, "中止旧进程失败：%s", err.Error())
		}
		if err := b.appCmd.Wait(); err != nil {
			println("wait:", err.Error())
		}
		b.logf(log.Success, "旧进程被终止!")
		b.appCmd = nil
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
