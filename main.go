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
		dbDrop(argv.PostProcess.Drop)

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

	cypherQ, err := cypher.NewQuery()
	if err != nil {
		log.Fatal(err.Error())
	}

	tx, err := cypherQ.Begin()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer tx.Rollback()

	iq := "CREATE BTREE INDEX IF NOT EXISTS FOR (n:%s) ON (n.%s)"
	tx.Run(fmt.Sprintf(iq, "Exe", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Exe", "path"), nil)
	tx.Run(fmt.Sprintf(iq, "Dll", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Dll", "path"), nil)
	tx.Run(fmt.Sprintf(iq, "Directory", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Directory", "path"), nil)

	tx.Run(fmt.Sprintf(iq, "Principal", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Principal", "name"), nil)
	tx.Run(fmt.Sprintf(iq, "Runner", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Runner", "name"), nil)
	tx.Run(fmt.Sprintf(iq, "Dep", "nid"), nil)
	tx.Run(fmt.Sprintf(iq, "Dep", "name"), nil)

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

func dbDrop(doIt bool) {
	if doIt {
		logerr.Info("dropping database")
		cypherQ, err := cypher.NewQuery()
		if err != nil {
			logerr.Fatalf("couldn't create neo4j session: %s", err.Error())
		}
		err = cypherQ.Append(`
			CALL apoc.periodic.iterate(
			'MATCH (n) RETURN n', 'DETACH DELETE n'
			, {batchSize:1000})
		`).ExecuteW()
		if err != nil {
			logerr.Fatalf("drop failed: %s", err.Error())
		}
	}
}
