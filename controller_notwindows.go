//go:build !windows

package main

import (
	"fmt"

	"github.com/audibleblink/pegopher/args"
)

func doCollectCmd(a args.ArgType) error {
	return fmt.Errorf("collect functionality is only available on Windows")
}
