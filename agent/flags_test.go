package agent

import (
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flags", func() {
	It("adds flags to parse", func() {
		Flags()
		flag.Parse()

		Expect(agentHost).To(Equal(defaultAgentHost))
		Expect(agentPort).To(Equal(defaultAgentPort))

		Expect(managerHost).To(Equal(defaultManagerHost))
		Expect(managerPort).To(Equal(defaultManagerPort))
	})
})
