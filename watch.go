// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

// 初始化监视器，paths指定需要监视的路径或是文件。
func initWatcher(arg *args) {
	log(info, "以下类型的文件将被忽略:", arg.exts)

	// path文件是否包含允许的扩展名。
	isEnabledExt := func(path string) bool {
		if arg.exts[0] == "*" {
			return true
		}

		for _, ext := range arg.exts {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}

		return false
	}

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log(erro, err)
		os.Exit(2)
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

				if event.Name == arg.outputName { // 过滤程序本身
					log(ignore, "watcher.Events:忽略程序本身的改变:", event)
					continue
				}

				if !isEnabledExt(event.Name) { // 不需要监视的扩展名
					log(ignore, "watcher.Events:忽略不被监视的文件:", event)
					continue
				}

				if time.Now().Unix()-buildTime <= 1 { // 已经记录
					log(ignore, "watcher.Events:该监控事件被忽略:", event)
					continue
				}

				buildTime = time.Now().Unix()
				log(info, "watcher.Events:触发编译事件:", event)

				go autoBuild(arg)
			case err := <-watcher.Errors:
				log(warn, "watcher.Errors", err)
			}
		}
	}()

	// 监视的路径，必定包含当前工作目录
	log(info, "初始化监视器...")
	log(info, "以下路径或是文件将被监视:", arg.paths)
	for _, path := range arg.paths {
		if err := watcher.Add(path); err != nil {
			log(erro, "watcher.Add:", err)
			os.Exit(2)
		}
	}
}
