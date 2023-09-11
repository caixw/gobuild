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

	out := new(bytes.Buffer)
	logs := NewConsoleLogger(true, out, nil, nil)
	a.NotNil(logs)

	logs.Output(System, Error, "error")
	time.Sleep(300 * time.Microsecond)
	a.NotEmpty(out.String())

	out.Reset()
	logs.Output(System, Ignore, "message")
	time.Sleep(300 * time.Microsecond)
	a.NotEmpty(out.String())

	// ignore=false
	out.Reset()
	logs = NewConsoleLogger(false, out, nil, nil)
	a.NotNil(logs)
	logs.Output(System, Ignore, "message")
	time.Sleep(300 * time.Microsecond)
	a.Empty(out.String())
}

func TestAsWriter(t *testing.T) {
	a := assert.New(t, false)
	out := &bytes.Buffer{}
	logs := NewConsoleLogger(true, out, nil, nil)
	a.NotNil(logs)

	w := asWriter(System, Info, logs)
	n, err := w.Write([]byte("abc"))
	a.NotError(err).Equal(n, 3)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc")

	w = asWriter(System, Error, logs)
	n, err = w.Write([]byte("defg"))
	a.NotError(err).Equal(n, 4)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc").Contains(out.String(), "defg")
}
