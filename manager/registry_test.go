package manager_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/chmking/horde/manager"
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
				agent := Registration{}
				for i := 0; i < 5; i++ {
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
			PContext("and is active", func() {

			})

			PContext("and it quarantined", func() {

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

			PIt("returns quarantined registered agents", func() {

			})

			It("returns active registered agents", func() {
				all := r.GetAll()
				Expect(all).To(BeEquivalentTo(expected))
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
			var expected []Registration

			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}

					expected = append(expected, agent)
					r.Add(agent)
				}
			})

			PIt("does not return quarantined registered agents", func() {

			})

			It("returns active registered agents", func() {
				active := r.GetActive()
				Expect(active).To(BeEquivalentTo(expected))
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
			var expected []Registration

			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					agent := Registration{Id: strconv.Itoa(i)}

					expected = append(expected, agent)
					r.Add(agent)
				}
			})

			It("does not return active registered agents", func() {

			})

			It("returns quarantined registered agents", func() {

			})
		})
	})
})
