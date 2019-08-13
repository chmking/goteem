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

	Describe("Stop", func() {
		var req StopRequest

		BeforeEach(func() {
			req = StopRequest{}
		})

		Context("when the status is IDLE", func() {
			BeforeEach(func() {
				agent.Status = Status_IDLE
			})

			It("leaves the status IDLE", func() {
				agent.Stop(ctx, &req)
				Expect(agent.Status).To(Equal(Status_IDLE))
			})

			It("returns a StopResponse", func() {
				resp, _ := agent.Stop(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Stop(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the status is SCALING", func() {
			BeforeEach(func() {
				agent.Status = Status_SCALING
			})

			Context("and stop completes", func() {
				It("sets the status to IDLE", func() {
					agent.Stop(ctx, &req)
					Expect(agent.Status).To(Equal(Status_IDLE))
				})

				It("returns a StopResponse", func() {
					resp, _ := agent.Stop(ctx, &req)
					Expect(resp).NotTo(BeNil())
				})

				It("does not return an error", func() {
					_, err := agent.Stop(ctx, &req)
					Expect(err).To(BeNil())
				})
			})

			Context("and stop is active", func() {
				BeforeEach(func() {
					agent.Session = &MockSession{}
				})

				It("sets the status to STOPPING", func() {
					agent.Stop(ctx, &req)
					Expect(agent.Status).To(Equal(Status_STOPPING))
				})

				It("returns a StopResponse", func() {
					resp, _ := agent.Stop(ctx, &req)
					Expect(resp).NotTo(BeNil())
				})

				It("does not return an error", func() {
					_, err := agent.Stop(ctx, &req)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("when the status is RUNNING", func() {
			BeforeEach(func() {
				agent.Status = Status_RUNNING
			})

			Context("and stop completes", func() {
				It("sets the status to IDLE", func() {
					agent.Stop(ctx, &req)
					Expect(agent.Status).To(Equal(Status_IDLE))
				})

				It("returns a StopResponse", func() {
					resp, _ := agent.Stop(ctx, &req)
					Expect(resp).NotTo(BeNil())
				})

				It("does not return an error", func() {
					_, err := agent.Stop(ctx, &req)
					Expect(err).To(BeNil())
				})
			})

			Context("and stop is active", func() {
				BeforeEach(func() {
					agent.Session = &MockSession{}
				})

				It("sets the status to STOPPING", func() {
					agent.Stop(ctx, &req)
					Expect(agent.Status).To(Equal(Status_STOPPING))
				})

				It("returns a StopResponse", func() {
					resp, _ := agent.Stop(ctx, &req)
					Expect(resp).NotTo(BeNil())
				})

				It("does not return an error", func() {
					_, err := agent.Stop(ctx, &req)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("when the agent is STOPPING", func() {
			BeforeEach(func() {
				agent.Status = Status_STOPPING
			})

			It("leaves the status STOPPING", func() {
				agent.Stop(ctx, &req)
				Expect(agent.Status).To(Equal(Status_STOPPING))
			})

			It("returns a StopResponse", func() {
				resp, _ := agent.Stop(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Stop(ctx, &req)
				Expect(err).To(BeNil())
			})
		})

		Context("when the agent is QUITTING", func() {
			BeforeEach(func() {
				agent.Status = Status_QUITTING
			})

			It("leaves the status QUITTING", func() {
				agent.Stop(ctx, &req)
				Expect(agent.Status).To(Equal(Status_QUITTING))
			})

			It("returns a StopResponse", func() {
				resp, _ := agent.Stop(ctx, &req)
				Expect(resp).NotTo(BeNil())
			})

			It("does not return an error", func() {
				_, err := agent.Stop(ctx, &req)
				Expect(err).To(BeNil())
			})
		})
	})
})
