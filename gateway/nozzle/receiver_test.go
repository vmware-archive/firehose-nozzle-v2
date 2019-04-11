package nozzle_test

import (
	"bytes"
	"errors"
	"rlp/nozzle"
	"rlp/nozzle/nozzlefakes"
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
