package app

// Spec defines the configuration relevant to the gRPC API.
type Spec struct {
	AppName         string          `yaml:"appName" envconfig:"APP_NAME"`
	Environment     string          `yaml:"environment"`
	Host            string          `yaml:"host"`
	Insecure        bool            `yaml:"insecure"`
	RestAPIProvider RestAPIProvider `yaml:"restAPIProvider"`
	Port            uint16          `yaml:"port"`
}

// Validate checks that the values defined in the config are valid.
func (s *Spec) Validate() error {
	err := s.RestAPIProvider.validate()
	if err != nil {
		return err
	}

	return nil
}
