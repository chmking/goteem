package goteem_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoteem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goteem Suite")
}
