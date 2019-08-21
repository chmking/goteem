package horde

import (
	"context"

	"github.com/chmking/horde/agent/recorder"
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
			var r *recorder.Recorder

			BeforeEach(func() {
				r = recorder.New()

				ctx = context.WithValue(ctx, recorderKey, r)
				Expect(ctx.Value(recorderKey).(*recorder.Recorder)).To(Equal(r))
			})

			It("returns the recorder", func() {
				result := RecorderFrom(ctx)
				Expect(result).To(Equal(r))
			})
		})

		Context("when there is no embedded recorder", func() {
			It("return a nil Recorder", func() {
				result := RecorderFrom(ctx)
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("WithRecorder", func() {
		var (
			r   *recorder.Recorder
			ctx context.Context
		)

		BeforeEach(func() {
			r = recorder.New()
			ctx = context.Background()
		})

		It("embeds the Recorder", func() {
			ctx = WithRecorder(ctx, r)
			result := RecorderFrom(ctx)
			Expect(result).To(Equal(r))
		})
	})
})
