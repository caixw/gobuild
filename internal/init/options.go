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

func initOptions(wd, base string) error {
	dir := path.Join(binBaseDir, base)
	o := &watch.Options{
		MainFiles:        path.Join(dir, "*.go"),
		OutputName:       path.Join(dir, base),
		Exts:             []string{".go", ".yaml", ".xml", ".yml", ".json"}, // 配置文件修改也重启
		Recursive:        true,
		Dirs:             []string{"./"},
		WatcherFrequency: watch.MinWatcherFrequency,
	}
	data, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(wd, ConfigFilename), data, fs.ModePerm)
}
