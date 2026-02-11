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

	return errors.Join(errs...)
}
