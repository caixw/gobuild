// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"strings"

	"github.com/issue9/term/v3/colors"

	"github.com/caixw/gobuild/watch"
)

type console struct {
	Logs       chan *watch.Log
	showIgnore bool
	writers    map[int8]*consoleWriter
	stop       chan struct{}
}

func newConsoleLogs(showIgnore bool, err, out io.Writer) *console {
	logs := &console{
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
func (logs *console) Stop() { logs.stop <- struct{}{} }

func (logs *console) output() {
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
