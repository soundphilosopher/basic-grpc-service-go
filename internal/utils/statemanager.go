package utils

import (
	"sync"

	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StateManager struct {
	mu       sync.Mutex
	state    map[string]*basicServiceV1.State
	start    map[string]*timestamppb.Timestamp
	complete map[string]*timestamppb.Timestamp
	errors   map[string]*[]error
}

func NewStateManager() *StateManager {
	return &StateManager{
		state:    make(map[string]*basicServiceV1.State),
		start:    make(map[string]*timestamppb.Timestamp),
		complete: make(map[string]*timestamppb.Timestamp),
		errors:   make(map[string]*[]error),
	}
}

func (m *StateManager) Start(hash string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state := basicServiceV1.State_STATE_PROCESS
	m.state[hash] = &state
	m.start[hash] = timestamppb.Now()
}

func (m *StateManager) Finish(hash string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var state basicServiceV1.State
	errors, exists := m.errors[hash]
	if exists && errors != nil && len(*errors) > 0 {
		state = basicServiceV1.State_STATE_COMPLETE_WITH_ERROR
	} else {
		state = basicServiceV1.State_STATE_COMPLETE
	}

	m.state[hash] = &state
	m.complete[hash] = timestamppb.Now()
}

func (m *StateManager) GetState(hash string) (*basicServiceV1.State, *timestamppb.Timestamp, *timestamppb.Timestamp) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state[hash], m.start[hash], m.complete[hash]
}

func (m *StateManager) SetError(hash string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		if _, exists := m.errors[hash]; !exists {
			m.errors[hash] = &[]error{}
		}
		*m.errors[hash] = append(*m.errors[hash], err)
		return
	}
}

func (m *StateManager) HasErrors(hash string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	errors, exists := m.errors[hash]
	if !exists || errors == nil {
		return false
	}
	return len(*errors) > 0
}

func (m *StateManager) GetErrors(hash string) []error {
	m.mu.Lock()
	defer m.mu.Unlock()

	errors, exists := m.errors[hash]
	if !exists || errors == nil {
		return []error{}
	}
	return *errors
}
