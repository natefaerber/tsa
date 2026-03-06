package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// stage groups a hook name with its steps.
type stage struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Usage       string     `json:"usage"`
	Steps       []StepInfo `json:"steps"`
}

var (
	flagFormat string
)

var stages = []stage{
	{Name: "pre-commit", Description: "Run pre-commit checks on staged changes", Usage: "tsa pre-commit [files...]", Steps: preCommitSteps},
	{Name: "commit-msg", Description: "Clean up commit messages", Usage: "tsa commit-msg <file>", Steps: commitMsgSteps},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stages and their steps",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch flagFormat {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "  ")
			return enc.Encode(stages)
		default:
			for i, s := range stages {
				fmt.Printf("%s  %s\n", colorStage(s.Name), colorDim(s.Description))
				fmt.Printf("  %s\n", colorDim(s.Usage))
				printStepTree(s.Steps)
				if i < len(stages)-1 {
					fmt.Println()
				}
			}
			return nil
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&flagFormat, "format", "f", "text", "output format (text, json)")
	rootCmd.AddCommand(listCmd)
}

func printStepTree(steps []StepInfo) {
	for i, s := range steps {
		connector := treeBranch()
		if i == len(steps)-1 {
			connector = treeLast()
		}
		fmt.Printf("%s%s  %s\n", connector, colorStep(s.Name), colorDim(s.Description))
	}
}
