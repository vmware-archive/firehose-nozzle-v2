// firehose-nozzle-v2
//
// Copyright (c) 2018-Present Pivotal Software, Inc. All Rights Reserved.
//
// This program and the accompanying materials are made available under the terms of the under the Apache License,
// Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may
// obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing permissions and
// limitations under the License.

package nozzle

import (
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UAAURL            string   `envconfig:"UAA_URL" required:"true"`
	UAAUser           string   `envconfig:"UAA_USER" required:"true"`
	UAAPass           string   `envconfig:"UAA_PASS" required:"true"`
	SkipSSLValidation bool     `envconfig:"SKIP_SSL_VALIDATION"`
	LogStreamUrl      string   `envconfig:"LOG_STREAM_URL" required:"true"`
	Envelopes         []string `envconfig:"envelopes" required:"true"`
}

func GetConfig() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		return nil, err
	}
	err = c.validEnvelopes()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validEnvelopes() error {
	for _, e := range c.Envelopes {
		if e != "log" && e != "counter" && e != "gauge" && e != "timer" && e != "event" {
			return errors.New(fmt.Sprintf(
				"'%s' is not a valid envelope type. Allowed values are: log, counter, gauge, timer, event", e,
			))
		}
	}

	return nil
}
