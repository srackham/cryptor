package config

import (
	"fmt"
	"os"

	"github.com/srackham/cryptor/internal/fsx"
	"gopkg.in/yaml.v3"
)

type Config struct {
	XratesURL string `yaml:"xrates-url"`
}

func LoadConfig(fileName string) (*Config, error) {
	if !fsx.FileExists(fileName) {
		return nil, fmt.Errorf("missing config file: %v", fileName)
	}
	confFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %s", fileName)
	}
	defer confFile.Close()

	var config Config
	yamlDecoder := yaml.NewDecoder(confFile)
	err = yamlDecoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %s", fileName)
	}

	return &config, nil
}
