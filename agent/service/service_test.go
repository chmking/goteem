package service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde"
	. "github.com/chmking/horde/agent/service"
	"github.com/chmking/horde/protobuf/private"
)

var _ = Describe("Service", func() {
	var service *Service

	BeforeEach(func() {
		service = New(horde.Config{})
	})

	Describe("Healthcheck", func() {
		var (
			req *private.HealthcheckRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &private.HealthcheckRequest{}
			ctx = context.Background()
		})

		Context("when the agent does not return an error", func() {
			It("returns a public.HealthcheckResponse", func() {
				resp, _ := service.Healthcheck(ctx, req)
				Expect(resp).To(BeEquivalentTo(&private.HealthcheckResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Healthcheck(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Scale", func() {
		var (
			req *private.Orders
			ctx context.Context
		)

		BeforeEach(func() {
			req = &private.Orders{}
			ctx = context.Background()
		})

		Context("when the agent does not return an error", func() {
			It("returns a public.ScaleResponse", func() {
				resp, _ := service.Scale(ctx, req)
				Expect(resp).To(BeEquivalentTo(&private.ScaleResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Scale(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Stop", func() {
		var (
			req *private.StopRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &private.StopRequest{}
			ctx = context.Background()
		})

		Context("when the agent does not return an error", func() {
			It("returns a public.StopResponse", func() {
				resp, _ := service.Stop(ctx, req)
				Expect(resp).To(BeEquivalentTo(&private.StopResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Stop(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Quit", func() {
		var (
			req *private.QuitRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &private.QuitRequest{}
			ctx = context.Background()
		})

		Context("when the agent does not return an error", func() {
			It("returns a public.QuitResponse", func() {
				resp, _ := service.Quit(ctx, req)
				Expect(resp).To(BeEquivalentTo(&private.QuitResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Quit(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})
})
