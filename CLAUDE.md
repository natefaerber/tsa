# CLAUDE.md

## Overview

**tsa** is a personal git hook utility written in Go. It provides a stage/step architecture for running checks and transformations during git hooks, with configurable patterns via YAML.

## Project Structure

```
.
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Cobra root command, global flags
│   ├── pre_commit.go       # pre-commit stage (nocommit check)
│   ├── commit_msg.go       # commit-msg stage (strip-attribution, collapse-blanks)
│   ├── list.go             # `tsa list` command, stage registry
│   ├── steps.go            # Step runner, filtering (--skip/--only), StepInfo type
│   ├── config.go           # YAML config loading, pattern compilation (text/regex)
│   ├── git.go              # GitClient interface (abstracts git for testability)
│   ├── color.go            # ANSI color helpers, auto-disabled when not a TTY
│   ├── *_test.go           # Tests alongside source files
├── Makefile                # build, test, lint, fmt, vet, install, clean
├── .github/workflows/
│   └── ci.yaml             # GitHub Actions: test, lint, cross-compile
└── go.mod
```

## Architecture

### Stages and Steps

Each git hook is a **stage** (cobra subcommand) containing one or more **steps**. Steps can be included/excluded at runtime via `--skip` and `--only` flags.

- Stages are registered in `cmd/list.go` (the `stages` slice)
- Step metadata lives in `StepInfo` structs (name + description)
- Step execution uses `RunSteps()` in `cmd/steps.go`

### Adding a New Stage

1. Create `cmd/<hook_name>.go` with a cobra command
2. Define `StepInfo` slice for the step metadata
3. In `RunE`, check `PrintStepInfo()` first, then build `[]Step` and call `RunSteps()`
4. Register the stage in `cmd/list.go`'s `stages` slice

### Adding a New Step to an Existing Stage

1. Add a `StepInfo` entry to the stage's metadata slice
2. Add the `Step` to the `[]Step` in the stage's `RunE`

### Config

- Location: `$XDG_CONFIG_HOME/tsa/config.yaml` (override with `--config`)
- When present, config patterns **replace** built-in defaults
- Pattern types: `text` (substring), `regex`, both support `case_insensitive: true`
- Config struct is in `cmd/config.go`

### Git Abstraction

`cmd/git.go` defines a `GitClient` interface. Production uses `execGitClient` (shells out to git). Tests inject fakes via the package-level `git` variable.

## Development

### Commands

```bash
make              # test + build
make test         # go test ./... -v
make build        # build with version/commit ldflags
make lint         # golangci-lint (must be installed)
make vet          # go vet
make fmt          # gofmt
make install      # go install
make clean        # remove binary
```

### Testing

- Tests are in `cmd/*_test.go`, run with `make test`
- Git operations are faked via `GitClient` interface — no real repo needed
- Config tests use temp files with `--config` flag override
- Global flags (`flagSkip`, `flagOnly`, `flagQuiet`, `flagConfig`) must be reset between tests via `resetFlags()`

### Build

Version info is injected via ldflags:

```bash
go build -ldflags "-s -w -X github.com/natefaerber/tsa/cmd.Version=v1.0.0 -X github.com/natefaerber/tsa/cmd.Commit=abc1234"
```

## Conventions

- YAML files use `.yaml` extension (not `.yml`)
- No dependencies on git libraries — system git via `GitClient` interface
- Colors auto-disable when stdout is not a TTY
- `--only` takes precedence over `--skip` if both are set
- Steps print "skipping X" to stderr unless `--quiet`
