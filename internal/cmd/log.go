// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/issue9/term/v3/colors"

	"github.com/caixw/gobuild/watch"
)

// Console 将日志输出到控制台
type Console struct {
	Logs       chan *watch.Log
	showIgnore bool
	writers    map[int8]*consoleWriter
	stop       chan struct{}
}

// NewConsole 声明 ConsoleLogs 实例
func NewConsole(showIgnore bool) *Console {
	return newConsoleLogs(showIgnore, os.Stderr, os.Stdout)
}

func newConsoleLogs(showIgnore bool, err, out io.Writer) *Console {
	logs := &Console{
		Logs:       make(chan *watch.Log, 100),
		showIgnore: showIgnore,
		stop:       make(chan struct{}, 1),
		writers: map[int8]*consoleWriter{
			watch.LogTypeSuccess: newWriter(out, colors.Green, "[SUCC] "),
			watch.LogTypeInfo:    newWriter(out, colors.Blue, "[INFO] "),
			watch.LogTypeWarn:    newWriter(err, colors.Magenta, "[WARN] "),
			watch.LogTypeError:   newWriter(err, colors.Red, "[ERRO] "),
			watch.LogTypeIgnore:  newWriter(out, colors.Default, "[IGNO] "),
			watch.LogTypeApp:     newWriter(out, colors.Default, "[APP] "),
			watch.LogTypeGo:      newWriter(out, colors.Default, "[GO] "),
		},
	}

	go logs.output()

	return logs
}

// Stop 停止输出
func (logs *Console) Stop() { logs.stop <- struct{}{} }

func (logs *Console) output() {
	for {
		select {
		case log := <-logs.Logs:
			if !logs.showIgnore && log.Type == watch.LogTypeIgnore {
				continue
			}

			w := logs.writers[log.Type]
			colors.Fprint(w.out, colors.Normal, w.color, colors.Default, w.prefix)
			msg := strings.TrimRight(log.Message, "\n")
			colors.Fprintln(w.out, colors.Normal, colors.Default, colors.Default, msg)
		case <-logs.stop:
			return
		}
	}
}

// 带色彩输出的控制台
type consoleWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func newWriter(out io.Writer, color colors.Color, prefix string) *consoleWriter {
	return &consoleWriter{
		out:    out,
		color:  color,
		prefix: prefix,
	}
}
