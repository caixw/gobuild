// SPDX-License-Identifier: MIT

// Package config 提供操作配置文件的能力
package config

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild/watch"
)

// 默认的 main 方法的父目录
const binBaseDir = "cmd"

// Filename 配置文件的文件名
const Filename = ".gobuild.yaml"

const fileHeader = `# 此文件由 gobuild<https://github.com/caixw/gobuild> 生成和使用

`

func initOptions(wd, base string) error {
	dir := path.Join(binBaseDir, base)
	o := &watch.Options{
		MainFiles:        path.Join("./", dir),
		Excludes:         []string{Filename},
		Exts:             []string{".go", ".yaml", ".xml", ".yml", ".json"}, // 配置文件修改也重启
		AutoTidy:         true,
		WatcherFrequency: watch.MinWatcherFrequency,
	}
	data, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	d := make([]byte, 0, len(fileHeader)+len(data))
	d = append(d, []byte(fileHeader)...)
	d = append(d, data...)
	return os.WriteFile(filepath.Join(wd, Filename), d, fs.ModePerm)
}
