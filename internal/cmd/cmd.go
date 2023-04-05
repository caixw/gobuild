// SPDX-License-Identifier: MIT

package cmd

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild/locales"
)

//go:embed *.yaml
var localeFS embed.FS

const (
	url     = "https://github.com/caixw/gobuild"
	license = "MIT"
)

func Exec() error {
	p := getPrinter()

	usage := p.Sprintf("cmd.usage %s %s", license, url)

	o := cmdopt.New(os.Stdout, flag.ExitOnError, usage, nil, func(s string) string { return p.Sprintf("未找到子命令 %s") })

	initVersion(o, p)
	initWatch(o, p)
	initInit(o, p)
	o.New("help", p.Sprintf("显示帮助信息"), p.Sprintf("显示帮助信息"), cmdopt.Help(o))
	return o.Exec(os.Args[1:])
}

func getPrinter() *message.Printer {
	tag, _ := localeutil.DetectUserLanguageTag()
	c := catalog.NewBuilder(catalog.Fallback(tag))

	if err := localeutil.LoadMessageFromFSGlob(c, &localeFS, "*.yaml", yaml.Unmarshal); err != nil {
		panic(err)
	}
	if err := localeutil.LoadMessageFromFSGlob(c, locales.Locales, "*.yaml", yaml.Unmarshal); err != nil {
		panic(err)
	}
	p, err := os.Executable()
	if err != nil { // 这里不退出
		fmt.Fprintln(os.Stderr, err)
	}
	if err := localeutil.LoadMessageFromFSGlob(c, os.DirFS(filepath.Dir(p)), "*.yaml", yaml.Unmarshal); err != nil {
		panic(err)
	}

	return message.NewPrinter(tag, message.Catalog(c))
}
