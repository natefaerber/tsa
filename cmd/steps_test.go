package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestRunStepsAllRun(t *testing.T) {
	var ran []string
	steps := []Step{
		{Name: "a", Run: func() error { ran = append(ran, "a"); return nil }},
		{Name: "b", Run: func() error { ran = append(ran, "b"); return nil }},
	}

	if err := RunSteps(steps, nil, nil, false); err != nil {
		t.Fatal(err)
	}
	if len(ran) != 2 || ran[0] != "a" || ran[1] != "b" {
		t.Errorf("expected [a b], got %v", ran)
	}
}

func TestRunStepsSkip(t *testing.T) {
	var ran []string
	steps := []Step{
		{Name: "a", Run: func() error { ran = append(ran, "a"); return nil }},
		{Name: "b", Run: func() error { ran = append(ran, "b"); return nil }},
	}

	if err := RunSteps(steps, []string{"a"}, nil, true); err != nil {
		t.Fatal(err)
	}
	if len(ran) != 1 || ran[0] != "b" {
		t.Errorf("expected [b], got %v", ran)
	}
}

func TestRunStepsOnly(t *testing.T) {
	var ran []string
	steps := []Step{
		{Name: "a", Run: func() error { ran = append(ran, "a"); return nil }},
		{Name: "b", Run: func() error { ran = append(ran, "b"); return nil }},
	}

	if err := RunSteps(steps, nil, []string{"b"}, true); err != nil {
		t.Fatal(err)
	}
	if len(ran) != 1 || ran[0] != "b" {
		t.Errorf("expected [b], got %v", ran)
	}
}

func TestRunStepsOnlyOverridesSkip(t *testing.T) {
	var ran []string
	steps := []Step{
		{Name: "a", Run: func() error { ran = append(ran, "a"); return nil }},
		{Name: "b", Run: func() error { ran = append(ran, "b"); return nil }},
	}

	// only takes precedence: run only "a", even though skip also says "a"
	if err := RunSteps(steps, []string{"a"}, []string{"a"}, true); err != nil {
		t.Fatal(err)
	}
	if len(ran) != 1 || ran[0] != "a" {
		t.Errorf("expected [a], got %v", ran)
	}
}

func TestRunStepsSkipPrintsMessage(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	steps := []Step{
		{Name: "a", Run: func() error { return nil }},
	}
	_ = RunSteps(steps, []string{"a"}, nil, false)

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	if got := buf.String(); got != "skipping a\n" {
		t.Errorf("expected 'skipping a\\n', got %q", got)
	}
}

func TestRunStepsQuietSuppressesMessage(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	steps := []Step{
		{Name: "a", Run: func() error { return nil }},
	}
	_ = RunSteps(steps, []string{"a"}, nil, true)

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	if got := buf.String(); got != "" {
		t.Errorf("expected no output, got %q", got)
	}
}

func TestRunStepsStopsOnError(t *testing.T) {
	var ran []string
	steps := []Step{
		{Name: "a", Run: func() error { ran = append(ran, "a"); return fmt.Errorf("fail") }},
		{Name: "b", Run: func() error { ran = append(ran, "b"); return nil }},
	}

	err := RunSteps(steps, nil, nil, true)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(ran) != 1 {
		t.Errorf("expected only [a] to run, got %v", ran)
	}
}
