package registry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chmking/horde/logger/log"
	"github.com/chmking/horde/protobuf/private"
)

var ErrInvalidAgent = errors.New("invalid agent")

type Registration struct {
	Id     string
	Client private.AgentClient
}

type Registry struct {
	agents     []*Metadata
	active     map[string]struct{}
	quarantine map[string]struct{}
	mtx        sync.Mutex
	cb         func()
}

type Metadata struct {
	Registration Registration
	Failed       int
}

func New() *Registry {
	return &Registry{
		active:     map[string]struct{}{},
		quarantine: map[string]struct{}{},
	}
}

func (r *Registry) RegisterCallback(cb func()) {
	r.cb = cb
}

func (r *Registry) Len() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return len(r.agents)
}

func (r *Registry) Add(regis Registration) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	_, isActive := r.active[regis.Id]
	_, isQuarantined := r.quarantine[regis.Id]

	if isActive || isQuarantined {
		return nil
	}

	r.agents = append(r.agents, &Metadata{Registration: regis})
	r.active[regis.Id] = struct{}{}

	if r.cb != nil {
		go r.cb()
	}

	return nil
}

func (r *Registry) Quarantine(id string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.active[id]; ok {
		delete(r.active, id)
		r.quarantine[id] = struct{}{}

		for _, metadata := range r.agents {
			if metadata.Registration.Id == id {
				metadata.Failed = 3
				break
			}
		}

		if r.cb != nil {
			go r.cb()
		}

		return nil
	}

	if _, ok := r.quarantine[id]; !ok {
		return ErrInvalidAgent
	}

	return nil
}

func (r *Registry) GetAll() []Registration {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []Registration

	if len(r.agents) == 0 {
		return result
	}

	result = make([]Registration, 0, len(r.agents))
	for _, metadata := range r.agents {
		result = append(result, metadata.Registration)
	}

	return result
}

func (r *Registry) GetActive() []Registration {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []Registration

	if len(r.active) == 0 {
		return result
	}

	result = make([]Registration, 0, len(r.active))
	for _, metadata := range r.agents {
		if _, ok := r.active[metadata.Registration.Id]; ok {
			result = append(result, metadata.Registration)
		}
	}

	return result
}

func (r *Registry) GetQuarantined() []Registration {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []Registration

	if len(r.quarantine) == 0 {
		return result
	}

	result = make([]Registration, 0, len(r.quarantine))
	for _, metadata := range r.agents {
		if _, ok := r.quarantine[metadata.Registration.Id]; ok {
			result = append(result, metadata.Registration)
		}
	}

	return result
}

func (r *Registry) Healthcheck() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	adjusted := false
	for _, metadata := range r.agents {
		client := metadata.Registration.Client
		_, err := client.Healthcheck(context.Background(), &private.HealthcheckRequest{})
		if err != nil {
			metadata.Failed = min(metadata.Failed+1, 3)
			if metadata.Failed == 3 {
				_, isActive := r.active[metadata.Registration.Id]
				if isActive {
					log.Info().Msg("Quarantining unhealthy client")
					delete(r.active, metadata.Registration.Id)
					r.quarantine[metadata.Registration.Id] = struct{}{}
					adjusted = true
				}
			}
		} else {
			metadata.Failed = max(metadata.Failed-1, 0)
			if metadata.Failed == 0 {
				_, isQuarantined := r.quarantine[metadata.Registration.Id]
				if isQuarantined {
					log.Info().Msg("Activating healthy client")
					delete(r.quarantine, metadata.Registration.Id)
					r.active[metadata.Registration.Id] = struct{}{}
					adjusted = true
				}
			}
		}
	}

	if adjusted && r.cb != nil {
		go r.cb()
	}
}

func (r *Registry) BeginHealthcheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				r.Healthcheck()
			}

			<-time.After(time.Second * 10)
		}
	}()
}

func min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

func max(lhs, rhs int) int {
	if lhs > rhs {
		return lhs
	}
	return rhs
}
