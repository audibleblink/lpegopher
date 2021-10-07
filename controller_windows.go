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
	defer logerr.ClearContext()

	switch {
	case args.Collect.All:
		all()
	case args.Collect.Dlls != nil:
		out, err := mkOutput(args.Collect.Dlls.File)
		if err != nil {
			return logerr.Wrap(err)
		}
		collectors.PEs(out, "*.dll", args.Collect.Dlls.Path)
	case args.Collect.Exes != nil:
		out, err := mkOutput(args.Collect.Exes.File)
		if err != nil {
			return logerr.Wrap(err)
		}
		collectors.PEs(out, "*.exe", args.Collect.Exes.Path)
	case args.Collect.Tasks != nil:
		out, err := mkOutput(args.Collect.Tasks.File)
		if err != nil {
			return logerr.Wrap(err)
		}
		collectors.Tasks(out)
	case args.Collect.Services != nil:
		out, err := mkOutput(args.Collect.Services.File)
		if err != nil {
			return logerr.Wrap(err)
		}
		collectors.Services(out)
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	return
}

func mkOutput(outfile string) (output io.Writer, err error) {
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
	wg.Add(4)

	go func(file, startPath string) {
		defer wg.Done()
		out, err := mkOutput(file)
		if err != nil {
			return
		}
		collectors.PEs(out, "*.dll", startPath)
	}("dlls.json", `C:\`)

	go func(file, startPath string) {
		defer wg.Done()
		out, err := mkOutput(file)
		if err != nil {
			return
		}
		collectors.PEs(out, "*.exe", startPath)
	}("exes.json", `C:\`)

	go func(file string) {
		defer wg.Done()
		out, err := mkOutput(file)
		if err != nil {
			return
		}
		collectors.Tasks(out)
	}("tasks.json")

	go func(file string) {
		defer wg.Done()
		out, err := mkOutput(file)
		if err != nil {
			return
		}
		collectors.Services(out)
	}("services.json")

	wg.Wait()

}

func getSystem(pid int) error {
	err := getsystem.InNewProcess(argv.GetSystem.PID, `c:\windows\system32\cmd.exe`, false)
	return logerr.DefaultLogger().Context("getsystem").Wrap(err)
}
