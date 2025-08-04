package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
)

func CallService(serviceName string, serviceType string) chan *basicServiceV1.SomeServiceResponse {
	response := make(chan *basicServiceV1.SomeServiceResponse)
	go func() {
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
		srvResp := &basicServiceV1.SomeServiceResponse{Id: uuid.NewString(), Name: serviceName, Version: "v0.1.0", Data: &basicServiceV1.SomeServiceData{Type: serviceType, Value: fmt.Sprintf("Some data from %s", serviceName)}}
		response <- srvResp

		close(response)
	}()

	return response
}

func MergeServiceResponses(responses ...chan *basicServiceV1.SomeServiceResponse) chan *basicServiceV1.SomeServiceResponses {
	var wg sync.WaitGroup
	output := make(chan *basicServiceV1.SomeServiceResponses)
	wg.Add(len(responses))
	for _, response := range responses {
		go func(response <-chan *basicServiceV1.SomeServiceResponse) {
			defer wg.Done()
			srvResponses := &basicServiceV1.SomeServiceResponses{}
			for srvResp := range response {
				srvResponses.Responses = append(srvResponses.Responses, srvResp)
			}
			output <- srvResponses
		}(response)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
