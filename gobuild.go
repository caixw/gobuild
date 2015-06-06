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

	"github.com/issue9/term/colors"
	"gopkg.in/fsnotify.v1"
)

// 当前程序的版本号
const version = "0.1.1.150605"

const usage = `gobuild 用于热编译Go程序。

用法:
 gobuild [options] [dependents]

 options:
  -h 显示当前帮助信息；
  -v 显示gobuild和go程序的版本信息；
  -o 执行编译后的可执行文件名。
  -main 指定需要编译的文件，默认为""。

 dependents
  指定其它依赖的目录，只能出现在命令的尾部。


常见用法:

 - gobuild
   监视当前目录，若有变动，则重新编译当前目录下的*.go文件；

 - gobuild -main=main.go
   监视当前目录，若有变动，则重新编译当前目录下的main.go文件；

 - gobuild dir1 dir2
   监视当前目录及dir1和dir2，若有变动，则重新编译当前目录下的*.go文件；
`

var (
	showHelp    = false
	showVersion = false
	mainFiles   = ""
	outputName  = ""

	watcher *fsnotify.Watcher

	wd string // 当前工作目录

	// outputName的命令
	cmd *exec.Cmd
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
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	// 初始化监视器
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log(erro, err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log(info, "watcher.Events:", event)
				if event.Name == outputName {
					continue
				}

				autoBuild()
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

	// 确定编译后的文件名
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}
	if runtime.GOOS == "windows" {
		outputName += ".exe"
	}
	cmd = exec.Command(outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

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

	autoBuild()
}

func autoBuild() {
	log(info, "编译代码...")

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

func kill() {
	defer func() {
		if err := recover(); err != nil {
			log(erro, "kill.defer:", err)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		log(info, "中止旧进程...")
		if err := cmd.Process.Kill(); err != nil {
			log(erro, "kill:", err)
		}
	}
}

func restart() {
	kill()

	if err := cmd.Run(); err != nil {
		log(erro, "启动进程时出错:", err)
	}
}
