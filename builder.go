// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

// 监视器的更新频率，只有文件更新的时长超过此值，才会被更新
const watcherFrequency = 1 * time.Second

type builder struct {
	exts      []string  // 需要监视的文件扩展名
	appName   string    // 输出的程序文件
	appCmd    *exec.Cmd // appName 的命令行包装引用，方便结束其进程。
	appArgs   []string  // 传递给 appCmd 的参数
	goCmdArgs []string  // 传递给 go build 的参数
	logs      chan *Log
}

func (b *builder) log(typ int8, msg ...interface{}) {
	b.logs <- &Log{
		Type:    typ,
		Message: fmt.Sprint(msg...),
	}
}

// 确定文件 path 是否属于被忽略的格式。
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

	b.restart()
}

// 重启被编译的程序
func (b *builder) restart() {
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
	b.appCmd.Dir = filepath.Dir(b.appName) // 确定程序的工作目录
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
	} // end for paths

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

	b.log(LogTypeInfo, "以下路径或是文件将被监视:")
	for _, path := range paths {
		b.log(LogTypeInfo, path)
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
func (b *builder) watch(watcher *fsnotify.Watcher) {
	go func() {
		var buildTime time.Time
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					b.log(LogTypeIgnore, "watcher.Events:忽略 CHMOD 事件:", event)
					continue
				}

				if b.isIgnore(event.Name) { // 不需要监视的扩展名
					b.log(LogTypeIgnore, "watcher.Events:忽略不被监视的文件:", event)
					continue
				}

				if time.Now().Sub(buildTime) <= watcherFrequency { // 已经记录
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
	}()
}
