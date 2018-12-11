// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"io"
	"os"

	"github.com/issue9/term/colors"
)

// Log 日志类型
type Log struct {
	Type    int8
	Message string
}

// 日志类型
const (
	LogTypeSuccess int8 = iota + 1
	LogTypeInfo
	LogTypeWarn
	LogTypeError
	LogTypeIgnore
	logTypeSize
)

// ConsoleLogs 将日志输出到控制台
type ConsoleLogs struct {
	Logs    chan *Log
	ignore  bool // 是否忽略 ignore 通道的日志
	writers map[int8]*logWriter
}

// NewConsoleLogs 声明 ConsoleLogs 实例
func NewConsoleLogs(ignore bool) *ConsoleLogs {
	logs := &ConsoleLogs{
		Logs:   make(chan *Log, 100),
		ignore: ignore,
		writers: map[int8]*logWriter{
			LogTypeSuccess: newWriter(os.Stdout, colors.Green, "[SUCC]"),
			LogTypeInfo:    newWriter(os.Stdout, colors.Blue, "[INFO]"),
			LogTypeWarn:    newWriter(os.Stderr, colors.Magenta, "[WARN]"),
			LogTypeError:   newWriter(os.Stderr, colors.Red, "[ERRO]"),
			LogTypeIgnore:  newWriter(os.Stderr, colors.Default, "[IGNO]"),
		},
	}

	go logs.output()

	return logs
}

func (logs *ConsoleLogs) output() {
	for log := range logs.Logs {
		w := logs.writers[log.Type]
		colors.Fprint(w.out, w.color, colors.Default, w.prefix)
		colors.Fprintln(w.out, colors.Default, colors.Default, log.Message)
	}
}

// 带色彩输出的控制台。
type logWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func newWriter(out io.Writer, color colors.Color, prefix string) *logWriter {
	return &logWriter{
		out:    out,
		color:  color,
		prefix: prefix,
	}
}
