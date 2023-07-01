package dataplane

import "context"

func (usecase *dataplane) Connect(ctx context.Context) error {
	if err := usecase.repos.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Connect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Connect(ctx); err != nil {
		return err
	}

	usecase.logger.Info("connected")
	return nil
}

func (usecase *dataplane) Disconnect(ctx context.Context) error {
	usecase.logger.Info("disconnected")

	if err := usecase.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.publisher.Disconnect(ctx); err != nil {
		return err
	}

	if err := usecase.cache.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
