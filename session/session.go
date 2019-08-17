package session

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chmking/horde"
)

type Callback func()

type ScaleOrder struct {
	Count int32
	Rate  float64
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
	count  int32
	cancel context.CancelFunc

	workers []context.CancelFunc
	mtx     sync.Mutex
}

func (s *Session) Count() int32 {
	return atomic.LoadInt32(&s.count)
}

func (s *Session) Scale(order ScaleOrder, cb Callback) {
	order.callback = cb

	if s.cancel != nil {
		s.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	log.Printf("Scaling with ScaleOrder: %+v", order)

	go s.doScale(ctx, order)
}

func (s *Session) doScale(ctx context.Context, order ScaleOrder) {

	current := s.Count()

	// Scale Down
	if order.Count < current {
		log.Println("Scaling Down")

		diff := current - order.Count

		s.mtx.Lock()

		// Cancel work
		for i := int32(len(s.workers)) - diff; i < int32(len(s.workers)); i++ {
			if cancel := s.workers[i]; cancel != nil {
				cancel()
			}
			s.workers[i] = nil
		}

		// Resize workers
		s.workers = s.workers[:int32(len(s.workers))-diff]

		// Update count
		current = int32(len(s.workers))
		atomic.StoreInt32(&s.count, current)

		s.mtx.Unlock()
	}

	if order.Count > current {
		log.Println("Scaling Up")

		// Wait is used to stagger scaling across agents
		<-time.After(time.Duration(order.Wait) * time.Millisecond)

		// Start workers
		for i := s.Count(); i < order.Count; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				ctx, cancel := context.WithCancel(context.Background())

				// Append worker handle
				s.mtx.Lock()
				s.workers = append(s.workers, cancel)
				s.mtx.Unlock()

				// Start worker
				go s.doWork(ctx, order.Work)

				// Wait for rate limit
				limit := time.Duration(float64(time.Second.Nanoseconds()) * order.Rate)
				<-time.After(limit)
			}
		}
	}

	if order.callback != nil {
		order.callback()
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
	defer cb()
}
