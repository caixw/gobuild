// SPDX-License-Identifier: MIT

package gobuild

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MinWatcherFrequency 监视器更新频率的最小值
const MinWatcherFrequency = 1 * time.Second

// Options 热编译的选项
type Options struct {
	// 为 go build 最后的文件参数，可以为空，表示当前目录。
	MainFiles string

	// 指定可执行文件输出的文件路径，为空表示默认值，
	// 若不带路径信息，会附加在 Dirs 的第一个路径上；
	//
	// windows 系统无须指定 .exe 扩展名，会自行添加。
	//
	// 如果带路径信息，则会使用该文件所在目录作为工作目录。
	OutputName string
	appName    string

	// 传递各个工具的参数，大致有以下几个，具体可参考 go build 的 xxflags 系列参数。
	//  - asm   --> asmflags
	//  - gccgo --> gccgoflags
	//  - gc    --> gcflags
	//  - ld    --> ldflags
	Flags map[string]string

	// 指定监视的文件扩展名，为空表示不监视任何文件，* 表示监视所有文件
	Exts string
	exts []string

	// 传递给编译成功后的程序的参数
	AppArgs string
	appArgs []string

	// 是否监视子目录
	Recursive bool

	// 表示需要监视的目录，至少指定一个目录，
	// 第一个目录被当作主目录，将编译其下的文件作为执行主体。
	//
	// 如果 OutputName 中未指定目录的话，第一个目录会被当作工作目录使用。
	Dirs  []string
	paths []string

	// 监视器的更新频率，只有文件更新的时长超过此值，才会被定义为更新。
	// 防止文件频繁修改导致的频繁编译调用。
	//
	// 此值不能小于 MinWatcherFrequency。
	//
	// 默认值为 MinWatcherFrequency 表示的值。
	WatcherFrequency time.Duration

	// 传递给 go 命令的参数
	goCmdArgs []string
}

func (opt *Options) sanitize() error {
	if len(opt.Dirs) < 1 {
		return errors.New("参数 dir 至少指定一个")
	}
	wd, err := filepath.Abs(opt.Dirs[0])
	if err != nil {
		return err
	}
	opt.Dirs[0] = wd

	if opt.appName, err = getAppName(opt.OutputName, wd); err != nil {
		return err
	}

	opt.exts = getExts(opt.Exts)

	opt.appArgs = splitArgs(opt.AppArgs)

	if opt.paths, err = recursivePaths(opt.Recursive, opt.Dirs); err != nil {
		return err
	}

	if opt.WatcherFrequency == 0 {
		opt.WatcherFrequency = MinWatcherFrequency
	} else if opt.WatcherFrequency < MinWatcherFrequency {
		return errors.New("watcherFrequency 值过小")
	}

	// 初始化 goCmd 的参数
	args := []string{"build", "-o", opt.appName}
	for k, v := range opt.Flags {
		args = append(args, "-"+k+"flags", v)
	}
	args = append(args, "-v")
	if len(opt.MainFiles) > 0 {
		args = append(args, opt.MainFiles)
	}
	opt.goCmdArgs = args

	return nil
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

	goexe := os.Getenv("GOEXE")
	if goexe != "" && !strings.HasSuffix(outputName, goexe) {
		outputName += goexe
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

		if fi.IsDir() && (strings.Index(path, "/.") < 0 || strings.Index(path, "\\.") < 0) { // 非隐藏的目录
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
