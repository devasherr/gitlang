package hooks

import (
	"fmt"

	"github.com/devasherr/gitlang/internal/config"
)

func CommitMsg(cfg config.CommitMsg) []error {
	fmt.Println("commit-msg hook ... ")
	return nil
}
