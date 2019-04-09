package nozzle_test

import (
	"bytes"
	"github.com/cf-platform-eng/firehose-nozzle-v2/src/nozzle"
	"github.com/cf-platform-eng/firehose-nozzle-v2/src/nozzle/nozzlefakes"
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
			w.Write([]byte("\n "))
			//bazaarAPIRequest = r
		})
		logStreamServer := httptest.NewServer(handler)

		c := nozzle.Config{
			LogStreamUrl: logStreamServer.URL,
		}
		uaaClient := &nozzlefakes.FakeUAA{}
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)
		err := nozzle.GatewayMain(&c, uaaClient, shipper)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("EOF"))
		Expect(string(writer.Bytes())).To(Equal(`data: {"content": "Something really interesting"}`))
	})
})
