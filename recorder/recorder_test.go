package recorder_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde/protobuf/public"
	. "github.com/chmking/horde/recorder"
)

var _ = Describe("New", func() {
	It("constructs a recorder", func() {
		r := New()
		Expect(r).ToNot(BeNil())
	})
})

var _ = Describe("Recorder", func() {
	var r *Recorder

	BeforeEach(func() {
		r = New()
	})

	Describe("Success", func() {
		It("records a success", func() {
			r.Success("foo", "bar", int64(42))
			results := r.Results()
			Expect(results[0].Code).To(Equal(public.Code_CODE_SUCCESS))
		})
	})

	Describe("Error", func() {
		It("records an error", func() {
			r.Error("foo", "bar", int64(42), errors.New("baz"))
			results := r.Results()
			Expect(results[0].Code).To(Equal(public.Code_CODE_ERROR))
		})
	})

	Describe("Panic", func() {
		It("records a panic", func() {
			r.Panic("foo", "bar", int64(42), errors.New("baz"))
			results := r.Results()
			Expect(results[0].Code).To(Equal(public.Code_CODE_PANIC))
		})
	})

	Describe("Results", func() {
		Context("when there are no results", func() {
			It("returns nil", func() {
				results := r.Results()
				Expect(results).To(BeNil())
			})
		})

		Context("when there are 5 results", func() {
			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					r.Success("foo", "bar", int64(42))
				}
			})

			It("returns all results", func() {
				results := r.Results()
				Expect(results).To(HaveLen(5))
			})
		})
	})
})
