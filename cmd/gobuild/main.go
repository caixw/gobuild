// SPDX-License-Identifier: MIT

// 一个简单的 Go 语言热编译工具
//
// 监视指定目录(可同时监视多个目录)下文件的变化，触发`go build`指令，
// 实时编译指定的 Go 代码，并在编译成功时运行该程序。
// 具体命令格式可使用`gobuild -h`来查看。
package main

import "github.com/caixw/gobuild/internal/cmd"

func main() {
	cmd.Exec()
}
