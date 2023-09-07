// SPDX-License-Identifier: MIT

package watch

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

func TestWatch(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{
		Logger:    NewConsoleLogger(true, io.Discard, io.Discard),
		MainFiles: "./testdir/main.go",
	}
	a.NotError(opt.sanitize())

	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan bool, 1)
	go func() {
		a.NotError(Watch(ctx, opt))
		exit <- true
	}()
	time.Sleep(500 * time.Millisecond) // 等待 go func() 启动
	cancel()
	time.Sleep(500 * time.Millisecond) // 等待 cancel 完成
	<-exit
	os.Remove("./testdir/testdir")
}
