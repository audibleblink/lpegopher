package main

import (
	"io/ioutil"
	"path/filepath"
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

	files, err := ioutil.ReadDir(args.Collect.Path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, f := range files {
		if f.IsDir() {
			path := filepath.Join(args.Collect.Path, f.Name())
			log.Infof("forking collection of %s", path)
			wg.Add(1)
			go func(startPath string) {
				defer wg.Done()
				collectors.PEs(startPath)
			}(path)

		}
	}

	wg.Add(1)
	log.Info("collecting tasks")
	go func() {
		defer wg.Done()
		collectors.Tasks()
	}()

	wg.Add(1)
	log.Info("collecting services")
	go func() {
		defer wg.Done()
		collectors.Services()
	}()

	wg.Wait()
	log.Info("flushing buffers and closing files")
	collectors.FlashAndClose()
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
