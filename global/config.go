package global

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type ConfigType struct {
	Editor         *[]string `yaml:"editor"`
	EditorPassword *string   `yaml:"editorPassword"`
	RootDir        *string   `yaml:"rootDir"`
}

var Config *ConfigType

func defaultConfig() *ConfigType {
	rootDir := "."
	return &ConfigType{
		Editor:  &[]string{"nano", "-R"},
		RootDir: &rootDir,
	}
}

func parseConfig(file string) *ConfigType {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return &ConfigType{}
	}
	data, err := os.ReadFile(file)
	if err != nil {
		log.Panicln(err)
	}
	ret := &ConfigType{}
	err = yaml.Unmarshal(data, ret)
	if err != nil {
		log.Panicln(err)
	}
	return ret
}

func mergeConfig(configs ...*ConfigType) *ConfigType {
	ret := &ConfigType{}
	for _, conf := range configs {
		if ret.Editor == nil {
			ret.Editor = conf.Editor
		}
		if ret.EditorPassword == nil {
			ret.EditorPassword = conf.EditorPassword
		}
		if ret.RootDir == nil {
			ret.RootDir = conf.RootDir
		}
	}
	return ret
}

func init() {
	userRC, _ := os.UserHomeDir()
	Config = mergeConfig(parseConfig(filepath.Join(userRC, ".blogrc")), defaultConfig())
}
