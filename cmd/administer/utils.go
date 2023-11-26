package administer

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/spf13/cobra"
)

var format = "2006-01-02"
var formatlong = format + " 15:04:05"

func daterange(cmd *cobra.Command) (*time.Time, *time.Time, error) {
	f, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, nil, err
	}
	from, err := time.Parse(format, f)
	if err != nil {
		from, err = time.Parse(formatlong, f)
		if err != nil {
			return nil, nil, err
		}
	}

	t, err := cmd.Flags().GetString("to")
	if err != nil {
		return nil, nil, err
	}
	to, err := time.Parse(format, t)
	if err != nil {
		to, err = time.Parse(formatlong, f)
		if err != nil {
			return nil, nil, err
		}
	}

	return &from, &to, nil
}

func timeout(cmd *cobra.Command) (context.Context, context.CancelFunc, error) {
	ctx := cmd.Context()
	t, err := cmd.Flags().GetInt64("timeout")
	if err != nil {
		return ctx, func() {}, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(t))
	return ctx, cancel, nil
}

func isValidAppIdArg(cmd *cobra.Command, args []string) error {
	if !suid.Valid(args[0]) {
		return fmt.Errorf("app id is not valid (you entered %s)", args[0])
	}

	if suid.Ns(args[0]) != entities.IdNsApp {
		return fmt.Errorf("app id must be started with %s (you entered %s)", entities.IdNsApp, args[0])
	}

	return nil
}
