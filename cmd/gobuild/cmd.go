// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
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

	cl "github.com/caixw/gobuild/cmd/locales"
	"github.com/caixw/gobuild/locales"
)

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
	p, err := os.Executable()
	if err != nil { // 这里不退出
		fmt.Fprintln(os.Stderr, err)
	}

	ls, err := serialize.LoadFSGlob(func(s string) serialize.UnmarshalFunc {
		return yaml.Unmarshal
	}, "*.yaml", cl.Locales, locales.Locales, os.DirFS(filepath.Dir(p)))
	if err != nil {
		panic(err)
	}

	tag, _ := localeutil.DetectUserLanguageTag()
	c := catalog.NewBuilder(catalog.Fallback(tag))

	for _, l := range ls {
		if err := l.Catalog(c); err != nil {
			panic(err)
		}
	}

	return message.NewPrinter(tag, message.Catalog(c))
}
