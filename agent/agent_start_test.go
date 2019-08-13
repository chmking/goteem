package agent_test

import (
	"context"
	"time"

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

	Describe("Start", func() {
		var req StartRequest

		BeforeEach(func() {
			req = StartRequest{}
		})

		Context("when the agent is IDLE", func() {
			BeforeEach(func() {
				agent.Status = Status_IDLE
			})

			Context("and scaling completes", func() {
				It("sets the status to RUNNING", func() {
					agent.Start(ctx, &req)
					<-time.After(time.Millisecond)

					Expect(agent.SafeStatus()).To(Equal(Status_RUNNING))
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

			Context("and scaling is active", func() {
				BeforeEach(func() {
					agent.Session = &MockSession{}
				})

				It("sets the status to SCALING", func() {
					agent.Start(ctx, &req)
					Expect(agent.SafeStatus()).To(Equal(Status_SCALING))
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
		})

		Context("when the status is SCALING", func() {
			BeforeEach(func() {
				agent.Status = Status_SCALING
			})

			Context("and scaling completes", func() {
				It("sets the status to RUNNING", func() {
					agent.Start(ctx, &req)
					<-time.After(time.Millisecond)
					Expect(agent.SafeStatus()).To(Equal(Status_RUNNING))
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

			Context("and scaling is active", func() {
				BeforeEach(func() {
					agent.Session = &MockSession{}
				})

				It("sets the status to SCALING", func() {
					agent.Start(ctx, &req)
					Expect(agent.SafeStatus()).To(Equal(Status_SCALING))
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
		})

		Context("when the status is RUNNING", func() {
			BeforeEach(func() {
				agent.Status = Status_RUNNING
			})

			Context("and scaling completes", func() {
				It("sets the status to RUNNING", func() {
					agent.Start(ctx, &req)
					<-time.After(time.Millisecond)
					Expect(agent.SafeStatus()).To(Equal(Status_RUNNING))
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

			Context("and scaling is active", func() {
				BeforeEach(func() {
					agent.Session = &MockSession{}
				})

				It("sets the status to SCALING", func() {
					agent.Start(ctx, &req)
					Expect(agent.SafeStatus()).To(Equal(Status_SCALING))
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
		})

		Context("when the agent is STOPPING", func() {
			BeforeEach(func() {
				agent.Status = Status_STOPPING
			})

			It("leaves the status STOPPING", func() {
				agent.Start(ctx, &req)
				Expect(agent.SafeStatus()).To(Equal(Status_STOPPING))
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
				agent.Status = Status_QUITTING
			})

			It("leaves the status QUITTING", func() {
				agent.Start(ctx, &req)
				Expect(agent.SafeStatus()).To(Equal(Status_QUITTING))
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
