// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package watch

import (
	"io"
	"testing"

	"github.com/issue9/assert/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestOptions_newBuilder(t *testing.T) {
	a := assert.New(t, false)
	l := NewConsoleLogger(false, io.Discard, nil, nil)
	p := message.NewPrinter(language.Und)

	opt := &Options{}
	a.NotError(opt.sanitize())
	b := opt.newBuilder(p, l)
	a.NotNil(b).False(b.anyExt)

	opt = &Options{Exts: []string{".a", "*", ".b"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder(p, l)
	a.NotNil(b).Equal(b.exts, []string{"*"}).True(b.anyExt)

	opt = &Options{Exts: []string{".a", ".b"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder(p, l)
	a.NotNil(b).Equal(b.exts, []string{".a", ".b"}).False(b.anyExt)
}

func TestBuilder_isIgnore(t *testing.T) {
	a := assert.New(t, false)
	l := NewConsoleLogger(false, io.Discard, nil, nil)
	p := message.NewPrinter(language.Und)

	// 未指定 exts，表示 *.go。
	opt := &Options{}
	a.NotError(opt.sanitize())
	b := opt.newBuilder(p, l)
	a.NotNil(b)
	a.False(b.isIgnore("./builder.go"))
	a.True(b.isIgnore("./go.mod"))

	// exts = "*"
	opt = &Options{Exts: []string{"*"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder(p, l)
	a.NotNil(b)
	a.False(b.isIgnore("builder.go")).
		False(b.isIgnore("not-exists.file"))

	opt = &Options{Exts: []string{"*"}, Excludes: []string{"builder.go", "*_test.go"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder(p, l)
	a.NotNil(b)
	a.True(b.isIgnore("builder.go")).
		False(b.isIgnore("not-exists.file")).
		False(b.isIgnore("log.go")).
		True(b.isIgnore("log_test.go"))
}
