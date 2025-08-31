package utils_test

import (
	"errors"
	"testing"
	"time"

	"github.com/soundphilosopher/basic-grpc-service-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateState(t *testing.T) {
	t.Parallel()
	t.Run("Initial state management", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		state, start, complete := sm.GetState(hash)
		assert.Nil(t, state)
		assert.Nil(t, start)
		assert.Nil(t, complete)
	})

	t.Run("should set initial state for starting state", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)

		state, start, complete := sm.GetState(hash)
		assert.NotNil(t, state)
		assert.NotNil(t, start)
		assert.Nil(t, complete)

		assert.Equal(t, "STATE_PROCESS", state.String())
		assert.IsType(t, time.Time{}, start.AsTime())
	})

	t.Run("should set complete when finished", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.Finish(hash)
		state, start, complete := sm.GetState(hash)
		assert.NotNil(t, state)
		assert.NotNil(t, start)
		assert.NotNil(t, complete)

		assert.Equal(t, "STATE_COMPLETE", state.String())
		assert.IsType(t, time.Time{}, start.AsTime())
		assert.IsType(t, time.Time{}, complete.AsTime())
	})

	t.Run("should set some errors", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.SetError(hash, errors.New("test error"))
		state, start, complete := sm.GetState(hash)
		errors := sm.GetErrors(hash)
		assert.NotNil(t, state)
		assert.NotNil(t, start)
		assert.Nil(t, complete)
		assert.NotNil(t, errors)

		assert.Equal(t, "STATE_PROCESS", state.String())
		assert.IsType(t, time.Time{}, start.AsTime())
		assert.Nil(t, complete)
		assert.Len(t, errors, 1)
	})

	t.Run("should return an complete with errors state when finished with errors", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.SetError(hash, errors.New("test error"))
		errors := sm.GetErrors(hash)
		sm.Finish(hash)
		state, start, complete := sm.GetState(hash)
		assert.NotNil(t, state)
		assert.NotNil(t, start)
		assert.NotNil(t, complete)
		assert.NotNil(t, errors)

		assert.Equal(t, "STATE_COMPLETE_WITH_ERROR", state.String())
		assert.IsType(t, time.Time{}, start.AsTime())
		assert.IsType(t, time.Time{}, complete.AsTime())
		assert.Len(t, errors, 1)
	})

	t.Run("should return false when no error is set to the state", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		assert.False(t, sm.HasErrors(hash))
	})

	t.Run("should return false when state caches an nil error", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.SetError(hash, nil)
		assert.False(t, sm.HasErrors(hash))
	})

	t.Run("should return true when state caches errors", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.SetError(hash, errors.New("test error"))
		assert.True(t, sm.HasErrors(hash))
	})

	t.Run("should return empty error set when no error is catched by the state", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		errors := sm.GetErrors(hash)
		assert.Empty(t, errors)
		assert.Len(t, errors, 0)
	})

	t.Run("should return errors if any are catched by the state", func(t *testing.T) {
		sm := utils.NewStateManager()
		hash := "test_hash"

		sm.Start(hash)
		sm.SetError(hash, errors.New("test error 1"))
		sm.SetError(hash, errors.New("test error 2"))
		errors := sm.GetErrors(hash)
		assert.NotEmpty(t, errors)
		assert.Len(t, errors, 2)
	})
}
