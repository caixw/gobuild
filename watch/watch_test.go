// SPDX-License-Identifier: MIT

package watch

import (
	"context"
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/gobuild/log"
)

func TestWatch(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{Dirs: []string{"./testdir"}}
	a.NotError(opt.sanitize())

	logs := log.NewConsole(true)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		a.Equal(Watch(ctx, logs.Logs, opt), context.Canceled)
	}()
	cancel()
	time.Sleep(500 * time.Microsecond) // 等待完成
}
