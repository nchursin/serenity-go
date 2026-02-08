package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

// TestJSONPlaceholderPosts demonstrates basic CRUD operations with JSONPlaceholder
func TestJSONPlaceholderPosts(t *testing.T) {
	// Create an actor with API ability
	actor := core.NewActor("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// CREATE: Create a new post
	newPost := map[string]interface{}{
		"title":  "Serenity-Go Test Post",
		"body":   "Testing Serenity-Go framework with real API",
		"userId": 1,
	}

	err := actor.AttemptsTo(
		core.NewInteraction("creates a new post", func(a core.Actor) error {
			req, err := api.Post("/posts").
				With(newPost).
				Build()
			if err != nil {
				return err
			}

			sendReq := api.SendRequest(req)
			return sendReq.PerformAs(a)
		}),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
	)

	require.NoError(t, err)
}

// TestJSONPlaceholderGetPosts demonstrates getting a list of posts
func TestJSONPlaceholderGetPosts(t *testing.T) {
	actor := core.NewActor("Reader").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	err := actor.AttemptsTo(
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)

	require.NoError(t, err)
}

// TestJSONPlaceholderErrorHandling demonstrates error scenarios
func TestJSONPlaceholderErrorHandling(t *testing.T) {
	actor := core.NewActor("ErrorTester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Test 404 - non-existent resource
	err := actor.AttemptsTo(
		api.GetRequest("/posts/99999"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(404)),
	)

	require.NoError(t, err)
}
