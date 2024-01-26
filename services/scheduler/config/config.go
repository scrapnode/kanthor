package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf.Scheduler, conf.Validate()
}

type Wrapper struct {
	Scheduler Config `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Scheduler.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Request SchedulerRequest `json:"request" yaml:"request" mapstructure:"request"`
}

func (conf *Config) Validate() error {
	if err := conf.Request.Validate(); err != nil {
		return err
	}

	return nil
}

type SchedulerRequest struct {
	Schedule SchedulerRequestSchedule `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
}

func (conf *SchedulerRequest) Validate() error {
	if err := conf.Schedule.Validate(); err != nil {
		return err
	}
	return nil
}

type SchedulerRequestSchedule struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *SchedulerRequestSchedule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("SCHEDULER.CONFIG.REQUEST.SCHEDULE.TIMEOUT", conf.Timeout, 1000),
	)
}
