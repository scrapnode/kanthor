package circuitbreaker

import "github.com/go-playground/validator/v10"

type Config struct {
	HalfOpenMaxPassThroughRequests     int     `json:"half_open_max_pass_through_requests" yaml:"half_open_max_pass_through_requests" mapstructure:"half_open_max_pass_through_requests" validate:"required,number,gte=1"`
	HalfOpenTriggerMinimumRequests     int     `json:"half_open_trigger_minimum_requests" yaml:"half_open_trigger_minimum_requests" mapstructure:"half_open_trigger_minimum_requests" validate:"required,number,gte=1"`
	HalfOpenTriggerErrorThresholdRatio float64 `json:"half_open_trigger_error_threshold_ratio" yaml:"half_open_trigger_error_threshold_ratio" mapstructure:"half_open_trigger_error_threshold_ratio" validate:"required,number,gt=0,lt=1"`
	CloseStateClearInterval            int     `json:"close_state_clear_interval" yaml:"close_state_clear_interval" mapstructure:"close_state_clear_interval" validate:"required,number,gte=0"`
	OpenStateDuration                  int     `json:"open_state_duration" yaml:"open_state_duration" mapstructure:"open_state_duration" validate:"required,number,gte=0"`
}

func (conf Config) Validate() error {
	return validator.New().Struct(conf)
}
