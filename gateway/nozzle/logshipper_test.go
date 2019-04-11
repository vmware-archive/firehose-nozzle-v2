package nozzle_test

import (
	"bytes"
	"rlp/nozzle"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogShipper", func() {
	It("Shipper ships", func() {
		writer := bytes.Buffer{}
		shipper := nozzle.NewSampleShipper(&writer)
		err := shipper.LogShip("YOLO")
		Expect(err).To(BeNil())
		Expect(string(writer.Bytes())).To(Equal("YOLO\n"))
	})
})
