package gitlabels

import (
	yaml "gopkg.in/yaml.v2"
)

// Config gitlabels config
type Config struct {
	Owner        string                 `yaml:"owner"`
	ORG          string                 `yaml:"org"`
	Labels       map[string]LabelConfig `yaml:"labels"`
	ProjectRegex string                 `yaml:"project-regex"`
	RenameLabels map[string]string      `yaml:"rename"`
	RemoveLabels []string               `yaml:"remove"`
}

// LabelConfig config for label
type LabelConfig struct {
	Color       string `yaml:"color"`
	Description string `yaml:"description"`
}

func (c Config) getUser() string {
	if c.ORG != "" {
		return c.ORG
	}
	return c.Owner
}

// ParseConfig parse config from []byte
func ParseConfig(content []byte) (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
