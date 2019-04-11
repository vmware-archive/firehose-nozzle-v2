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
	"bytes"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"encoding/base64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"rlp/nozzle"
)

var _ = Describe("LogShipper", func() {
	It("Shipper ships", func() {
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)
		err := shipper.LogShip(&loggregator_v2.Envelope{
			SourceId: "unit-test",
			Message: &loggregator_v2.Envelope_Log{
				Log: &loggregator_v2.Log{
					Payload: []byte("Hello test"),
				},
			},
		})

		Expect(err).To(BeNil())

		messageString := string(writer.Bytes())
		Expect(messageString).To(ContainSubstring("unit-test"))

		encoded := base64.StdEncoding.EncodeToString([]byte("Hello test"))
		Expect(messageString).To(ContainSubstring(encoded))
	})
})
