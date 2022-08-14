// SPDX-License-Identifier: MIT

package log

import (
	"io"
	"os"

	"github.com/issue9/term/v3/colors"
)

// Console 将日志输出到控制台
type Console struct {
	Logs       chan *Log
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
		Logs:       make(chan *Log, 100),
		showIgnore: showIgnore,
		stop:       make(chan struct{}, 1),
		writers: map[int8]*consoleWriter{
			Success: newWriter(out, colors.Green, "[SUCC] "),
			Info:    newWriter(out, colors.Blue, "[INFO] "),
			Warn:    newWriter(err, colors.Magenta, "[WARN] "),
			Error:   newWriter(err, colors.Red, "[ERRO] "),
			Ignore:  newWriter(out, colors.Default, "[IGNO] "),
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
			if !logs.showIgnore && log.Type == Ignore {
				continue
			}

			w := logs.writers[log.Type]
			colors.Fprint(w.out, colors.Normal, w.color, colors.Default, w.prefix)
			colors.Fprintln(w.out, colors.Normal, colors.Default, colors.Default, log.Message)
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
