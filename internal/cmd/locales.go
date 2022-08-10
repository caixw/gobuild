// SPDX-License-Identifier: MIT

package cmd

import "embed"

const usage = `gobuild 是 Go 的热编译工具，监视文件变化，并编译和运行程序。

命令行语法：
 gobuild [options] [dependents]

 options:

%s

 dependents:
  指定其它依赖的目录，只能出现在命令的尾部。


常见用法:

 gobuild
   监视当前目录，若有变动，则重新编译当前目录下的 *.go 文件；

 gobuild -main=main.go
   监视当前目录，若有变动，则重新编译当前目录下的 main.go 文件；

 gobuild -main="main.go" dir1 dir2
   监视当前目录及 dir1 和 dir2，若有变动，则重新编译当前目录下的 main.go 文件；


NOTE: 不会监视隐藏文件和隐藏目录下的文件。

源代码采用 MIT 开源许可证，并发布于 https://github.com/caixw/gobuild`

//go:embed *.yaml
var localeFS embed.FS
