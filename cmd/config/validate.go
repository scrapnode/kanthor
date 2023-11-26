package config

import (
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	infrastructure "github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func NewValidate(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "validate",
		ValidArgs: append([]string{services.ALL}, services.SERVICES...),
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			infra, err := infrastructure.New(provider)
			if err != nil {
				return err
			}
			if err := infra.Validate(); err != nil {
				return err
			}

			confs, err := Services(provider, args[0])
			if err != nil {
				return err
			}

			var returning error
			for _, conf := range confs {
				if err := conf.Validate(); err != nil {
					returning = errors.Join(returning, err)
				}
			}

			defer func() {
				if verbose, err := cmd.Flags().GetBool("verbose"); err == nil && verbose {
					fmt.Println("--- infrastructure ---")
					fmt.Println(utils.StringifyIndent(infra))

					for name, conf := range confs {
						fmt.Println("--- " + name + " ---")
						fmt.Println(utils.StringifyIndent(conf))
					}
				}
			}()

			return returning
		},
	}
	return command
}
