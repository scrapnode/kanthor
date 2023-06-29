package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

func NewShow(provider configuration.Provider, conf *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "show",
		ValidArgs: []string{"config"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "config" {
				return showConfig(provider, conf)
			}

			return nil
		},
	}
	return command
}

func showConfig(provider configuration.Provider, conf *config.Config) error {
	bytes, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	title := "SOURCES (lower priority will be overridden by higher)"
	fmt.Println(strings.Repeat("-", len(title)+2))
	fmt.Println(title)
	t := table.NewWriter()
	t.AppendHeader(table.Row{"file", "found", "priority"})

	sources := provider.Sources()
	for priority, source := range sources {
		t.AppendRow([]interface{}{source.Source, source.Found, priority})
	}
	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}
