// SPDX-License-Identifier: MIT

package init

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/caixw/gobuild/watch"
)

const ConfigFilename = ".gobuild.yaml"

const fileheader = `# 此文件由 gobuild<https://github.com/caixw/gobuild> 生成和使用
`

func initOptions(wd, base string) error {
	dir := path.Join(binBaseDir, base)
	o := &watch.Options{
		MainFiles:        path.Join("./", dir),
		OutputName:       path.Join(dir, base),
		Excludes:         []string{ConfigFilename},
		Exts:             []string{".go", ".yaml", ".xml", ".yml", ".json"}, // 配置文件修改也重启
		Recursive:        true,
		Dirs:             []string{"./"},
		AutoTidy:         true,
		WatcherFrequency: watch.MinWatcherFrequency,
	}
	data, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	d := make([]byte, 0, len(fileheader)+len(data))
	d = append(d, []byte(fileheader)...)
	d = append(d, data...)
	return os.WriteFile(filepath.Join(wd, ConfigFilename), d, fs.ModePerm)
}
