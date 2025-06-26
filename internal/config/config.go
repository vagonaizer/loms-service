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
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
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
