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
	"github.com/issue9/localeutil/message/serialize"
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

const helpUsage = localeutil.StringPhrase("显示帮助信息")

func Exec() error {
	p := getPrinter()

	usage := p.Sprintf("cmd.usage %s %s", license, url)

	o := cmdopt.New(os.Stdout, flag.ExitOnError, usage, nil, func(s string) string {
		return localeutil.Phrase("未找到子命令 %s").LocaleString(p)
	})

	initVersion(o, p)
	initWatch(o, p)
	initInit(o, p)
	cmdopt.Help(o, "help", helpUsage.LocaleString(p), helpUsage.LocaleString(p))
	return o.Exec(os.Args[1:])
}

func getPrinter() *localeutil.Printer {
	tag, _ := localeutil.DetectUserLanguageTag()
	c := catalog.NewBuilder(catalog.Fallback(tag))

	l1, err := serialize.LoadFSGlob(&localeFS, "*.yaml", yaml.Unmarshal)
	if err != nil {
		panic(err)
	}
	for _, l := range l1 {
		if err := l.Catalog(c); err != nil {
			panic(err)
		}
	}

	l2, err := serialize.LoadFSGlob(&locales.Locales, "*.yaml", yaml.Unmarshal)
	if err != nil {
		panic(err)
	}
	for _, l := range l2 {
		if err := l.Catalog(c); err != nil {
			panic(err)
		}
	}

	p, err := os.Executable()
	if err != nil { // 这里不退出
		fmt.Fprintln(os.Stderr, err)
	}

	l3, err := serialize.LoadFSGlob(os.DirFS(filepath.Dir(p)), "*.yaml", yaml.Unmarshal)
	if err != nil {
		panic(err)
	}
	for _, l := range l3 {
		if err := l.Catalog(c); err != nil {
			panic(err)
		}
	}

	return message.NewPrinter(tag, message.Catalog(c))
}
