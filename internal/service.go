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

func (s *BasicServiceV1) Background(ctx context.Context, req *connect.Request[basicServiceV1.BackgroundRequest], stream *connect.ServerStream[basicServiceV1.BackgroundResponse]) error {
	hash := uuid.NewString()
	state, _, _, _ := s.StateManager.GetState(hash)

	data := []*basicServiceV1.SomeServiceResponse{}

	if state == nil {
		s.StateManager.Start(hash, basicServiceV1.State_PROCESS_STATE_PROCESS)
		go func() {
			// fan-out
			s1 := utils.CallService("service-1", "rest")
			s2 := utils.CallService("service-2", "rpc")
			s3 := utils.CallService("service-3", "grpc")
			s4 := utils.CallService("service-4", "rest")
			s5 := utils.CallService("service-5", "grpc")

			// fan-in
			for response := range utils.MergeServiceResponses(s1, s2, s3, s4, s5) {
				log.Printf("Received response: %v", response)
				data = append(data, response.Responses...)
			}

			s.StateManager.Finish(hash, basicServiceV1.State_PROCESS_STATE_COMPLETE)
		}()
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			current_state, start, finish, _ := s.StateManager.GetState(hash)

			if *current_state != basicServiceV1.State_PROCESS_STATE_PROCESS {
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
