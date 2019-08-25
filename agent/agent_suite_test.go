//go:generate mockgen -package=agent -destination=mock_agent_test.go github.com/chmking/horde/agent Session,StateMachine
package agent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAgent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agent Suite")
}
