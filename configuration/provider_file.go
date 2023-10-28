package configuration

import (
	"fmt"
	"path"
	"strings"

	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"
	"github.com/spf13/viper"
)

var FileLookingDirs = []string{"$KANTHOR_HOME/", "$HOME/.kanthor/", "./"}
var FileName = "configs"
var FileExt = "yaml"

func NewFile(dirs []string) (Provider, error) {
	instance := viper.New()
	instance.SetConfigName(FileName) // name of config file (without extension)
	instance.SetConfigType(FileExt)  // extension

	sources := []Source{}
	for _, dir := range dirs {
		dir = strings.Trim(dir, " ")
		filename := FileName + "." + FileExt
		sources = append(sources, Source{Looking: path.Join(dir, filename), Found: path.Join(utils.AbsPathify(dir), filename)})
		instance.AddConfigPath(dir)
	}

	if err := instance.MergeInConfig(); err != nil {
		// ignore not found files, otherwise return error
		if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
			return nil, fmt.Errorf("config.viper.MergeInConfig(): %v", err)
		}
	}

	for index, source := range sources {
		source.Used = instance.ConfigFileUsed() != "" && instance.ConfigFileUsed() == source.Found
		sources[index] = source
	}

	instance.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	instance.SetEnvPrefix(project.Namespace())
	instance.AutomaticEnv()

	return &file{viper: instance, sources: sources}, nil
}

type file struct {
	viper   *viper.Viper
	sources []Source
}

func (provider *file) Unmarshal(dest interface{}) error {
	return provider.viper.Unmarshal(dest)
}

func (provider *file) Sources() []Source {
	return provider.sources
}
