package nozzle_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNozzle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nozzle Suite")
}
