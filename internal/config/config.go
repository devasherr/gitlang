package config

import (
	"gopkg.in/yaml.v3"
)

type Conventions struct {
	Naming []string `yaml:"naming"`
	Casing []string `yaml:"casing"`
}

type PreCommit struct {
	Enabled             bool     `yaml:"enabled"`
	MaxFileSizeKb       float64  `yaml:"max_file_size_kb"`
	ForbiddenExtensions []string `yaml:"forbidden_extensions"`
	NamingConventions   struct {
		Folder Conventions `yaml:"folder"`
		File   Conventions `yaml:"file"`
	} `yaml:"naming_conventions"`
}

type CommitMsg struct {
	Enabled          bool     `yaml:"enabled"`
	MinLength        int      `yaml:"min_length"`
	NoTrailingPeriod bool     `yaml:"no_trailing_period"`
	ForbiddenWords   []string `yaml:"forbidden_words"`
}

type Config struct {
	CommitMsg CommitMsg `yaml:"commit-msg"`
	PreCommit PreCommit `yaml:"pre-commit"`
}

func LoadConfig(body []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
