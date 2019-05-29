// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gobuild

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestLogs(t *testing.T) {
	a := assert.New(t)
	erro := new(bytes.Buffer)
	out := new(bytes.Buffer)

	logs := newConsoleLogs(true, erro, out)
	a.NotNil(logs)

	logs.Logs <- &Log{Type: LogTypeError, Message: "error"}
	time.Sleep(300 * time.Microsecond)
	a.NotEmpty(erro.String())
	a.Empty(out.String())

	erro.Reset()
	out.Reset()
	logs.Logs <- &Log{Type: LogTypeIgnore, Message: "message"}
	time.Sleep(300 * time.Microsecond)
	a.Empty(erro.String())
	a.NotEmpty(out.String())

	// ignore=false
	erro.Reset()
	out.Reset()
	logs = newConsoleLogs(false, erro, out)
	a.NotNil(logs)
	logs.Logs <- &Log{Type: LogTypeIgnore, Message: "message"}
	time.Sleep(300 * time.Microsecond)
	a.Empty(out.String())
}
