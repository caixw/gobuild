// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert"
)

func TestGetExts(t *testing.T) {
	a := assert.New(t)

	a.Equal(getExts(""), []string{})
	a.Equal(getExts(",, ,"), []string{})
	a.Equal(getExts(",.go, ,.php,"), []string{".go", ".php"})
	a.Equal(getExts(",go,.php,"), []string{".go", ".php"})
	a.Equal(getExts(",go , .php,"), []string{".go", ".php"})
}

func pathsEqual(a *assert.Assertion, paths1, paths2 []string) {
	a.Equal(len(paths1), len(paths2))

LOOP:
	for _, p1 := range paths1 {
		p1 = filepath.ToSlash(p1)
		var p2 string
		for _, p2 = range paths2 {
			if p1 == p2 {
				continue LOOP
			}
		}

		a.True(false, "路径 %s %s 不相等", p1, p2)
	}
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	paths, err := recursivePaths(false, []string{"./testdir"})
	a.NotError(err)
	pathsEqual(a, paths, []string{
		"./testdir",
	})

	paths, err = recursivePaths(true, []string{"./testdir"})
	a.NotError(err)
	pathsEqual(a, paths, []string{
		"./testdir",
		"testdir/testdir1",
		"testdir/testdir2",
		"testdir/testdir2/testdir3",
	})

	paths, err = recursivePaths(true, []string{"./testdir/testdir1", "./testdir/testdir2"})
	a.NotError(err)
	pathsEqual(a, paths, []string{
		"./testdir/testdir1",
		"./testdir/testdir2",
		"testdir/testdir2/testdir3",
	})

	paths, err = recursivePaths(true, []string{"./testdir/testdir2"})
	a.NotError(err)
	pathsEqual(a, paths, []string{
		"./testdir/testdir2",
		"testdir/testdir2/testdir3",
	})
}

func TestSplitArgs(t *testing.T) {
	a := assert.New(t)

	a.Equal(splitArgs("x=5 y=6"), []string{"x", "5", "y", "6"})
	a.Equal(splitArgs("xxx=5 -yy=6 -bool"), []string{"xxx", "5", "-yy", "6", "-bool"})
	a.Equal(splitArgs("xxx=5 yy=6 bool="), []string{"xxx", "5", "yy", "6", "bool"})
}

func TestGetAppName(t *testing.T) {
	a := assert.New(t)
	goexe := os.Getenv("GOEXE")

	name, err := getAppName("", "./testdir")
	a.NotError(err).True(strings.HasSuffix(name, "testdir"+goexe), name)

	name, err = getAppName("a", "./testdir/a")
	a.NotError(err).True(strings.HasSuffix(name, "a"+goexe), name)

	name, err = getAppName("a.exe", "./testdir")
	a.NotError(err).True(strings.HasSuffix(name, "a.exe"), name)
}
