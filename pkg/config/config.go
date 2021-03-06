package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DSN            string `yaml:"dsn"`
	Database       string `yaml:"database"`
	GameCollection string `yaml:"gameCollection"`
	NameCollection string `yaml:"nameCollection"`
}

func LoadConfig(path string) (*Config, error) {
	var config *Config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, err
	}
	return config, nil
}
