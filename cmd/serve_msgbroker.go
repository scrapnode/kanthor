package cmd

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/ioc"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func NewServeMsgBroker() *cobra.Command {
	command := &cobra.Command{
		Use:       "msgbroker [-d data] [-c count] {pub | sub}",
		ValidArgs: []string{"pub", "sub"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		Short:     "command to play with our msgbroker",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.New()
			if err != nil {
				return err
			}

			broker, err := ioc.InitializeMsgBroker(conf)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			if err := broker.Connect(ctx); err != nil {
				return err
			}
			defer func() {
				_ = broker.Disconnect(ctx)
			}()

			if args[0] == "pub" {
				count, err := cmd.Flags().GetInt("count")
				if err != nil {
					return err
				}

				data, err := cmd.Flags().GetString("data")
				if err != nil {
					return err
				}
				if data == "" {
					data = fmt.Sprintf("ts:%d", time.Now().UTC().Unix())
				}

				for i := 0; i < count; i++ {
					e := &msgbroker.Event{
						Tier:  "starter",
						AppId: "cli",
						Type:  "testing.go",
						Id:    utils.ID("msg"),
						Data:  []byte(data),
						Metadata: map[string]string{
							"count": strconv.Itoa(i),
						},
					}

					if err := broker.Pub(ctx, e); err != nil {
						return err
					}
				}

				cancel()
				return nil
			}

			if err := broker.Sub(ctx, sub); err != nil {
				return err
			}
			// Listen for the interrupt signal.
			<-ctx.Done()
			// make sure once we stop process, we cancel all the execution
			cancel()

			return nil
		},
	}

	command.Flags().StringP("data", "d", "", "--data=hello")
	command.Flags().IntP("count", "c", 1, "--count=1")
	return command
}

func sub(event *msgbroker.Event) error {
	log.Printf("id:%s data:%s", event.Id, event.String())
	return nil
}
