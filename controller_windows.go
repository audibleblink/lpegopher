package main

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/logerr"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	logerr.Context("doCollectCmd")
	logerr.Info("collection started")
	defer logerr.ClearContext()

	switch {
	case args.Collect.All:
		all()
	case args.Collect.PEs != nil:
		out, err := setOutputWriter(args.Collect.PEs.File)
		if err != nil {
			return logerr.Wrap(err)
		}

		collectors.PEs(out, args.Collect.PEs.Path)
	case args.Collect.Runners != nil:
		out, err := setOutputWriter(args.Collect.Runners.File)
		if err != nil {
			return logerr.Wrap(err)
		}

		collectors.Tasks(out)
		collectors.Services(out)
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return
}

func setOutputWriter(outfile string) (output io.Writer, err error) {
	if outfile == "stdOut" {
		output = os.Stdout
		return
	}

	absFilePath, err := filepath.Abs(outfile)
	if err != nil {
		err = logerr.Wrap(err)
		return
	}

	return os.OpenFile(absFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	}("pes.json", `C:\`)

	go func(file string) {
		defer wg.Done()
		out, err := setOutputWriter(file)
		if err != nil {
			return
		}
		collectors.Tasks(out)
		collectors.Services(out)
	}("runners.json")

	wg.Wait()
}

func getSystem(pid int) error {
	err := getsystem.InNewProcess(argv.GetSystem.PID, `c:\windows\system32\cmd.exe`, false)
	return logerr.Add("getsystem").Wrap(err)
}
