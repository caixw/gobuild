// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/issue9/term/colors"
)

var (
	succ   = log.New(&logWriter{out: os.Stdout, color: colors.Green, prefix: "[SUCC]"}, "", log.Ltime)
	info   = log.New(&logWriter{out: os.Stdout, color: colors.Blue, prefix: "[INFO]"}, "", log.Ltime)
	warn   = log.New(&logWriter{out: os.Stderr, color: colors.Magenta, prefix: "[WARN]"}, "", log.Ltime)
	erro   = log.New(&logWriter{out: os.Stderr, color: colors.Red, prefix: "[ERRO]"}, "", log.Ltime)
	ignore = log.New(ioutil.Discard, "", log.Ltime) // 默认情况下不显示此类信息，全部发送到 Discard
)

// 带色彩输出的控制台。
type logWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func (w *logWriter) Write(bs []byte) (int, error) {
	colors.Fprint(w.out, w.color, colors.Default, w.prefix)
	return colors.Fprint(w.out, colors.Default, colors.Default, string(bs))
}
