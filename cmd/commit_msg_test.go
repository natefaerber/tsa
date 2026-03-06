package cmd

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "commit-msg-test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(got)
}

func resetFlags() {
	flagSkip = nil
	flagOnly = nil
	flagQuiet = false
	flagConfig = ""
}

func TestCommitMsgStripsAttribution(t *testing.T) {
	resetFlags()
	input := "feat: add feature\n\nSome description.\n\nCo-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>\n"
	want := "feat: add feature\n\nSome description.\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgStripsEmoji(t *testing.T) {
	resetFlags()
	input := "fix: bug\n\n🤖 Generated with Claude Code\n"
	want := "fix: bug\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgCollapsesBlankLines(t *testing.T) {
	resetFlags()
	input := "\n\nfeat: thing\n\n\n\nBody here.\n\n\n"
	want := "feat: thing\n\nBody here.\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgLeavesCleanMessageAlone(t *testing.T) {
	resetFlags()
	input := "feat: clean commit\n\nNo attribution here.\n"
	want := "feat: clean commit\n\nNo attribution here.\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgStripsBothAttributions(t *testing.T) {
	resetFlags()
	input := "chore: stuff\n\nDetails.\n\n🤖 Generated with Claude Code\n\nCo-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>\n"
	want := "chore: stuff\n\nDetails.\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgSkipStripAttribution(t *testing.T) {
	resetFlags()
	flagSkip = []string{"strip-attribution"}
	flagQuiet = true

	input := "feat: thing\n\nCo-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>\n"
	// Attribution kept, but blank lines still collapsed
	want := "feat: thing\n\nCo-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgSkipCollapseBlanks(t *testing.T) {
	resetFlags()
	flagSkip = []string{"collapse-blanks"}
	flagQuiet = true

	input := "feat: thing\n\n\n\nBody here.\n\n\n"
	// Attribution stripped (nothing to strip), blanks NOT collapsed.
	// The trailing \n from Split + Join means the file keeps its structure.
	want := "feat: thing\n\n\n\nBody here.\n\n\n\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgOnlyStripAttribution(t *testing.T) {
	resetFlags()
	flagOnly = []string{"strip-attribution"}
	flagQuiet = true

	input := "\n\nfeat: thing\n\n\n\n🤖 Generated with Claude Code\n\n\n"
	// Attribution stripped, but blanks NOT collapsed
	want := "\n\nfeat: thing\n\n\n\n\n\n\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgCaseInsensitiveDefault(t *testing.T) {
	resetFlags()
	input := "feat: thing\n\nco-authored-by: Claude <noreply@anthropic.com>\n"
	want := "feat: thing\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "tsa-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestCommitMsgConfigTextPattern(t *testing.T) {
	cfg := writeConfig(t, `
commit-msg:
  strip-attribution:
    patterns:
      - text: "Signed-off-by: Bot"
`)
	resetFlags()
	flagConfig = cfg

	input := "feat: thing\n\nSigned-off-by: Bot\n"
	want := "feat: thing\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgConfigCaseInsensitive(t *testing.T) {
	cfg := writeConfig(t, `
commit-msg:
  strip-attribution:
    patterns:
      - text: "signed-off-by"
        case_insensitive: true
`)
	resetFlags()
	flagConfig = cfg

	input := "feat: thing\n\nSIGNED-OFF-BY: Someone\n"
	want := "feat: thing\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgConfigRegex(t *testing.T) {
	cfg := writeConfig(t, `
commit-msg:
  strip-attribution:
    patterns:
      - regex: "Co-[Aa]uthored-[Bb]y:.*bot"
`)
	resetFlags()
	flagConfig = cfg

	input := "feat: thing\n\nCo-Authored-By: some-bot <bot@example.com>\n"
	want := "feat: thing\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgConfigRegexCaseInsensitive(t *testing.T) {
	cfg := writeConfig(t, `
commit-msg:
  strip-attribution:
    patterns:
      - regex: "generated with"
        case_insensitive: true
`)
	resetFlags()
	flagConfig = cfg

	input := "feat: thing\n\nGENERATED WITH some tool\n"
	want := "feat: thing\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCommitMsgConfigReplacesDefaults(t *testing.T) {
	// Config with custom pattern should NOT strip the default patterns
	cfg := writeConfig(t, `
commit-msg:
  strip-attribution:
    patterns:
      - text: "CUSTOM-MARKER"
`)
	resetFlags()
	flagConfig = cfg

	input := "feat: thing\n\n🤖 Generated with Claude Code\nCUSTOM-MARKER\n"
	// Default emoji pattern should NOT be stripped, only CUSTOM-MARKER
	want := "feat: thing\n\n🤖 Generated with Claude Code\n"

	path := writeTemp(t, input)
	if err := runCommitMsg(nil, []string{path}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, path); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
