package core

import (
	"fmt"
)

// task implements the Task interface
type task struct {
	description string
	activities  []Activity
}

// NewTask creates a new task with the given description and activities
func NewTask(description string, activities ...Activity) Task {
	return &task{
		description: description,
		activities:  activities,
	}
}

// Description returns the task description
func (t *task) Description() string {
	return t.description
}

// PerformAs executes the task as the given actor
func (t *task) PerformAs(actor Actor) error {
	for _, activity := range t.activities {
		if err := activity.PerformAs(actor); err != nil {
			return fmt.Errorf("task '%s' failed during activity '%s': %w",
				t.Description(), activity.Description(), err)
		}
	}
	return nil
}

// FailureMode returns the failure mode for tasks (default: FailFast)
func (t *task) FailureMode() FailureMode {
	return FailFast
}

// Where creates a new task with the given description and activities
// This is a convenience function similar to Serenity/JS Task.where
func Where(description string, activities ...Activity) Task {
	return NewTask(description, activities...)
}

// interaction implements the Interaction interface
type interaction struct {
	description string
	perform     func(actor Actor) error
}

// Do creates a new interaction with the given description and perform function
func Do(description string, perform func(actor Actor) error) Interaction {
	return &interaction{
		description: description,
		perform:     perform,
	}
}

// Description returns the interaction description
func (i *interaction) Description() string {
	return fmt.Sprintf("#actor %s", i.description)
}

// PerformAs executes the interaction as the given actor
func (i *interaction) PerformAs(actor Actor) error {
	return i.perform(actor)
}

// FailureMode returns the failure mode for interactions (default: FailFast)
func (i *interaction) FailureMode() FailureMode {
	return FailFast
}
