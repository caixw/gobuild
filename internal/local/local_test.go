// SPDX-License-Identifier: MIT

package local

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestGoVersion(t *testing.T) {
	a := assert.New(t, false)
	v, err := GoVersion()
	a.NotError(err).NotEmpty(v)
}
