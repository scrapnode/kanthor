package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Endeavour struct {
	Planner  EndeavourPlanner  `json:"planner" yaml:"planner" mapstructure:"planner"`
	Executor EndeavourExecutor `json:"executor" yaml:"executor" mapstructure:"executor"`
}

func (conf *Endeavour) Validate() error {
	if err := conf.Planner.Validate(); err != nil {
		return fmt.Errorf("attempt.endeavour.planner: %v", err)
	}
	if err := conf.Executor.Validate(); err != nil {
		return fmt.Errorf("attempt.endeavour.executor: %v", err)
	}
	return nil
}

type EndeavourPlanner struct {
	Schedule string `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Timeout  int64  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size     int    `json:"size" yaml:"size" mapstructure:"size"`

	ScanStart int64 `json:"scan_start" yaml:"scan_start" mapstructure:"scan_start"`
	ScanEnd   int64 `json:"scan_end" yaml:"scan_end" mapstructure:"scan_end"`
}

func (conf *EndeavourPlanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("attempt.endeavour.planner.schedule", conf.Schedule),
		validator.NumberGreaterThanOrEqual("attempt.endeavour.planner.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("attempt.endeavour.planner.size", conf.Size, 0),
		validator.NumberLessThan("attempt.endeavour.planner.scan_end", conf.ScanEnd, 0),
		validator.NumberLessThan("attempt.endeavour.planner.scan_start", conf.ScanStart, conf.ScanEnd),
	)
}

type EndeavourExecutor struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size    int   `json:"size" yaml:"size" mapstructure:"size"`
}

func (conf *EndeavourExecutor) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("attempt.endeavour.executor.timeout", int(conf.Timeout), 1000),
		validator.NumberGreaterThan("attempt.endeavour.executor.size", conf.Size, 0),
	)
}
