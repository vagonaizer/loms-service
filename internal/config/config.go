package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPC struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	LOMS struct {
		Address string `yaml:"address"`
	} `yaml:"loms"`
}

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig(path string) (*Config, error) {
	config := NewConfig()

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}
