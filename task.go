package main

import (
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// NewTask creates a new Task instance.
func NewTask(id int, description string) Task {
	return Task{
		ID:          id,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
