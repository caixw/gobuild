// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

type builder struct {
	exts      []string  // 不需要监视的文件扩展名
	appCmd    *exec.Cmd // 被编译的程序
	goCmdArgs []string  // 传递给go build的参数
}

// 确定文件path是否属于被忽略的格式。
func (b *builder) isIgnore(path string) bool {
	if b.appCmd != nil && b.appCmd.Path == path { // 忽略程序本身的监视
		return true
	}

	for _, ext := range b.exts {
		if len(ext) == 0 {
			continue
		}
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
	log(info, "编译代码...")

	goCmd := exec.Command("go", b.goCmdArgs...)
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		log(erro, "编译失败:", err)
		return
	}

	log(succ, "编译成功!")

	b.restart()
}

// 重启被编译的程序
func (b *builder) restart() {
	defer func() {
		if err := recover(); err != nil {
			log(erro, "restart.defer:", err)
		}
	}()

	// kill process
	if b.appCmd != nil && b.appCmd.Process != nil {
		log(info, "中止旧进程...")
		if err := b.appCmd.Process.Kill(); err != nil {
			log(erro, "kill:", err)
		}
	}

	if err := b.appCmd.Run(); err != nil {
		log(erro, "启动进程时出错:", err)
	}
}

// 开始监视paths中指定的目录或文件。
func (b *builder) watch(paths []string) {
	log(info, "初始化监视器...")

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log(erro, err)
		os.Exit(2)
	}

	// 监视的路径，必定包含当前工作目录
	log(info, "以下路径或是文件将被监视:", paths)
	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			log(erro, "watcher.Add:", err)
			os.Exit(2)
		}
	}

	go func() {
		var buildTime int64
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log(ignore, "watcher.Events:忽略CHMOD事件:", event)
					continue
				}

				if b.isIgnore(event.Name) { // 不需要监视的扩展名
					log(ignore, "watcher.Events:忽略不被监视的文件:", event)
					continue
				}

				if time.Now().Unix()-buildTime <= 1 { // 已经记录
					log(ignore, "watcher.Events:该监控事件被忽略:", event)
					continue
				}

				buildTime = time.Now().Unix()
				log(info, "watcher.Events:触发编译事件:", event)

				go b.build()
			case err := <-watcher.Errors:
				log(warn, "watcher.Errors", err)
			}
		}
	}()
}
