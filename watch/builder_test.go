// SPDX-License-Identifier: MIT

package watch

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestOptions_newBuilder(t *testing.T) {
	a := assert.New(t, false)

	opt := &Options{Dirs: []string{"./"}}
	a.NotError(opt.sanitize())

	b := opt.newBuilder(nil)
	a.NotNil(b)
}
