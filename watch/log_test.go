// SPDX-License-Identifier: MIT

package watch

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

type logs struct {
	Logs chan *Log
	erro *bytes.Buffer
	out  *bytes.Buffer
}

func (l *logs) close() { close(l.Logs) }

func newLogs() *logs {
	l := &logs{
		Logs: make(chan *Log, 10),
		erro: &bytes.Buffer{},
		out:  &bytes.Buffer{},
	}
	go func() {
		for log := range l.Logs {
			if log == nil {
				return
			}
			switch log.Type {
			case LogTypeInfo:
				l.out.WriteString(log.Message)
			case LogTypeError:
				l.erro.WriteString(log.Message)
			}

		}
	}()

	return l
}

func TestAsWriter(t *testing.T) {
	a := assert.New(t, false)
	logs := newLogs()
	a.NotNil(logs)
	defer logs.close()

	w := asWriter(LogTypeInfo, logs.Logs)
	n, err := w.Write([]byte("abc"))
	a.NotError(err).Equal(n, 3)
	time.Sleep(300 * time.Microsecond)
	a.Contains(logs.out.String(), "abc").NotContains(logs.erro.String(), "abc")

	w = asWriter(LogTypeError, logs.Logs)
	n, err = w.Write([]byte("defg"))
	a.NotError(err).Equal(n, 4)
	time.Sleep(300 * time.Microsecond)
	a.Contains(logs.out.String(), "abc").Contains(logs.erro.String(), "defg")
}
