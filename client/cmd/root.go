package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd is the base command for your CLI.
var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A CLI application for quizzes",
	Long:  "quiz is a terminal-based application that tests your knowledge by interacting with an API.",
}

// Execute is the entry point for Cobra to run the root command.
func Execute() error {
	return rootCmd.Execute()
}
