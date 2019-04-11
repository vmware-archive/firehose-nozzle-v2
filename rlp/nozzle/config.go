package nozzle

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CACertPath       string `envconfig:"CA_CERT_PATH" required:"true"`
	CertPath      string `envconfig:"CERT_PATH" required:"true"`
	KeyPath      string `envconfig:"KEY_PATH" required:"true"`
	LogsAPIAddr string `envconfig:"LOGS_API_ADDR" required:"true"`
	ShardID string `envconfig:"SHARD_ID" required:"true"`
}

func GetConfig() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	return c, err
}
