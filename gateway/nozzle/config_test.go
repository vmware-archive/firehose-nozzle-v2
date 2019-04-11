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
	"gateway/nozzle"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {
	It("GetConfig should return an error when missing required", func() {
		_, err := nozzle.GetConfig()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("UAA_URL"))
	})

	It("GetConfig should build config from environment", func() {
		os.Setenv("UAA_URL", "my_url")
		os.Setenv("UAA_USER", "my_user")
		os.Setenv("UAA_PASS", "my_pass")
		os.Setenv("LOG_STREAM_URL", "log-stream.sys.cf.example.com")
		c, err := nozzle.GetConfig()
		Expect(err).To(BeNil())
		Expect(c.UAAURL).To(Equal("my_url"))
		Expect(c.UAAUser).To(Equal("my_user"))
		Expect(c.LogStreamUrl).To(Equal("log-stream.sys.cf.example.com"))
		Expect(c.UAAPass).To(Equal("my_pass"))
	})
})
