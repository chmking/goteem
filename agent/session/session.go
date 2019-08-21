package session

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chmking/horde"
)

type Callback func()

type ScaleOrder struct {
	Count int32
	Rate  int64
	Wait  int64
	Work  Work

	callback Callback
}

type Work struct {
	Tasks   []*horde.Task
	WaitMin int64
	WaitMax int64

	weightSum int
}

type Session struct {
	cancel context.CancelFunc

	workers []context.CancelFunc
	mtx     sync.Mutex
}

func (s *Session) Count() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return len(s.workers)
}

func (s *Session) Scale(ctx context.Context, order ScaleOrder, cb Callback) {
	if s.cancel != nil {
		s.cancel()
	}

	order.callback = cb

	log.Printf("Scaling with ScaleOrder: %+v", order)

	go s.doScale(ctx, order)
}

func (s *Session) doScale(ctx context.Context, order ScaleOrder) {
	s.scaleDown(ctx, order)

	// Wait is used to stagger scaling across agents
	<-time.After(time.Duration(order.Wait) * time.Millisecond)

	s.scaleUp(ctx, order)

	if order.callback != nil {
		order.callback()
	}
}

func (s *Session) scaleDown(ctx context.Context, order ScaleOrder) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	current := len(s.workers)

	if int(order.Count) >= current {
		return
	}

	diff := current - int(order.Count)

	// Cancel work on workers
	for i := current - diff; i < current; i++ {
		if cancel := s.workers[i]; cancel != nil {
			cancel()
		}
		s.workers[i] = nil
	}

	// Resize workers
	s.workers = s.workers[:current-diff]
}

func (s *Session) scaleUp(ctx context.Context, order ScaleOrder) {
	for {
		s.mtx.Lock()

		if len(s.workers) >= int(order.Count) {
			s.mtx.Unlock()
			return
		}

		select {
		case <-ctx.Done():
			s.mtx.Unlock()
			return
		default:
			workerCtx, cancel := context.WithCancel(ctx)

			// Append worker handle
			s.workers = append(s.workers, cancel)
			s.mtx.Unlock()

			// Start worker
			log.Println("Spawning worker")
			go s.doWork(workerCtx, order.Work)

			// Wait for rate limit
			limit := time.Duration(order.Rate)
			<-time.After(limit)
		}
	}
}

func (s *Session) doWork(ctx context.Context, work Work) {
	if len(work.Tasks) == 0 {
		panic("empty task list provided to worker")
	}

	for _, task := range work.Tasks {
		work.weightSum = work.weightSum + task.Weight
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var selected *horde.Task

			if work.weightSum == 0 {
				// Tasks were not given weights. Randomly select
				// amogst the tasks provided.
				selected = work.Tasks[rand.Int63n(int64(len(work.Tasks)))]
			} else {
				// At least one task was weighted. Randomly select
				// a weight to determine task.
				index := rand.Int63n(int64(work.weightSum))
				var taskSum int
				for _, task := range work.Tasks {
					taskSum = taskSum + task.Weight
					if index < int64(taskSum) {
						selected = task
						break
					}
				}
			}

			selected.Func.Exec(ctx)

			if work.WaitMax >= 0 {
				min := maxInt64(0, work.WaitMin)
				min = minInt64(min, work.WaitMax)

				var wait time.Duration
				if min == work.WaitMax {
					wait = time.Millisecond * time.Duration(min)
				} else {
					val := rand.Int63n(work.WaitMax-min) + min
					wait = time.Millisecond * time.Duration(val)
				}

				<-time.After(wait)
			}
		}
	}
}

func maxInt64(lhs, rhs int64) int64 {
	if lhs > rhs {
		return lhs
	}

	return rhs
}

func minInt64(lhs, rhs int64) int64 {
	if lhs < rhs {
		return lhs
	}

	return rhs
}

func (s *Session) Stop(cb Callback) {
	s.mtx.Lock()
	for _, cancel := range s.workers {
		cancel()
	}
	s.mtx.Unlock()

	defer cb()
}
