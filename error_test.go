package horde_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde"
)

var _ = Describe("Error", func() {
	Describe("Error", func() {
		It("returns the string", func() {
			err := Error("foo")
			result := err.Error()

			Expect(result).To(Equal("foo"))
		})
	})
})
