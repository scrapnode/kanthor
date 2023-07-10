package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

func NewShow(provider configuration.Provider, conf *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "show",
		ValidArgs: []string{"config", "version"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}
			validating, err := cmd.Flags().GetBool("validate")
			if err != nil {
				return err
			}

			if name == "config" {
				return showConfig(conf, provider.Sources(), verbose, validating)
			}
			if name == "version" {
				return showVersion(conf, verbose)
			}

			return nil
		},
	}

	command.Flags().BoolP("validate", "", false, "should validate the output we show you")
	return command
}

func showConfig(conf *config.Config, sources []configuration.Source, verbose, validating bool) error {
	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	if verbose {
		t := table.NewWriter()
		t.AppendHeader(table.Row{"looking", "found", "used"})
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
		return conf.Validate(services.ALL)
	}

	return nil
}

func showVersion(conf *config.Config, verbose bool) error {
	fmt.Println(conf.Version)
	return nil
}
