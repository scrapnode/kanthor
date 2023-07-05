package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

func NewShow(provider configuration.Provider, conf *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "show",
		ValidArgs: []string{"config"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")

			if args[0] == "config" {
				return showConfig(provider, conf, verbose)
			}

			return nil
		},
	}

	return command
}

func showConfig(provider configuration.Provider, conf *config.Config, verbose bool) error {
	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	if verbose {
		t := table.NewWriter()
		t.AppendHeader(table.Row{"origin", "found", "used"})
		sources := provider.Sources()
		for _, source := range sources {
			var check string
			if source.Used {
				check = "x"
			}
			t.AppendRow([]interface{}{source.Origin, source.Found, check})
		}
		t.SetOutputMirror(os.Stdout)
		t.Render()
	}

	return nil
}
