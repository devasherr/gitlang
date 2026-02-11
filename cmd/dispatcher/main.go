package main

import (
	"fmt"
	"os"

	"github.com/devasherr/gitlang/internal/config"
	"github.com/devasherr/gitlang/internal/hooks"
)

const (
	ERROR = "error"
)

func logger(status, msg string) {
	writer := os.Stdout
	if status == "error" {
		writer = os.Stderr
	}

	fmt.Fprintf(writer, "gitlang(%s): %s\n", status, msg)
}

func main() {
	if len(os.Args) < 2 {
		logger(ERROR, "required arguments not passed")
		os.Exit(1)
	}

	body, err := os.ReadFile(".gitlang.yaml")
	if err != nil {
		logger(ERROR, ".gitlang.yaml does not exist")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(body)
	if err != nil {
		logger(ERROR, err.Error())
		os.Exit(1)
	}

	hook := os.Args[1]

	var errs []error

	switch hook {
	case "commit-msg":
		errs = hooks.CommitMsg(cfg.CommitMsg)
	default:
		logger(ERROR, "unknown hook: "+hook)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			logger(ERROR, err.Error())
		}
		os.Exit(1)
	}

	os.Exit(1) // remove this
}
