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
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CACertPath string `envconfig:"CA_CERT_PATH"`
	CertPath   string `envconfig:"CERT_PATH"`
	KeyPath    string `envconfig:"KEY_PATH"`
	CommonName string `envconfig:"COMMON_NAME"`

	CACert string `envconfig:"CA_CERT"`
	Cert   string `envconfig:"CERT"`
	Key    string `envconfig:"KEY"`

	LogsAPIAddr string `envconfig:"LOGS_API_ADDR"`
	ShardID     string `envconfig:"SHARD_ID" required:"true"`
	PrintStats  bool   `envconfig:"PRINT_STATS"`
}

func GetConfig() (*Config, error) {
	c := &Config{
		CommonName: "reverselogproxy",
	}
	err := envconfig.Process("", c)

	if c.CACert == "" && c.CACertPath == "" {
		return nil, errors.New("one of CA_CERT or CA_CERT_PATH is required")
	}

	if c.Cert == "" && c.CertPath == "" {
		return nil, errors.New("one of CERT or CERT_PATH is required")
	}

	if c.Key == "" && c.KeyPath == "" {
		return nil, errors.New("one of KEY or KEY_PATH is required")
	}

	return c, err
}
