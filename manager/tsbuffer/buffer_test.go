package tsbuffer

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chmking/horde/protobuf/private"
)

var MockInitMillisecond = func() int64 {
	return int64(42000)
}

var MockNowMillisecond = func() int64 {
	return int64(44000)
}

var _ = Describe("New", func() {
	It("constructs a new Buffer", func() {
		result := New(0)
		Expect(result).ToNot(BeNil())
	})

	It("sets the buffer window", func() {
		result := New(time.Millisecond * 100)
		Expect(result.window).To(BeNumerically("==", 100))
	})
})

var _ = Describe("Buffer", func() {
	var buf *Buffer

	BeforeEach(func() {
		buf = newInjected(time.Millisecond*100, MockInitMillisecond)
	})

	Describe("Len", func() {
		Context("when there are no items", func() {
			It("returns 0", func() {
				value := buf.Len()
				Expect(value).To(BeNumerically("==", 0))
			})
		})

		Context("when there is an item", func() {
			BeforeEach(func() {
				result := &private.Result{Millisecond: MockInitMillisecond()}
				buf.Add(result)
			})

			It("returns 1", func() {
				value := buf.Len()
				Expect(value).To(BeNumerically("==", 1))
			})
		})
	})

	Describe("Add", func() {
		var result *private.Result

		BeforeEach(func() {
			result = &private.Result{}
		})

		Context("when the item is in range", func() {
			BeforeEach(func() {
				result.Millisecond = MockInitMillisecond()
			})

			It("adds the item", func() {
				Expect(buf.results).To(HaveLen(0))
				buf.Add(result)
				Expect(buf.results).To(HaveLen(1))
			})
		})

		Context("when the item is out of range", func() {
			BeforeEach(func() {
				result.Millisecond = MockInitMillisecond() - 1000
			})

			It("drops the item", func() {
				Expect(buf.results).To(HaveLen(0))
				buf.Add(result)
				Expect(buf.results).To(HaveLen(0))
			})
		})

		Context("when items share same timestamp", func() {
			BeforeEach(func() {
				result.Millisecond = MockInitMillisecond()
			})

			It("only creates a single index", func() {
				Expect(buf.results).To(HaveLen(0))
				for i := 0; i < 5; i++ {
					buf.Add(result)
				}
				Expect(buf.results).To(HaveLen(1))
			})

			It("appends items to the same index", func() {
				for i := 0; i < 5; i++ {
					buf.Add(result)
				}
				values := buf.results[result.Millisecond]
				Expect(values).To(HaveLen(5))
			})
		})
	})

	Describe("Collect", func() {
		Context("when there are no results", func() {
			It("returns nil results", func() {
				results := buf.Collect()
				Expect(results).To(BeNil())
			})
		})

		Context("when there are results to return", func() {
			var expectedReturn map[int64][]*private.Result

			BeforeEach(func() {
				expectedReturn = make(map[int64][]*private.Result, 0)
				// Currently the pointer is at MockInitMillisecond - 100.
				// Write expected results starting at MockInitMillisecond.
				start := MockInitMillisecond()
				for i := 0; i < 5; i++ {
					value := &private.Result{Millisecond: start + int64(i)}
					expectedReturn[value.Millisecond] = []*private.Result{value}
					buf.results[value.Millisecond] = []*private.Result{value}
				}

				// Assing fast-forwarded time 2000 millisecond in future.
				buf.millisecond = MockNowMillisecond
			})

			It("returns results older than buffer", func() {
				results := buf.Collect()
				Expect(results).To(Equal(expectedReturn))
			})

			It("deletes the returned results", func() {
				buf.Collect()
				Expect(buf.results).To(HaveLen(0))
			})

			Context("when there are buffering results", func() {
				var remaining map[int64][]*private.Result

				BeforeEach(func() {
					remaining = make(map[int64][]*private.Result, 0)
					// Destination pointed it MockNowMillisecond - 100.
					// Results written there or after should remain.
					start := MockNowMillisecond() - 100
					for i := 0; i < 5; i++ {
						value := &private.Result{Millisecond: start + int64(i)}
						remaining[value.Millisecond] = []*private.Result{value}
						buf.results[value.Millisecond] = []*private.Result{value}
					}
				})

				It("returns results older than buffer", func() {
					results := buf.Collect()
					Expect(results).To(Equal(expectedReturn))
				})

				It("leaves the buffering results", func() {
					buf.Collect()
					Expect(buf.results).To(HaveLen(5))
				})
			})
		})
	})
})
