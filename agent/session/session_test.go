package session

import (
	"context"

	"github.com/chmking/horde"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var (
		session *Session
		order   ScaleOrder
		done    chan struct{}
		ctx     context.Context
	)

	BeforeEach(func() {
		session = &Session{}
		order = ScaleOrder{}
		done = make(chan struct{}, 1)
		ctx = context.Background()
	})

	Describe("Scale", func() {
		It("calls the Callback when scaled", func() {
			called := false

			session.Scale(ctx, order, func() {
				called = true
				close(done)
			})

			<-done

			Expect(called).To(BeTrue())
		})
	})

	Describe("Stop", func() {
		It("calls the Callback when stopped", func() {
			called := false
			session.Stop(func() { called = true })
			Expect(called).To(BeTrue())
		})
	})

	Describe("scaleDown", func() {
		Context("when the target is greater than current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 2
			})

			It("does not scale the workers", func() {
				session.scaleDown(ctx, order)
				Expect(session.workers).To(HaveLen(1))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleDown(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})

		Context("when the target is equal to current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 1
			})

			It("does not scale the workers", func() {
				session.scaleDown(ctx, order)
				Expect(session.workers).To(HaveLen(1))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleDown(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})

		Context("when the target is less than current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 0
			})

			It("scales workers", func() {
				session.scaleDown(ctx, order)
				Expect(session.workers).To(HaveLen(0))
			})

			It("cancels the work", func(done Done) {
				ctx, cancel := context.WithCancel(context.Background())
				session.workers[0] = cancel

				session.scaleDown(ctx, order)
				_, ok := <-ctx.Done()
				Expect(ok).To(BeFalse())
				close(done)
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleDown(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})
	})

	Describe("scaleUp", func() {
		Context("when the target is less than current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 0
			})

			It("does not scale workers", func() {
				session.scaleUp(ctx, order)
				Expect(session.workers).To(HaveLen(1))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleUp(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})

		Context("when the target is equal to current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 1
			})

			It("does not scale workers", func() {
				session.scaleUp(ctx, order)
				Expect(session.workers).To(HaveLen(1))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleUp(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})

		Context("when the target is greater than current", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 2
				order.Work = Work{
					Tasks: []*horde.Task{{Func: func(ctx context.Context) {}}},
				}
			})

			It("scales workers", func() {
				session.scaleUp(ctx, order)
				Expect(session.workers).To(HaveLen(2))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleUp(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})

		Context("when the scale is cancelled", func() {
			BeforeEach(func() {
				session.workers = make([]context.CancelFunc, 1, 1)
				order.Count = 2

				cancelled, cancel := context.WithCancel(ctx)
				cancel()

				ctx = cancelled
			})

			It("does not scale workers", func() {
				session.scaleUp(ctx, order)
				Expect(session.workers).To(HaveLen(1))
			})

			It("unlocks the mutex", func(done Done) {
				session.scaleUp(ctx, order)
				session.mtx.Lock()
				close(done)
			})
		})
	})
})
