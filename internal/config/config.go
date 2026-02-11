package config

import (
	"gopkg.in/yaml.v3"
)

type CommitMsg struct {
	Enabled          bool     `yaml:"enabled"`
	MinLength        int      `yaml:"min_length"`
	NoTrailingPeriod bool     `yaml:"no_trailing_period"`
	ForbiddenWords   []string `yaml:"forbidden_words"`
}

type Config struct {
	CommitMsg CommitMsg `yaml:"commit-msg"`
}

func LoadConfig(body []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
