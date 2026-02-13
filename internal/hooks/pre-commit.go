package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/devasherr/gitlang/internal/config"
)

func PreCommit(c config.Config) error {
	if c.Branch.Enabled && len(c.Branch.Protected) > 0 {
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		branch_name := strings.TrimRight(string(output), "\n")

		// regex check
		if c.Branch.Pattern != "" {
			matched, _ := regexp.MatchString(c.Branch.Pattern, branch_name)
			if !matched {
				// temporary form of logging
				fmt.Fprintf(os.Stderr, "gitlang(warn): branch name does not satisfy pattern `%s`\n", c.Branch.Pattern)
			}
		}

		if slices.Contains(c.Branch.Protected, branch_name) {
			return fmt.Errorf("protected branch `%s`. can not make direct commits to this branch", branch_name)
		}
	}

	cfg := c.PreCommit
	if !cfg.Enabled {
		return nil
	}

	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var errs []error

	folderDirSet := map[string]bool{}
	for path := range strings.SplitSeq(strings.TrimRight(string(output), "\n"), "\n") {
		fileInfo, err := os.Stat(path)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		size := float64(fileInfo.Size())
		if cfg.MaxFileSizeKb > 0 {
			if size > cfg.MaxFileSizeKb {
				errs = append(errs, fmt.Errorf("%s of size %.2fKB exceeds max size limit", path, size))
			}
		}

		ext := filepath.Ext(path)

		// linux/mac don't care about extension for exectuable
		// but for the sake of consistency ext will default to .exe
		if ext == "" && fileInfo.Mode()&0111 != 0 {
			ext = ".exe"
		}

		// TODO: add .env and other sensitive entension here
		// since hashmap will be used no need to worry about creating duplicates

		if len(cfg.ForbiddenExtensions) > 0 {
			// TODO: replace with hashmap for faster lookup
			// also for removing duplicate extensions that might have been given
			if ext != "" && slices.Contains(cfg.ForbiddenExtensions, ext) {
				errs = append(errs, fmt.Errorf("%s extension is forbidden", ext))
			}
		}

		errs = append(errs, handleNamingConventions(fileInfo.Name(), cfg.NamingConventions.File))
		folderName := strings.TrimRight(path, fileInfo.Name())
		if !folderDirSet[folderName] {
			folderDirSet[folderName] = true
			errs = append(errs, handleNamingConventions(folderName, cfg.NamingConventions.Folder))

			// feature or bug ??
			// Top/nested/ <- only lowercase letters allowed
			// Top/nested/deeper <- only lowercase letters allowed
			// should it only fail/warn once since Top is the root cause, or warn on every folder path
		}
	}

	return errors.Join(errs...)
}

func checkNoSpace(s string) error {
	if strings.Contains(s, " ") {
		return fmt.Errorf("%s: spaces are not allowed", s)
	}
	return nil
}

func checkLowercase(s string) error {
	for _, c := range s {
		if !unicode.IsLower(c) && unicode.IsLetter(c) {
			return fmt.Errorf("%s: only lowercase letters allowed", s)
		}
	}
	return nil
}

func handleNamingConventions(s string, cfg config.Conventions) error {
	var errs []error

	for _, nameing := range cfg.Naming {
		if nameing == "no_spaces" {
			errs = append(errs, checkNoSpace(s))
		}
		if nameing == "lowercase" {
			errs = append(errs, checkLowercase(s))
		}
	}

	return errors.Join(errs...)
}
