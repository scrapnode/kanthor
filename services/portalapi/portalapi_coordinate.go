package portalapi

func (service *portalapi) coordinate() error {
	return service.infra.Coordinator.Receive(func(cmd string, data []byte) error {
		service.logger.Debugw("coordinating", "cmd", cmd, "data", data)

		//ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		//defer cancel()

		return nil
	})
}
