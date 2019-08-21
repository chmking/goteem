package service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/manager/service"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
)

var _ = Describe("Service", func() {
	var service *Service

	BeforeEach(func() {
		service = New()
	})

	Describe("Start", func() {
		var (
			req *public.StartRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &public.StartRequest{}
			ctx = context.Background()
		})

		Context("when the manager does not return an error", func() {
			It("returns a public.StartResponse", func() {
				resp, _ := service.Start(ctx, req)
				Expect(resp).To(BeEquivalentTo(&public.StartResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Start(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Status", func() {
		var (
			req *public.StatusRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &public.StatusRequest{}
			ctx = context.Background()
		})

		Context("when the manager does not return an error", func() {
			It("returns a public.StatusResponse", func() {
				resp, _ := service.Status(ctx, req)
				Expect(resp).To(BeEquivalentTo(&public.StatusResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Status(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Stop", func() {
		var (
			req *public.StopRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &public.StopRequest{}
			ctx = context.Background()
		})

		Context("when the manager does not return an error", func() {
			It("returns a public.StopResponse", func() {
				resp, _ := service.Stop(ctx, req)
				Expect(resp).To(BeEquivalentTo(&public.StopResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Stop(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Quit", func() {
		var (
			req *public.QuitRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &public.QuitRequest{}
			ctx = context.Background()
		})

		Context("when the manager does not return an error", func() {
			It("returns a public.QuitResponse", func() {
				resp, _ := service.Quit(ctx, req)
				Expect(resp).To(BeEquivalentTo(&public.QuitResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Quit(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("Register", func() {
		var (
			req *private.RegisterRequest
			ctx context.Context
		)

		BeforeEach(func() {
			req = &private.RegisterRequest{}
			ctx = context.Background()
		})

		Context("when the manager does not return an error", func() {
			It("returns a private.RegisterResponse", func() {
				resp, _ := service.Register(ctx, req)
				Expect(resp).To(BeEquivalentTo(&private.RegisterResponse{}))
			})

			It("does not return an error", func() {
				_, err := service.Register(ctx, req)
				Expect(err).To(BeNil())
			})
		})
	})
})
