// SPDX-License-Identifier: MIT

//go:generate web locale -l=und -m -f=yaml ./
//go:generate web update-locale -src=./locales/und.yaml -dest=./locales/zh-CN.yaml,./locales/zh-TW.yaml

// Package gobuild 热编译 Go 代码
package gobuild

import (
	"context"
	"io"

	"golang.org/x/text/message"

	"github.com/caixw/gobuild/internal/config"
	"github.com/caixw/gobuild/watch"
	"github.com/issue9/term/v3/colors"
)

type (
	WatchOptions = watch.Options
	Logger       = watch.Logger
)

// Watch 监视文件变化执行热编译服务
//
// p 用于处理本地化的错误信息；
// l 用于输出错误信息；
// 如果初始化参数有误，则反错误信息，如果是编译过程中出错，将直接将错误内容输出到 [Logger]。
func Watch(ctx context.Context, p *message.Printer, l Logger, o *WatchOptions) error {
	return watch.Watch(ctx, p, l, o)
}

// Init 初始化一个空的项目
//
// wd 为工作目录，将在此目录下初始化项目；
// configFilename 为配置文件的文件名；
// name 为 go.mod 中定义的模块的名称。
// name 的最后一个元素会作为名称在 wd 指定的目录下创建子目录，
// 同时在子目录下会添加以下内容：
//   - go.mod 以 name 作为模块名；
//   - configFilename 指定的文件名作为 gobuild 的配置文件；
//   - cmd/{base}/{base}.go 程序入口 main 函数，base 为 name 的最后一个元素；
func Init(wd, name, configFilename string) error { return config.Init(wd, name, configFilename) }

// NewConsoleLogger 将日志输出到控制台的 Logger 实现
//
// colors 和 prefixes 可以为 nil，会采用默认值。
func NewConsoleLogger(showIgnore bool, out io.Writer, colors map[int8]colors.Color, prefixes map[int8]string) Logger {
	return watch.NewConsoleLogger(showIgnore, out, colors, prefixes)
}
