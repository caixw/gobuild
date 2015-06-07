// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/issue9/assert"
)

func TestBuilder_isIgnore(t *testing.T) {
	a := assert.New(t)

	fn := func(isIgnore bool, exts []string, path string) {
		b := newBuilder("*.go", "o.out", exts, []string{})
		a.NotNil(b)

		if isIgnore {
			a.True(b.isIgnore(path))
		} else {
			a.False(b.isIgnore(path))
		}

	}

	fn(true, []string{}, "abc.go")
	fn(true, []string{}, "abc.go/go")
	fn(true, []string{"."}, "abc.go/go")
	fn(true, []string{""}, "abc.go/go")
	fn(true, []string{"go"}, "abc.go/abcgo")
	fn(true, []string{"Go"}, "go/abc/abcgo")
	fn(true, []string{"go", "php"}, "abcgo/abc.html")
	fn(true, []string{"go", "php"}, "abc.go/abc.html")
	fn(true, []string{".Go"}, "abc.go/abcgo")
	fn(true, []string{".go"}, "go/abc/abcgo")
	fn(true, []string{".go", "php"}, "abcgo/abc.html")
	fn(true, []string{".go", "php"}, "abc.go/abc.html")
	fn(true, []string{".Go"}, "abc/.go")

	fn(false, []string{".go"}, "abc/.go")
	fn(false, []string{"go"}, "abc/.go")
	fn(false, []string{"go", ".php"}, "abc/.php")
}
