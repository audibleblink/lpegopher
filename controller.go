package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {
	switch {

	case args.Process.PEs != nil:
		fileProcessor := newFileProcessor(processor.CreatePEFromJSON)

		logerr.Info("processing pEs")
		err = fileProcessor(args.Process.PEs.File)
		if err != nil {
			return
		}

		logerr.Info("creating pe relationships")
		err = processor.RelatePEs(args.Process.PEs.File)
		return

	case args.Process.Runners != nil:
		fileProcessor := newFileProcessor(processor.CreateRunnerFromJSON)
		logerr.Info("processing runners")
		err = fileProcessor(args.Process.Runners.File)
		if err != nil {
			return
		}

		logerr.Info("creating runners' relationships")
		err = processor.RelateRunners(args.Process.Runners.File)
		return

	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return
}

type jsonProcessor func([]byte) error

func newFileProcessor(jp jsonProcessor) func(file string) error {
	return func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		count := 0
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 8*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			count += 1
			text := scanner.Bytes()
			err = jp(text)
			if err != nil {
				switch err.(type) {
				case *json.SyntaxError:
					fmt.Fprintf(os.Stderr, "malformed json at line %d", count)
					continue
				default:
					return err
				}
			}
			if count%500 == 0 {
				logerr.Infof("Checkpoint. Processed %d lines", count)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "scanner quit:", err)
		}

		logerr.Infof("Done. Processed %d lines", count)
		return nil

	}
}
