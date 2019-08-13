package horde

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("RecorderFrom", func() {
		Context("when there is an embedded recorder", func() {
			var recorder *Recorder

			BeforeEach(func() {
				recorder = &Recorder{}

				ctx = context.WithValue(ctx, recorderKey, recorder)
				Expect(ctx.Value(recorderKey).(*Recorder)).To(Equal(recorder))
			})

			It("returns the recorder", func() {
				result := RecorderFrom(ctx)
				Expect(result).To(Equal(recorder))
			})
		})

		Context("when there is no embedded recorder", func() {
			It("return a nil Recorder", func() {
				result := RecorderFrom(ctx)
				Expect(result).To(BeNil())
			})
		})
	})
})
