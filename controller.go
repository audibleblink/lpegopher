package main

import (
	"bufio"
	"os"

	"github.com/alexflint/go-arg"
)

func handleProcess(cli *arg.Parser, args argType) (err error) {
	switch {
	case args.Process.Dlls != nil:
	case args.Process.Exes != nil:
	case args.Process.Tasks != nil:
	case args.Process.Services != nil:
		err = processServices(args)
	}
	return
}

func processDlls(args argType)  {}
func processExes(args argType)  {}
func processTasks(args argType) {}

func processServices(args argType) (err error) {
	file, err := os.Open(args.Process.Services.File)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		runner := &Runner{}
		runner, err = NewRunnerFromJson(scanner.Bytes())
		if err != nil {
			return
		}
		err = runner.Save()
		if err != nil {
			return
		}
	}
	return
}
