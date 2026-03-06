package cmd

import (
	"fmt"
	"testing"
)

// fakeGitClient returns a canned diff response.
type fakeGitClient struct {
	diff string
	err  error
}

func (f fakeGitClient) StagedDiff(files []string) (string, error) {
	return f.diff, f.err
}

func withFakeGit(t *testing.T, client GitClient) {
	t.Helper()
	orig := git
	git = client
	t.Cleanup(func() { git = orig })
}

func TestCheckNoCommitClean(t *testing.T) {
	withFakeGit(t, fakeGitClient{diff: "+func hello() {}\n"})
	resetFlags()

	if err := checkNoCommit(nil); err != nil {
		t.Fatal(err)
	}
}

func TestCheckNoCommitWithFiles(t *testing.T) {
	var captured []string
	withFakeGit(t, fakeGitClientFunc(func(files []string) (string, error) {
		captured = files
		return "", nil
	}))
	resetFlags()

	if err := checkNoCommit([]string{"foo.go", "bar.go"}); err != nil {
		t.Fatal(err)
	}
	if len(captured) != 2 || captured[0] != "foo.go" || captured[1] != "bar.go" {
		t.Errorf("expected [foo.go bar.go], got %v", captured)
	}
}

func TestCheckNoCommitGitError(t *testing.T) {
	withFakeGit(t, fakeGitClient{err: fmt.Errorf("git failed")})
	resetFlags()

	err := checkNoCommit(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "git failed" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunPreCommitSkipsNocommit(t *testing.T) {
	withFakeGit(t, fakeGitClient{diff: "// !nocommit\n"})
	resetFlags()
	flagSkip = []string{"nocommit"}
	flagQuiet = true

	// With nocommit skipped, even a diff containing !nocommit should pass
	if err := runPreCommit(nil, nil); err != nil {
		t.Fatal(err)
	}
}

// fakeGitClientFunc adapts a function to the GitClient interface.
type fakeGitClientFunc func(files []string) (string, error)

func (f fakeGitClientFunc) StagedDiff(files []string) (string, error) {
	return f(files)
}
