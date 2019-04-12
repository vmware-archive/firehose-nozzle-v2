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
	"errors"
	"gateway/nozzle"
	"gateway/nozzle/nozzlefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Receiver", func() {
	It("Receiver runs", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			str := `data: {"content": "Something really interesting"}`
			w.Write([]byte(str))
			w.Write([]byte("\n"))
		})
		logStreamServer := httptest.NewServer(handler)
		c := nozzle.Config{
			LogStreamUrl: logStreamServer.URL,
		}
		uaaClient := &nozzlefakes.FakeUAA{}
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)

		err := nozzle.Receive(&c, uaaClient, shipper)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("EOF"))
		Expect(string(writer.Bytes())).To(Equal(`data: {"content": "Something really interesting"}` + "\n"))
	})

	It("Receiver sends auth", func() {
		var streamerRequest *http.Request
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			streamerRequest = r
		})
		logStreamServer := httptest.NewServer(handler)
		c := nozzle.Config{
			LogStreamUrl: logStreamServer.URL,
		}
		uaaClient := &nozzlefakes.FakeUAA{}
		uaaClient.GetAuthTokenReturns("MyCrazyAuth", nil)
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)

		nozzle.Receive(&c, uaaClient, shipper)

		Expect(streamerRequest.Header.Get("Authorization")).To(Equal("MyCrazyAuth"))
	})

	It("Receiver requests multiple envelope types", func() {
		var streamerRequest *http.Request
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			streamerRequest = r
		})
		logStreamServer := httptest.NewServer(handler)
		c := nozzle.Config{
			LogStreamUrl: logStreamServer.URL,
			Envelopes:    []string{"counter", "log", "gauge"},
		}
		uaaClient := &nozzlefakes.FakeUAA{}
		uaaClient.GetAuthTokenReturns("MyCrazyAuth", nil)
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)

		nozzle.Receive(&c, uaaClient, shipper)

		Expect(streamerRequest.RequestURI).To(Equal("/v2/read?counter&log&gauge"))
	})

	It("Receiver surfaces error from UAA GetAuthToken", func() {
		uaaClient := &nozzlefakes.FakeUAA{}
		uaaClient.GetAuthTokenReturns("", errors.New("bad authorization"))
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)

		err := nozzle.Receive(&nozzle.Config{}, uaaClient, shipper)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("bad authorization"))
	})
})
