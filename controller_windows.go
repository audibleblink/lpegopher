package main

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/processor"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	switch {
	case args.Collect.Dlls != nil:
	case args.Collect.Exes != nil:
		proc := newFileProcessor(processor.NewExeFromJson)
		err = proc(args.Process.Exes.File)
	case args.Collect.Tasks != nil:
		proc := newFileProcessor(processor.NewRunnerFromJson)
		err = proc(args.Process.Tasks.File)
	case args.Collect.Services != nil:
		proc := newFileProcessor(processor.NewRunnerFromJson)
		err = proc(args.Process.Services.File)
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return
}
