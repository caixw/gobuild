// SPDX-License-Identifier: MIT

package config

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestGetRoot(t *testing.T) {
	a := assert.New(t, false)

	wd, err := getRoot("./testdir/mod4/mod44")
	a.NotError(err).
		Equal(filepath.Base(wd), "mod4")

	wd, err = getRoot("./testdir/mod4")
	a.NotError(err).
		Equal(filepath.Base(wd), "mod4")

	wd, err = getRoot("./") // 当前项目并未设置 gobuild.yaml 文件
	a.Error(err).Empty(wd)
}
