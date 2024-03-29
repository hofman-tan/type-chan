package cmd

import (
	"fmt"
	"typechan/app"

	"github.com/spf13/cobra"
)

// timedCmd launches the typing test in timed mode.
var timedCmd = &cobra.Command{
	Use:   "timed",
	Short: "Begins the test in timed mode",
	Long:  `Begins the test in timed mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if app.Timeout <= 0 {
			return fmt.Errorf("timeout must be larger than 0")
		}

		a := app.New()
		a.Start(app.Timed)
		return nil
	},
}

func init() {
	timedCmd.PersistentFlags().DurationVarP(&app.Timeout, "seconds", "s", app.Timeout, "Timer timeout e.g. 30s, 5m")
	rootCmd.AddCommand(timedCmd)
}
