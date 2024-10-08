// SPDX-FileCopyrightText: 2015-2024 caixw
//
// SPDX-License-Identifier: MIT

package watch

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{}
	a.NotError(opt.sanitize())

	opt.MainFiles = "./"
	a.NotError(opt.sanitize()).
		Equal(opt.WatcherFrequency, MinWatcherFrequency).
		When(runtime.GOOS == "windows", func(a *assert.Assertion) {
			a.True(strings.HasSuffix(opt.appName, "watch.exe"))
		}).
		When(runtime.GOOS != "windows", func(a *assert.Assertion) {
			a.True(strings.HasSuffix(opt.appName, "watch"))
		})

	opt.Excludes = []string{}
	a.NotError(opt.sanitize())

	opt.Excludes = []string{"/abc/*"}
	a.NotError(opt.sanitize())

	opt.Excludes = []string{"/abc/*****/def", "abc/[]/def"}
	a.Error(opt.sanitize())
}

func pathsEqual(a *assert.Assertion, paths1, paths2 []string) {
	a.TB().Helper()
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

func TestOptions_sanitizeExts(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{}
	opt.sanitizeExts()
	a.Equal(opt.Exts, []string{".go"}).False(opt.anyExts)

	opt = &Options{Exts: []string{".go", " ", "java", " .php"}}
	opt.sanitizeExts()
	a.Equal(opt.Exts, []string{".go", ".java", ".php"}).False(opt.anyExts)

	opt = &Options{Exts: []string{".go", " ", "java", "*"}}
	opt.sanitizeExts()
	a.True(opt.anyExts)
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t, false)

	abs := func(s string) string {
		ss, err := filepath.Abs(s)
		a.NotError(err).NotEmpty(ss)
		return filepath.ToSlash(ss)
	}

	paths, err := recursivePaths("./testdir/testdir1")
	a.NotError(err)
	pathsEqual(a, paths, []string{
		abs("./testdir/testdir1"),
	})

	paths, err = recursivePaths("./testdir/testdir2")
	a.NotError(err)
	pathsEqual(a, paths, []string{
		abs("./testdir"),
		// "testdir/testdir1",  // 有 go.mod
		abs("testdir/testdir2"),
		abs("testdir/testdir2/testdir3"),
	})
}

func TestSplitArgs(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(splitArgs("x=5    y=6"), []string{"x", "5", "y", "6"})
	a.Equal(splitArgs("xxx=5 -yy=6 -bool"), []string{"xxx", "5", "-yy", "6", "-bool"})
	a.Equal(splitArgs("xxx=5 yy=6  bool="), []string{"xxx", "5", "yy", "6", "bool"})
	a.Equal(splitArgs(`aa=1 bb "xxx=5 yy=6 bool="`), []string{"aa", "1", "bb", "xxx=5 yy=6 bool="})
	a.Equal(splitArgs(`aa=1 bb "x"`), []string{"aa", "1", "bb", "x"})
	a.Equal(splitArgs(`aa=1 bb ""`), []string{"aa", "1", "bb"})
	a.Equal(splitArgs(`  ""`), []string{})
}

func TestGetWD(t *testing.T) {
	a := assert.New(t, false)

	wd, err := getWD("./")
	a.NotError(err).NotEmpty(wd)

	wd, err = getWD("./testdir")
	a.NotError(err).True(strings.HasSuffix(wd, "testdir"))

	wd, err = getWD("./testdir/go.mod")
	a.NotError(err).True(strings.HasSuffix(wd, "testdir"))

	wd, err = getWD("./testdir/*.go")
	a.NotError(err).True(strings.HasSuffix(wd, "testdir"))
}
