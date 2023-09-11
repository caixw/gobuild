// SPDX-License-Identifier: MIT

package watch

import (
	"io"
	"strings"

	"github.com/issue9/localeutil"
	"github.com/issue9/term/v3/colors"
)

// 日志类型
const (
	System = "system" // gobuild 系统信息
	Go     = "go"     // go 编译器的信息
	App    = "app"    // 被编译程序的信息
)

const (
	Success int8 = iota
	Info
	Warn
	Error
	Ignore // 默认情况下被忽略的信息，一般内容比较多，且价值不高的内容会显示在此通道。
)

var defaultColors = map[int8]colors.Color{
	Success: colors.Green,
	Info:    colors.Blue,
	Warn:    colors.Yellow,
	Error:   colors.Red,
	Ignore:  colors.Default,
}

type (
	// Logger 热编译过程中的日志接收对象
	Logger interface {
		// Output 输出日志内容
		//
		// source 表示信息来源；
		// t 表示信息类型；
		Output(source string, t int8, message string)
	}

	loggerWriter struct {
		s string
		t int8
		w Logger
	}

	consoleLogger struct {
		out        io.Writer
		showIgnore bool
		colors     map[int8]colors.Color
		sources    map[string]string
	}
)

func (w *loggerWriter) Write(bs []byte) (int, error) {
	w.w.Output(w.s, w.t, string(bs))
	return len(bs), nil
}

func asWriter(s string, t int8, w Logger) io.Writer {
	return &loggerWriter{s: s, t: t, w: w}
}

func (c *consoleLogger) Output(source string, t int8, msg string) {
	if !c.showIgnore && t == Ignore {
		return
	}

	if s, found := c.sources[source]; found {
		source = s
	}
	source = "[" + source + "] "

	colors.Fprint(c.out, colors.Normal, c.colors[t], colors.Default, source)
	msg = strings.TrimRight(msg, "\n")
	colors.Fprintln(c.out, colors.Normal, colors.Default, colors.Default, msg)
}

// NewConsoleLogger 将日志输出到控制台的 Logger 实现
//
// colors 表示各类日志的颜色值；
// sources 表示各类信息源的名称；
// colors 和 prefixes 可以为 nil，会采用默认值。
func NewConsoleLogger(showIgnore bool, out io.Writer, colors map[int8]colors.Color, sources map[string]string) Logger {
	if colors == nil {
		colors = defaultColors
	}

	return &consoleLogger{
		out:        out,
		showIgnore: showIgnore,
		colors:     colors,
		sources:    sources,
	}
}

func (b *builder) systemLog(typ int8, msg localeutil.LocaleStringer) {
	b.logs.Output(System, typ, msg.LocaleString(b.p))
}
