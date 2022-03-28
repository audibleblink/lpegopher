package main

import (
	"net/http"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {
	log := logerr.Add("postprocessing")

	if args.PostProcess.HTTP != "" {
		log.Info("starting fileserver")
		fileServer := serveFiles(args.PostProcess.HTTP)
		defer fileServer.Close()
	}

	log.Info("creating file and principal nodes")
	err = processor.InsertAllNodes(args.PostProcess.HTTP)
	if err != nil {
		return
	}

	log.Info("creating runner nodes")
	err = processor.InsertAllRunners(args.PostProcess.HTTP)
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
	err = processor.RelateDependecies(args.PostProcess.HTTP)
	if err != nil {
		return
	}

	log.Info("creating ACL relationships")
	err = processor.RelateACLs(args.PostProcess.HTTP)

	log.Info("postprocessing complete")
	return
}

func serveFiles(server string) *http.Server {
	log := logerr.Add("fileserver")

	srv := &http.Server{Addr: server}
	srv.Handler = http.FileServer(http.Dir("."))

	log.Infof("http server starting on %s", server)
	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
