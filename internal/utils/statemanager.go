// Package utils provides utility functions and state management for the basic service.
package utils

import (
	"sync"

	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StateManager tracks the lifecycle of background operations using unique hash identifiers.
// It maintains state, timestamps, and errors for concurrent operations in a thread-safe manner.
type StateManager struct {
	mu       sync.Mutex
	state    map[string]*basicServiceV1.State
	start    map[string]*timestamppb.Timestamp
	complete map[string]*timestamppb.Timestamp
	errors   map[string]*[]error
}

// NewStateManager creates a new StateManager with initialized internal maps.
func NewStateManager() *StateManager {
	return &StateManager{
		state:    make(map[string]*basicServiceV1.State),
		start:    make(map[string]*timestamppb.Timestamp),
		complete: make(map[string]*timestamppb.Timestamp),
		errors:   make(map[string]*[]error),
	}
}

// Start marks the beginning of an operation by setting its state to processing
// and recording the start timestamp.
func (m *StateManager) Start(hash string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state := basicServiceV1.State_STATE_PROCESS
	m.state[hash] = &state
	m.start[hash] = timestamppb.Now()
}

// Finish completes an operation by setting the final state based on error conditions
// and recording the completion timestamp. Operations with errors are marked as
// STATE_COMPLETE_WITH_ERROR, otherwise STATE_COMPLETE.
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

// GetState returns the current state, start time, and completion time for the given hash.
// Returns nil values for times that haven't been set yet.
func (m *StateManager) GetState(hash string) (*basicServiceV1.State, *timestamppb.Timestamp, *timestamppb.Timestamp) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state[hash], m.start[hash], m.complete[hash]
}

// SetError adds an error to the operation's error list. If err is nil, no action is taken.
func (m *StateManager) SetError(hash string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		if _, exists := m.errors[hash]; !exists {
			m.errors[hash] = &[]error{}
		}
		*m.errors[hash] = append(*m.errors[hash], err)
	}
}

// HasErrors returns true if the operation has recorded any errors.
func (m *StateManager) HasErrors(hash string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	errors, exists := m.errors[hash]
	if !exists || errors == nil {
		return false
	}
	return len(*errors) > 0
}

// GetErrors returns all errors recorded for the operation, or an empty slice if none exist.
func (m *StateManager) GetErrors(hash string) []error {
	m.mu.Lock()
	defer m.mu.Unlock()

	errors, exists := m.errors[hash]
	if !exists || errors == nil {
		return []error{}
	}
	return *errors
}
