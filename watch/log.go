// SPDX-License-Identifier: MIT

package watch

import "io"

// 日志类型
const (
	LogTypeSuccess int8 = iota
	LogTypeInfo
	LogTypeWarn
	LogTypeError
	LogTypeIgnore
	LogTypeApp
	LogTypeGo
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

func asWriter(t int8, w chan<- *Log) io.Writer { return &writer{t: t, w: w} }
