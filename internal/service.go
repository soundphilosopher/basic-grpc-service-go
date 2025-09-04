// Package internal implements the BasicService gRPC handlers with Cloud Events support.
package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/soundphilosopher/basic-grpc-service-go/internal/talk"
	"github.com/soundphilosopher/basic-grpc-service-go/internal/utils"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

// BasicServiceV1 implements the gRPC BasicService interface providing
// Hello, Talk, and Background operations with state management capabilities.
type BasicServiceV1 struct {
	StateManager *utils.StateManager
}

// NewBasicServiceV1 creates a new BasicServiceV1 instance with an initialized StateManager.
// The StateManager tracks the lifecycle of background operations.
func NewBasicServiceV1() *BasicServiceV1 {
	return &BasicServiceV1{
		StateManager: utils.NewStateManager(),
	}
}

// Hello handles simple greeting requests and returns a Cloud Event response.
// The greeting message is formatted with the provided input message.
func (s *BasicServiceV1) Hello(ctx context.Context, req *connect.Request[basicServiceV1.HelloRequest]) (*connect.Response[basicServiceV1.HelloResponse], error) {
	event, err := anypb.New(&basicServiceV1.HelloResponseEvent{Greeting: fmt.Sprintf("Hello, %s", req.Msg.Message)})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cloudevent, err := utils.CreateCloudEvent(req, event)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&basicServiceV1.HelloResponse{CloudEvent: cloudevent})
	resp.Header().Set("Basic-Service-Version", "v1")

	return resp, nil
}

// Talk handles bidirectional streaming conversation using the talk module.
// The stream continues until the client closes it or the talk module signals end.
func (s *BasicServiceV1) Talk(ctx context.Context, stream *connect.BidiStream[basicServiceV1.TalkRequest, basicServiceV1.TalkResponse]) error {
	for {
		if err := ctx.Err(); err != nil {
			return connect.NewError(connect.CodeAborted, err)
		}

		receive, err := stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return connect.NewError(connect.CodeCanceled, err)
		}

		reply, end := talk.Reply(receive.Message)
		if err := stream.Send(&basicServiceV1.TalkResponse{Answer: reply}); err != nil {
			return connect.NewError(connect.CodeCanceled, err)
		}
		if end {
			return nil
		}
	}
}

// Background handles long-running operations by orchestrating multiple service calls
// and streaming periodic status updates. Uses fan-out/fan-in pattern to call
// multiple services concurrently and reports progress every 2 seconds.
func (s *BasicServiceV1) Background(ctx context.Context, req *connect.Request[basicServiceV1.BackgroundRequest], stream *connect.ServerStream[basicServiceV1.BackgroundResponse]) error {
	hash := uuid.NewString()
	state, _, _ := s.StateManager.GetState(hash)

	data := []*basicServiceV1.SomeServiceResponse{}

	// Start background processing if not already running
	if state == nil {
		s.StateManager.Start(hash)
		go func() {
			// Fan-out: call multiple services concurrently
			s1 := utils.CallService("service-1", "rest")
			s2 := utils.CallService("service-2", "rpc")
			s3 := utils.CallService("service-3", "grpc")
			s4 := utils.CallService("service-4", "rest")
			s5 := utils.CallService("service-5", "grpc")

			// Fan-in: collect responses as they arrive
			for response := range utils.MergeServiceResponses(s1, s2, s3, s4, s5) {
				log.Printf("Received response: %v", response)
				data = append(data, response.Responses...)
			}

			s.StateManager.Finish(hash)
		}()
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Stream status updates until processing completes
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			current_state, start, finish := s.StateManager.GetState(hash)

			// Send final response when processing is complete
			if *current_state != basicServiceV1.State_STATE_PROCESS {
				event, err := anypb.New(&basicServiceV1.BackgroundResponseEvent{State: *current_state, StartedAt: start, CompletedAt: finish, Responses: data})
				if err != nil {
					return connect.NewError(connect.CodeInternal, err)
				}

				cloudevent, err := utils.CreateCloudEvent(req, event)
				if err != nil {
					return connect.NewError(connect.CodeInternal, err)
				}

				if err := stream.Send(&basicServiceV1.BackgroundResponse{CloudEvent: cloudevent}); err != nil {
					return connect.NewError(connect.CodeCanceled, err)
				}

				return nil
			}

			// Send progress update
			event, err := anypb.New(&basicServiceV1.BackgroundResponseEvent{State: *current_state, StartedAt: start, CompletedAt: finish, Responses: data})
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}

			cloudevent, err := utils.CreateCloudEvent(req, event)
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}

			if err := stream.Send(&basicServiceV1.BackgroundResponse{CloudEvent: cloudevent}); err != nil {
				return connect.NewError(connect.CodeCanceled, err)
			}
		}
	}
}
