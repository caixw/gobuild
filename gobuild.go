// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild // import "github.com/caixw/gobuild"

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Build 执行热编译操作
func Build(mainFiles, outputName, exts string, recursive bool, appArgs string, in, su, wa, er, ign *log.Logger) error {
	info = in
	succ = su
	erro = er
	warn = wa
	ignore = ign

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 初始化 goCmd 的参数
	args := []string{"build", "-o", outputName}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}

	b := &builder{
		exts:      getExts(exts),
		appName:   getAppName(outputName, wd),
		appArgs:   splitArgs(appArgs),
		goCmdArgs: args,
	}

	w, err := b.initWatcher(recursivePaths(recursive, append(flag.Args(), wd)))
	if err != nil {
		return err
	}
	defer w.Close()

	b.watch(w)
	go b.build()

	<-make(chan bool)
	return nil
}

func splitArgs(args string) []string {
	ret := make([]string, 0, 10)
	var state byte
	var start, index int

	for index = 0; index < len(args); index++ {
		b := args[index]
		if b == ' ' {
			if state != ' ' {
				ret = append(ret, args[start:index])
				state = ' '
			}
			start = index + 1
			continue
		}

		if b == '=' {
			if state != '=' {
				ret = append(ret, args[start:index])
				state = '='
			}
			start = index + 1
			continue
		}

		state = 0
	} // end for

	if start < len(args) {
		ret = append(ret, args[start:len(args)])
	}

	info.Println("给程序传递了以下参数：", ret)

	return ret
}

// 根据 recursive 值确定是否递归查找 paths 每个目录下的子目录。
func recursivePaths(recursive bool, paths []string) []string {
	if !recursive {
		return paths
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			erro.Println("在遍历监视目录时，发生以下错误:", err)
		}

		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			erro.Println("在遍历监视目录时，发生以下错误:", err)
		}
	}

	return ret
}

// 将 extString 分解成数组，并清理掉无用的内容，比如空字符串
func getExts(extString string) []string {
	exts := strings.Split(extString, ",")
	ret := make([]string, 0, len(exts))

	for _, ext := range exts {
		ext = strings.TrimSpace(ext)

		if len(ext) == 0 {
			continue
		}
		if ext[0] != '.' {
			ext = "." + ext
		}
		ret = append(ret, ext)
	}

	switch {
	case len(ret) == 0: // 允许不监视任意文件，但输出一信息来警告
		warn.Println("将 ext 设置为空值，意味着不监视任何文件的改变！")
	case len(ret) > 0:
		info.Println("系统将监视以下类型的文件:", ret)
	}

	return ret
}

func getAppName(outputName, wd string) string {
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = wd + string(filepath.Separator) + outputName
	}

	// 转成绝对路径
	outputName, err := filepath.Abs(outputName)
	if err != nil {
		erro.Println(err)
	}

	info.Println("输出文件为:", outputName)

	return outputName
}
