// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"runtime"
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

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	paths, err := recursivePaths(false, []string{"./testdir"})
	a.NotError(err).
		Equal(paths, []string{
			"./testdir",
		})

	paths, err = recursivePaths(true, []string{"./testdir"})
	a.NotError(err).Equal(paths, []string{
		"./testdir",
		"testdir/testdir1",
		"testdir/testdir2",
		"testdir/testdir2/testdir3",
	})

	paths, err = recursivePaths(true, []string{"./testdir/testdir1", "./testdir/testdir2"})
	a.NotError(err).
		Equal(paths, []string{
			"./testdir/testdir1",
			"./testdir/testdir2",
			"testdir/testdir2/testdir3",
		})

	paths, err = recursivePaths(true, []string{"./testdir/testdir2"})
	a.NotError(err).Equal(paths, []string{
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

	name, err := getAppName("", "./testdir")
	a.NotError(err)
	if runtime.GOOS != "windows" {
		a.True(strings.HasSuffix(name, "testdir"), name)
	} else {
		a.True(strings.HasSuffix(name, "testdir.exe"), name)
	}

	name, err = getAppName("a", "./testdir")
	a.NotError(err)
	if runtime.GOOS != "windows" {
		a.True(strings.HasSuffix(name, "testdir/a"), name)
	} else {
		a.True(strings.HasSuffix(name, "testdir/a.exe"), name)
	}

	name, err = getAppName("a.exe", "./testdir")
	a.NotError(err).
		True(strings.HasSuffix(name, "testdir/a.exe"), name)
}
