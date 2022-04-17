package ops

import (
	"errors"
)

// Spec defines the configuration relevant to the Ops server (used for profilers, tracing, etc).
type Spec struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

// Validate checks that the values defined in the Spec config are valid.
func (s *Spec) Validate() error {

	if s.Port <= 0 {
		return errors.New("Config OpsSpec Port must be higher than 0")
	}

	return nil
}
