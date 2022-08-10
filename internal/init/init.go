// SPDX-License-Identifier: MIT

// Package init 提供初始化项目的功能
package init

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/issue9/source"
)

const binBaseDir = "cmd"

func Init(wd, name string) error {
	base := path.Base(name)

	wd = filepath.Join(wd, base)
	fi, err := os.Stat(wd)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		if err := os.Mkdir(wd, fs.ModePerm); err != nil {
			return err
		}
	case err != nil:
		return err
	default: // 不存在错误，说明存在文件夹或是同名文件，判断其是否为空。
		if !fi.IsDir() {
			return fs.ErrExist
		}

		dirs, err := os.ReadDir(wd)
		if err != nil {
			return err
		}
		if len(dirs) > 0 {
			return initOptions(wd, base)
		}
	}

	// 生成 go.mod
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = wd // go mod init 不会创建目录，而是在工作目录下直接创建。
	if err := cmd.Run(); err != nil {
		return err
	}

	// 创建 cmd/{base}/{base}.go 代码文件
	if err := initCmd(wd, base); err != nil {
		return err
	}

	return initOptions(wd, base)
}

func initCmd(wd, base string) error {
	cmd := filepath.Join(wd, binBaseDir, base)
	if err := os.MkdirAll(cmd, fs.ModePerm); err != nil {
		return err
	}
	return source.DumpGoSource(filepath.Join(cmd, "main.go"), []byte(code))
}

const code = `package main

func main() {
	// TODO
}`
