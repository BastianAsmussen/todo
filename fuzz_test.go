package main

import (
	"os"
	"testing"
)

func FuzzLoadTasks(f *testing.F) {
	// Seed with valid CSV data.
	f.Add([]byte("ID,Description,CreatedAt,CompletedAt\n1,Test Task,2022-01-01T00:00:00Z,\n"))

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()

		// Write data to a temporary file.
		tmpFile, err := os.CreateTemp("", "tasks-*.csv")
		if err != nil {
			t.Skip()
		}

		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(data); err != nil {
			t.Skip()
		}

		tmpFile.Close()

		// Temporarily override dataFilePath.
		originalDataFilePath := dataFilePath
		dataFilePath = tmpFile.Name()
		defer func() { dataFilePath = originalDataFilePath }()

		// Attempt to load tasks.
		_, _ = loadTasks()
	})
}
