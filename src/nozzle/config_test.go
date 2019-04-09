package nozzle_test

import (
	"github.com/cf-platform-eng/firehose-nozzle-v2/src/nozzle"
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
		c, err := nozzle.GetConfig()
		Expect(err).To(BeNil())
		Expect(c.UaaURL).To(Equal("my_url"))
	})
})
