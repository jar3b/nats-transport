package nats_transport

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"net/http"
	"time"
)

// SubjectResolveFunc resolves the NATS subject (based on http.Request struct)
type SubjectResolveFunc func(r *http.Request) string

// ModifyRequestHookFunc is the hook to modify Request struct (protobuf)
type ModifyRequestHookFunc = func(req *http.Request, r *Request)

type NatsTransport struct {
	// NatsConnection connection to NATS server
	NatsConnection *nats.Conn

	// Subject NATS subject to push wrapped HTTP request on
	Subject string

	// SubjectResolver used only if Subject is not specified
	SubjectResolver SubjectResolveFunc

	// ModifyRequestHook calls after Request struct parsed from HTTP and before sent to NATS
	ModifyRequestHook ModifyRequestHookFunc

	// Timeout NATS request timeout
	Timeout time.Duration
}

func (nt NatsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	request := NewRequest()

	if err := request.FromHTTP(r); err != nil {
		return nil, fmt.Errorf("nats_transport: cannot parse HTTP request: %v", err)
	}

	// call modify request hook
	if nt.ModifyRequestHook != nil {
		nt.ModifyRequestHook(r, request)
	}

	// Serialize the request.
	requestBytes, err := proto.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("nats_transport: cannot serialize request")
	}

	// get the outgoing NATS subject
	natsSubject := nt.Subject
	if natsSubject == "" {
		natsSubject = nt.SubjectResolver(r)
	}
	if natsSubject == "" {
		return nil, fmt.Errorf("nats_transport: cannot detect NATS subject")
	}

	// Send Request to NATS server
	msg, err := nt.NatsConnection.Request(
		natsSubject,
		requestBytes,
		nt.Timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("nats_transport: cannot send NATS request")
	}

	// Get Response object from NATS message
	response := NewResponse()
	if err := response.ReadFrom(msg.Data); err != nil {
		return nil, fmt.Errorf("nats_transport: %v", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("nats_transport: %s", response.Error)
	}

	// prepare HTTP response
	httpResponse, err := response.ToHTTPResponse(r)
	if err != nil {
		return nil, fmt.Errorf("nats_transport: %v", err)
	}

	return httpResponse, nil
}
