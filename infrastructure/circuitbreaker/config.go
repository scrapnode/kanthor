package circuitbreaker

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	HalfOpenMaxPassThroughRequests     int     `json:"half_open_max_pass_through_requests" yaml:"half_open_max_pass_through_requests" mapstructure:"half_open_max_pass_through_requests"`
	HalfOpenTriggerMinimumRequests     int     `json:"half_open_trigger_minimum_requests" yaml:"half_open_trigger_minimum_requests" mapstructure:"half_open_trigger_minimum_requests"`
	HalfOpenTriggerErrorThresholdRatio float64 `json:"half_open_trigger_error_threshold_ratio" yaml:"half_open_trigger_error_threshold_ratio" mapstructure:"half_open_trigger_error_threshold_ratio"`
	CloseTriggerMinimumRequests        int     `json:"close_trigger_minimum_requests" yaml:"close_trigger_minimum_requests" mapstructure:"close_trigger_minimum_requests"`
	CloseTriggerErrorThresholdRatio    float64 `json:"close_trigger_error_threshold_ratio" yaml:"close_trigger_error_threshold_ratio" mapstructure:"close_trigger_error_threshold_ratio"`
	CloseStateClearInterval            int     `json:"close_state_clear_interval" yaml:"close_state_clear_interval" mapstructure:"close_state_clear_interval"`
	OpenStateDuration                  int     `json:"open_state_duration" yaml:"open_state_duration" mapstructure:"open_state_duration"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("circuitbreaker.conf.half_open_max_pass_through_requests", conf.HalfOpenMaxPassThroughRequests, 1),
		validator.NumberGreaterThanOrEqual("circuitbreaker.conf.half_open_trigger_minimum_requests", conf.HalfOpenTriggerMinimumRequests, 1),
		validator.NumberInRange("circuitbreaker.conf.half_open_trigger_minimum_requests", conf.HalfOpenTriggerErrorThresholdRatio, float64(0), float64(1)),
		validator.NumberGreaterThanOrEqual("circuitbreaker.conf.close_trigger_minimum_requests", conf.HalfOpenTriggerMinimumRequests, 1),
		validator.NumberInRange("circuitbreaker.conf.close_trigger_error_threshold_ratio", conf.CloseTriggerErrorThresholdRatio, float64(0), float64(1)),
		validator.NumberGreaterThanOrEqual("circuitbreaker.conf.close_state_clear_interval", conf.CloseStateClearInterval, 0),
		validator.NumberGreaterThanOrEqual("circuitbreaker.conf.open_state_duration", conf.OpenStateDuration, 0),
	)
}
