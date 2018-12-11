// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

// Log 日志类型
type Log struct {
	Type    int8
	Message string
}

// 日志类型
const (
	LogTypeSuccess int8 = iota + 1
	LogTypeInfo
	LogTypeWarn
	LogTypeError
	LogTypeIgnore
)
