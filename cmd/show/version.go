package show

import (
	"fmt"
	"github.com/scrapnode/kanthor/config"
)

func Version(conf *config.Config, verbose bool) error {
	fmt.Println(conf.Version)
	return nil
}
