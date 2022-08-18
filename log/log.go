// SPDX-License-Identifier: MIT

// Package log 输出的日志管理
package log

import "io"

// 日志类型
const (
	Success int8 = iota
	Info
	Warn
	Error
	Ignore
	App
	Go
)

type Log struct {
	Type    int8
	Message string
}

type writer struct {
	t int8
	w chan<- *Log
}

func (w *writer) Write(bs []byte) (int, error) {
	w.w <- &Log{
		Type:    w.t,
		Message: string(bs),
	}
	return len(bs), nil
}

func AsWriter(t int8, w chan<- *Log) io.Writer { return &writer{t: t, w: w} }
