package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nchursin/serenity-go/serenity/core"
)

// LastResponseStatus returns the status code of the last response
type LastResponseStatus struct{}

// AnsweredBy returns the status code from the last HTTP response
func (lr LastResponseStatus) AnsweredBy(actor core.Actor) (int, error) {
	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return 0, fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)
	resp := callAbility.LastResponse()
	if resp == nil {
		return 0, fmt.Errorf("no response available")
	}

	return resp.StatusCode, nil
}

// Description returns the question description
func (lr LastResponseStatus) Description() string {
	return "asks for the last response status code"
}

// LastResponseBody returns the body of the last response
type LastResponseBody struct{}

// AnsweredBy returns the body from the last HTTP response
func (lr LastResponseBody) AnsweredBy(actor core.Actor) (string, error) {
	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return "", fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)
	resp := callAbility.LastResponse()
	if resp == nil {
		return "", fmt.Errorf("no response available")
	}

	defer func() {
		_ = resp.Body.Close() // Ignore cleanup error
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Restore body for potential re-reading
	resp.Body = io.NopCloser(strings.NewReader(string(body)))

	return string(body), nil
}

// Description returns the question description
func (lr LastResponseBody) Description() string {
	return "asks for the last response body"
}

// ResponseHeader returns a specific header from the last response
type ResponseHeader struct {
	key string
}

// NewResponseHeader creates a new question for a specific header
func NewResponseHeader(key string) ResponseHeader {
	return ResponseHeader{key: key}
}

// AnsweredBy returns the header value from the last HTTP response
func (rh ResponseHeader) AnsweredBy(actor core.Actor) (string, error) {
	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return "", fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)
	resp := callAbility.LastResponse()
	if resp == nil {
		return "", fmt.Errorf("no response available")
	}

	return resp.Header.Get(rh.key), nil
}

// Description returns the question description
func (rh ResponseHeader) Description() string {
	return fmt.Sprintf("asks for the response header '%s'", rh.key)
}

// ResponseBodyAsJSON returns the response body parsed as JSON
type ResponseBodyAsJSON[T any] struct{}

// NewResponseBodyAsJSON creates a new question for JSON response body
func NewResponseBodyAsJSON[T any]() ResponseBodyAsJSON[T] {
	return ResponseBodyAsJSON[T]{}
}

// AnsweredBy returns the response body parsed as JSON
func (rbaj ResponseBodyAsJSON[T]) AnsweredBy(actor core.Actor) (T, error) {
	var result T

	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return result, fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)
	resp := callAbility.LastResponse()
	if resp == nil {
		return result, fmt.Errorf("no response available")
	}

	defer func() {
		_ = resp.Body.Close() // Ignore cleanup error
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}

	// Restore body for potential re-reading
	resp.Body = io.NopCloser(strings.NewReader(string(body)))

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return result, nil
}

// Description returns the question description
func (rbaj ResponseBodyAsJSON[T]) Description() string {
	return "asks for the last response body as JSON"
}

// JSONPath represents a JSON path query on the response body
type JSONPath struct {
	path string
}

// NewJSONPath creates a new JSON path question
func NewJSONPath(path string) JSONPath {
	return JSONPath{path: path}
}

// AnsweredBy returns the value at the specified JSON path
func (jp JSONPath) AnsweredBy(actor core.Actor) (any, error) {
	ability, err := actor.AbilityTo(&callAnAPI{})
	if err != nil {
		return nil, fmt.Errorf("actor does not have the ability to call an API: %w", err)
	}

	callAbility := ability.(CallAnAPI)
	resp := callAbility.LastResponse()
	if resp == nil {
		return nil, fmt.Errorf("no response available")
	}

	defer func() {
		_ = resp.Body.Close() // Ignore cleanup error
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Restore body for potential re-reading
	resp.Body = io.NopCloser(strings.NewReader(string(body)))

	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return jp.extractValue(data, strings.Split(jp.path, "."))
}

// Description returns the question description
func (jp JSONPath) Description() string {
	return fmt.Sprintf("asks for JSON path '%s'", jp.path)
}

// extractValue recursively extracts value from JSON data using path segments
func (jp JSONPath) extractValue(data any, path []string) (any, error) {
	if len(path) == 0 {
		return data, nil
	}

	switch v := data.(type) {
	case map[string]any:
		if nextValue, exists := v[path[0]]; exists {
			return jp.extractValue(nextValue, path[1:])
		}
		return nil, fmt.Errorf("path '%s' not found", strings.Join(path, "."))

	case []any:
		if path[0] == "*" {
			var results []any
			for _, item := range v {
				value, err := jp.extractValue(item, path[1:])
				if err == nil {
					results = append(results, value)
				}
			}
			return results, nil
		}

		index, err := strconv.Atoi(path[0])
		if err != nil {
			return nil, fmt.Errorf("invalid array index '%s': %w", path[0], err)
		}

		if index < 0 || index >= len(v) {
			return nil, fmt.Errorf("array index '%d' out of bounds", index)
		}

		return jp.extractValue(v[index], path[1:])

	default:
		return nil, fmt.Errorf("cannot traverse non-object/array with path '%s'", path[0])
	}
}

// ResponseTime returns the response time of the last request
type ResponseTime struct{}

// AnsweredBy returns the response time (currently returns 0 as timing needs to be implemented in interactions)
func (rt ResponseTime) AnsweredBy(actor core.Actor) (int64, error) {
	// This would need timing implementation in interactions
	return 0, nil
}

// Description returns the question description
func (rt ResponseTime) Description() string {
	return "asks for the last request response time"
}

// Convenience variables for common questions
var (
	LastResponseStatusQ = LastResponseStatus{}
	LastResponseBodyQ   = LastResponseBody{}
	ResponseTimeQ       = ResponseTime{}
)
