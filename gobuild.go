// SPDX-License-Identifier: MIT

package gobuild

import (
	"context"

	"github.com/caixw/gobuild/log"
	"github.com/caixw/gobuild/watch"
)

type (
	WatchOptions = watch.Options
	Log          = log.Log
)

// Watch 监视文件变化执行热编译服务
func Watch(ctx context.Context, logs chan *Log, opt *WatchOptions) error {
	return watch.Watch(ctx, logs, opt)
}
