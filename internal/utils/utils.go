package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	basicServiceV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/service/v1"
)

// CallService simulates an asynchronous service call with random delay.
// Returns a channel that will receive a single response after 0-9 seconds
// and then close. Used for testing fan-out patterns.
func CallService(serviceName string, serviceType string) chan *basicServiceV1.SomeServiceResponse {
	response := make(chan *basicServiceV1.SomeServiceResponse)
	go func() {
		// Simulate variable response time
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)

		srvResp := &basicServiceV1.SomeServiceResponse{
			Id:      uuid.NewString(),
			Name:    serviceName,
			Version: "v0.1.0",
			Data: &basicServiceV1.SomeServiceData{
				Type:  serviceType,
				Value: fmt.Sprintf("Some data from %s", serviceName),
			},
		}
		response <- srvResp

		close(response)
	}()

	return response
}

// MergeServiceResponses implements a fan-in pattern by collecting responses from
// multiple service call channels. Each input channel's responses are grouped
// and sent as a batch on the output channel. The output channel closes when
// all input channels have been processed.
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

	// Close output channel when all responses are processed
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
