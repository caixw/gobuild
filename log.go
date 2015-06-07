// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/issue9/term/colors"
)

type logLevel int

const (
	succ logLevel = iota
	info
	warn
	erro
	ignore
	max // 永远在最后，用于判断logLevel的值有没有超标
)

var levelStrings = map[logLevel]string{
	succ:   "SUCCESS",
	info:   "INFO",
	warn:   "WARINNG",
	erro:   "ERROR",
	ignore: "IGNORE",
}

var levelColors = map[logLevel]colors.Color{
	succ:   colors.Green,
	info:   colors.Blue,
	warn:   colors.Magenta,
	erro:   colors.Red,
	ignore: colors.Default,
}

// 输出日志信息
func log(level logLevel, msg ...interface{}) {
	if level < 0 || level >= max {
		panic("log:无效的level值")
	}

	data := time.Now().Format("2006-01-02 15:04:05 ")
	colors.Print(colors.Stdout, colors.Default, colors.Default, data)
	colors.Print(colors.Stdout, levelColors[level], colors.Default, "[", levelStrings[level], "] ")
	colors.Println(colors.Stdout, levelColors[level], colors.Default, msg...)
}
