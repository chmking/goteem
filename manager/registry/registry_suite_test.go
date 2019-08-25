//go:generate mockgen -package=registry_test -destination=mock_private_test.go github.com/chmking/horde/protobuf/private AgentClient
package registry_test

import (
	"testing"

	"github.com/chmking/horde/logger"
	"github.com/chmking/horde/logger/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Suite")
}

var _ = BeforeSuite(func() {
	log.Logger = logger.NewStdLogger(GinkgoWriter)
})
