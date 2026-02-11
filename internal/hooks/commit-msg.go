package hooks

import (
	"fmt"
	"os"

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

	msg := string(data)
	if cfg.MinLength > 0 {
		if len(msg) < cfg.MinLength {
			return fmt.Errorf("commit message too short")
		}
	}

	return nil
}
