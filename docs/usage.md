# Usage

## Show Current Identity

Display the AWS caller identity for your current credentials:

```sh
whoiam
```

Output includes the account name (if configured), account ID, and ARN.

---

## Initialize Config

Create a config file before adding account mappings.

**Project-local** (creates `.whoiam/whoiam.yaml` in the current directory):

```sh
whoiam init
```

This also creates a `.whoiam/.gitignore` that excludes `expected-env` from version control. Commit `.whoiam/whoiam.yaml` to share account mappings with your team.

**Global** (creates `~/.whoiam/whoiam.yaml`):

```sh
whoiam init --global
```

---

## Set the Expected Environment

Tell `whoiam` which account you expect to be authenticated with. This saves you from passing `--env` on every command.

```sh
whoiam set production           # write to .whoiam/expected-env (project-local)
whoiam set --global staging     # write to ~/.whoiam/expected-env (global, all projects)
whoiam set                      # clear the local expected environment
whoiam set --global             # clear the global expected environment
```

The local setting takes precedence over the global one. Use `whoiam status` to see what is currently set.

---

## Check Status

Show the current expected environment and whether you are authenticated:

```sh
whoiam status
```

Example output:

```
Expected env: production (local)
Authenticated:  yes
Account:        production (123456789012)
ARN:            arn:aws:iam::123456789012:role/my-role
```

---

## Validate

Assert that the current AWS credentials match the expected account. Exits non-zero on mismatch, making it suitable as a pre-flight check.

```sh
whoiam validate                    # uses the expected env set by 'whoiam set'
whoiam validate --env production   # explicit environment
```

Use this in Taskfiles, CI pipelines, or `mise` hooks to fail fast before a destructive operation:

```yaml
# Taskfile.yml example
tasks:
  deploy:
    cmds:
      - whoiam validate --env production
      - terraform apply
```

---

## Exec

Verify the expected account and then run a command. If the account matches, the command runs; if not, it exits with an error before anything executes.

```sh
whoiam exec --env production -- terraform apply
whoiam exec --env staging -- aws s3 ls
```

If no command is provided, `whoiam exec` opens an interactive subshell with the account already verified:

```sh
whoiam exec --env production
# Opens a subshell — type 'exit' to return to the parent shell
```

If you have already set the expected environment with `whoiam set`, you can omit `--env`:

```sh
whoiam set production
whoiam exec -- terraform apply
```

---

## View Config

Print the effective merged configuration (global + project-local), showing the source of each account:

```sh
whoiam config
```

---

## Command Reference

| Command | Description |
|---------|-------------|
| `whoiam` | Show current AWS caller identity |
| `whoiam init` | Initialize project-local config |
| `whoiam init --global` | Initialize global config |
| `whoiam set [env]` | Set (or clear) the local expected environment |
| `whoiam set --global [env]` | Set (or clear) the global expected environment |
| `whoiam status` | Show expected env and current auth state |
| `whoiam validate [--env <env>]` | Assert current account matches expected |
| `whoiam exec [--env <env>] [-- cmd]` | Verify account then run command or open subshell |
| `whoiam config` | Print merged config with sources |
