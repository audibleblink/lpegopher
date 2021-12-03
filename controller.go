package main

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("postprocessing")

	log.Info("creating file and principal nodes")
	err = processor.InsertAllNodes()
	if err != nil {
		return
	}

	log.Info("creating runner nodes")
	err = processor.InsertAllRunners()
	if err != nil {
		return
	}

	log.Info("creating filetree relationships")
	err = processor.BulkRelateFileTree()
	if err != nil {
		return
	}

	log.Info("creating ownership relationships")
	err = processor.RelateOwnership()
	if err != nil {
		return
	}

	log.Info("creating runner relationships")
	err = processor.BulkRelateRunners()
	if err != nil {
		return
	}

	log.Info("creating imports relationships")
	err = processor.RelateDependecies()
	if err != nil {
		return
	}

	log.Info("creating ACL relationships")
	err = processor.RelateACLs()

	log.Info("postprocessing complete")
	return
}
