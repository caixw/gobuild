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
	LogTypeIgnore // 默认情况下被忽略的信息，一般内容比较多，且价格不高的内容会显示在此通道。
	LogTypeApp    // 被编译程序返回的信息
	LogTypeGo     // Go 编译器返回的信息
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
		showIgnore bool
		writers    map[int8]*consoleWriter
	}

	consoleWriter struct {
		out    io.Writer
		color  colors.Color
		prefix string
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

	w := c.writers[t]
	colors.Fprint(w.out, colors.Normal, w.color, colors.Default, w.prefix)
	msg = strings.TrimRight(msg, "\n")
	colors.Fprintln(w.out, colors.Normal, colors.Default, colors.Default, msg)
}

// NewConsoleLogger 返回将日志输出到控制台的 Logger 接口实现
func NewConsoleLogger(showIgnore bool, err, out io.Writer) Logger {
	newCW := func(out io.Writer, color colors.Color, prefix string) *consoleWriter {
		return &consoleWriter{
			out:    out,
			color:  color,
			prefix: prefix,
		}
	}

	return &consoleLogger{
		showIgnore: showIgnore,
		writers: map[int8]*consoleWriter{
			LogTypeSuccess: newCW(out, colors.Green, "[SUCC] "),
			LogTypeInfo:    newCW(out, colors.Blue, "[INFO] "),
			LogTypeWarn:    newCW(err, colors.Magenta, "[WARN] "),
			LogTypeError:   newCW(err, colors.Red, "[ERRO] "),
			LogTypeIgnore:  newCW(out, colors.Default, "[IGNO] "),
			LogTypeApp:     newCW(out, colors.Default, "[APP] "),
			LogTypeGo:      newCW(out, colors.Default, "[GO] "),
		},
	}
}
