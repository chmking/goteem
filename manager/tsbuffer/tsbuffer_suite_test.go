package tsbuffer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTsbuffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tsbuffer Suite")
}
