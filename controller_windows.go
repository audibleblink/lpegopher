package main

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/collectors"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	switch {
	case args.Collect.Dlls != nil:
		setPrinter(args.Collect.Dlls.File)
		path := args.Collect.Dlls.Path
		collectors.PEs("*.dll", path)
	case args.Collect.Exes != nil:
		setPrinter(args.Collect.Exes.File)
		path := args.Collect.Exes.Path
		collectors.PEs("*.exe", path)
	case args.Collect.Tasks != nil:
		collectors.Tasks()
	case args.Collect.Services != nil:
		collectors.Services()
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return
}

func setPrinter(outfile string) (err error) {
	if outfile == "stdOut" {
		collectors.Printer = os.Stdout
		return
	}

	absFilePath, err := filepath.Abs(outfile)
	if err != nil {
		return
	}

	f, err := os.OpenFile(absFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	collectors.Printer = f
	return
}
