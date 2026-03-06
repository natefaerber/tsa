package cmd

import (
	"os"

	"golang.org/x/term"
)

var (
	bold  = "\033[1m"
	dim   = "\033[2m"
	cyan  = "\033[36m"
	reset = "\033[0m"
)

func init() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		bold = ""
		dim = ""
		cyan = ""
		reset = ""
	}
}

func colorStage(s string) string { return bold + cyan + s + reset }
func colorStep(s string) string  { return cyan + s + reset }
func colorDim(s string) string   { return dim + s + reset }
func treeBranch() string         { return colorDim("├── ") }
func treeLast() string           { return colorDim("└── ") }
