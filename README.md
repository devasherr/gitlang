# Gitlang

language-agnostic Git hook manager

> funfact: I actually had the project name first and then started thinking “what could I build that fits `gitlang`?”

Gitlang helps enforce clean Git workflows using a single configuration file.

- Enforce branch naming rules
- Validate commit messages
- Prevent unwanted files from being committed
- Keep repositories consistent across teams

---

# Installation

## Download Binary (once)

```bash
sudo curl -L https://github.com/devasherr/gitlang/releases/download/v0.1.1/gitlang-cli-linux-amd64 -o /usr/local/bin/gitlang && sudo chmod +x /usr/local/bin/gitlang
```

---

# Quick Start

## 1. Initialize Gitlang

```bash
gitlang init
```

Creates a default `.gitlang.yaml` configuration file.

And thats it!! Gitlang automatically validates your git actvities.

---

# Features

### Examples

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
      match: ["*.js", "*.ts"]

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
