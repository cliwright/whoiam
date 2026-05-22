# Configuration

## Config Files

`whoiam` uses YAML config files to map account names to AWS account IDs.

| Scope | Location | Purpose |
|-------|----------|---------|
| Global | `~/.whoiam/whoiam.yaml` | Shared across all projects |
| Project-local | `.whoiam/whoiam.yaml` | Per-project overrides, committable to version control |

Create them with `whoiam init` (project-local) or `whoiam init --global` (global).

## Config Format

```yaml
accounts:
  production: "123456789012"
  staging: "210987654321"
  development: "345678901234"
```

Each key is a name you choose; the value is the 12-digit AWS account ID.

## Config Merging

When both a global and a project-local config exist, `whoiam` merges them. Project-local definitions take precedence — if the same account name appears in both files, the local value wins.

Run `whoiam config` to see the merged result with source labels:

```
+---------------+----------------+--------+
| ACCOUNT NAME  | ACCOUNT NUMBER | SOURCE |
+---------------+----------------+--------+
| production    | 123456789012   | local  |
| staging       | 210987654321   | global |
+---------------+----------------+--------+
```

## Session State (Expected Environment)

`whoiam set` and `whoiam exec`/`whoiam validate` use a session file to track which account you expect to be authenticated with.

| Scope | Location |
|-------|----------|
| Project-local | `.whoiam/expected-env` |
| Global | `~/.whoiam/expected-env` |

`whoiam init` automatically creates a `.whoiam/.gitignore` that excludes `expected-env` — this file is personal session state and should not be committed.

## Expected Environment Resolution Order

When `whoiam exec` or `whoiam validate` resolves which account to check against, it uses this priority:

1. `--env <name>` flag (highest priority)
2. `WHOIAM_EXPECTED_ENV` environment variable
3. `.whoiam/expected-env` (project-local session file)
4. `~/.whoiam/expected-env` (global session file)

If none of these are set, the command fails with an error asking you to use `--env` or run `whoiam set`.

## Environment Variable

You can override the expected environment without touching any files:

```sh
WHOIAM_EXPECTED_ENV=production whoiam validate
WHOIAM_EXPECTED_ENV=staging whoiam exec -- terraform plan
```

This is useful in CI pipelines where you want to inject the expected account from a secret or pipeline variable.

## AWS Credentials

`whoiam` uses the AWS SDK for Go and reads credentials from the standard locations: environment variables, shared credentials file (`~/.aws/credentials`), IAM instance profiles, and so on. It does not manage credentials itself.
