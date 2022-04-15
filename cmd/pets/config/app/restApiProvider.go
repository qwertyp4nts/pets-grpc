package app

import (
	"errors"
	"time"
)

// RestAPIProvider defines the configuration specific to the connection to RestAPIProvider.
type RestAPIProvider struct {
	Host     string        `yaml:"host"`
	Insecure bool          `yaml:"insecure"`
	Timeout  time.Duration `yaml:"timeout"`
}

func (r RestAPIProvider) validate() error {
	if r.Host == "" {
		return errors.New("RestAPIProvider Host is required")
	}

	return nil
}
