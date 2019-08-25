package state_test

import (
	"testing"

	"github.com/chmking/horde/logger"
	"github.com/chmking/horde/logger/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State Suite")
}

var _ = BeforeSuite(func() {
	log.Logger = logger.NewStdLogger(GinkgoWriter)
})
