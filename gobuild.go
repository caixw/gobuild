// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 一个简单的Go语言热编译工具。
//
// gobuild会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，
// 一旦触发，就会调用`go build`编译Go源文件并执行。
//  // 监视当前目录下的文件，若发生变化，则触发go build -main="*.go"
//  gobuild
//
//  // 监视当前目录和term目录下的文件，若发生变化，则触发go build -main="main.go"
//  gobuild -main=main.go ~/Go/src/github.com/issue9/term
package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/issue9/term/colors"
)

// 当前程序的版本号
const version = "0.1.3.150607"

var cmd *exec.Cmd // outputName的命令

func init() {
	// 基本环境检测
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		log(erro, "未设置环境变量GOPATH")
		os.Exit(2)
	}
}

func main() {
	initFlag()

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

	wd, err := os.Getwd()
	if err != nil {
		log(erro, "获取当前工作目录时，发生以下错误:", err)
		os.Exit(2)
	}
	// 确定编译后的文件名
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}

	// 监视的路径，必定包含当前工作目录
	initWatcher(append(flag.Args(), wd), extString)

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
