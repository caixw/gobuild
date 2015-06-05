// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/issue9/term/colors"
)

var (
	succ = colors.New(colors.Stdout, colors.Green, colors.Default)
	info = colors.New(colors.Stdout, colors.Default, colors.Default)
	def  = colors.New(colors.Stdout, colors.Default, colors.Default)
	erro = colors.New(colors.Stdout, colors.Red, colors.Default)
	warn = colors.New(colors.Stdout, colors.Magenta, colors.Default)
)
