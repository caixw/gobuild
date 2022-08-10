// SPDX-License-Identifier: MIT

package watch

import (
	"runtime"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestOptions_newBuilder(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{Dirs: []string{"./"}}
	a.NotError(opt.sanitize())

	b, err := opt.newBuilder(nil)
	a.NotError(err).NotNil(b)
	a.Contains(b.env, runtime.Version())
}
