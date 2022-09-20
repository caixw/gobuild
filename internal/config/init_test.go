// SPDX-License-Identifier: MIT

package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestInit(t *testing.T) {
	t.Run("存在同名文件", func(t *testing.T) {
		a := assert.New(t, false)
		a.ErrorIs(Init("testdir", "mod1"), fs.ErrExist)
	})

	t.Run("正常创建", func(t *testing.T) {
		a := assert.New(t, false)
		a.NotError(Init("testdir/mod2", "mod"))
		a.FileExists("testdir/mod2/mod/go.mod")
		a.FileExists(filepath.Join("testdir/mod2/mod/", Filename))
		a.FileExists(filepath.Join("testdir/mod2/mod/", binBaseDir, "mod/main.go"))
		os.RemoveAll("testdir/mod2/mod")
	})

	t.Run("仅创建.gobuild.yaml", func(t *testing.T) {
		a := assert.New(t, false)
		a.NotError(Init("testdir", "mod3"))
		a.FileNotExists("testdir/mod3/go.mod")
		a.FileExists(filepath.Join("testdir/mod3", Filename))
		a.FileNotExists(filepath.Join("testdir/mod3/", binBaseDir, "mod/main.go"))
		os.RemoveAll(filepath.Join("testdir/mod3", Filename))
	})
}
