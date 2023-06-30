package dispatcher

import "context"

func (service *dispatcher) Connect(ctx context.Context) error {
	if err := service.repos.Connect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("connected")
	return nil
}

func (service *dispatcher) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.repos.Disconnect(ctx); err != nil {
		return err
	}

	if err := service.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
