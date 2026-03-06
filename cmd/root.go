package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Set via ldflags at build time.
	Version = "dev"
	Commit  = "unknown"
)

var (
	flagSkip      []string
	flagOnly      []string
	flagQuiet     bool
	flagListSteps bool
)

var rootCmd = &cobra.Command{
	Use:     "tsa",
	Short:   "Personal git hook utilities",
	Version: Version + " (" + Commit + ")",
}

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&flagSkip, "skip", nil, "steps to skip (comma-separated)")
	rootCmd.PersistentFlags().StringSliceVar(&flagOnly, "only", nil, "run only these steps (comma-separated)")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "suppress skip messages")
	rootCmd.PersistentFlags().BoolVar(&flagListSteps, "list-steps", false, "list available steps and exit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
