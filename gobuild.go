// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 一个简单的Go语言热编译工具
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/term/colors"
	"gopkg.in/fsnotify.v1"
)

// 当前程序的版本号
const version = "0.1.1.150606"

var (
	showHelp    = false
	showVersion = false
	mainFiles   = ""
	outputName  = ""
	extString   = "go"
	exts        []string

	watcher *fsnotify.Watcher
	wd      string    // 当前工作目录
	cmd     *exec.Cmd // outputName的命令
)

func init() {
	// 基本环境检测
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		log(erro, "未设置环境变量GOPATH")
		os.Exit(2)
	}

	// 获取所有被监视的路径
	var err error
	wd, err = os.Getwd()
	if err != nil {
		log(erro, "获取当前工作目录时，发生以下错误:", err)
		os.Exit(2)
	}

	// 初始化flag相关设置
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&outputName, "o", "", "指定输出名称")
	flag.StringVar(&extString, "ext", "go", "指定监视的文件扩展名")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	// 初始化监视器
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log(erro, err)
	}

	go func() {
		var buildTime int64
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}

				if event.Name == outputName { // 过滤程序本身
					continue
				}

				if !isEnabledExt(event.Name) { // 不需要监视的扩展名
					continue
				}

				if time.Now().Unix()-buildTime <= 1 { // 已经记录
					log(info, "该监控事件被忽略:", event)
					continue
				}

				buildTime = time.Now().Unix()
				log(info, "watcher.Events:", event)
				go autoBuild()
			case err := <-watcher.Errors:
				log(erro, "watcher.Errors", err)
			}
		}
	}()
}

func main() {
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		colors.Print(colors.Stdout, colors.Green, colors.Default, "gobuild: ")
		colors.Println(colors.Stdout, colors.Default, colors.Default, version)

		colors.Print(colors.Stdout, colors.Green, colors.Default, "Go: ")
		goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
		colors.Println(colors.Stdout, colors.Default, colors.Default, goVersion)

		return
	}

	if len(extString) > 0 {
		exts = strings.Split(extString, ",")
	}

	// 确定编译后的文件名
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}

	// 监视的路径，必定包含当前工作目录
	paths := append(flag.Args(), wd)
	log(info, "初始化监视器...")
	log(info, "以下路径或是文件将被监视:", paths)
	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			log(erro, "watcher.Add:", err)
			os.Exit(2)
		}
	}

	// 初始化cmd变量
	cmd = exec.Command(outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// 首次编译程序。
	autoBuild()

	done := make(chan bool)
	<-done
}

func autoBuild() {

	log(info, "编译代码...")

	// 初始化goCmd变量
	args := []string{"build", "-o", outputName}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}
	goCmd := exec.Command("go", args...)
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	if err := goCmd.Run(); err != nil {
		log(erro, "编译失败:", err)
		return
	}

	log(succ, "编译成功!")

	restart()
}

// 重启cmd程序
func restart() {
	defer func() {
		if err := recover(); err != nil {
			log(erro, "restart.defer:", err)
		}
	}()

	// kill process
	if cmd != nil && cmd.Process != nil {
		log(info, "中止旧进程...")
		if err := cmd.Process.Kill(); err != nil {
			log(erro, "kill:", err)
		}
	}

	if err := cmd.Run(); err != nil {
		log(erro, "启动进程时出错:", err)
	}
}

// path文件是否包含允许的扩展名。
func isEnabledExt(path string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
