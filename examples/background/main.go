package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/v1/basicV1connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	client := basicV1connect.NewBasicServiceClient(http.DefaultClient, "https://127.0.0.1:8999", connect.WithGRPC())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	stream, err := client.Background(ctx, connect.NewRequest(&basicServiceV1.BackgroundRequest{Processes: 4}))
	if err != nil {
		log.Fatalf("error calling Background: %v\n", err)
	}

	for stream.Receive() {
		response := stream.Msg()

		data := &basicServiceV1.BackgroundResponseEvent{}
		if err := proto.Unmarshal(response.CloudEvent.GetProtoData().Value, data); err != nil {
			log.Fatalf("error unmarshalling response: %v\n", err)
		}

		j, err := protojson.Marshal(data)
		if err != nil {
			log.Fatalf("error marshaling response: %v\n", err)
		}

		log.Printf("Received: %s\n", string(j))
	}

	if err := stream.Err(); err != nil {
		log.Fatalf("receive error: %v\n", err)
	}
}
