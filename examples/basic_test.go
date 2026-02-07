package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nchursin/serenity-go/serenity/api"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

// TestJSONPlaceholderBasics demonstrates basic API testing with JSONPlaceholder
func TestJSONPlaceholderBasics(t *testing.T) {
	actor := core.NewActor("APITester").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	// Test GET posts - should return existing posts
	err := actor.AttemptsTo(
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
	)
	require.NoError(t, err)

	// Test GET users - should return existing users
	err = actor.AttemptsTo(
		api.GetRequest("/users"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("email")),
	)
	require.NoError(t, err)

	// Test GET specific post
	err = actor.AttemptsTo(
		api.GetRequest("/posts/1"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("sunt aut facere")),
	)
	require.NoError(t, err)
}

// TestJSONPlaceholderErrors demonstrates error scenarios
func TestJSONPlaceholderErrors(t *testing.T) {
	actor := core.NewActor("ErrorTester").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	// Test 404 for non-existent post
	err := actor.AttemptsTo(
		api.GetRequest("/posts/99999"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(404)),
	)
	require.NoError(t, err)

	// Test 404 for non-existent endpoint
	err = actor.AttemptsTo(
		api.GetRequest("/nonexistent"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(404)),
	)
	require.NoError(t, err)
}

// TestJSONPlaceholderPostRequest demonstrates POST request functionality
func TestJSONPlaceholderPostRequest(t *testing.T) {
	actor := core.NewActor("PostTester").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	// Create a new post (JSONPlaceholder will return the data with an ID)
	newPost := map[string]interface{}{
		"title":  "Test Post",
		"body":   "This is a test post",
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
			return api.SendRequest(req).PerformAs(a)
		}),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("Test Post")),
	)
	require.NoError(t, err)
}

// TestJSONPlaceholderHeaders demonstrates header assertions
func TestJSONPlaceholderHeaders(t *testing.T) {
	actor := core.NewActor("HeaderTester").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	err := actor.AttemptsTo(
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.NewResponseHeader("content-type"), expectations.Contains("json")),
	)
	require.NoError(t, err)
}

// TestMultipleActors demonstrates using multiple actors
func TestMultipleActors(t *testing.T) {
	// Different actors for different roles
	admin := core.NewActor("Admin").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))
	user := core.NewActor("RegularUser").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	// Both actors can read posts
	err := admin.AttemptsTo(
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)
	require.NoError(t, err)

	err = user.AttemptsTo(
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)
	require.NoError(t, err)
}

// TestTaskComposition demonstrates creating reusable tasks
func TestTaskComposition(t *testing.T) {
	actor := core.NewActor("TaskUser").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

	// Define a reusable task for checking API availability
	checkApiAvailable := core.Where(
		"checks if API is available",
		api.GetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)

	// Use the task
	err := actor.AttemptsTo(
		checkApiAvailable,
	)
	require.NoError(t, err)
}
