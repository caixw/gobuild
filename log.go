// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"io"
	"log"
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
)

// ConsoleLogs 将日志输出到控制台
type ConsoleLogs struct {
	Logs   chan *Log
	ignore bool // 是否忽略 ignore 通道的日志
	succ   *log.Logger
	info   *log.Logger
	warn   *log.Logger
	erro   *log.Logger
	igno   *log.Logger
}

// NewConsoleLogs 声明 ConsoleLogs 实例
func NewConsoleLogs(ignore bool) *ConsoleLogs {
	logs := &ConsoleLogs{
		succ:   log.New(newWriter(os.Stdout, colors.Green, "[SUCC]"), "", log.Ltime),
		info:   log.New(newWriter(os.Stdout, colors.Blue, "[INFO]"), "", log.Ltime),
		warn:   log.New(newWriter(os.Stderr, colors.Magenta, "[WARN]"), "", log.Ltime),
		erro:   log.New(newWriter(os.Stderr, colors.Red, "[ERRO]"), "", log.Ltime),
		igno:   log.New(newWriter(os.Stderr, colors.Default, "[IGNO]"), "", log.Ltime),
		Logs:   make(chan *Log, 100),
		ignore: ignore,
	}

	go logs.output()

	return logs
}

func (logs *ConsoleLogs) output() {
	for log := range logs.Logs {
		switch log.Type {
		case LogTypeError:
			logs.erro.Println(log.Message)
		case LogTypeIgnore:
			if !logs.ignore {
				logs.igno.Println(log.Message)
			}
		case LogTypeInfo:
			logs.info.Println(log.Message)
		case LogTypeSuccess:
			logs.succ.Println(log.Message)
		case LogTypeWarn:
			logs.warn.Println(log.Message)
		default:
			panic("无效的日志类型")
		}
	}
}

// 带色彩输出的控制台。
type logWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func newWriter(out io.Writer, color colors.Color, prefix string) io.Writer {
	return &logWriter{
		out:    out,
		color:  color,
		prefix: prefix,
	}
}

func (w *logWriter) Write(bs []byte) (int, error) {
	colors.Fprint(w.out, w.color, colors.Default, w.prefix)
	return colors.Fprint(w.out, colors.Default, colors.Default, string(bs))
}
