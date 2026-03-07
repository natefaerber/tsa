package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

// Pattern represents a single match rule for filtering lines.
type Pattern struct {
	Text            string `yaml:"text" json:"text,omitempty"`
	Regex           string `yaml:"regex,omitempty" json:"regex,omitempty"`
	CaseInsensitive bool   `yaml:"case_insensitive,omitempty" json:"case_insensitive,omitempty"`
}

// CommitMsgConfig holds configuration for the commit-msg stage.
type CommitMsgConfig struct {
	StripAttribution struct {
		Patterns []Pattern `yaml:"patterns"`
	} `yaml:"strip-attribution"`
}

// Config is the top-level configuration file structure.
type Config struct {
	CommitMsg CommitMsgConfig `yaml:"commit-msg"`
}

var flagConfig string

func init() {
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "config file path (default $XDG_CONFIG_HOME/tsa/config.yaml)")
}

// loadConfig reads the config file from --config flag or the default XDG path.
// Returns a zero Config (no error) if no config file exists.
func loadConfig() (Config, error) {
	path := flagConfig
	if path == "" {
		path = defaultConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, nil
}

func defaultConfigPath() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		home, _ := os.UserHomeDir()
		xdg = filepath.Join(home, ".config")
	}
	return filepath.Join(xdg, "tsa", "config.yaml")
}

// matcher is a compiled pattern that can test a line.
type matcher interface {
	matches(line string) bool
}

type textMatcher struct {
	text            string
	caseInsensitive bool
}

func (m textMatcher) matches(line string) bool {
	if m.caseInsensitive {
		return strings.Contains(strings.ToLower(line), strings.ToLower(m.text))
	}
	return strings.Contains(line, m.text)
}

type regexMatcher struct {
	re *regexp.Regexp
}

func (m regexMatcher) matches(line string) bool {
	return m.re.MatchString(line)
}

// compilePatterns turns config patterns into compiled matchers.
func compilePatterns(patterns []Pattern) ([]matcher, error) {
	matchers := make([]matcher, 0, len(patterns))
	for _, p := range patterns {
		switch {
		case p.Text != "" && p.Regex != "":
			return nil, fmt.Errorf("pattern has both text and regex: %q / %q", p.Text, p.Regex)
		case p.Text != "":
			matchers = append(matchers, textMatcher{text: p.Text, caseInsensitive: p.CaseInsensitive})
		case p.Regex != "":
			expr := p.Regex
			if p.CaseInsensitive {
				expr = "(?i)" + expr
			}
			re, err := regexp.Compile(expr)
			if err != nil {
				return nil, fmt.Errorf("compiling regex %q: %w", p.Regex, err)
			}
			matchers = append(matchers, regexMatcher{re: re})
		default:
			return nil, fmt.Errorf("pattern has neither text nor regex")
		}
	}
	return matchers, nil
}

// matchesAny returns true if any matcher matches the line.
func matchesAny(line string, matchers []matcher) bool {
	for _, m := range matchers {
		if m.matches(line) {
			return true
		}
	}
	return false
}
