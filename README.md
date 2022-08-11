# gobuild

[![Latest Release](https://img.shields.io/github/release/caixw/gobuild.svg?style=flat-square)](https://github.com/caixw/gobuild/releases/latest)
[![Test](https://github.com/caixw/gobuild/workflows/Test/badge.svg)](https://github.com/caixw/gobuild/actions?query=workflow%3ATest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/caixw/gobuild)](https://pkg.go.dev/github.com/caixw/gobuild)
![Go version](https://img.shields.io/github/go-mod/go-version/caixw/gobuild)
![License](https://img.shields.io/github/license/caixw/gobuild)
[![codecov](https://codecov.io/gh/caixw/gobuild/branch/master/graph/badge.svg)](https://codecov.io/gh/caixw/gobuild)

gobuild 是一个简单的 Go 代码热编译工具。
会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，并编译和运行程序。

## 命令行语法

主要包含了 watch 和 init 两个子命令。具体的子命令可以通过 `gobuild help` 查看。

### init

初始化项目，添加项目的必备的文件，比如热编译的配置文件 `.gobuild.yaml`。
如果是空目录，还会顺带初始化 `go.mod` 等文件。

```shell
gobuild init github.com/owner/mod
```

### watch

监视文件并进行热编译，热编译的配置项从当前目录下的 `.gobuild.yaml` 加载。

```shell
gobuild watch [options]
```

### 配置文件

配置文件为当前目录下的 `.gobuild.yaml`，可由 `gobuild init` 子命令生成，包含了以下字段：

 字段       | 类型         | 描述
------------|--------------|-------------------------------------
 main       | string       | 指定需要编译的文件，如果为空表示当前目录。
 output     | string       | 指定可执行文件输出的文件路径
 flags      | map          | 指定可执行文件输出的文件路径
 exts       | []string     | 指定监视的文件扩展名
 args       | string       | 传递给编译成功后的程序的参数
 recursive  | bool         | 是否监视子目录
 dirs       | dirs         | 表示需要监视的目录
 freq       | duration     | 监视器的更新频率

## 支持平台

平台支持依赖 [colors](https://github.com/issue9/term) 与 [fsnotify](https://github.com/fsnotify/fsnotify) 两个模块，
目前支持以下平台：windows, linux, macOS, BSD。

## 安装

macOS 和 linux 用户可以直接使用 brew 进行安装：

```shell
brew tap caixw/brew
brew install caixw/brew/gobuild
```

常用平台可以从 <https://github.com/caixw/gobuild/releases> 下载，并将二进制文件放入 `PATH` 即可。

如果不存在你当前平台的二进制，可以自己编译：

```shell
git clone https://github.com/caixw/gobuild.git
cd gobuild
./build.sh
```

## 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
