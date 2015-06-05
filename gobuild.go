// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/fsnotify.v1"
)

const version = "0.1.0.150604"

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
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本号")
	flag.StringVar(&outputName, "o", "", "指定程序名称")
	flag.StringVar(&mainFiles, "main", "", "指定需要编译的文件")

	flag.Usage = func() {
		def.Println(usage)
	}

	var err error
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		erro.Println(err)
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				info.Println("watcher:", event)
			case err := <-watcher.Errors:
				erro.Println(err)
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
		succ.Print("gobuild:")
		def.Println(version)

		succ.Print("Go:")
		def.Println(runtime.Version(), runtime.GOOS, runtime.GOARCH)

		return
	}

	autoBuild()

	if err := watcher.Add("./"); err != nil {
		erro.Println(err)
	}
	watcher.Close()
}

func autoBuild() {
	info.Println("============开始编译===============")

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
		erro.Println(err)
		erro.Println("============编译失败===============")
		return
	}
	succ.Println("============编译成功===============")

	restart(outputName)
}

var cmd *exec.Cmd

func kill() {
	info.Println("中止旧进程...")
	defer func() {
		if err := recover(); err != nil {
			def.Println(err)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			def.Println(err)
		}
	}
}

func start(outputName string) {
	info.Println("准备启动进程:", outputName, "...")

	if strings.IndexByte(outputName, '/') < 0 && strings.IndexByte(outputName, filepath.Separator) < 0 {
		wd, err := os.Getwd()
		if err != nil {
			erro.Println(err)
			return
		}

		outputName = wd + string(filepath.Separator) + outputName
	}

	cmd = exec.Command(outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		erro.Println(err)
	}
}

func restart(outputName string) {
	info.Println("重启中...")

	kill()
	start(outputName)
}
