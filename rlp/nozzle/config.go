package nozzle

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UAAURL       string `envconfig:"UAA_URL" required:"true"`
	UAAUser      string `envconfig:"UAA_USER" required:"true"`
	UAAPass      string `envconfig:"UAA_PASS" required:"true"`
	LogStreamUrl string `envconfig:"LOG_STREAM_URL" required:"true"`
}

func GetConfig() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	return c, err
}
