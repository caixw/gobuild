// SPDX-License-Identifier: MIT

// Package local 宿主机的一些操作
package local

import (
	"bytes"
	"os/exec"
	"strings"
)

// GoVersion 返回本地 Go 的版本信息
func GoVersion() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.TrimPrefix(buf.String(), "go version ")), nil
}
