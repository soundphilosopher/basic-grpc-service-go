package utils

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	cloudeventsV1 "github.com/soundphilosopher/basic-grpc-service-go/sdk/io/cloudevents/v1"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateCloudEvent(request any, event *anypb.Any) (*cloudeventsV1.CloudEvent, error) {
	if req, ok := request.(connect.AnyRequest); ok {
		ce := &cloudeventsV1.CloudEvent{
			Id:          uuid.New().String(),
			SpecVersion: "1.0",
			Type:        string(event.MessageName()),
			Source:      req.Header().Get("Host") + req.Spec().Procedure,
			Data: &cloudeventsV1.CloudEvent_ProtoData{
				ProtoData: event,
			},
			Attributes: map[string]*cloudeventsV1.CloudEvent_CloudEventAttributeValue{},
		}

		ce.Attributes["time"] = &cloudeventsV1.CloudEvent_CloudEventAttributeValue{
			Attr: &cloudeventsV1.CloudEvent_CloudEventAttributeValue_CeTimestamp{
				CeTimestamp: timestamppb.Now(),
			},
		}

		return ce, nil
	}

	return nil, fmt.Errorf("cannot convert request to cloudevent. Req: %v, Event: %v", request, event)
}
