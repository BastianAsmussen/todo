package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	homeDir, _   = os.UserHomeDir()
	dataFilePath = filepath.Join(homeDir, ".todo", "tasks.csv")
	walFilePath  = filepath.Join(homeDir, ".todo", "tasks.wal")
)

var lock sync.Mutex

// loadTasks loads tasks from the CSV file.
func loadTasks() ([]Task, error) {
	file, err := os.Open(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil // No file means no tasks yet.
		}

		return nil, fmt.Errorf("Failed to open data file: %w", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to read data file: %w", err)
	}

	var tasks []Task
	for _, record := range records[1:] { // Skip header row.
		if len(record) < 4 {
			fmt.Fprintf(os.Stderr, "Warning: Skipping invalid record: %v\n", record)
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("Invalid task ID '%s': %w", record[0], err)
		}
		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			return nil, fmt.Errorf("Invalid creation date '%s': %w", record[2], err)
		}

		var completedAt *time.Time
		if record[3] != "" {
			t, err := time.Parse(time.RFC3339, record[3])
			if err != nil {
				return nil, fmt.Errorf("Invalid completion date '%s': %w", record[3], err)
			}

			completedAt = &t
		}

		tasks = append(tasks, Task{
			ID:          id,
			Description: record[1],
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		})
	}

	return tasks, nil
}

// saveTasks saves tasks to the CSV file.
func saveTasks(tasks []Task) error {
	lock.Lock()
	defer lock.Unlock()

	// Ensure the directory exists.
	if err := os.MkdirAll(filepath.Dir(dataFilePath), 0755); err != nil {
		return fmt.Errorf("Failed to create data directory: %w", err)
	}

	// First write to WAL.
	if err := writeWAL(tasks); err != nil {
		return fmt.Errorf("Failed to write WAL: %w", err)
	}

	file, err := os.Create(dataFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create data file: %w", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header.
	if err := writer.Write([]string{"ID", "Description", "CreatedAt", "CompletedAt"}); err != nil {
		return fmt.Errorf("Failed to write header: %w", err)
	}

	// Write task data.
	for _, task := range tasks {
		var completedAt string
		if task.CompletedAt != nil {
			completedAt = task.CompletedAt.Format(time.RFC3339)
		}

		if err := writer.Write([]string{
			strconv.Itoa(task.ID),
			task.Description,
			task.CreatedAt.Format(time.RFC3339),
			completedAt,
		}); err != nil {
			return fmt.Errorf("Failed to write task data: %w", err)
		}
	}

	return nil
}

// writeWAL writes tasks to the WAL file before committing to the main file.
func writeWAL(tasks []Task) error {
	file, err := os.Create(walFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create WAL file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("Failed to encode WAL data: %w", err)
	}

	return nil
}
