package utils

import (
	"context"
	"errors"
	"time"
)

type Stoppable interface {
	Stop(ctx context.Context) error
}

func Stop(instances ...Stoppable) error {
	// wait a little to stop our service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errc := make(chan error, 1)
	go func() {
		var returning error
		for _, instance := range instances {
			if err := instance.Stop(ctx); err != nil {
				returning = errors.Join(returning, err)
			}
		}

		errc <- returning
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
