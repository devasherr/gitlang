package main

import (
	"fmt"
	"os"
	"strings"

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

	var errs error

	switch hook {
	case "pre-commit":
		errs = hooks.PreCommit(cfg.PreCommit)
	case "commit-msg":
		os.Exit(1) // WARN: remove this
		errs = hooks.CommitMsg(cfg.CommitMsg, os.Args[2:])
	default:
		logger(ERROR, "unknown hook: "+hook)
		os.Exit(1)
	}

	if errs != nil {
		for err := range strings.SplitSeq(errs.Error(), "\n") {
			logger(ERROR, err)
		}
		os.Exit(1)
	}

	os.Exit(1) // WARN: REMOVE THIS
}
