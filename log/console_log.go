// SPDX-License-Identifier: MIT

package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

func TestLogs(t *testing.T) {
	a := assert.New(t, false)

	erro := new(bytes.Buffer)
	out := new(bytes.Buffer)
	logs := newConsoleLogs(true, erro, out)
	a.NotNil(logs)
	defer logs.Stop()

	logs.Logs <- &Log{Type: Error, Message: "error"}
	time.Sleep(300 * time.Microsecond)
	a.NotEmpty(erro.String())
	a.Empty(out.String())

	erro.Reset()
	out.Reset()
	logs.Logs <- &Log{Type: Ignore, Message: "message"}
	time.Sleep(300 * time.Microsecond)
	a.Empty(erro.String())
	a.NotEmpty(out.String())

	// ignore=false
	erro.Reset()
	out.Reset()
	logs = newConsoleLogs(false, erro, out)
	a.NotNil(logs)
	logs.Logs <- &Log{Type: Ignore, Message: "message"}
	time.Sleep(300 * time.Microsecond)
	a.Empty(out.String())
}
