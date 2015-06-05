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

	"gopkg.in/fsnotify.v1"
)

const version = "0.1.1.150605"

const usage = `gobuild 用于热编译Go程序。

用法:
 gobuild [options] [dependents]

 options:
  -h 显示当前帮助信息；
  -v 显示gobuild和go程序的版本信息；
  -o 指定编译后的文件名；
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

 - gobuild -o="/var/main" -main="main.go" dir1 dir2
   监视当前目录及dir1和dir2，若有变动，则重新编译当前目录下的main.go文件并保存到/var/main；
`

var (
	showHelp    = false
	showVersion = false
	outputName  = ""
	mainFiles   = ""

	watcher *fsnotify.Watcher
)

func init() {
	// 初始化flag相关设置
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&outputName, "o", "", "指定程序名称")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	// 初始化监视器
	var err error
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log(erro, err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log(info, "watcher:", event)
				autoBuild()
			case err := <-watcher.Errors:
				log(erro, err)
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
		fmt.Println("gobuild:", version)
		fmt.Println("Go:", runtime.Version(), runtime.GOOS, runtime.GOARCH)

		return
	}

	log(info, "初始化监视器...")
	log(info, "以下路径或是文件将被监视:")
	// TODO 输出监视文件
	if err := watcher.Add("./"); err != nil {
		log(erro, err)
		os.Exit(2)
	}

	autoBuild()
	//watcher.Close()
}

func autoBuild() {
	log(info, "编译代码...")

	if len(outputName) > 0 && runtime.GOOS == "windows" {
		outputName += ".exe"
	}

	args := []string{"build"}
	if len(outputName) > 0 {
		args = append(args, "-o", outputName)
	}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}

	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log(erro, "编译失败:", err)
		return
	}

	log(succ, "编译成功!")
	restart(outputName)
}

var cmd *exec.Cmd

func kill() {
	log(info, "中止旧进程...")
	defer func() {
		if err := recover(); err != nil {
			log(erro, err)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			log(erro, err)
		}
	}
}

func start(outputName string) {
	log(info, "准备启动进程:", outputName, "...")

	if strings.IndexByte(outputName, '/') < 0 && strings.IndexByte(outputName, filepath.Separator) < 0 {
		wd, err := os.Getwd()
		if err != nil {
			log(erro, err)
			return
		}

		outputName = wd + string(filepath.Separator) + outputName
	}

	cmd = exec.Command(outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log(erro, err)
	}
}

func restart(outputName string) {
	log(info, "重启进程...")

	kill()
	start(outputName)
}
