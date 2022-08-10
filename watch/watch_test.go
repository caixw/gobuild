// SPDX-License-Identifier: MIT

package watch

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/gobuild/log"
)

func TestWatch(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{
		Dirs:       []string{"./testdir"},
		MainFiles:  "./testdir/main.go",
		OutputName: "outputName",
	}
	a.NotError(opt.sanitize())

	logs := log.NewConsole(true)
	defer logs.Stop()

	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan bool, 1)
	go func() {
		a.NotError(Watch(ctx, logs.Logs, opt))
		exit <- true
	}()
	time.Sleep(500 * time.Millisecond) // 等待 go func() 启动
	cancel()
	time.Sleep(500 * time.Millisecond) // 等待 cancel 完成
	<-exit
	a.NotError(os.Remove("./testdir/outputName"))
}
