package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

func completeTask(taskID int) error {
	tasks, err := loadTasks()
	if err != nil {
		return fmt.Errorf("Error loading tasks: %v", err)
	}

	for i, task := range tasks {
		if task.ID == taskID {
			now := time.Now()
			tasks[i].CompletedAt = &now
			if err := saveTasks(tasks); err != nil {
				return fmt.Errorf("Error saving tasks: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("Task with ID %d not found!", taskID)
}

// newAddCmd returns the add command.
func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [description]",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			tasks, err := loadTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
				return
			}

			task := NewTask(generateID(tasks), description)

			tasks = append(tasks, task)
			if err := saveTasks(tasks); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
				return
			}
			fmt.Printf("Task added: %s\n", task.Description)
		},
	}
}

// newListCmd returns the list command.
func newListCmd() *cobra.Command {
	var showAll bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := loadTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
				return
			}

			// Initialize tabwriter
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

			// Write header
			fmt.Fprintln(w, "ID\tDescription\tCreated At\tCompleted At")

			// Write tasks
			for _, task := range tasks {
				if !showAll && task.CompletedAt != nil {
					continue // Skip completed tasks unless --all is specified
				}

				createdStr := timediff.TimeDiff(task.CreatedAt)
				var completedStr string
				if task.CompletedAt != nil {
					completedStr = timediff.TimeDiff(*task.CompletedAt)
				} else {
					completedStr = "Incomplete..."
				}

				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					task.ID,
					task.Description,
					createdStr,
					completedStr)
			}

			// Flush the writer.
			w.Flush()
		},
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks, including completed ones.")

	return cmd
}

// newCompleteCmd returns the complete command.
func newCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "complete [taskID]",
		Short: "Mark a task as complete",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid task ID '%s'. Please provide a valid numeric ID!\n", args[0])
				return
			}

			if err := completeTask(taskID); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Printf("Task ID %d marked as complete.\n", taskID)
			}
		},
	}
}

// newDeleteCmd returns the delete command.
func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [taskID]",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid task ID '%s'. Please provide a valid numeric ID.\n", args[0])
				return
			}

			tasks, err := loadTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
				return
			}

			var updatedTasks []Task
			var found bool
			for _, task := range tasks {
				if task.ID != taskID {
					updatedTasks = append(updatedTasks, task)
				} else {
					found = true
				}
			}

			if !found {
				fmt.Fprintf(os.Stderr, "Task with ID %d not found!\n", taskID)
				return
			}

			if err := saveTasks(updatedTasks); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving tasks: %v\n", err)
			} else {
				fmt.Printf("Task ID %d deleted.\n", taskID)
			}
		},
	}
}

// generateID generates a new unique ID for tasks.
func generateID(tasks []Task) int {
	idSet := make(map[int]struct{})
	for _, task := range tasks {
		idSet[task.ID] = struct{}{}
	}

	newID := 1
	for {
		if _, exists := idSet[newID]; !exists {
			break
		}

		newID++
	}

	return newID
}
