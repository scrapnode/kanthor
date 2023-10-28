package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Scheduler Scheduler `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
}

func (conf *Config) Validate() error {
	if err := conf.Scheduler.Validate(); err != nil {
		return err
	}
	return nil
}

type Scheduler struct {
	Request SchedulerRequest `json:"request" yaml:"request" mapstructure:"request"`
}

func (conf *Scheduler) Validate() error {
	if err := conf.Request.Validate(); err != nil {
		return fmt.Errorf("scheduler.request: %v", err)
	}

	return nil
}

type SchedulerRequest struct {
	Schedule SchedulerRequestSchedule `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
}

func (conf *SchedulerRequest) Validate() error {
	if err := conf.Schedule.Validate(); err != nil {
		return fmt.Errorf("scheduler.request.schedule: %v", err)
	}
	return nil
}

type SchedulerRequestSchedule struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *SchedulerRequestSchedule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("scheduler.request.schedule.timeout", conf.Timeout, 1000),
	)
}
