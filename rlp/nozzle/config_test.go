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
		Expect(err.Error()).To(ContainSubstring("CA_CERT_PATH"))
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
})
