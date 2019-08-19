package manager

import "errors"

var ErrInvalidAgent = errors.New("invalid agent")

type Registration struct {
	Id string
}

type Registry struct {
	agents     []Registration
	active     map[string]struct{}
	quarantine map[string]struct{}
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
	r.agents = append(r.agents, regis)
	r.active[regis.Id] = struct{}{}
	return nil
}

func (r *Registry) Quarantine(id string) error {
	if _, ok := r.active[id]; ok {
		delete(r.active, id)
		r.quarantine[id] = struct{}{}
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
	for _, regis := range r.agents {
		result = append(result, regis)
	}

	return result
}

func (r *Registry) GetActive() []Registration {
	var result []Registration

	if len(r.active) == 0 {
		return result
	}

	result = make([]Registration, 0, len(r.active))
	for _, regis := range r.agents {
		if _, ok := r.active[regis.Id]; ok {
			result = append(result, regis)
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
	for _, regis := range r.agents {
		if _, ok := r.quarantine[regis.Id]; ok {
			result = append(result, regis)
		}
	}

	return result
}
