package answerable

import (
	"context"

	"github.com/nchursin/serenity-go/serenity/core"
)

// functionQuestion[T] implements core.Question[T] for functions.
// It executes the provided function when asked by any actor.
type functionQuestion[T any] struct {
	description string
	function    func(core.Actor, context.Context) (T, error)
}

// AnsweredBy executes the function and returns its result.
// If the function returns an error, that error is returned.
func (f *functionQuestion[T]) AnsweredBy(actor core.Actor, ctx context.Context) (T, error) {
	return f.function(actor, ctx)
}

// Description returns the provided description.
func (f *functionQuestion[T]) Description() string {
	return f.description
}
