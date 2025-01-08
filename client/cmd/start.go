package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/matheuspolitano/quiz-go/client/internal/config"
	"github.com/matheuspolitano/quiz-go/client/internal/quiz"
)

// startCmd defines the command that starts the quiz flow.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the quiz",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load config from the current directory (".") or wherever your app.env is.
		cfg, err := config.LoadConfig(".")
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			return
		}

		// 2. Use the loaded config's APIURL
		quiz.RunQuizFlow(cfg.API_URL)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
