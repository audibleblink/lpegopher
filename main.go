package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"github.com/audibleblink/logerr"
	"github.com/audibleblink/lpegopher/args"
	"github.com/audibleblink/lpegopher/cypher"
	"github.com/audibleblink/lpegopher/node"
)

// Global variables
var (
	argv = args.ArgType{}       // Command line arguments
	cli  = arg.MustParse(&argv) // Parsed command line
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

	case argv.Process != nil:
		dbInit()
		if argv.Process.Drop {
			err := dbDrop()
			if err != nil {
				logerr.Fatalf("db drop failed: %v", err)
			}
		}
		err := dbCreateIndices()
		if err != nil {
			logerr.Fatalf("index creation failed: %v", err)
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
	host := fmt.Sprintf("%s://%s", argv.Process.Protocol, argv.Process.Host)

	var err error
	cypher.Driver, err = neo4j.NewDriver(
		host,
		neo4j.BasicAuth(argv.Process.Username, argv.Process.Password, ""),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func dbCreateIndices() error {
	log := logerr.Add("db indices")
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		return log.Wrap(err)
	}

	tx, err := cypherQ.Begin()
	if err != nil {
		return log.Wrap(err)
	}
	defer tx.Rollback()

	log.Debug("creating indices")

	// Use the node package to create all schema constraints and indices
	nodeSchema := node.NewNodeSchema(tx)

	// Create unique constraints
	log.Debug("creating unique constraints")
	if err := nodeSchema.CreateUniqueConstraints(); err != nil {
		return log.Wrap(err)
	}

	// Create btree indices
	log.Debug("creating btree indices")
	if err := nodeSchema.CreateBTreeIndices(); err != nil {
		return log.Wrap(err)
	}

	err = tx.Commit()
	if err != nil {
		switch e := err.(type) {
		case *neo4j.Neo4jError:
			if e.Code == "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists" {
				log.Debug("node constraints already in place, skipping")
				return nil
			} else {
				return log.Wrap(err)
			}
		default:
			log.Errorf("tx commit failed %s", err)
			return log.Wrap(err)
		}
	}
	return nil
}

func dbDrop() error {
	log := logerr.Add("drop")
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		return log.Add("session creation failed").Wrap(err)
	}

	log.Info("dropping graph")
	err = cypherQ.Append(`
			CALL apoc.periodic.iterate(
			'MATCH (n) RETURN n', 'DETACH DELETE n'
			, {batchSize: 5000, parallel: true})
		`).ExecuteW()
	if err != nil {
		return log.Add("couldn't drop database").Wrap(err)
	}

	log.Debug("dropping schema")
	cypherQ, _ = cypher.NewQuery()
	err = cypherQ.Append(`
			CALL apoc.schema.assert({},{},true) YIELD label, key RETURN *;
		`).ExecuteW()
	if err != nil {
		return log.Add("couldn't reset schema").Wrap(err)
	}
	return nil
}
