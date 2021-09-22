package nats_transport

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
)

// NewResponse creates blank initialized Response object.
func NewResponse() *Response {
	return &Response{
		StatusCode: int32(200),
		Header:     make(map[string]*Values, 0),
		Body:       make([]byte, 0),
	}
}

func (resp *Response) ReadFrom(responseData []byte) error {
	if responseData == nil || len(responseData) == 0 {
		return errors.New("response content is empty")
	}
	if err := proto.Unmarshal(responseData, resp); err != nil {
		return err
	}
	return nil
}

func (resp *Response) ToHTTPResponse(r *http.Request) (*http.Response, error) {
	httpResponse := http.Response{
		Status:        fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(int(resp.StatusCode))),
		StatusCode:    int(resp.StatusCode),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBuffer(resp.Body)),
		ContentLength: int64(len(resp.Body)),
		Request:       r,
		Header:        make(http.Header, 0),
	}

	// copy headers
	for headerName, headerValues := range resp.Header {
		for _, headerValue := range headerValues.Arr {
			httpResponse.Header.Add(headerName, headerValue)
		}
	}

	return &httpResponse, nil
}
