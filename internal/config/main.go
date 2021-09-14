package config

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config interface {
	GetAuth() Auth
	Logger
}

type config struct {
	Auth     Auth   `yaml:"auth"`
	LogLevel string `yaml:"log"`

	Logger
}

type Auth struct {
	Proxy     string `yaml:"proxy"`
	CSRFToken string `yaml:"csrf_token"`
	Cookie    string `yaml:"cookie"`
}

func New(path string) Config {
	cfg := config{}

	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to read config: %s", path)))
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to unmarshal config: %s", path)))
	}

	cfg.Logger = NewLogger(cfg.LogLevel)

	return &cfg
}

func (c *config) GetAuth() Auth {
	return c.Auth
}
