package main

import (
	"bufio"
	"os"

	"github.com/alexflint/go-arg"
)

type jsonProcessor func([]byte) error

func handleProcess(cli *arg.Parser, args argType) (err error) {
	switch {
	case args.Process.Dlls != nil:
	case args.Process.Exes != nil:
	case args.Process.Tasks != nil:
		processor := newFileProcessor(NewRunnerFromJson)
		err = processor(args)
	case args.Process.Services != nil:
		processor := newFileProcessor(NewRunnerFromJson)
		err = processor(args)
	}
	return
}

func processDlls(args argType)  {}
func processExes(args argType)  {}
func processTasks(args argType) {}

func newFileProcessor(jp jsonProcessor) func(args argType) error {
	return func(args argType) error {
		file, err := os.Open(args.Process.Services.File)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			err = jp(scanner.Bytes())
			if err != nil {
				return err
			}
		}
		return nil

	}
}
