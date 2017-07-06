// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
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

	a.Equal(recursivePaths(false, []string{"./testdir"}), []string{
		"./testdir",
	})

	a.Equal(recursivePaths(true, []string{"./testdir"}), []string{
		"./testdir",
		"testdir/testdir1",
		"testdir/testdir2",
		"testdir/testdir2/testdir3",
	})

	a.Equal(recursivePaths(true, []string{"./testdir/testdir1", "./testdir/testdir2"}), []string{
		"./testdir/testdir1",
		"./testdir/testdir2",
		"testdir/testdir2/testdir3",
	})

	a.Equal(recursivePaths(true, []string{"./testdir/testdir2"}), []string{
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
