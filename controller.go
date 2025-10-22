package main

import (
	"net/http"

	"github.com/alexflint/go-arg"

	"github.com/audibleblink/logerr"
	"github.com/audibleblink/lpegopher/args"
	"github.com/audibleblink/lpegopher/processor"
)

func doProcessCmd(args args.ArgType, cli *arg.Parser) (err error) {
	_ = cli
	log := logerr.Add("postprocessing")

	if args.Process.HTTP != "" {
		log.Info("starting fileserver")
		fileServer := serveFiles(args.Process.HTTP)
		defer fileServer.Close()
	}

	log.Info("creating file and principal nodes")
	err = processor.InsertAllNodes(args.Process.HTTP)
	if err != nil {
		return
	}

	log.Info("creating runner nodes")
	err = processor.InsertAllRunners(args.Process.HTTP)
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
	err = processor.RelateDependecies(args.Process.HTTP)
	if err != nil {
		return
	}

	log.Info("creating ACL relationships")
	err = processor.RelateACLs(args.Process.HTTP)
	if err != nil {
		return
	}

	log.Info("creating user/group memberships")
	err = processor.RelateMembership()
	if err != nil {
		return
	}

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
