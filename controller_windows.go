package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/memutils"

	"github.com/audibleblink/lpegopher/args"
	"github.com/audibleblink/lpegopher/collectors"
	"github.com/audibleblink/lpegopher/logerr"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("doCollectCmd")
	log.Info("collection started")

	collectors.InitOutputFiles()

	log.Info("collecting system principals")
	collectors.CreateGroupPrincipals()

	var wg sync.WaitGroup

	files, err := os.ReadDir(args.Collect.Root)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, f := range files {
		if f.IsDir() {
			path := filepath.Join(args.Collect.Root, f.Name())
			log.Debugf("forking collection of %s", path)
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

	wg.Add(1)
	log.Info("collecting autoruns")
	go func() {
		defer wg.Done()
		collectors.Autoruns()
	}()

	wg.Add(1)
	log.Info("collecting processes")
	go func() {
		defer wg.Done()
		collectors.Processes()
	}()

	wg.Wait()
	log.Info("flushing buffers and closing files")
	collectors.FlushAndClose()
	log.Info("collection complete")
	log.Warn(
		"=============================================================================================",
	)
	log.Warn(
		"don't forget to upload/move *.csv to neo4j's `import` directory before running postprocessing",
	)
	log.Warn(
		"=============================================================================================",
	)
	return
}

func getSystem() error {
	pid := argv.GetSystem.PID
	if pid == 0 {
		pid = memutils.PidForName("winlogon.exe")
		logerr.Infof("stealing winlogon token from pid %d", pid)
	}
	return getsystem.InNewProcess(pid, `c:\windows\system32\cmd.exe`, false)
}
