package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadBinary(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status)
	}

	// create dir if not exist
	if err = os.Mkdir(".gbin", 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = os.Chmod(path, 0755); err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	return err
}

func routeGitHooks() error {
	hooks := []string{"commit-msg", "pre-commit"}
	var errs []error
	for _, hook := range hooks {
		path := ".git/hooks/" + hook
		_ = os.Rename(path+".sample", path)

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer file.Close()

		content := fmt.Sprintf(
			`#!/bin/sh
exec ./.gbin/dispatcher %s "$@"`, hook)
		_, err = file.WriteString(content)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}

func generateConfigFile() error {
	file, err := os.OpenFile(".gitlang.yaml", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	// default config
	content := []byte(`# ==========================================
# This is default generate config file by gitlang cli
# Feel free to change the config as you wish
# ==========================================

# ------------------------------------------
# Branch Rules
# ------------------------------------------
branch:
  enabled: true # If false, branch validation is skipped

  # List of branches that cannot receive direct commits.
  # Commits on these branches will fail.
  protected:
    - main
    - test

  # Regex pattern that branch names must follow (excluding master/main)
  pattern: "^(feature|bugfix|hotfix)/[A-Z]+-[0-9]+"

# ------------------------------------------
# Pre-Commit Rules
# ------------------------------------------
pre-commit:
  enabled: true # If false, pre-commit checks are skipped

  # Maximum allowed file size (in KB) for staged files
  # Fails larger than this will cause commit to fail
  max_file_size_kb: 5000

  # File extendsions that are not allowed to be commited
  forbidden_extensions: 
    - .log
    - .tmp
    - .swp
    # in unix system executable are considered to have .exe extension
    # so to skip exectables just specify .exe
    - .exe

  naming_conventions:
    folder:
      naming:
        - no_spaces # Folder names cannot contain spaces
        - lowercase # Folder names must be lowercase only
    file:
      naming:
        - no_spaces # File names cannot contain spaces
        - lowercase # File names must be lowercase only

  # Command to run before commit
  # These commands only execute if staged files match the given patterns
  run:
    # Run ESLint only if staged JS or TS files exist
    - cmd: "npm run lint"
      match:
        - "*.js"
        - "*.ts"

    # Run go vet only if staged Go files exist
    - cmd: "go vet ./..."
      match:
        - "*.go"

# ------------------------------------------
# Commit Message Rules
# ------------------------------------------
commit-msg:
  enabled: true # If false, commit message validation is skipped
  min_length: 10 # Minimum number of characters required in commit message

  # commit message must NOT end with a period
  # invalid example: "fix login bug."
  no_trailing_period: true

  # Words thats are not allowed inside commit messages
  # Matching is case-insensitive
  forbidden_words:
    - tmp
    - stuff
    - wip
`)
	_, err = file.Write(content)
	return err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments provided")
	}

	switch os.Args[1] {
	case "init":
		_, err := os.ReadDir(".git")
		if err != nil {
			log.Fatal("unable to locate .git, make sure current project is tracked by git")
		}

		if err := downloadBinary("https://github.com/devasherr/gitlang/releases/download/v0.1.0/gitlang-dispatcher-linux-amd64", "./.gbin/dispatcher"); err != nil {
			log.Fatalf("failed to download dispatcher: %s", err.Error())
		}

		// route git hooks to use dispatcher
		if err = routeGitHooks(); err != nil {
			log.Fatalf("failed to reroute git hook to dispatcher: %s", err.Error())
		}

		// make .gitlang.yaml file
		if err = generateConfigFile(); err != nil {
			log.Fatalf("failed to generate default config file: %s", err.Error())
		}
	default:
		log.Fatalf("unknown argument: %s", os.Args[1])
	}
}
