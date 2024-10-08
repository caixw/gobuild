// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package watch

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/source"
)

// MinWatcherFrequency 监视器更新频率的最小值
const MinWatcherFrequency = time.Second

// Options 热编译的选项
type Options struct {
	XMLName struct{} `xml:"gobuild" json:"-" yaml:"-"`

	// 指定编译的文件
	//
	// 为 go build 最后的文件参数，可以为空，表示当前目录。
	//
	// 同时也会根据此值查找 go.mod 的项目根目录。
	MainFiles string `xml:"main,omitempty" json:"main,omitempty" yaml:"main,omitempty"`
	appName   string
	paths     []string
	wd        string

	// 传递给编译器的参数
	//
	// 即传递给 go build 命令的参数，但是不应该包含 -o 等参数
	Args []string `xml:"args,omitempty" json:"args,omitempty" yaml:"args,omitempty"`

	// 指定监视的文件扩展名
	//
	// 如果指定了 *，表示所有文件类型，包括没有扩展名的文件。默认为 .go。
	Exts    []string `xml:"exts,omitempty" json:"exts,omitempty" yaml:"exts,omitempty"`
	anyExts bool

	// 忽略的文件
	//
	// 采用 [path.Match] 作为匹配方式。_test.go 始终被忽略，不需要在此指定。默认为空。
	Excludes []string `xml:"excludes>glob,omitempty" json:"excludes,omitempty" yaml:"excludes,omitempty"`

	// 传递给编译成功后的程序的参数
	AppArgs string `xml:"appArgs,omitempty" yaml:"appArgs,omitempty" json:"appArgs,omitempty"`
	appArgs []string

	// 监视器的更新频率
	//
	// 只有文件更新的时长超过此值，才会被定义为更新。防止文件频繁修改导致的频繁编译调用。
	//
	// 此值不能小于 [MinWatcherFrequency]。默认值为 [MinWatcherFrequency]。
	WatcherFrequency time.Duration `xml:"freq,omitempty" yaml:"freq,omitempty" json:"freq,omitempty"`

	// 传递给 go 命令的参数
	goCmdArgs []string
}

func (opt *Options) sanitize() (err error) {
	// 检测 glob 语法
	for _, p := range opt.Excludes {
		if _, err := filepath.Match(p, "abc"); err != nil {
			return err
		}
	}

	if opt.MainFiles == "" {
		opt.MainFiles = "./"
	}

	// 根据 MainFiles 拿到 wd 和 appName

	// MainFiles 可能是 *.go 等非正常的目录结构，根据最后一个字符作简单判断。
	opt.wd, err = getWD(opt.MainFiles)
	if err != nil {
		return err
	}
	// BUG: 如果获得的 opt.wd == /，那么 appName 将是个非法值。
	opt.appName = filepath.Join(opt.wd, filepath.Base(opt.wd))
	if runtime.GOOS == "windows" {
		opt.appName += ".exe"
	}

	// 根据 wd 获取项目根目录，从而拿到需要监视的列表
	//
	// TODO 处理 go.work 中的内容
	if opt.paths, err = recursivePaths(opt.wd); err != nil {
		return err
	}

	opt.sanitizeExts()

	opt.appArgs = splitArgs(opt.AppArgs)

	if opt.WatcherFrequency == 0 {
		opt.WatcherFrequency = MinWatcherFrequency
	} else if opt.WatcherFrequency < MinWatcherFrequency {
		return errors.New("watcherFrequency 值过小")
	}

	// 初始化 goCmd 的参数
	args := []string{"build", "-o", opt.appName}
	args = append(args, opt.Args...)
	if len(opt.MainFiles) > 0 {
		args = append(args, opt.MainFiles)
	}
	opt.goCmdArgs = args

	return nil
}

func getWD(mainFiles string) (wd string, err error) {
	if wd, err = filepath.Abs(mainFiles); err != nil {
		return "", err
	}

	// 不以路径分隔符结尾的，wd 可能表示的是文件，而不是目录。
	if last := wd[len(wd)-1]; last != '/' && last != filepath.Separator {
		stat, err := os.Stat(wd)
		if err != nil || !stat.IsDir() { // err!=nil 可能 wd 不是一个正常的文件表示，比如 ./*.go
			wd = filepath.Dir(wd)
		}
	}

	return wd, nil
}

func (opt *Options) sanitizeExts() {
	if len(opt.Exts) == 0 {
		opt.Exts = []string{".go"}
	}

	exts := make([]string, 0, len(opt.Exts))
	for _, ext := range opt.Exts {
		ext = strings.TrimSpace(ext)
		if len(ext) == 0 {
			continue
		}

		if ext == "*" {
			opt.anyExts = true
			return
		}

		if ext[0] != '.' {
			ext = "." + ext
		}
		exts = append(exts, ext)
	}
	opt.Exts = exts
}

// 根据 recursive 值确定是否递归查找 paths 每个目录下的子目录
func recursivePaths(wd string) ([]string, error) {
	p, mod, err := source.ModFile(wd)
	if err != nil {
		return nil, err
	}
	root := filepath.Dir(p)

	dirs := make([]string, 0, len(mod.Replace)+1)
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		if root != p {
			if d.Name()[0] == '.' { // 隐藏目录
				return fs.SkipDir
			}

			stat, err := os.Stat(filepath.Join(p, "go.mod"))
			if err == nil && !stat.IsDir() {
				return fs.SkipDir
			}
		}

		dirs = append(dirs, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func splitArgs(args string) []string {
	ret := make([]string, 0, 10)
	var state byte
	var start, index int

	for index = 0; index < len(args); index++ {
		b := args[index]
		switch b {
		case ' ':
			if state == '"' {
				break
			}

			if state != ' ' {
				ret = appendArg(ret, args[start:index])
				state = ' '
			}
			start = index + 1
		case '=':
			if state == '"' {
				break
			}

			if state != '=' {
				ret = appendArg(ret, args[start:index])
				state = '='
			}
			start = index + 1
			state = 0
		case '"':
			if state == '"' {
				ret = appendArg(ret, args[start:index])
				state = 0
				start = index + 1
				break
			}

			if start != index {
				ret = appendArg(ret, args[start:index])
			}
			state = '"'
			start = index + 1
		default:
			if state == ' ' {
				state = 0
				start = index
			}
		}
	} // end for

	if start < len(args) {
		ret = appendArg(ret, args[start:])
	}

	return ret
}

func appendArg(args []string, arg string) []string {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return args
	}

	return append(args, arg)
}
