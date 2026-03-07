package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the default configuration",
	Long:  "Print the built-in default config as YAML. Redirect to a file to start customizing.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := Config{
			CommitMsg: CommitMsgConfig{},
		}
		cfg.CommitMsg.StripAttribution.Patterns = defaultStripPatterns

		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf,
			yaml.Indent(2),
			yaml.IndentSequence(false),
		)
		if err := enc.Encode(cfg); err != nil {
			return fmt.Errorf("marshalling config: %w", err)
		}
		fmt.Fprint(os.Stdout, buf.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
