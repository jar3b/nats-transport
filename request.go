package nats_transport

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (r *Request) FromHTTP(req *http.Request) error {
	if req == nil {
		return errors.New("nats_transport: request cannot be nil")
	}

	defer req.Body.Close()
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("nats_transport: cannot read request body")
	}
	// bufReader := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Proto = req.Proto
	r.Scheme = req.URL.Scheme
	r.Host = req.Host
	r.URL = req.URL.String()
	r.Method = req.Method
	r.Header = copyMap(req.Header)
	r.RemoteAddr = req.RemoteAddr
	r.Body = buf

	return nil
}

func NewRequest() *Request {
	return &Request{
		Header: make(map[string]*Values),
		Body:   make([]byte, 0),
	}
}
