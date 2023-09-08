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
		writers    map[int8]*consoleWriter
	}

	consoleWriter struct {
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
	colors.Fprint(c.out, colors.Normal, w.color, colors.Default, w.prefix)
	msg = strings.TrimRight(msg, "\n")
	colors.Fprintln(c.out, colors.Normal, colors.Default, colors.Default, msg)
}

// NewConsoleLogger 将日志输出到控制台的 Logger 实现
func NewConsoleLogger(showIgnore bool, out io.Writer) Logger {
	newCW := func(out io.Writer, color colors.Color, prefix string) *consoleWriter {
		return &consoleWriter{
			color:  color,
			prefix: prefix,
		}
	}

	return &consoleLogger{
		out:        out,
		showIgnore: showIgnore,
		writers: map[int8]*consoleWriter{
			LogTypeSuccess: newCW(out, colors.Green, "[SUCC] "),
			LogTypeInfo:    newCW(out, colors.Blue, "[INFO] "),
			LogTypeWarn:    newCW(out, colors.Yellow, "[WARN] "),
			LogTypeError:   newCW(out, colors.Red, "[ERRO] "),
			LogTypeIgnore:  newCW(out, colors.Default, "[IGNO] "),
			LogTypeApp:     newCW(out, colors.Magenta, "[APP] "),
			LogTypeGo:      newCW(out, colors.Cyan, "[GO] "),
		},
	}
}
