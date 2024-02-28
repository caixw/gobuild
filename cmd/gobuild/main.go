// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:generate web locale -l=und -m -f=yaml ./
//go:generate web update-locale -src=./locales/und.yaml -dest=./locales/zh-CN.yaml,./locales/zh-TW.yaml

// 一个简单的 Go 语言热编译工具
//
// 监视指定目录(可同时监视多个目录)下文件的变化，触发 go build 指令，
// 实时编译指定的 Go 代码，并在编译成功时运行该程序。
// 具体命令格式可使用 gobuild help 来查看。
package main

func main() {
	if err := Exec(); err != nil {
		panic(err)
	}
}
