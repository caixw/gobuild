// SPDX-License-Identifier: MIT

package watch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/message"
)

type builder struct {
	exts        []string // 需要监视的文件扩展名
	anyExt      bool
	excludes    []string
	appName     string // 输出的程序文件
	logs        Logger
	watcherFreq time.Duration
	p           *message.Printer

	// 被编译程序的执行环境
	appWD      string    // 工作目录
	appCmd     *exec.Cmd // appName 的命令行包装引用，方便结束其进程。
	appArgs    []string  // 传递给 appCmd 的参数
	appKillMux sync.Mutex

	// go build 或是 go mod tidy 的运行环境
	goTidy    bool      // 自动运行 go mod tidy
	goCmd     *exec.Cmd // go build 或是 go mod 进程
	goArgs    []string  // 传递给 go build 的参数，go mod tidy 忽略此参数。
	goKillMux sync.Mutex
}

func (opt *Options) newBuilder() *builder {
	exts := opt.Exts
	if opt.anyExts {
		exts = []string{"*"}
	}
	return &builder{
		exts:        exts,
		anyExt:      opt.anyExts,
		excludes:    opt.Excludes,
		appName:     opt.appName,
		logs:        opt.Logger,
		watcherFreq: opt.WatcherFrequency,
		p:           opt.Printer,

		appWD:   filepath.Dir(opt.appName),
		appArgs: opt.appArgs,

		goTidy: opt.AutoTidy,
		goArgs: opt.goCmdArgs,
	}
}

// 输出翻译后的内容
func (b *builder) logf(typ int8, key message.Reference, msg ...interface{}) {
	b.logs.Output(typ, b.p.Sprintf(key, msg...))
}

// 确定文件 path 是否属于被忽略的格式
func (b *builder) isIgnore(path string) bool {
	if b.appCmd != nil && b.appCmd.Path == path { // 忽略程序本身的监视
		return true
	}

	for _, p := range b.excludes {
		matched, err := filepath.Match(p, path)
		if err != nil {
			b.logf(LogTypeError, b.p.Sprintf(err.Error))
			return false
		}
		if matched {
			return true
		}
	}

	if b.anyExt {
		return false
	}
	for _, ext := range b.exts {
		if strings.HasSuffix(path, ext) {
			return false
		}
	}

	return true
}

func (b *builder) tidy() {
	b.killGo()

	b.logf(LogTypeInfo, "执行 go mod tidy...")

	b.goCmd = exec.Command("go", "mod", "tidy")
	b.goCmd.Stderr = asWriter(LogTypeGo, b.logs)
	b.goCmd.Stdout = asWriter(LogTypeGo, b.logs)
	if err := b.goCmd.Run(); err != nil {
		b.logf(LogTypeError, "go mod tidy 失败：%s", err.Error())
		return
	}

	b.logf(LogTypeSuccess, "go mod tidy 完成!")
	b.goCmd = nil
}

// 开始编译代码
func (b *builder) build() {
	b.killGo()

	b.logf(LogTypeInfo, "编译代码...")

	b.goCmd = exec.Command("go", b.goArgs...)
	b.goCmd.Stderr = asWriter(LogTypeGo, b.logs)
	b.goCmd.Stdout = asWriter(LogTypeGo, b.logs)
	if err := b.goCmd.Run(); err != nil {
		b.logf(LogTypeError, "编译失败：%s", err.Error())
		return
	}

	b.logf(LogTypeSuccess, "编译成功!")
	b.goCmd = nil

	b.restartApp()
}

func (b *builder) killGo() {
	b.goKillMux.Lock()
	defer b.goKillMux.Unlock()

	if b.goCmd != nil && b.goCmd.Process != nil {
		b.logf(LogTypeInfo, "中止旧的编译进程")
		if err := b.goCmd.Process.Kill(); err != nil {
			b.logf(LogTypeError, "中止旧的编译进程失败：%s", err.Error())
		}
		if err := b.goCmd.Wait(); err != nil {
			b.logf(LogTypeError, "被中止编译进程非正常退出：%s", err.Error())
		}
		b.logf(LogTypeSuccess, "旧的编译进程被终止!")
		b.goCmd = nil
	}
}

// 重启被编译的程序
func (b *builder) restartApp() {
	defer func() {
		if err := recover(); err != nil {
			b.logf(LogTypeError, "重启失败：%v", err)
		}
	}()

	b.killApp()

	b.logf(LogTypeInfo, "启动新进程：%s", b.appName)
	b.appCmd = exec.Command(b.appName, b.appArgs...)
	b.appCmd.Dir = b.appWD
	b.appCmd.Stderr = asWriter(LogTypeApp, b.logs)
	b.appCmd.Stdout = asWriter(LogTypeApp, b.logs)
	if err := b.appCmd.Start(); err != nil {
		b.logf(LogTypeError, "启动进程时出错：%s", err)
	}
}

func (b *builder) killApp() {
	b.appKillMux.Lock()
	defer b.appKillMux.Unlock()

	if b.appCmd != nil && b.appCmd.Process != nil {
		b.logf(LogTypeInfo, "中止旧进程：%s", b.appName)
		if err := b.appCmd.Process.Kill(); err != nil {
			b.logf(LogTypeError, "中止旧进程失败：%s", err.Error())
		}
		if err := b.appCmd.Wait(); err != nil {
			b.logf(LogTypeError, "被中止进程非正常退出：%s", err.Error())
		}
		b.logf(LogTypeSuccess, "旧进程被终止!")
		b.appCmd = nil
	}
}

// 过滤掉不需要监视的目录。以下目录会被过滤掉：
// - 整个目录下都没需要监视的文件；
func (b *builder) filterPaths(paths []string) []string {
	ret := make([]string, 0, len(paths))

	for _, dir := range paths {
		fs, err := os.ReadDir(dir)
		if err != nil {
			b.logf(LogTypeError, err)
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
