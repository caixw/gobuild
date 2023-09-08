// SPDX-License-Identifier: MIT

package watch

import (
	"io"
	"strings"

	"github.com/issue9/term/v3/colors"
)

// 日志类型
const (
	LogTypeSuccess int8 = iota
	LogTypeInfo
	LogTypeWarn
	LogTypeError
	LogTypeIgnore // 默认情况下被忽略的信息，一般内容比较多，且价值不高的内容会显示在此通道。
	LogTypeApp    // 被编译程序返回的信息
	LogTypeGo     // Go 编译器返回的信息
)

var (
	defaultColors = map[int8]colors.Color{
		LogTypeSuccess: colors.Green,
		LogTypeInfo:    colors.Blue,
		LogTypeWarn:    colors.Yellow,
		LogTypeError:   colors.Red,
		LogTypeIgnore:  colors.Default,
		LogTypeApp:     colors.Magenta,
		LogTypeGo:      colors.Cyan,
	}

	defaultPrefixes = map[int8]string{
		LogTypeSuccess: "[SUCC] ",
		LogTypeInfo:    "[INFO] ",
		LogTypeWarn:    "[WARN] ",
		LogTypeError:   "[ERRO] ",
		LogTypeIgnore:  "[IGNO] ",
		LogTypeApp:     "[APP] ",
		LogTypeGo:      "[GO] ",
	}
)

type (
	// Logger 热编译过程中的日志接收对象
	Logger interface {
		// Output 输出日志内容
		//
		// t 表示日志类型，一般表示日志的重要程度或是日志的来源信息。
		Output(t int8, message string)
	}

	loggerWriter struct {
		t int8
		w Logger
	}

	consoleLogger struct {
		out        io.Writer
		showIgnore bool
		colors     map[int8]colors.Color
		prefixes   map[int8]string
	}
)

func (w *loggerWriter) Write(bs []byte) (int, error) {
	w.w.Output(w.t, string(bs))
	return len(bs), nil
}

func asWriter(t int8, w Logger) io.Writer { return &loggerWriter{t: t, w: w} }

func (c *consoleLogger) Output(t int8, msg string) {
	if !c.showIgnore && t == LogTypeIgnore {
		return
	}

	colors.Fprint(c.out, colors.Normal, c.colors[t], colors.Default, c.prefixes[t])
	msg = strings.TrimRight(msg, "\n")
	colors.Fprintln(c.out, colors.Normal, colors.Default, colors.Default, msg)
}

// NewConsoleLogger 将日志输出到控制台的 Logger 实现
//
// colors 和 prefixes 可以为 nil，会采用默认值。
func NewConsoleLogger(showIgnore bool, out io.Writer, colors map[int8]colors.Color, prefixes map[int8]string) Logger {
	if colors == nil {
		colors = defaultColors
	}
	if prefixes == nil {
		prefixes = defaultPrefixes
	}

	return &consoleLogger{
		out:        out,
		showIgnore: showIgnore,
		colors:     colors,
		prefixes:   prefixes,
	}
}
