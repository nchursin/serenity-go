package api

import (
	"net/http"

	"github.com/nchursin/serenity-go/serenity/core"
)

// SendRequest creates a SendRequest interaction (exported function)
func SendRequest(req *http.Request) core.Activity {
	return a(req)
}

// Convenience functions for HTTP methods that return activities
func GetRequest(url string) core.Activity {
	req, err := NewRequestBuilder("GET", url).Build()
	if err != nil {
		return core.NewInteraction("get request", func(actor core.Actor) error {
			return err
		})
	}
	return SendRequest(req)
}

// TODO: научится констурировать PostRequest с апишкой типа SendPostRequest(url).WithBody(jsonMarshable)
func PostRequest(url string) core.Activity {
	req, err := NewRequestBuilder("POST", url).Build()
	if err != nil {
		return core.NewInteraction("post request", func(actor core.Actor) error {
			return err
		})
	}
	return SendRequest(req)
}

// TODO: научится констурировать PutRequest с апишкой типа SendPutRequest(url).WithBody(jsonMarshable)
func PutRequest(url string) core.Activity {
	req, err := NewRequestBuilder("PUT", url).Build()
	if err != nil {
		return core.NewInteraction("put request", func(actor core.Actor) error {
			return err
		})
	}
	return SendRequest(req)
}

func DeleteRequest(url string) core.Activity {
	req, err := NewRequestBuilder("DELETE", url).Build()
	if err != nil {
		return core.NewInteraction("delete request", func(actor core.Actor) error {
			return err
		})
	}
	return SendRequest(req)
}
