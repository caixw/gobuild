// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 一个简单的Go语言热编译工具
//
// USAGE:
//  gobuild [options] path
//  options:
//  -appname 给定一个程序名称，若不指定，则使用统一使用main，windows会加上.exe后缀名；
//  -help 显示帮助内容；
//  -version 显示版本信息；
package main
