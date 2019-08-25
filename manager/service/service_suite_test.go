//go:generate mockgen -package=service -destination=mock_manager_test.go github.com/chmking/horde/manager/service Manager
package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}
