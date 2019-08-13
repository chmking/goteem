package agent_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde"
	. "github.com/chmking/horde/agent"
	. "github.com/chmking/horde/protobuf/private"
)

var _ = Describe("Agent", func() {
	var (
		ctx   context.Context
		agent *Agent
	)

	BeforeEach(func() {
		ctx = context.Background()
		agent = New(Config{})
	})

	Describe("Quit", func() {
		var req QuitRequest

		BeforeEach(func() {
			req = QuitRequest{}
		})

		Context("when the agent is IDLE", func() {
			BeforeEach(func() {
				agent.Status = Status_IDLE
			})

			It("sets the status to QUITTING", func() {
				agent.Quit(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a QuitResponse", func() {
				resp, _ := agent.Quit(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Quit(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the status is SCALING", func() {
			BeforeEach(func() {
				agent.Status = Status_SCALING
			})

			It("sets the status to QUITTING", func() {
				agent.Quit(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a QuitResponse", func() {
				resp, _ := agent.Quit(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Quit(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the status is RUNNING", func() {
			BeforeEach(func() {
				agent.Status = Status_RUNNING
			})

			It("sets the status to QUITTING", func() {
				agent.Quit(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a QuitResponse", func() {
				resp, _ := agent.Quit(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Quit(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the agent is STOPPING", func() {
			BeforeEach(func() {
				agent.Status = Status_STOPPING
			})

			It("sets the status to QUITTING", func() {
				agent.Quit(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a QuitResponse", func() {
				resp, _ := agent.Quit(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Quit(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the agent is QUITTING", func() {
			BeforeEach(func() {
				agent.Status = Status_QUITTING
			})

			It("leaves the status QUITTING", func() {
				agent.Quit(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a QuitResponse", func() {
				resp, _ := agent.Quit(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Quit(ctx, &req)
				Expect(err).To(BeNil())
			})
		})
	})
})
