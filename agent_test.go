package goteem_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/goteem"
)

const address = ":5557"

var _ = Describe("Agent", func() {

	var agent Agent

	BeforeEach(func() {
		agent = Agent{}
		agent.Listen(address)
	})

	Describe("Teem", func() {
		var req TeemRequest
	})
})
