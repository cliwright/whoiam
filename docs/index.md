[![Sketch fonts](https://see.fontimg.com/api/rf5/BWWo5/YTZkZTMxNDlhNDEwNDZhZmFiZThhODFhNjA5N2U3NTgub3Rm/d2hvaWFt/typo-draft-demo.png?r=fs&h=250&w=2000&fg=5A3922&bg=FFFFFF&tb=1&s=125)](https://www.fontspace.com/typo-draft-font-f41179)

`whoiam` is a CLI tool that prevents fat-finger deployments to the wrong AWS account — because nothing ruins your morning quite like realising you just ran `terraform apply` against production.

<div align="center"><a href="https://imgflip.com/i/ashxye"><img src="https://i.imgflip.com/ashxye.jpg" title="made at imgflip.com"/></a></div>

You know that sinking feeling. You get dizzy and the walls start closing in on you. Was your session pointed to dev... or prod?
You know your team shouldn't have local production credentials, but hey... startups. We've all been there.

This has happened to me, and teams I have worked on more times than I care to admit. And that's why I built `whoiam`.

---

## How it works

Most tools focus on getting credentials *into* your shell. `whoiam` asks a different question: **are those credentials pointing at the right account?**

Before any command runs, `whoiam` calls `sts:GetCallerIdentity` and compares the result against the account you declared you expected to be on. If they don't match, it exits immediately — before a single byte of infrastructure changes.

It doesn't store credentials, manage profiles, or replace anything you already use. It sits in front of your existing workflow as a single verification step.

---

## Works with whatever you already use

`whoiam` is credential-agnostic. It works with `aws-vault`, AWS SSO, raw `~/.aws/credentials`, instance profiles — anything the AWS SDK can resolve.

**With aws-vault:**
```sh
aws-vault exec production -- whoiam exec -- terraform apply
```

**With AWS SSO / profiles — set `AWS_PROFILE` and whoiam picks it up automatically if the profile name matches an account in your config:**
```sh
AWS_PROFILE=production whoiam exec -- terraform apply
```

**With nothing special — just environment variables or instance credentials:**
```sh
whoiam exec --env production -- terraform apply
```

`aws-vault` ensures you *have* credentials. `whoiam` ensures they're pointing at the *right account*.

---

## Features

- **Fail fast, before anything runs** — exits non-zero immediately on account mismatch, so you never get halfway through a deployment on the wrong environment
- **Works with any credential source** — aws-vault, SSO, instance profiles, environment variables; if the AWS SDK can see it, whoiam can verify it
- **AWS_PROFILE support** — if your profile name matches an account in your config, whoiam uses it automatically; no extra flags needed
- **Pin an expected environment per project** — `whoiam set production` saves your intent so you don't repeat `--env` on every command
- **Shareable account map** — commit `.whoiam/whoiam.yaml` to your repo so the whole team uses the same account IDs; personal session state stays out of git
- **CI and Taskfile friendly** — `whoiam validate` exits non-zero on mismatch, making it a drop-in pre-flight check for any pipeline

---

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

---

## Documentation

- [Installation](installation.md) — Homebrew, binary download, or build from source
- [Usage](usage.md) — all commands with examples
- [Configuration](configuration.md) — config files, merging, session state, and environment variables
