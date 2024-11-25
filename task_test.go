package main

import (
	"os"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	// Setup.
	cleanupDataFiles()

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Test.
	description := "Test task"
	task := NewTask(generateID(tasks), description)
	tasks = append(tasks, task)
	if err := saveTasks(tasks); err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Verify.
	tasks, err = loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d!", len(tasks))
	}

	if tasks[0].Description != description {
		t.Errorf("Expected description '%s', got '%s'!", description, tasks[0].Description)
	}
}

func TestCompleteTask(t *testing.T) {
	// Setup.
	cleanupDataFiles()

	tasks := []Task{NewTask(1, "Test Completion")}
	if err := saveTasks(tasks); err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Test.
	err := completeTask(1)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	// Verify.
	tasks, err = loadTasks()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	if tasks[0].CompletedAt == nil {
		t.Errorf("Expected task to be completed!")
	}
}

func TestListTasks(t *testing.T) {
	// Setup.
	cleanupDataFiles()

	tasks := []Task{
		NewTask(1, "Incomplete Task"),
		func() Task {
			task := NewTask(2, "Completed Task")
			now := time.Now()
			task.CompletedAt = &now

			return task
		}(),
	}

	if err := saveTasks(tasks); err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Test.
	uncompletedTasks := filterTasks(tasks, false)
	if len(uncompletedTasks) != 1 {
		t.Errorf("Expected 1 uncompleted task, got %d!", len(uncompletedTasks))
	}
	if uncompletedTasks[0].Description != "Incomplete Task" {
		t.Errorf("Unexpected task in uncompleted tasks: %s", uncompletedTasks[0].Description)
	}

	allTasks := filterTasks(tasks, true)
	if len(allTasks) != 2 {
		t.Errorf("Expected 2 tasks when showAll is true, got %d!", len(allTasks))
	}
}

func cleanupDataFiles() {
	os.Remove(dataFilePath)
	os.Remove(walFilePath)
}

// Helper function to filter tasks as per the `list` command logic.
func filterTasks(tasks []Task, showAll bool) []Task {
	var filtered []Task
	for _, task := range tasks {
		if !showAll && task.CompletedAt != nil {
			continue
		}

		filtered = append(filtered, task)
	}

	return filtered
}
