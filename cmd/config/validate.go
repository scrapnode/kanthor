package config

import (
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
		ValidArgs: services.SERVICES,
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			infra, err := infrastructure.New(provider)
			if err != nil {
				return err
			}
			if err := infra.Validate(); err != nil {
				return err
			}

			serviceName := args[0]
			service, err := Service(provider, serviceName)
			if err != nil {
				return err
			}

			if err := service.Validate(); err != nil {
				return err
			}

			fmt.Println("--- infrastructure ---")
			fmt.Println(utils.StringifyIndent(infra, ""))

			fmt.Println("--- " + serviceName + " ---")
			fmt.Println(utils.StringifyIndent(service, ""))
			return nil
		},
	}
	return command
}
