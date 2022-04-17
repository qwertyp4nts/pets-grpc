package config

import (
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config/app"
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config/ops"
)

const (
	defaultPathConfigFile = "./config/app/local.yaml"
)

// Config defines the configuration for the application.
type Config struct {
	AppSpec app.Spec `yaml:"spec" envconfig:"SPEC"`
	OpsSpec ops.Spec `yaml:"ops,omitempty" envconfig:"OPS"`
}

// Validate checks that the values defined in the config are valid.
func (s *Config) Validate() error {
	err := s.AppSpec.Validate()
	if err != nil {
		return err
	}

	err = s.OpsSpec.Validate()
	if err != nil {
		return err
	}

	return nil
}
