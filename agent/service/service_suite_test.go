//go:generate mockgen -package=service -destination=mock_agent_test.go github.com/chmking/horde/agent/service Agent
package service_test

import (
	"testing"

	"github.com/chmking/horde/logger"
	"github.com/chmking/horde/logger/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = BeforeSuite(func() {
	log.Logger = logger.NewStdLogger(GinkgoWriter)
})
