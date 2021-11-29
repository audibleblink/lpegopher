package main

import (
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/rpcls/pkg/procs"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("doCollectCmd")
	log.Info("collection started")

	var wg sync.WaitGroup
	wg.Add(2)

	go func(startPath string) {
		defer wg.Done()
		collectors.PEs(startPath)
	}(args.Collect.Path)

	go func() {
		defer wg.Done()
		collectors.Tasks()
		collectors.Services()
	}()

	wg.Wait()

	log.Info("collection complete")
	return
}

func getSystem() error {

	pid := argv.GetSystem.PID
	if pid == 0 {
		pid = procs.PidForName("winlogon.exe")
		logerr.Infof("stealing winlogon token from pid %d", pid)
	}
	return getsystem.InNewProcess(pid, `c:\windows\system32\cmd.exe`, false)
}
