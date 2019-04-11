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
