// SPDX-License-Identifier: MIT

package gobuild

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestBuild(t *testing.T) {
	a := assert.New(t)

	opt := &Options{Dirs: []string{"./testdir"}}
	a.NotError(opt.sanitize())

	logs := NewConsoleLogs(true)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		a.Equal(Build(ctx, logs.Logs, opt), context.Canceled)
	}()
	cancel()
	time.Sleep(500 * time.Microsecond) // 等待完成
}

func TestOptions_newBuilder(t *testing.T) {
	a := assert.New(t)

	opt := &Options{Dirs: []string{"./"}}
	a.NotError(opt.sanitize())

	b, err := opt.newBuilder(nil)
	a.NotError(err).NotNil(b)
	a.Contains(b.env, runtime.Version())
}
