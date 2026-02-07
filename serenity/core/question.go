package core

import (
	"fmt"
)

// question implements the Question interface
type question[T any] struct {
	description string
	ask         func(actor Actor) (T, error)
}

// NewQuestion creates a new question with the given description and ask function
func NewQuestion[T any](description string, ask func(actor Actor) (T, error)) Question[T] {
	return &question[T]{
		description: description,
		ask:         ask,
	}
}

// Description returns the question description
func (q *question[T]) Description() string {
	return fmt.Sprintf("asks %s", q.description)
}

// AnsweredBy returns the answer when asked by the given actor
func (q *question[T]) AnsweredBy(actor Actor) (T, error) {
	return q.ask(actor)
}

// Of creates a new question with the given description and ask function
// This is a convenience function for creating questions
func Of[T any](description string, ask func(actor Actor) (T, error)) Question[T] {
	return NewQuestion(description, ask)
}
