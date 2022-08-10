// SPDX-License-Identifier: MIT

// Package watch 监视文件变化并编译
package watch

import (
	"context"

	"github.com/caixw/gobuild/log"
)

// Watch 执行热编译服务
//
// 如果初始化参数有误，则反错误信息，如果是编译过程中出错，将直接将错误内容输出到 logs。
func Watch(ctx context.Context, logs chan<- *log.Log, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}

	b, err := opt.newBuilder(logs)
	if err != nil {
		return err
	}

	b.logf(log.Info, "当前环境参数如下：%s", b.env)

	b.logf(log.Info, "给程序传递了以下参数：%s", b.appArgs) // 输出提示信息

	switch { // 提示扩展名
	case len(b.exts) == 0: // 允许不监视任意文件，但输出一信息来警告
		b.logf(log.Warn, "将 ext 设置为空值，意味着不监视任何文件的改变！")
	case len(b.exts) > 0:
		b.logf(log.Info, "系统将监视以下类型的文件：%s", b.exts)
	}

	b.logf(log.Info, "输出文件为：%s", b.appName) // 提示 appName

	w, err := b.initWatcher(opt.paths)
	if err != nil {
		return err
	}
	defer w.Close()

	go b.build() // 第一次主动编译程序，后续的才是监视变化。

	b.watch(ctx, w)
	return nil
}
