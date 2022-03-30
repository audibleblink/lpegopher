//go:build !windows

package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/lpegopher/args"
	"github.com/audibleblink/lpegopher/logerr"
)

func doCollectCmd(a args.ArgType, cli *arg.Parser) error {
	return fmt.Errorf("collect functionality is only available on Windows")
}

func getSystem() error {
	return logerr.Wrap(fmt.Errorf("only available on Windows"))
}
