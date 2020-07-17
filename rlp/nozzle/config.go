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
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CACertPath  string `envconfig:"CA_CERT_PATH" required:"true"`
	CertPath    string `envconfig:"CERT_PATH" required:"true"`
	KeyPath     string `envconfig:"KEY_PATH" required:"true"`
	CommonName  string `envconfig:"COMMON_NAME"`
	LogsAPIAddr string `envconfig:"LOGS_API_ADDR" required:"true"`
	ShardID     string `envconfig:"SHARD_ID" required:"true"`
}

func GetConfig() (*Config, error) {
	c := &Config{
		CommonName: "reverselogproxy",
	}
	err := envconfig.Process("", c)
	return c, err
}
