package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Scheduler struct {
	Request SchedulerRequest `json:"request" yaml:"request" mapstructure:"request"`
}

func (conf *Scheduler) Validate() error {
	if err := conf.Request.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.request: %v", err)
	}

	return nil
}

type SchedulerRequest struct {
	Schedule SchedulerRequestSchedule `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
}

func (conf *SchedulerRequest) Validate() error {
	if err := conf.Schedule.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.request.schedule: %v", err)
	}
	return nil
}

type SchedulerRequestSchedule struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *SchedulerRequestSchedule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.scheduler.request.schedule.timeout", conf.Timeout, 1000),
	)
}
