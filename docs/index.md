[![Sketch fonts](https://see.fontimg.com/api/rf5/BWWo5/YTZkZTMxNDlhNDEwNDZhZmFiZThhODFhNjA5N2U3NTgub3Rm/d2hvaWFt/typo-draft-demo.png?r=fs&h=250&w=2000&fg=5A3922&bg=FFFFFF&tb=1&s=125)](https://www.fontspace.com/typo-draft-font-f41179)

`whoiam` is a CLI tool that prevents accidental deployments to the wrong AWS account. Before running a command, it verifies that your current credentials match the account you expect — protecting you from "fat finger" mistakes when working across multiple environments.

## Features

- **Account verification** — assert that current credentials match a named account before running any command
- **Safe exec** — wrap any command with `whoiam exec` to fail fast if the wrong account is active
- **Session state** — use `whoiam set` to pin an expected environment for the current project or globally, avoiding repetitive `--env` flags
- **Config merging** — combine a global account list with per-project overrides; local definitions take precedence
- **Pre-flight checks** — `whoiam validate` exits non-zero on mismatch, suitable for Taskfiles, CI pipelines, and `mise` hooks
- **Identity display** — run `whoiam` with no arguments to see your current AWS account and ARN

## Quick Start

```sh
# 1. Initialize a project-local config
whoiam init

# 2. Edit .whoiam/whoiam.yaml to add your account mappings
# accounts:
#   production: "123456789012"
#   staging:    "210987654321"

# 3. Pin the expected environment for this project
whoiam set production

# 4. Run commands safely — whoiam verifies the account first
whoiam exec -- terraform apply
```

## Documentation

- [Installation](installation.md) — Homebrew, binary download, or build from source
- [Usage](usage.md) — all commands with examples
- [Configuration](configuration.md) — config files, merging, session state, and environment variables
