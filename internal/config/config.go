package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port             string   `yaml:"port"`
	MaxCurrentTasks  int      `yaml:"max_current_tasks"`
	AllowedMIMETypes []string `yaml:"allowed_mime_types"`
}

func MustLoad(configPath string) *Config {
	cfg := new(Config)
	yml, err := os.ReadFile(configPath)
	if err != nil {
		panic("Can not read config: " + err.Error())
	}
	err = yaml.Unmarshal(yml, cfg)
	if err != nil {
		panic("Can not unmarshall config" + err.Error())
	}
	return cfg
}
