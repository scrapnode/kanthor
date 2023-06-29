package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"strings"
)

var FileLookingDirs = []string{"./", "$HOME/.kanthor/", "$KANTHOR_CONFIG_DIR/"}
var FileName = "configs"
var FileExt = "yaml"

func NewFile(dirs []string) (Provider, error) {
	provider := &file{viper: viper.New()}
	provider.viper.SetConfigName(FileName) // name of config file (without extension)
	provider.viper.SetConfigType(FileExt)  // extension

	for _, dir := range dirs {
		source := Source{Source: path.Join(dir, fmt.Sprintf("%s.%s", FileName, FileExt)), Found: true}

		provider.viper.AddConfigPath(dir)
		if err := provider.viper.MergeInConfig(); err != nil {
			// ignore not found files, otherwise return error

			if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
				return nil, fmt.Errorf("config.viper.MergeInConfig(): %v", err)
			}
		}

		source.Found = true
		provider.sources = append(provider.sources, source)
	}

	provider.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	provider.viper.SetEnvPrefix("KANTHOR")
	provider.viper.AutomaticEnv()

	return provider, nil
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
