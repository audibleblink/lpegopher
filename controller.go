package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/processor"
	"github.com/audibleblink/pegopher/util"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {
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
			text, err = util.DecodeUTF16(text)
			if err != nil {
				return err
			}

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
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "scanner quit:", err)
		}

		fmt.Println(count)
		return nil

	}
}
