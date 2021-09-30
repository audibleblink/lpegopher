//go:build !windows

package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
)

func doCollectCmd(a args.ArgType, cli *arg.Parser) error {
	return fmt.Errorf("collect functionality is only available on Windows")
}
