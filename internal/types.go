package internal

import "github.com/soundphilosopher/basic-grpc-service-go/internal/utils"

type BasicServiceV1 struct {
	StateManager *utils.StateManager
}

func NewBasicServiceV1() *BasicServiceV1 {
	return &BasicServiceV1{
		StateManager: utils.NewStateManager(),
	}
}
