package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	flag "github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

// Load configuration from file, env and flags and return compiled and validated config.
func Load() (*Config, error) {
	f := flag.CommandLine

	flags(f)

	flag.Parse()

	return create(f)
}

// Create validated configuration from file, env, and flags.
func create(f *flag.FlagSet) (*Config, error) {
	var config Config

	configFile, _ := f.GetString("config")

	if configFile == "" {
		configFile = defaultPathConfigFile
	}

	cfg, err := os.Open(configFile)
	if err != nil {
		log.Println("error while trying to open configuration file")

		return nil, err
	}
	defer cfg.Close()

	decoder := yaml.NewDecoder(cfg)

	err = decoder.Decode(&config)
	if err != nil {
		log.Println("error while trying to open configuration file")

		return nil, err
	}

	err = envconfig.Process("", &config)
	if err != nil {
		log.Println("error while trying to bind environment variables to configuration file")

		return nil, err
	}

	// Validate app config
	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Flags enables shorthand flags for specific config values.
func flags(f *flag.FlagSet) {
	f.StringP("config", "c", "", "The configuration file to use to configure this application")
}
