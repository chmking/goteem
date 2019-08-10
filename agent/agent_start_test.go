package agent_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde"
	. "github.com/chmking/horde/agent"
)

var _ = Describe("Agent", func() {
	var (
		agent *Agent
		ctx   context.Context
	)

	BeforeEach(func() {
		agent = New(Config{})
		ctx = context.Background()
	})

	Describe("Start", func() {
		var req StartRequest

		BeforeEach(func() {
			req = StartRequest{}
		})

		Context("when the agent is IDLE", func() {
			BeforeEach(func() {
				Expect(agent.Status()).To(Equal(Status_IDLE))
			})

			It("sets the status to RUNNING", func() {
				agent.Start(ctx, &req)
				Expect(agent.Status()).To(Equal(Status_RUNNING))
			})

			It("returns a StartResponse", func() {
				resp, _ := agent.Start(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Start(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the status is RUNNING", func() {
			BeforeEach(func() {
				agent.Start(ctx, &StartRequest{})
				Expect(agent.Status()).To(Equal(Status_RUNNING))
			})

			It("leaves the status RUNNING", func() {
				agent.Start(ctx, &req)
				Expect(agent.Status()).To(Equal(Status_RUNNING))
			})

			It("returns a StartResponse", func() {
				resp, _ := agent.Start(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Start(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the agent is STOPPING", func() {
			BeforeEach(func() {
				agent.Start(ctx, &StartRequest{})
				agent.Stop(ctx, &StopRequest{})
				Expect(agent.Status()).To(Equal(Status_STOPPING))
			})

			It("leaves the status STOPPING", func() {
				agent.Start(ctx, &req)
				Expect(agent.Status()).To(Equal(Status_STOPPING))
			})

			It("does not return a StartResponse", func() {
				resp, _ := agent.Start(ctx, &req)
				Expect(resp).To(BeNil())
			})

			It("returns ErrStatusStopping", func() {
				_, err := agent.Start(ctx, &req)
				Expect(err).To(Equal(ErrStatusStopping))
			})
		})

		Context("when the agent is QUITTING", func() {
			BeforeEach(func() {
				agent.Quit(ctx, &QuitRequest{})
				Expect(agent.Status()).To(Equal(Status_QUITTING))
			})

			It("leaves the status QUITTING", func() {
				agent.Start(ctx, &req)
				Expect(agent.Status()).To(Equal(Status_QUITTING))
			})

			It("does not return a StartResponse", func() {
				resp, _ := agent.Start(ctx, &req)
				Expect(resp).To(BeNil())
			})

			It("returns ErrStatusQuitting", func() {
				_, err := agent.Start(ctx, &req)
				Expect(err).To(Equal(ErrStatusQuitting))
			})
		})
	})
})
