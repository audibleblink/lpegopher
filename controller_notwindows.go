//go:build !windows

package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
)

func doCollectCmd(cli *arg.Parser, a args.ArgType) error {
	return fmt.Errorf("collect functionality is only available on Windows")
}
