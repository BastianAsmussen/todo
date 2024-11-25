package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "todo",
		Short: "Manage your tasks.",
	}

	// Add commands.
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newCompleteCmd())
	rootCmd.AddCommand(newDeleteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
