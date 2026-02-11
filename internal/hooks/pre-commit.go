package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/devasherr/gitlang/internal/config"
)

func PreCommit(cfg config.PreCommit) error {
	if !cfg.Enabled {
		return nil
	}

	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var errs []error

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

		if len(cfg.ForbiddenExtensions) > 0 {
			if ext != "" && slices.Contains(cfg.ForbiddenExtensions, ext) {
				errs = append(errs, fmt.Errorf("%s extension is forbidden", ext))
			}
		}
	}

	return errors.Join(errs...)
}
