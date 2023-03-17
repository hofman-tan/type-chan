package cmd

import (
	"type-chan/app"

	"github.com/spf13/cobra"
)

// timedCmd launches the test in timed mode.
var timedCmd = &cobra.Command{
	Use:   "timed",
	Short: "Begins the test in timed mode",
	Long:  `Begins the test in timed mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New()
		a.Start(app.Timed)
	},
}

func init() {
	rootCmd.AddCommand(timedCmd)
}
