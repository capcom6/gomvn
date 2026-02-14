package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type App struct {
	Name        string
	Server      *Server
	Repository  []string
	Permissions Permissions
	Debug       bool
	Database    Database `yaml:"database"`
	Storage     Storage  `yaml:"storage"`
}

func NewAppConfig(file string) (*App, error) {
	if file == "" {
		file = "config.yml"
	}

	var conf App
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	if ymlErr := yaml.Unmarshal(yamlFile, &conf); ymlErr != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", ymlErr)
	}

	if conf.Permissions.Index == nil {
		v := true
		conf.Permissions.Index = &v
	}

	return &conf, nil
}
