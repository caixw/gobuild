// SPDX-License-Identifier: MIT

package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestInit(t *testing.T) {
	a := assert.New(t, false)

	a.Run("存在同名文件", func(a *assert.Assertion) {
		a.ErrorIs(Init("testdir", "mod1"), fs.ErrExist)
	})

	a.Run("正常创建", func(a *assert.Assertion) {
		a.NotError(Init("testdir/mod2", "mod"))
		a.FileExists("testdir/mod2/mod/go.mod")
		a.FileExists(filepath.Join("testdir/mod2/mod/", ConfigFilename))
		a.FileExists(filepath.Join("testdir/mod2/mod/", binBaseDir, "mod/main.go"))
		os.RemoveAll("testdir/mod2/mod")
	})

	a.Run("仅创建.gobuild.yaml", func(a *assert.Assertion) {
		a.NotError(Init("testdir", "mod3"))
		a.FileNotExists("testdir/mod3/go.mod")
		a.FileExists(filepath.Join("testdir/mod3", ConfigFilename))
		a.FileNotExists(filepath.Join("testdir/mod3/", binBaseDir, "mod/main.go"))
		os.RemoveAll(filepath.Join("testdir/mod3", ConfigFilename))
	})
}
