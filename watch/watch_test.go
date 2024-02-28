// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package watch

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestWatch(t *testing.T) {
	a := assert.New(t, false)
	l := NewConsoleLogger(true, io.Discard, nil, nil)

	opt := &Options{
		MainFiles: "./testdir/main.go",
	}
	a.NotError(opt.sanitize())

	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan bool, 1)
	go func() {
		a.NotError(Watch(ctx, message.NewPrinter(language.SimplifiedChinese), l, opt))
		exit <- true
	}()
	time.Sleep(500 * time.Millisecond) // 等待 go func() 启动
	cancel()
	time.Sleep(500 * time.Millisecond) // 等待 cancel 完成
	<-exit
	os.Remove("./testdir/testdir")
}
