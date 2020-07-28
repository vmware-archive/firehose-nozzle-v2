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

package nozzle_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"rlp/nozzle"
)

var _ = Describe("Config", func() {
	It("GetConfig should return an error when missing required", func() {
		_, err := nozzle.GetConfig()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("CA_CERT_PATH"))
	})

	AfterEach(func() {
		os.Clearenv()
	})

	It("GetConfig should build config from environment", func() {
		os.Setenv("CA_CERT_PATH", "my_ca_cert_path")
		os.Setenv("CERT_PATH", "my_cert_path")
		os.Setenv("KEY_PATH", "my_key_path")
		os.Setenv("LOGS_API_ADDR", "my_log_api_addr")
		os.Setenv("SHARD_ID", "my_shard_id")

		c, err := nozzle.GetConfig()

		Expect(err).To(BeNil())
		Expect(c.CACertPath).To(Equal("my_ca_cert_path"))
		Expect(c.CertPath).To(Equal("my_cert_path"))
		Expect(c.KeyPath).To(Equal("my_key_path"))
		Expect(c.LogsAPIAddr).To(Equal("my_log_api_addr"))
		Expect(c.ShardID).To(Equal("my_shard_id"))
	})

	It("one method of config is required", func() {
		_, err := nozzle.GetConfig()

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("CA_CERT"))
		Expect(err.Error()).To(ContainSubstring("CA_CERT_PATH"))
	})
})
