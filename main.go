package main

import (
	"fmt"
	"log"
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
		Level:   logerr.LogLevelInfo,
		Output:  os.Stderr,
		NoColor: argv.NoColor,
	}

	if argv.Debug {
		l.Level = logerr.LogLevelDebug
	}

	l.SetContext("lpegopher").SetAsGlobal()
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
		if argv.PostProcess.Drop {
			err := dbDrop()
			if err != nil {
				log.Fatalf("db drop failed: %v", err)
			}
		}
		err := dbCreateIndices()
		if err != nil {
			log.Fatalf("index creation failed: %v", err)
		}

		err = doProcessCmd(argv, cli)
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
}

func dbCreateIndices() error {
	log := logerr.Add("db indices")
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		log.Fatal(err.Error())
	}

	tx, err := cypherQ.Begin()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer tx.Rollback()

	iq := "CREATE CONSTRAINT ON (a:%s) ASSERT a.%s IS UNIQUE;"
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
				return nil
			}
		default:
			log.Errorf("tx commit failed %s", err)
			return err
		}
	}
	return nil
}

func dbDrop() error {
	log := logerr.Add("db drop")
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		return log.Add("session creation failed").Wrap(err)
	}
	err = cypherQ.Append(`
			CALL apoc.periodic.iterate(
			'MATCH (n) RETURN n', 'DETACH DELETE n'
			, {batchSize:10000})
		`).ExecuteW()
	if err != nil {
		return log.Add("couldn't drop database").Wrap(err)
	}

	log.Info("dropping schema")
	cypherQ, _ = cypher.NewQuery()
	err = cypherQ.Append(`
			CALL apoc.schema.assert({},{},true) YIELD label, key RETURN *;
		`).ExecuteW()
	if err != nil {
		return log.Add("couldn't reset schema").Wrap(err)
	}
	log.Info("database dropped")
	return nil
}
