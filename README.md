# tsa

Personal git hook utilities. Screens your commits before they're allowed through.

## Install

### mise

```bash
mise use github:natefaerber/tsa
```

Or add to your `.mise.toml`:

```toml
[tools]
"github:natefaerber/tsa" = "latest"
```

### go install

```bash
go install github.com/natefaerber/tsa@latest
```

### GitHub releases

Download pre-built binaries from [Releases](https://github.com/natefaerber/tsa/releases).

### From source

```bash
git clone https://github.com/natefaerber/tsa.git
cd tsa
make install
```

## Stages

tsa organizes hooks into **stages**, each containing one or more **steps**:

```
$ tsa list
pre-commit  Run pre-commit checks on staged changes
  tsa pre-commit [files...]
└── nocommit  Check for !nocommit markers in staged changes

commit-msg  Clean up commit messages
  tsa commit-msg <file>
├── strip-attribution  Remove AI attribution and co-author lines
└── collapse-blanks  Collapse consecutive blank lines and trim leading/trailing blanks
```

## Usage

### As git hooks

Add to your git hooks (e.g. `.git/hooks/pre-commit`):

```bash
#!/bin/sh
tsa pre-commit
```

```bash
#!/bin/sh
tsa commit-msg "$1"
```

### Pre-commit

Checks staged changes for `!nocommit` markers. Optionally scope to specific files:

```bash
tsa pre-commit                    # check all staged files
tsa pre-commit foo.go bar.go      # check only these files
```

### Commit-msg

Cleans up commit messages by stripping AI attribution lines and collapsing blank lines:

```bash
tsa commit-msg .git/COMMIT_EDITMSG
```

## Controlling Steps

Skip or isolate specific steps within a stage:

```bash
tsa commit-msg "$1" --skip=collapse-blanks
tsa commit-msg "$1" --only=strip-attribution
tsa pre-commit --skip=nocommit
```

List available steps for a stage:

```bash
tsa pre-commit --list-steps
tsa commit-msg --list-steps
```

Suppress "skipping X" messages:

```bash
tsa commit-msg "$1" --skip=collapse-blanks --quiet
```

## Configuration

tsa reads an optional config file from `$XDG_CONFIG_HOME/tsa/config.yaml` (defaults to `~/.config/tsa/config.yaml`). Override with `--config`:

```bash
tsa commit-msg "$1" --config=./my-config.yaml
```

### Strip-attribution patterns

When a config file is present, its patterns **replace** the built-in defaults. Each pattern is either a substring (`text`) or a regular expression (`regex`), with optional case-insensitive matching:

```yaml
commit-msg:
  strip-attribution:
    patterns:
      - text: "Generated with"
      - text: "Co-Authored-By: Claude"
        case_insensitive: true
      - regex: "Signed-off-by:.*bot"
        case_insensitive: true
```

### Built-in defaults

When no config file exists, tsa uses these defaults:

| Type | Pattern | Case Insensitive |
|------|---------|:----------------:|
| text | `🤖 Generated with` | no |
| text | `Co-Authored-By: Claude` | yes |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--skip=step1,step2` | | Skip named steps |
| `--only=step1,step2` | | Run only named steps |
| `--quiet` | `-q` | Suppress skip messages |
| `--list-steps` | | List steps for a stage and exit |
| `--config=path` | | Config file path |
| `--version` | | Print version |

## Output formats

```bash
tsa list              # colored tree output
tsa list -f json      # JSON array of stages with steps
```

## Development

```bash
make              # test + build
make test         # run tests with verbose output
make lint         # run golangci-lint
make build        # compile binary with version info
make install      # install to $GOPATH/bin
```

### Releasing

Tag a version to trigger the release workflow:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This runs [goreleaser](https://goreleaser.com/) via GitHub Actions, which builds cross-platform binaries (linux/darwin, amd64/arm64) and creates a GitHub release with checksums and changelog.

## License

MIT
