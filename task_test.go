package horde_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde"
)

var _ = Describe("Task", func() {
	Describe("Exec", func() {
		var (
			execed bool
			ctx    context.Context
		)

		BeforeEach(func() {
			execed = false
			ctx = context.Background()
		})

		It("executes the function", func() {
			fn := TaskFunc(func(ctx context.Context) {
				execed = true
			})

			fn.Exec(ctx)
			Expect(execed).To(BeTrue())
		})
	})
})
