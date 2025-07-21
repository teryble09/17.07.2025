package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Port             string   `yaml:"port"`
	MaxCurrentTasks  int      `yaml:"max_current_tasks"`
	MaxURLsInTask    int      `yaml:"max_urls_in_task"`
	AllowedMIMETypes []string `yaml:"allowed_mime_types"`

	//In milliseconds
	RetryWaitTime  int `yaml:"retry_wait_time"`
	MaxRetryAmount int `yaml:"max_retry_amount"`
	//In milliseconds
	HttpClientTimeout int `yaml:"http_client_timeout"`
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
