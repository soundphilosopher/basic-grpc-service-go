package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
	"github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/v1/basicV1connect"
)

func main() {
	client := basicV1connect.NewBasicServiceClient(http.DefaultClient, "https://127.0.0.1:8999", connect.WithGRPC())

	resp, err := client.Hello(context.Background(), connect.NewRequest(&basicServiceV1.HelloRequest{Message: "You"}))
	if err != nil {
		log.Fatalf("error calling Hello: %v\n", err)
		return
	}

	log.Printf("Response: %+v\n", resp.Msg)
}
