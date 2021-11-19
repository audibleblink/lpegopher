package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	argv = args.ArgType{}
	cli  = arg.MustParse(&argv)
)

func init() {
	l := &logerr.Logger{
		Level:  logerr.LogLevelInfo,
		Output: os.Stderr,
	}

	if argv.Debug {
		l.Level = logerr.LogLevelDebug
	}

	l.Context("lpegopher").SetAsGlobal()
}

func main() {
	args.Args = argv

	switch {
	case argv.GetSystem != nil:
		err := getSystem()
		if err != nil {
			logerr.Fatalf("getsystem failed: %v", err)
		}

	case argv.Collect != nil:
		err := doCollectCmd(argv, cli)
		if err != nil {
			logerr.Fatalf("collection failed: %v", err)
		}
	case argv.PostProcess != nil:
		dbInit()
		p := argv.PostProcess
		if p.All != nil {
			dbDrop(p.All.Drop)
		} else if p.Runners != nil {
			dbDrop(p.Runners.Drop)
		} else if p.PEs != nil {
			dbDrop(p.PEs.Drop)
		} else if p.Relationships != nil {
			dbDrop(p.Relationships.Drop)
		} else {
			cli.WriteHelp(os.Stderr)
			logerr.Fatal("you must choose a post-processing task")
		}

		// dbDrop(true)
		err := doProcessCmd(argv, cli)
		if err != nil {
			logerr.Fatalf("processing failed: %v", err)
		}
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
}

func dbInit() {
	log := logerr.Add("db init")
	host := fmt.Sprintf("%s://%s", argv.PostProcess.Protocol, argv.PostProcess.Host)

	var err error
	cypher.Driver, err = neo4j.NewDriver(
		host,
		neo4j.BasicAuth(argv.PostProcess.Username, argv.PostProcess.Password, ""),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// cypherQ, err := cypher.NewQuery()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	//
	// tx, err := cypherQ.Begin()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer tx.Rollback()

	// tx.Run("CREATE CONSTRAINT ON (a:Exe) ASSERT a.path IS UNIQUE;", nil)
	// tx.Run("CREATE CONSTRAINT ON (a:Dll) ASSERT a.path IS UNIQUE;", nil)
	// tx.Run("CREATE CONSTRAINT ON (a:Directory) ASSERT a.path IS UNIQUE;", nil)
	// tx.Run("CREATE CONSTRAINT ON (a:Principal) ASSERT a.name IS UNIQUE;", nil)
	// tx.Run("CREATE CONSTRAINT ON (a:Runner) ASSERT a.name IS UNIQUE;", nil)
	// tx.Run("CREATE CONSTRAINT ON (a:Dep) ASSERT a.name IS UNIQUE;", nil)
	//
	// err = tx.Commit()
	// if err != nil {
	// 	switch e := err.(type) {
	// 	case *neo4j.Neo4jError:
	// 		if e.Code == "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists" {
	// 			log.Debug("node constraints already  in place, skipping")
	// 		}
	// 	default:
	// 		log.Errorf("tx commit failed %s", err)
	// 	}
	// }

}

func dbDrop(doIt bool) {
	if doIt {
		logerr.Info("dropping database")
		cypherQ, err := cypher.NewQuery()
		if err != nil {
			logerr.Fatalf("drop failed: %s", err.Error())
		}
		cypherQ.Append(`
			CALL apoc.periodic.iterate(
			'MATCH (n) RETURN n', 'DETACH DELETE n'
			, {batchSize:1000})
		`).ExecuteW()
	}
}
