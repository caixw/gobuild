// SPDX-License-Identifier: MIT

package gobuild

import (
	"context"

	"golang.org/x/text/message"

	"github.com/caixw/gobuild/internal/config"
	"github.com/caixw/gobuild/watch"
)

type WatchOptions = watch.Options

// Watch 监视文件变化执行热编译服务
//
// 如果初始化参数有误，则反错误信息，如果是编译过程中出错，将直接将错误内容输出到 [watch.Logger]。
func Watch(ctx context.Context, o *WatchOptions) error { return watch.Watch(ctx, o) }

// Init 初始化一个空的项目
//
// wd 为工作目录，将在此目录下初始化项目；
// name 为 go.mod 中定义的模块的名称。
// name 的最后一个元素会作为名称在 wd 指定的目录下创建子目录，
// 同时在子目录下会添加以下内容：
//   - go.mod 以 name 作为模块名；
//   - .gobuild.yaml 为 gobuild 的配置文件；
//   - cmd/{base}/{base}.go 程序入口 main 函数，base 为 name 的最后一个元素；
func Init(wd, name string) error { return config.Init(wd, name) }

// WatchConfig 监视配置文件
//
// 如果配置文件发生变化，那么重启热编译程序；
// 如果 wd 中不存在配置文件，则向上一级查找，如果一直未找到，将返回 fs.ErrNotExists 错误。
func WatchConfig(wd string, p *message.Printer, logs watch.Logger) error {
	return config.Watch(wd, p, logs)
}
