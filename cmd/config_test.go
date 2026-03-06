package cmd

import (
	"os"
	"testing"
)

func TestCompilePatternsText(t *testing.T) {
	patterns := []Pattern{{Text: "hello"}}
	matchers, err := compilePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}
	if !matchesAny("say hello world", matchers) {
		t.Error("expected match")
	}
	if matchesAny("say goodbye", matchers) {
		t.Error("expected no match")
	}
}

func TestCompilePatternsTextCaseInsensitive(t *testing.T) {
	patterns := []Pattern{{Text: "hello", CaseInsensitive: true}}
	matchers, err := compilePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}
	if !matchesAny("HELLO world", matchers) {
		t.Error("expected case-insensitive match")
	}
}

func TestCompilePatternsRegex(t *testing.T) {
	patterns := []Pattern{{Regex: `Co-[Aa]uthored-[Bb]y:.*bot`}}
	matchers, err := compilePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}
	if !matchesAny("Co-Authored-By: some-bot", matchers) {
		t.Error("expected regex match")
	}
	if matchesAny("Co-Authored-By: human", matchers) {
		t.Error("expected no match")
	}
}

func TestCompilePatternsRegexCaseInsensitive(t *testing.T) {
	patterns := []Pattern{{Regex: "generated with", CaseInsensitive: true}}
	matchers, err := compilePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}
	if !matchesAny("GENERATED WITH tool", matchers) {
		t.Error("expected case-insensitive regex match")
	}
}

func TestCompilePatternsErrorBothTextAndRegex(t *testing.T) {
	patterns := []Pattern{{Text: "foo", Regex: "bar"}}
	_, err := compilePatterns(patterns)
	if err == nil {
		t.Error("expected error for pattern with both text and regex")
	}
}

func TestCompilePatternsErrorNeitherTextNorRegex(t *testing.T) {
	patterns := []Pattern{{CaseInsensitive: true}}
	_, err := compilePatterns(patterns)
	if err == nil {
		t.Error("expected error for pattern with neither text nor regex")
	}
}

func TestCompilePatternsErrorInvalidRegex(t *testing.T) {
	patterns := []Pattern{{Regex: "[invalid"}}
	_, err := compilePatterns(patterns)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestLoadConfigMissing(t *testing.T) {
	flagConfig = "/tmp/tsa-nonexistent-config.yaml"
	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.CommitMsg.StripAttribution.Patterns) != 0 {
		t.Error("expected empty config for missing file")
	}
}

func TestLoadConfigValid(t *testing.T) {
	content := `
commit-msg:
  strip-attribution:
    patterns:
      - text: "test-pattern"
      - regex: "foo.*bar"
        case_insensitive: true
`
	f, err := os.CreateTemp("", "tsa-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(content)
	f.Close()

	flagConfig = f.Name()
	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}

	patterns := cfg.CommitMsg.StripAttribution.Patterns
	if len(patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(patterns))
	}
	if patterns[0].Text != "test-pattern" {
		t.Errorf("expected text 'test-pattern', got %q", patterns[0].Text)
	}
	if patterns[1].Regex != "foo.*bar" || !patterns[1].CaseInsensitive {
		t.Errorf("unexpected second pattern: %+v", patterns[1])
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	f, err := os.CreateTemp("", "tsa-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("commit-msg:\n  strip-attribution:\n    patterns: not-a-list\n")
	f.Close()

	flagConfig = f.Name()
	_, err = loadConfig()
	if err == nil {
		t.Error("expected error for invalid YAML structure")
	}
}

func TestMatchesAnyMultipleMatchers(t *testing.T) {
	patterns := []Pattern{
		{Text: "alpha"},
		{Text: "beta"},
		{Regex: "gam+a"},
	}
	matchers, err := compilePatterns(patterns)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		line string
		want bool
	}{
		{"has alpha in it", true},
		{"has beta in it", true},
		{"has gamma in it", true},
		{"has delta in it", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := matchesAny(tt.line, matchers); got != tt.want {
			t.Errorf("matchesAny(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}
