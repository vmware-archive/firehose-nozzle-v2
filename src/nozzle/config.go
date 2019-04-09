package nozzle

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UaaURL string `envconfig:"UAA_URL" required:"true"`
}

func GetConfig() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	return c, err
}