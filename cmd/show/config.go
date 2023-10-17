package show

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/services"
	"gopkg.in/yaml.v3"
)

func Config(conf *config.Config, sources []configuration.Source, verbose, validating bool) error {
	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	if verbose {
		t := table.NewWriter()
		t.AppendHeader(table.Row{"looking", "path", "used"})
		for _, source := range sources {
			var check string
			if source.Used {
				check = "x"
			}
			t.AppendRow([]interface{}{source.Looking, source.Found, check})
		}
		t.SetOutputMirror(os.Stdout)
		t.Render()
	}

	if validating {
		return conf.Validate(services.SERVICE_ALL)
	}

	return nil
}
