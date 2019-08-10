package horde_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHorde(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Horde Suite")
}
