package main

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {

	peProcess := processor.QueryBuilder(processor.CreatePEFromJSON)
	logerr.Info("processing pes")
	err = peProcess(args.Process.PEs)
	if err != nil {
		return
	}

	runnerProcess := processor.QueryBuilder(processor.CreateRunnerFromJSON)
	logerr.Info("processing runners")
	err = runnerProcess(args.Process.Runners)
	if err != nil {
		return
	}

	logerr.Info("creating pe relationships")
	peRelBuilder := processor.QueryBuilder(processor.RelatePEs)
	err = peRelBuilder(args.Process.PEs)

	logerr.Info("creating runner relationships")
	runnerRelBuilder := processor.QueryBuilder(processor.RelateRunners)
	err = runnerRelBuilder(args.Process.PEs)
	return

	return
}
