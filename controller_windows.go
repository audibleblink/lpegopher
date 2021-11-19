package main

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/logerr"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("doCollectCmd")
	log.Info("collection started")

	collectors.DoJSON = args.Collect.JSON

	switch {
	case args.Collect.All:
		all()
	case args.Collect.PEs != nil:
		out, err := setOutputWriter(args.Collect.PEs.File)
		if err != nil {
			return log.Wrap(err)
		}

		collectors.PEs(out, args.Collect.PEs.Path)
	case args.Collect.Runners != nil:
		out, err := setOutputWriter(args.Collect.Runners.File)
		if err != nil {
			return log.Wrap(err)
		}

		collectors.Tasks(out)
		collectors.Services(out)
	default:
		cli.WriteHelp(os.Stderr)
		log.Fatal("you must choose a post-processing task")
	}

	log.Info("collection complete")
	return
}

func setOutputWriter(outfile string) (output *bufio.Writer, err error) {
	if outfile == "stdOut" {
		output = bufio.NewWriter(os.Stdout)
		return
	}

	absFilePath, err := filepath.Abs(outfile)
	if err != nil {
		err = logerr.Wrap(err)
		return
	}

	f, err := os.OpenFile(absFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return bufio.NewWriter(f), err
}

func all() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(file, startPath string) {
		defer wg.Done()
		out, err := setOutputWriter(file)
		if err != nil {
			return
		}
		collectors.PEs(out, startPath)
	}("inodes.csv", `C:\`)

	go func(file string) {
		defer wg.Done()
		out, err := setOutputWriter(file)
		if err != nil {
			return
		}
		collectors.Tasks(out)
		collectors.Services(out)
	}("runners.csv")

	wg.Wait()
}
