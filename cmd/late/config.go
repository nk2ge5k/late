package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres struct {
		URL string `yaml:"url"`
	} `yaml:"postgres"`
	GRPC struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc"`
}

func LoadConfig(filename string) (cfg Config, err error) {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return cfg, fmt.Errorf("could not read config file: %w", err)
	}

	if err := yaml.Unmarshal(bb, &cfg); err != nil {
		return cfg, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return
}
