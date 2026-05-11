# Gitlang

A fast, language-agnostic Git hook manager

> funfact: I actually had the project name first and then started thinking:
> “what could I build that fits `gitlang`?”

Gitlang helps enforce clean Git workflows using a single configuration file.

- Enforce branch naming rules
- Protect important branches
- Validate commit messages
- Run commands on staged files
- Prevent unwanted files from being committed
- Keep repositories consistent across teams

---

# Features

## Branch Rules

Protect important branches and enforce naming conventions.

```yaml
branch:
  enabled: true
  protected: [main, test]
  pattern: "^(feature|bugfix|hotfix)/[A-Z]+-[0-9]+"
```

### Examples

✅ Valid branch names:

```txt
feature/PROJ-123-login
bugfix/PROJ-88-auth
hotfix/PROJ-991-panic-fix
```

❌ Invalid branch names:

```txt
login-fix
feature/login
main
```

---

## Pre-Commit Checks

Prevent bad commits before they happen.

```yaml
pre-commit:
  enabled: true
  max_file_size_kb: 5000
  forbidden_extensions: [.log, .tmp, .swp, .exe]
```

### Supported Checks

- Maximum file size
- Forbidden file extensions
- File and folder naming conventions
- Run commands on staged files

---

## Run Commands on Staged Files

Run tools only when matching files are staged.

```yaml
pre-commit:
  run:
    - cmd: "npm run lint"
      match: ["*.js", "*.ts"]

    - cmd: "go vet ./..."
      match: ["*.go"]
```

This keeps hooks fast and avoids unnecessary commands.

---

## Commit Message Validation

Enforce clean and readable commit messages.

```yaml
commit-msg:
  enabled: true
  min_length: 10
  no_trailing_period: true
  forbidden_words: [tmp, stuff, wip]
```

### Example

✅ Valid:

```txt
fix: resolve login timeout issue
```

❌ Invalid:

```txt
wip.
```

---

# Installation

## Download Binary (once)

```bash
sudo curl -L https://github.com/devasherr/gitlang/releases/download/v0.1.0/gitlang-cli-linux-amd64 -o /usr/local/bin/gitlang && sudo chmod +x /usr/local/bin/gitlang
```

---

# Quick Start

## 1. Initialize Gitlang

```bash
gitlang init
```

Creates a default `gitlang.yml` configuration file.

And thats it!! Gitlang automatically validates your git actvities.

---

# Example Configuration

```yaml
branch:
  enabled: true
  protected: [main, test]
  pattern: "^(feature|bugfix|hotfix)/[A-Z]+-[0-9]+"

pre-commit:
  enabled: true
  max_file_size_kb: 5000
  forbidden_extensions: [.log, .tmp, .swp, .exe]

  naming_conventions:
    folder:
      naming: [no_spaces, lowercase]

    file:
      naming: [no_spaces, lowercase]

  run:
    - cmd: "npm run lint"
      match: ["*.js", "*.ts"] # if match is empty it applies to all staged files

    - cmd: "go vet ./..."
      match: ["*.go"]

commit-msg:
  enabled: true
  min_length: 10
  no_trailing_period: true
  forbidden_words: [tmp, stuff, wip]
```
A complete, well-commented example configuration is available at:

[examples/.example.gitlang.yaml](./examples/.example.gitlang.yaml)
