package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/processor"
	"github.com/audibleblink/pegopher/util"
)

func doProcessCmd(args args.ArgType) (err error) {
	switch {
	case args.Process.Dlls != nil:
	case args.Process.Exes != nil:
		proc := newFileProcessor(processor.NewExeFromJson)
		err = proc(args.Process.Exes.File)
	case args.Process.Tasks != nil:
		proc := newFileProcessor(processor.NewRunnerFromJson)
		err = proc(args.Process.Tasks.File)
	case args.Process.Services != nil:
		proc := newFileProcessor(processor.NewRunnerFromJson)
		err = proc(args.Process.Services.File)
	}
	return
}

func processDlls(args args.ArgType) {}
func processExes(args args.ArgType) {}

type jsonProcessor func([]byte) error

func newFileProcessor(jp jsonProcessor) func(file string) error {
	return func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		count := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			count += 1
			text := scanner.Bytes()
			text, err = util.DecodeUTF16(text)
			if err != nil {
				return err
			}
			err = jp(text)
			if err != nil {
				return err
			}
		}
		fmt.Println(count)
		return nil

	}
}
