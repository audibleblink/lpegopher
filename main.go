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
		Level:            logerr.LogLevelInfo,
		Output:           os.Stderr,
		LogWrappedErrors: true,
	}

	l.Context("lpegopher").SetAsGlobal()
}

func main() {

	switch {
	case argv.GetSystem != nil:
		err := getSystem(argv.GetSystem.PID)
		if err != nil {
			logerr.Fatalf("getsystem failed:", err)
		}
	case argv.Collect != nil:
		err := doCollectCmd(argv, cli)
		if err != nil {
			logerr.Fatalf("collection failed:", err)
		}
	case argv.Process != nil:
		dbInit()
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
