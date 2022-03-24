package main

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/memutils"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/logerr"
)

func doCollectCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("doCollectCmd")
	log.Info("collection started")

	collectors.InitOutputFiles()

	var wg sync.WaitGroup

	files, err := ioutil.ReadDir(args.Collect.Path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, f := range files {
		if f.IsDir() {
			path := filepath.Join(args.Collect.Path, f.Name())
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
	collectors.FlashAndClose()
	log.Info("collection complete")
	log.Warn("=============================================================================================")
	log.Warn("don't forget to upload/move *.csv to neo4j's `import` directory before running postprocessing")
	log.Warn("=============================================================================================")
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
