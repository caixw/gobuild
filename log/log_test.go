// SPDX-License-Identifier: MIT

package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

func TestAsWriter(t *testing.T) {
	a := assert.New(t, false)
	erro := new(bytes.Buffer)
	out := new(bytes.Buffer)
	logs := newConsoleLogs(true, erro, out)
	a.NotNil(logs)
	defer logs.Stop()

	w := AsWriter(Info, logs.Logs)
	n, err := w.Write([]byte("abc"))
	a.NotError(err).Equal(n, 3)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc").NotContains(erro.String(), "abc")

	w = AsWriter(Error, logs.Logs)
	n, err = w.Write([]byte("defg"))
	a.NotError(err).Equal(n, 4)
	time.Sleep(300 * time.Microsecond)
	a.Contains(out.String(), "abc").Contains(erro.String(), "defg")
}
