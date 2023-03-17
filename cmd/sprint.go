package cmd

import (
	"type-chan/app"

	"github.com/spf13/cobra"
)

// sprintCmd launches the test in sprint mode.
var sprintCmd = &cobra.Command{
	Use:   "sprint",
	Short: "Begins the test in sprint mode",
	Long:  `Begins the test in sprint mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New()
		a.Start(app.Sprint)
	},
}

func init() {
	rootCmd.AddCommand(sprintCmd)
}
