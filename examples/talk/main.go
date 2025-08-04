package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/v1/basicV1connect"
)

func main() {
	client := basicV1connect.NewBasicServiceClient(http.DefaultClient, "https://127.0.0.1:8999", connect.WithGRPC())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	stream := client.Talk(ctx)

	go func() {
		requests := []*basicServiceV1.TalkRequest{
			{Message: "Hello"},
			{Message: "How are you?"},
			{Message: "Bye"},
		}

		for _, req := range requests {
			if err := stream.Send(req); err != nil {
				log.Fatalf("send error: %v", err)
			}
			time.Sleep(time.Second)
		}

		stream.CloseRequest()
	}()

	for {
		resp, err := stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("receive error: %v", err)
		}

		log.Printf("Received: %s\n", resp.Answer)
	}
}
