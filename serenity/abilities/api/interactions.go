package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nchursin/serenity-go/serenity/core"
)

// sendRequest is an interaction that sends an HTTP request
type sendRequest struct {
	request *http.Request
}

// a creates a new SendRequest interaction
func a(req *http.Request) core.Activity {
	return &sendRequest{request: req}
}

// Description returns the interaction description
func (s *sendRequest) Description() string {
	if s.request == nil {
		return "#actor sends an HTTP request"
	}

	return fmt.Sprintf("#actor sends %s request to %s", s.request.Method, s.request.URL.String())
}

// PerformAs executes the send request interaction
func (s *sendRequest) PerformAs(actor core.Actor) error {
	if s.request == nil {
		return fmt.Errorf("request is nil")
	}

	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)

	_, err = callAbility.SendRequest(s.request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	return nil
}

// RequestBuilder helps build HTTP requests with fluent interface
type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(method, url string) *RequestBuilder {
	return &RequestBuilder{
		method:  method,
		url:     url,
		headers: make(map[string]string),
	}
}

// WithHeader adds a header to the request
func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

// WithHeaders adds multiple headers to the request
func (rb *RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	for k, v := range headers {
		rb.headers[k] = v
	}
	return rb
}

// WithBody sets the request body
func (rb *RequestBuilder) WithBody(body io.Reader) *RequestBuilder {
	rb.body = body
	return rb
}

// WithJSONBody sets a JSON body by marshaling the provided data
func (rb *RequestBuilder) WithJSONBody(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON body: %w", err)
	}

	rb.body = bytes.NewBuffer(jsonData)
	if rb.headers == nil {
		rb.headers = make(map[string]string)
	}
	rb.headers["Content-Type"] = "application/json"

	return nil
}

// With sets the request body (convenience method for interface{} values)
func (rb *RequestBuilder) With(data interface{}) *RequestBuilder {
	if data == nil {
		return rb
	}

	switch v := data.(type) {
	case io.Reader:
		rb.body = v
	case []byte:
		rb.body = bytes.NewBuffer(v)
	case string:
		rb.body = bytes.NewBufferString(v)
	default:
		// Try to marshal as JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			// Fall back to string representation
			rb.body = bytes.NewBufferString(fmt.Sprintf("%v", data))
		} else {
			rb.body = bytes.NewBuffer(jsonData)
			if rb.headers == nil {
				rb.headers = make(map[string]string)
			}
			rb.headers["Content-Type"] = "application/json"
		}
	}

	return rb
}

// Build creates the HTTP request
func (rb *RequestBuilder) Build() (*http.Request, error) {
	req, err := http.NewRequest(rb.method, rb.url, rb.body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// Convenience functions for different HTTP methods

// Get creates a GET request builder
func Get(url string) *RequestBuilder {
	return NewRequestBuilder(http.MethodGet, url)
}

// Post creates a POST request builder
func Post(url string) *RequestBuilder {
	return NewRequestBuilder(http.MethodPost, url)
}

// Put creates a PUT request builder
func Put(url string) *RequestBuilder {
	return NewRequestBuilder(http.MethodPut, url)
}

// Delete creates a DELETE request builder
func Delete(url string) *RequestBuilder {
	return NewRequestBuilder(http.MethodDelete, url)
}

// Patch creates a PATCH request builder
func Patch(url string) *RequestBuilder {
	return NewRequestBuilder(http.MethodPatch, url)
}
