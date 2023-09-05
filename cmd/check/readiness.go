package check

import (
	"fmt"
	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func NewReadiness() *cobra.Command {
	command := &cobra.Command{
		Use: "readiness",
		ValidArgs: []string{
			services.PORTAL_API,
			services.SDK_API,
			services.SCHEDULER,
			services.DISPATCHER,
			services.STORAGE,
		},
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client := background.NewClient(healthcheck.DefaultConfig(fmt.Sprintf("kanthor.%s", serviceName)))
			if err := client.Readiness(); err != nil {
				return err
			}
			fmt.Println("ready")
			return nil
		},
	}
	return command
}