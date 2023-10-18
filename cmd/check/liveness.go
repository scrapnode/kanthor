package check

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/healthcheck"
	"github.com/scrapnode/kanthor/pkg/healthcheck/background"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func NewLiveness() *cobra.Command {
	command := &cobra.Command{
		Use: "liveness",
		ValidArgs: []string{
			services.SERVICE_SCHEDULER,
			services.SERVICE_DISPATCHER,
			services.SERVICE_STORAGE,
			services.SERVICE_ATTEMPT_TRIGGER,
		},
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client := background.NewClient(healthcheck.DefaultConfig(serviceName))
			if err := client.Liveness(); err != nil {
				return err
			}
			fmt.Println("live")
			return nil
		},
	}
	return command
}
