package manager

import (
	"context"
	"errors"
	"time"

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
}

type Metadata struct {
	Registration Registration
	Failed       int
}

func NewRegistry() *Registry {
	return &Registry{
		active:     map[string]struct{}{},
		quarantine: map[string]struct{}{},
	}
}

func (r *Registry) Len() int {
	return len(r.agents)
}

func (r *Registry) Add(regis Registration) error {
	_, isActive := r.active[regis.Id]
	_, isQuarantined := r.quarantine[regis.Id]

	if isActive || isQuarantined {
		return nil
	}

	r.agents = append(r.agents, &Metadata{Registration: regis})
	r.active[regis.Id] = struct{}{}
	return nil
}

func (r *Registry) Quarantine(id string) error {
	if _, ok := r.active[id]; ok {
		delete(r.active, id)
		r.quarantine[id] = struct{}{}

		for _, metadata := range r.agents {
			if metadata.Registration.Id == id {
				metadata.Failed = 3
			}
		}

		return nil
	}

	if _, ok := r.quarantine[id]; !ok {
		return ErrInvalidAgent
	}

	return nil
}

func (r *Registry) GetAll() []Registration {
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
	for _, metadata := range r.agents {
		client := metadata.Registration.Client
		_, err := client.Heartbeat(context.Background(), &private.HeartbeatRequest{})
		if err != nil {
			metadata.Failed = min(metadata.Failed+1, 3)
			if metadata.Failed == 3 {
				_, isActive := r.active[metadata.Registration.Id]
				if isActive {
					delete(r.active, metadata.Registration.Id)
					r.quarantine[metadata.Registration.Id] = struct{}{}
				}
			}
		} else {
			metadata.Failed = max(metadata.Failed-1, 0)
			if metadata.Failed == 0 {
				_, isQuarantined := r.quarantine[metadata.Registration.Id]
				if isQuarantined {
					delete(r.quarantine, metadata.Registration.Id)
					r.active[metadata.Registration.Id] = struct{}{}
				}
			}
		}
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

			<-time.After(time.Second * 30)
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
