// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/issue9/term/colors"
	"time"
)

type logLevel int

const (
	succ logLevel = iota
	info
	warn
	erro
)

func (l logLevel) String() string {
	switch l {
	case succ:
		return "SUCC"
	case info:
		return "INFO"
	case warn:
		return "WARN"
	case erro:
		return "ERROR"
	default:
		return "<UNKNOWN>"
	}
}

func log(level logLevel, msg ...interface{}) {
	data := time.Now().Format("2006-01-02 15:04:05 ")
	colors.Print(colors.Stdout, colors.Default, colors.Default, data)
	colors.Print(colors.Stdout, colors.Red, colors.Default, "[", level, "] ")
	colors.Println(colors.Stdout, colors.Default, colors.Default, msg...)
}
