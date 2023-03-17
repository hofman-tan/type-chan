package cmd

import (
	"fmt"
	"type-chan/app"

	"github.com/spf13/cobra"
)

// timedCmd launches the test in timed mode.
var timedCmd = &cobra.Command{
	Use:   "timed",
	Short: "Begins the test in timed mode",
	Long:  `Begins the test in timed mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if app.Countdown <= 0 {
			return fmt.Errorf("countdown seconds must be larger than 0")
		}

		a := app.New()
		a.Start(app.Timed)
		return nil
	},
}

func init() {
	timedCmd.PersistentFlags().IntVarP(&app.Countdown, "countdown", "c", app.Countdown, "Timer countdown in seconds")
	rootCmd.AddCommand(timedCmd)
}
