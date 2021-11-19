package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/rpcls/pkg/procs"
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
		pid := argv.GetSystem.PID
		if pid == 0 {
			pid = procs.PidForName("winlogon.exe")
			logerr.Infof("stealing winlogon token from pid %d", pid)
		}
		err := getsystem.InNewProcess(pid, `c:\windows\system32\cmd.exe`, false)
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
		if argv.PostProcess.Drop {
			dbDrop()
		}
		if argv.PostProcess.Runners == "" {
			logerr.Fatal("supply runner file")
		} else if argv.PostProcess.Runners != "" {
			logerr.Fatal("--all and  --pe/--runner are mutually exclusive")
		}
		err := doProcessCmd(argv, cli)
		if err != nil {
			logerr.Fatalf("processing failed:", err)
		}
	default:
		cli.WriteHelp(os.Stderr)
		os.Exit(1)
	}
}

func dbInit() {
	log := logerr.Add("db init")
	host := fmt.Sprintf("%s://%s", argv.Protocol, argv.Host)

	var err error
	cypher.Driver, err = neo4j.NewDriver(
		host,
		neo4j.BasicAuth(argv.Username, argv.Password, ""),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	cypherQ, err := cypher.NewQuery()
	if err != nil {
		log.Fatal(err.Error())
	}

	tx, err := cypherQ.Begin()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer tx.Rollback()

	tx.Run("CREATE CONSTRAINT ON (a:Exe) ASSERT a.path IS UNIQUE;", nil)
	tx.Run("CREATE CONSTRAINT ON (a:Dll) ASSERT a.path IS UNIQUE;", nil)
	tx.Run("CREATE CONSTRAINT ON (a:Directory) ASSERT a.path IS UNIQUE;", nil)
	tx.Run("CREATE CONSTRAINT ON (a:Principal) ASSERT a.name IS UNIQUE;", nil)
	tx.Run("CREATE CONSTRAINT ON (a:Runner) ASSERT a.name IS UNIQUE;", nil)

	err = tx.Commit()
	if err != nil {
		switch e := err.(type) {
		case *neo4j.Neo4jError:
			if e.Code == "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists" {
				log.Debug("node constraints already  in place, skipping")
			}
		default:
			log.Errorf("tx commit failed %s", err)
		}
	}

}

func dbDrop() {
	logerr.Debug("dropping database")
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
