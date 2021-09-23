package main

import (
	"github.com/alexflint/go-arg"
)

func handleProcess(cli *arg.Parser, args argType) {
	switch {
	case args.Process.Dlls != nil:
	case args.Process.Exes != nil:
	case args.Process.Tasks != nil:
	case args.Process.Services != nil:
	}
}

func processDlls(args argType)     {}
func processExes(args argType)     {}
func processTasks(args argType)    {}
func processServices(args argType) {}
