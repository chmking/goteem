package service

import (
	"context"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
)

var _ = Describe("Service", func() {
	var (
		service *Service
		ctx     context.Context

		mockCtrl    *gomock.Controller
		mockManager *MockManager
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockManager = NewMockManager(mockCtrl)

		service = New()
		service.manager = mockManager

		ctx = context.Background()
	})

	Describe("Start", func() {
		var req *public.StartRequest

		BeforeEach(func() {
			req = &public.StartRequest{}
		})

		Context("when the manager does not return an error", func() {
			BeforeEach(func() {
				mockManager.EXPECT().Start(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			})

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

	Describe("Stop", func() {
		var req *public.StopRequest

		BeforeEach(func() {
			req = &public.StopRequest{}
		})

		Context("when the manager does not return an error", func() {
			BeforeEach(func() {
				mockManager.EXPECT().Stop().Return(nil).AnyTimes()
			})

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
		var req *public.QuitRequest

		BeforeEach(func() {
			req = &public.QuitRequest{}
		})

		Context("when the manager does not return an error", func() {
			BeforeEach(func() {
				mockManager.EXPECT().Stop().Return(nil).AnyTimes()
			})

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
		var req *private.RegisterRequest

		BeforeEach(func() {
			req = &private.RegisterRequest{}
		})

		Context("when the manager does not return an error", func() {
			BeforeEach(func() {
				mockManager.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			})

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
