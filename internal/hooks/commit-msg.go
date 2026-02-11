package hooks

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/devasherr/gitlang/internal/config"
)

func CommitMsg(cfg config.CommitMsg, args []string) error {
	if !cfg.Enabled {
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("commit message file missing")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	var errs []error

	msg := strings.TrimRight(string(data), "\n")
	if cfg.MinLength > 0 {
		if len(msg) < cfg.MinLength {
			errs = append(errs, fmt.Errorf("commit message too short"))
		}
	}

	if cfg.NoTrailingPeriod {
		if strings.HasSuffix(msg, ".") {
			errs = append(errs, fmt.Errorf("commit messge can not end with a period"))
		}
	}

	if len(cfg.ForbiddenWords) > 0 {
		for _, word := range cfg.ForbiddenWords {
			if strings.Contains(msg, word) {
				errs = append(errs, fmt.Errorf("`%s` is forbidden in commit message", word))
			}
		}
	}

	return errors.Join(errs...)
}
