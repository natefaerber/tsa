package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Default patterns used when no config file is present.
var defaultStripPatterns = []Pattern{
	{Text: "🤖 Generated with"},
	{Text: "Co-Authored-By: Claude", CaseInsensitive: true},
}

var commitMsgCmd = &cobra.Command{
	Use:   "commit-msg <file>",
	Short: "Clean up commit messages (strip AI attribution, collapse blanks)",
	Long:  "Steps: strip-attribution, collapse-blanks",
	Args: func(cmd *cobra.Command, args []string) error {
		if flagListSteps {
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: runCommitMsg,
}

func init() {
	rootCmd.AddCommand(commitMsgCmd)
}

var commitMsgSteps = []StepInfo{
	{Name: "strip-attribution", Description: "Remove AI attribution and co-author lines"},
	{Name: "collapse-blanks", Description: "Collapse consecutive blank lines and trim leading/trailing blanks"},
}

func runCommitMsg(cmd *cobra.Command, args []string) error {
	if PrintStepInfo(commitMsgSteps) {
		return nil
	}
	path := args[0]

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	patterns := defaultStripPatterns
	if len(cfg.CommitMsg.StripAttribution.Patterns) > 0 {
		patterns = cfg.CommitMsg.StripAttribution.Patterns
	}

	matchers, err := compilePatterns(patterns)
	if err != nil {
		return fmt.Errorf("strip-attribution: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading commit message: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	steps := []Step{
		{Name: "strip-attribution", Description: commitMsgSteps[0].Description, Run: func() error {
			lines = filterLines(lines, matchers)
			return nil
		}},
		{Name: "collapse-blanks", Description: commitMsgSteps[1].Description, Run: func() error {
			lines = collapseBlankLines(lines)
			return nil
		}},
	}

	if err := RunSteps(steps, flagSkip, flagOnly, flagQuiet); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("writing commit message: %w", err)
	}

	return nil
}

func filterLines(lines []string, matchers []matcher) []string {
	var filtered []string
	for _, line := range lines {
		if matchesAny(line, matchers) {
			continue
		}
		filtered = append(filtered, line)
	}
	return filtered
}

func collapseBlankLines(lines []string) []string {
	var result []string
	contentStarted := false
	pendingBlank := false

	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			if !contentStarted {
				contentStarted = true
			} else if pendingBlank {
				result = append(result, "")
			}
			result = append(result, line)
			pendingBlank = false
		} else if contentStarted {
			pendingBlank = true
		}
	}

	return result
}
