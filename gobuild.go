// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild // import "github.com/caixw/gobuild"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Build 执行热编译服务。
//
// logs 日志输出通道
// mainFiles 为 go build 最后的文件参数，可以为空，表示当前目录；
// outputName 指定可执行文件输出的文件路径，为空表示默认值；
// exts 指定监视的文件扩展名，为空表示不监视任何文件，* 表示监视所有文件；
// recursive 是否监视子目录；
// appArgs 传递给程序的参数；
// dir 表示需要监视的目录，至少指定一个目录，第一个目录被当作子目录，将编译其下的文件。
func Build(logs chan *Log, mainFiles, outputName, exts string, recursive bool, appArgs string, dir ...string) error {
	if len(dir) < 1 {
		return errors.New("参数 dir 至少指定一个")
	}
	wd, err := filepath.Abs(dir[0])
	if err != nil {
		return err
	}
	dir[0] = wd

	appName, err := getAppName(outputName, wd)
	if err != nil {
		return err
	}

	// 初始化 goCmd 的参数
	args := []string{"build", "-o", appName}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}

	b := &builder{
		exts:      getExts(exts),
		appName:   appName,
		appArgs:   splitArgs(appArgs),
		goCmdArgs: args,
		logs:      logs,
	}

	// 输出提示信息
	logs <- &Log{
		Type:    LogTypeInfo,
		Message: fmt.Sprint("给程序传递了以下参数：", b.appArgs),
	}

	// 提示扩展名
	switch {
	case len(b.exts) == 0: // 允许不监视任意文件，但输出一信息来警告
		logs <- &Log{
			Type:    LogTypeWarn,
			Message: "将 ext 设置为空值，意味着不监视任何文件的改变！",
		}
	case len(b.exts) > 0:
		logs <- &Log{
			Type:    LogTypeInfo,
			Message: fmt.Sprint("系统将监视以下类型的文件:", b.exts),
		}
	}

	// 提示 appName
	logs <- &Log{
		Type:    LogTypeInfo,
		Message: fmt.Sprint("输出文件为:", b.appName),
	}

	paths, err := recursivePaths(recursive, dir)
	if err != nil {
		return err
	}
	w, err := b.initWatcher(paths)
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

	return ret
}

// 根据 recursive 值确定是否递归查找 paths 每个目录下的子目录。
func recursivePaths(recursive bool, paths []string) ([]string, error) {
	if !recursive {
		return paths, nil
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			return nil, err
		}
	}

	return ret, nil
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

	return ret
}

func getAppName(outputName, wd string) (string, error) {
	if outputName == "" {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}

	// 没有分隔符，表示仅有一个文件名，需要加上 wd
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = filepath.Join(wd, outputName)
	}

	// 转成绝对路径
	outputName, err := filepath.Abs(outputName)
	if err != nil {
		return "", err
	}

	return outputName, nil
}
