// SPDX-License-Identifier: MIT

package watch

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
)

var (
	_ xml.Marshaler   = Flags{}
	_ xml.Unmarshaler = &Flags{}
)

func TestFlags_Marshal(t *testing.T) {
	a := assert.New(t, false)
	f := Flags{"k1": "v1"}
	data, err := xml.Marshal(f)
	a.NotError(err).Equal(string(data), "<Flags><k1>v1</k1></Flags>")

	f2 := Flags{}
	a.NotError(xml.Unmarshal(data, &f2))
	a.Equal(f, f2)
}

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{}
	a.Error(opt.sanitize()).
		NotNil(opt.Printer)

	opt.Dirs = []string{"./"}
	a.NotError(opt.sanitize()).
		Equal(opt.WatcherFrequency, MinWatcherFrequency)

	opt.Excludes = []string{}
	a.NotError(opt.sanitize())

	opt.Excludes = []string{"/abc/*"}
	a.NotError(opt.sanitize())

	opt.Excludes = []string{"/abc/*****/def", "abc/[]/def"}
	a.Error(opt.sanitize())
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

func TestOptions_sanitizeExts(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{}
	opt.sanitizeExts()
	a.Empty(opt.Exts).False(opt.anyExts)

	opt = &Options{Exts: []string{".go", " ", "java", " .php"}}
	opt.sanitizeExts()
	a.Equal(opt.Exts, []string{".go", ".java", ".php"}).False(opt.anyExts)

	opt = &Options{Exts: []string{".go", " ", "java", "*"}}
	opt.sanitizeExts()
	a.True(opt.anyExts)
}

func TestGetAppName(t *testing.T) {
	a := assert.New(t, false)
	goexe := os.Getenv("GOEXE")

	name, err := getAppName("", "./testdir")
	a.NotError(err).True(strings.HasSuffix(name, "testdir"+goexe), name)

	name, err = getAppName("a", "./testdir/a")
	a.NotError(err).True(strings.HasSuffix(name, "a"+goexe), name)

	name, err = getAppName("a.exe", "./testdir")
	a.NotError(err).True(strings.HasSuffix(name, "a.exe"), name)
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t, false)

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
	a := assert.New(t, false)

	a.Equal(splitArgs("x=5    y=6"), []string{"x", "5", "y", "6"})
	a.Equal(splitArgs("xxx=5 -yy=6 -bool"), []string{"xxx", "5", "-yy", "6", "-bool"})
	a.Equal(splitArgs("xxx=5 yy=6  bool="), []string{"xxx", "5", "yy", "6", "bool"})
	a.Equal(splitArgs(`aa=1 bb "xxx=5 yy=6 bool="`), []string{"aa", "1", "bb", "xxx=5 yy=6 bool="})
	a.Equal(splitArgs(`aa=1 bb "x"`), []string{"aa", "1", "bb", "x"})
	a.Equal(splitArgs(`aa=1 bb ""`), []string{"aa", "1", "bb"})
	a.Equal(splitArgs(`  ""`), []string{})
}
