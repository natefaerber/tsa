package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var preCommitCmd = &cobra.Command{
	Use:   "pre-commit [files...]",
	Short: "Run pre-commit checks on staged changes",
	Long:  "Steps: nocommit\n\nWhen files are given, only those files are checked. Otherwise all staged files are checked.",
	RunE:  runPreCommit,
}

func init() {
	rootCmd.AddCommand(preCommitCmd)
}

var preCommitSteps = []StepInfo{
	{Name: "nocommit", Description: "Check for !nocommit markers in staged changes"},
}

func runPreCommit(cmd *cobra.Command, args []string) error {
	if PrintStepInfo(preCommitSteps) {
		return nil
	}
	steps := []Step{
		{Name: "nocommit", Description: preCommitSteps[0].Description, Run: func() error { return checkNoCommit(args) }},
	}
	return RunSteps(steps, flagSkip, flagOnly, flagQuiet)
}

func checkNoCommit(files []string) error {
	diff, err := git.StagedDiff(files)
	if err != nil {
		return err
	}

	if strings.Contains(diff, "!nocommit") {
		fmt.Fprintln(os.Stderr, "Trying to commit non-committable code.")
		fmt.Fprintln(os.Stderr, "Remove the !nocommit string and try again.")
		os.Exit(1)
	}

	return nil
}
