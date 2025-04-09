// SPDX-FileCopyrightText: 2015-2025 caixw
//
// SPDX-License-Identifier: MIT

// Package config 提供操作配置文件的能力
package config

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/caixw/gobuild/watch"
)

// 默认的 main 方法的父目录
const binBaseDir = "cmd"

const fileHeader = `# 此文件由 gobuild<https://github.com/caixw/gobuild> 生成和使用

`

func initOptions(wd, base, configFilename string) error {
	o := &watch.Options{
		MainFiles:        "./" + path.Join(binBaseDir, base), // path.Join 会将 ./ 去除，所以 ./ 不能在 path.Join 中。
		Excludes:         []string{configFilename},
		Exts:             []string{".go", ".yaml", ".xml", ".yml", ".json"}, // 配置文件修改也重启
		WatcherFrequency: watch.MinWatcherFrequency,
	}
	data, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	d := make([]byte, 0, len(fileHeader)+len(data))
	d = append(d, []byte(fileHeader)...)
	d = append(d, data...)
	return os.WriteFile(filepath.Join(wd, configFilename), d, fs.ModePerm)
}
