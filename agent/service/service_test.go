package service_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/chmking/horde"
	agent "github.com/chmking/horde/agent"
	. "github.com/chmking/horde/agent/service"
	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
)

var _ = Describe("Service", func() {
	var (
		service *Service
		ctx     context.Context

		mockCtrl  *gomock.Controller
		mockAgent *MockAgent
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAgent = NewMockAgent(mockCtrl)

		service = New(horde.Config{})
		service.Agent = mockAgent

		ctx = context.Background()
	})

	Describe("Healthcheck", func() {
		var req *private.HealthcheckRequest

		BeforeEach(func() {
			req = &private.HealthcheckRequest{}

			mockAgent.EXPECT().Status().Return(agent.Status{
				State: public.Status_STATUS_SCALING,
				Count: 1,
			}).AnyTimes()
		})

		It("returns a public.HealthcheckResponse", func() {
			resp, _ := service.Healthcheck(ctx, req)
			Expect(resp).NotTo(BeNil())
			Expect(*resp).To(MatchFields(IgnoreExtras, Fields{
				"State": Equal(public.Status_STATUS_SCALING),
				"Count": BeNumerically("==", 1),
			}))
		})

		It("does not return an error", func() {
			_, err := service.Healthcheck(ctx, req)
			Expect(err).To(BeNil())
		})
	})

	Describe("Scale", func() {
		var req *private.ScaleRequest

		BeforeEach(func() {
			req = &private.ScaleRequest{Orders: &private.Orders{}}
		})

		Context("when the agent does not return an error", func() {
			BeforeEach(func() {
				mockAgent.EXPECT().Scale(gomock.Any()).Return(nil).AnyTimes()
			})

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
		var req *private.StopRequest

		BeforeEach(func() {
			req = &private.StopRequest{}
		})

		Context("when the agent does not return an error", func() {
			BeforeEach(func() {
				mockAgent.EXPECT().Stop().Return(nil).AnyTimes()
			})

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
		var req *private.QuitRequest

		BeforeEach(func() {
			req = &private.QuitRequest{}
		})

		Context("when the agent does not return an error", func() {
			BeforeEach(func() {
				mockAgent.EXPECT().Stop().Return(nil).AnyTimes()
			})

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
