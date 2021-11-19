package main

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {

	runnerFile := args.PostProcess.Runners

	if runnerFile != "" {
		runnerProcess := processor.QueryBuilder(processor.CreateRunnerFromJSON)
		logerr.Info("creating runners")
		err = runnerProcess(runnerFile)
		if err != nil {
			return
		}
	}

	logerr.Info("creating filetree relationships")
	err = processor.BulkRelateFileTree()
	if err != nil {
		return
	}

	logerr.Info("creating runner relationships")
	err = processor.BulkRelateRunners()
	if err != nil {
		return
	}

	return

}
