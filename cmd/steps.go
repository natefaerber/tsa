package cmd

import (
	"fmt"
	"os"
	"slices"
)

// Step represents a named unit of work within a hook.
type Step struct {
	Name        string
	Description string
	Run         func() error
}

// StepInfo holds static metadata for a step (used for --list-steps).
type StepInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RunSteps executes steps, filtering by skip/only flags.
// skip and only are mutually exclusive; if both are set, only takes precedence.
func RunSteps(steps []Step, skip, only []string, quiet bool) error {
	for _, s := range steps {
		if shouldSkip(s.Name, skip, only) {
			if !quiet {
				fmt.Fprintf(os.Stderr, "skipping %s\n", s.Name)
			}
			continue
		}
		if err := s.Run(); err != nil {
			return err
		}
	}
	return nil
}

func shouldSkip(name string, skip, only []string) bool {
	if len(only) > 0 {
		return !slices.Contains(only, name)
	}
	return slices.Contains(skip, name)
}

// PrintStepInfo prints step names and descriptions as a tree.
// Returns true if it printed (i.e. flagListSteps was set).
func PrintStepInfo(steps []StepInfo) bool {
	if !flagListSteps {
		return false
	}
	printStepTree(steps)
	return true
}
