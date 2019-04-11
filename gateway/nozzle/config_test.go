package nozzle_test

import (
	"rlp/nozzle"
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
