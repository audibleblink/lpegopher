package main

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {

	peProcess := processor.QueryBuilder(processor.CreatePEFromJSON)
	logerr.Info("creating inodes")
	err = peProcess(args.Process.PEs)
	if err != nil {
		return
	}

	runnerProcess := processor.QueryBuilder(processor.CreateRunnerFromJSON)
	logerr.Info("creating runners")
	err = runnerProcess(args.Process.Runners)
	if err != nil {
		return
	}

	logerr.Info("creating inode relationships")
	err = processor.BulkRelateFileTree()
	if err != nil {
		return
	}

	logerr.Info("creating runner relationships")
	err = processor.BulkRelateRunners()
	if err != nil {
		return
	}

	// logerr.Info("creating ACL relationships")
	// peRelBuilder := processor.QueryBuilder(processor.RelatePEs)
	// err = peRelBuilder(args.Process.PEs)
	return

}
