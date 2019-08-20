package manager_test

import (
	"errors"
	"strconv"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/manager"
	"github.com/chmking/horde/protobuf/private"
)

var _ = Describe("Registry", func() {
	var r *Registry

	BeforeEach(func() {
		r = NewRegistry()
	})

	Describe("Len", func() {
		Context("when the registry is empty", func() {
			It("returns 0", func() {
				value := r.Len()
				Expect(value).To(BeNumerically("==", 0))
			})
		})

		Context("when there are registered agents", func() {
			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}
					r.Add(agent)
				}
			})

			It("returns the number of agents", func() {
				value := r.Len()
				Expect(value).To(BeNumerically("==", 5))
			})
		})
	})

	Describe("Add", func() {
		var agent Registration

		BeforeEach(func() {
			agent = Registration{Id: "foo"}
		})

		Context("when the agent is not in the registry", func() {
			It("adds the agent to the registry", func() {
				r.Add(agent)
				Expect(r.GetAll()[0]).To(BeEquivalentTo(agent))
			})

			It("marks the agent as active", func() {
				r.Add(agent)
				Expect(r.GetActive()[0]).To(BeEquivalentTo(agent))
			})

			It("does not return an error", func() {
				err := r.Add(agent)
				Expect(err).To(BeNil())
			})
		})

		Context("when the agent is already in the registry", func() {
			BeforeEach(func() {
				r.Add(agent)
			})

			Context("and is active", func() {
				It("does not add a duplicate", func() {
					r.Add(agent)
					Expect(r.Len()).To(BeNumerically("==", 1))
				})
			})

			Context("and is quarantined", func() {
				BeforeEach(func() {
					r.Quarantine(agent.Id)
				})

				It("does not add a duplicate", func() {
					r.Add(agent)
					Expect(r.Len()).To(BeNumerically("==", 1))
				})
			})
		})
	})

	Describe("Quarantine", func() {
		Context("when the agent doesn't exist", func() {
			It("return an ErrInvalidAgent", func() {
				err := r.Quarantine("foo")
				Expect(err).To(Equal(ErrInvalidAgent))
			})
		})

		Context("when the agent does exist", func() {
			var agent Registration

			BeforeEach(func() {
				agent = Registration{Id: "foo"}
				r.Add(agent)
			})

			Context("and is already quarantined", func() {
				BeforeEach(func() {
					err := r.Quarantine(agent.Id)
					Expect(err).To(BeNil())
				})

				It("leaves the agent in quarantine", func() {
					quarantined := r.GetQuarantined()
					Expect(quarantined).To(ContainElement(agent))
				})

				It("does not add the agent to active", func() {
					active := r.GetActive()
					Expect(active).NotTo(ContainElement(agent))
				})

				It("does not return and error", func() {
					err := r.Quarantine("foo")
					Expect(err).To(BeNil())
				})
			})

			Context("and is active", func() {
				It("moves the agent to quarantine", func() {
					quarantined := r.GetQuarantined()
					Expect(quarantined).NotTo(ContainElement(agent))
				})

				It("does not leave the agent active", func() {
					active := r.GetActive()
					Expect(active).To(ContainElement(agent))
				})

				It("does not return and error", func() {
					err := r.Quarantine("foo")
					Expect(err).To(BeNil())
				})
			})
		})
	})

	Describe("GetAll", func() {
		Context("when the registry is empty", func() {
			It("returns nil", func() {
				all := r.GetAll()
				Expect(all).To(BeNil())
			})
		})

		Context("when there are registered agents", func() {
			var expected []Registration

			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}

					expected = append(expected, agent)
					r.Add(agent)
				}
			})

			It("returns quarantined registered agents", func() {
				all := r.GetAll()
				for _, value := range r.GetQuarantined() {
					Expect(all).To(ContainElement(value))
				}
			})

			It("returns active registered agents", func() {
				all := r.GetAll()
				for _, value := range r.GetActive() {
					Expect(all).To(ContainElement(value))
				}
			})
		})
	})

	Describe("GetActive", func() {
		Context("when the registry is empty", func() {
			It("returns nil", func() {
				active := r.GetActive()
				Expect(active).To(BeNil())
			})
		})

		Context("when there are registered agents", func() {
			var active []Registration

			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}
					active = append(active, agent)
					r.Add(agent)
				}
			})

			It("returns active registered agents", func() {
				result := r.GetActive()
				Expect(result).To(BeEquivalentTo(active))
			})

			Context("and there are quarantined agents", func() {
				var quarantined []Registration

				BeforeEach(func() {
					for i := 5; i < 10; i++ {
						agent := Registration{Id: strconv.Itoa(i)}
						quarantined = append(quarantined, agent)
						r.Add(agent)
						r.Quarantine(agent.Id)
					}
				})

				It("does not return quarantined registered agents", func() {
					result := r.GetActive()
					for _, value := range quarantined {
						Expect(result).NotTo(ContainElement(value))
					}
				})
			})
		})
	})

	Describe("GetQuarantined", func() {
		Context("when the registry is empty", func() {
			It("returns nil", func() {
				quarantined := r.GetQuarantined()
				Expect(quarantined).To(BeNil())
			})
		})

		Context("when there are registered agents", func() {
			var quarantined []Registration

			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}

					quarantined = append(quarantined, agent)
					r.Add(agent)
					r.Quarantine(agent.Id)
				}
			})

			It("returns quarantined registered agents", func() {
				result := r.GetQuarantined()
				Expect(result).To(BeEquivalentTo(quarantined))
			})

			Context("and there are quarantined agents", func() {
				var active []Registration

				BeforeEach(func() {
					for i := 5; i < 10; i++ {
						agent := Registration{Id: strconv.Itoa(i)}
						active = append(active, agent)
						r.Add(agent)
					}
				})

				It("does not return active registered agents", func() {
					result := r.GetQuarantined()
					for _, value := range active {
						Expect(result).NotTo(ContainElement(value))
					}
				})
			})
		})
	})

	Describe("Healthcheck", func() {
		var (
			agent           Registration
			mockCtrl        *gomock.Controller
			mockAgentClient *MockAgentClient
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockAgentClient = NewMockAgentClient(mockCtrl)

			agent = Registration{Id: "foo"}
		})

		JustBeforeEach(func() {
			r.Add(agent)
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		Context("when an active agent Healthcheck returns errors", func() {
			BeforeEach(func() {
				mockAgentClient.EXPECT().Heartbeat(gomock.Any(), gomock.Any()).
					Return(&private.HeartbeatResponse{}, errors.New("foo")).AnyTimes()
				agent.Client = mockAgentClient
			})

			Context("one time", func() {
				It("leaves the agent active", func() {
					r.Healthcheck()
					active := r.GetActive()
					Expect(active).To(ContainElement(agent))
				})
			})

			Context("two times", func() {
				JustBeforeEach(func() {
					r.Healthcheck()
				})

				It("leaves the agent active", func() {
					r.Healthcheck()
					active := r.GetActive()
					Expect(active).To(ContainElement(agent))
				})
			})

			Context("three times", func() {
				JustBeforeEach(func() {
					r.Healthcheck()
					r.Healthcheck()
				})

				It("quarantines the agent", func() {
					r.Healthcheck()
					quarantined := r.GetQuarantined()
					Expect(quarantined).To(ContainElement(agent))
				})
			})
		})

		Context("when an quarantined agent Healthcheck passes", func() {
			BeforeEach(func() {
				mockAgentClient.EXPECT().Heartbeat(gomock.Any(), gomock.Any()).
					Return(&private.HeartbeatResponse{}, nil).AnyTimes()
				agent.Client = mockAgentClient
			})

			JustBeforeEach(func() {
				r.Quarantine(agent.Id)
			})

			Context("one time", func() {
				It("leaves the agent quarantined", func() {
					r.Healthcheck()
					quarantined := r.GetQuarantined()
					Expect(quarantined).To(ContainElement(agent))
				})
			})

			Context("two times", func() {
				JustBeforeEach(func() {
					r.Healthcheck()
				})

				It("leaves the agent quarantined", func() {
					r.Healthcheck()
					quarantined := r.GetQuarantined()
					Expect(quarantined).To(ContainElement(agent))
				})
			})

			Context("three times", func() {
				JustBeforeEach(func() {
					r.Healthcheck()
					r.Healthcheck()
				})

				It("activates the agent", func() {
					r.Healthcheck()
					active := r.GetActive()
					Expect(active).To(ContainElement(agent))
				})
			})
		})
	})
})
