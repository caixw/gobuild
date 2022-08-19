// SPDX-License-Identifier: MIT

package watch

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

func TestLogger(t *testing.T) {
	a := assert.New(t, false)

	erro := new(bytes.Buffer)
	out := new(bytes.Buffer)
	logs := NewConsoleLogger(true, erro, out)
	a.NotNil(logs)

	logs.Output(LogTypeError, "error")
	time.Sleep(300 * time.Microsecond)
	a.NotEmpty(erro.String())
	a.Empty(out.String())

	erro.Reset()
	out.Reset()
	logs.Output(LogTypeIgnore, "message")
	time.Sleep(300 * time.Microsecond)
	a.Empty(erro.String())
	a.NotEmpty(out.String())

	// ignore=false
	erro.Reset()
	out.Reset()
	logs = NewConsoleLogger(false, erro, out)
	a.NotNil(logs)
	logs.Output(LogTypeIgnore, "message")
	time.Sleep(300 * time.Microsecond)
	a.Empty(out.String())
}

func TestAsWriter(t *testing.T) {
	a := assert.New(t, false)
	erro := &bytes.Buffer{}
	out := &bytes.Buffer{}
	logs := NewConsoleLogger(true, erro, out)
	a.NotNil(logs)

	w := asWriter(LogTypeInfo, logs)
	n, err := w.Write([]byte("abc"))
	a.NotError(err).Equal(n, 3)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc").NotContains(erro.String(), "abc")

	w = asWriter(LogTypeError, logs)
	n, err = w.Write([]byte("defg"))
	a.NotError(err).Equal(n, 4)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc").Contains(erro.String(), "defg")
}
