//go:generate mockgen -package=manager_test -destination=mock_private_test.go github.com/chmking/horde/protobuf/private AgentClient
//go:generate mockgen -package=manager_test -destination=mock_manager_test.go github.com/chmking/horde/manager Registry,StateMachine
package manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manager Suite")
}
