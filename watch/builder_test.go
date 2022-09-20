// SPDX-License-Identifier: MIT

package watch

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestOptions_newBuilder(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{Dirs: []string{"./"}}
	a.NotError(opt.sanitize())
	b := opt.newBuilder()
	a.NotNil(b).False(b.anyExt)

	opt = &Options{Dirs: []string{"./"}, Exts: []string{".a", "*", ".b"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder()
	a.NotNil(b).Equal(b.exts, []string{"*"}).True(b.anyExt)

	opt = &Options{Dirs: []string{"./"}, Exts: []string{".a", ".b"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder()
	a.NotNil(b).Equal(b.exts, []string{".a", ".b"}).False(b.anyExt)
}

func TestBuilder_isIgnore(t *testing.T) {
	a := assert.New(t, false)

	// 未指定 exts，忽略所有。
	opt := &Options{Dirs: []string{"./"}}
	a.NotError(opt.sanitize())
	b := opt.newBuilder()
	a.NotNil(b)
	a.True(b.isIgnore("./builder.go"))
	a.True(b.isIgnore("./go.mod"))

	// AutoTidy 自动监视 go.mod
	opt = &Options{Dirs: []string{"./"}, AutoTidy: true}
	a.NotError(opt.sanitize())
	b = opt.newBuilder()
	a.NotNil(b)
	a.True(b.isIgnore("./builder.go"))
	a.False(b.isIgnore("./go.mod"))

	// exts = "*"
	opt = &Options{Dirs: []string{"./"}, Exts: []string{"*"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder()
	a.NotNil(b)
	a.False(b.isIgnore("builder.go")).
		False(b.isIgnore("not-exists.file"))

	opt = &Options{Dirs: []string{"./"}, Exts: []string{"*"}, Excludes: []string{"builder.go", "*_test.go"}}
	a.NotError(opt.sanitize())
	b = opt.newBuilder()
	a.NotNil(b)
	a.True(b.isIgnore("builder.go")).
		False(b.isIgnore("not-exists.file")).
		False(b.isIgnore("log.go")).
		True(b.isIgnore("log_test.go"))
}
