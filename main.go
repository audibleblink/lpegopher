package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
	"github.com/mindstand/gogm/v2"
)

var (
	argv = args.ArgType{}
	cli  = arg.MustParse(&argv)
)

func logInit() {
	l := &logerr.Logger{
		Level:            logerr.LogLevelInfo,
		Output:           os.Stderr,
		LogWrappedErrors: true,
	}
	l.SetAsGlobal()
}

func main() {

	logInit()

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

	config := gogm.Config{
		IndexStrategy: gogm.IGNORE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:      50,
		Port:          argv.Process.Port,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          argv.Process.Host,
		Password:      argv.Process.Password,
		Username:      argv.Process.Username,
		Protocol:      argv.Process.Protocol,
		// UseSystemCertPool: true,
		CAFileLocation:  "./ca.crt",
		EnableLogParams: false,
		Logger:          logerr.DefaultLogger(),
	}

	driver, err := gogm.New(
		&config,
		gogm.DefaultPrimaryKeyStrategy,
		&node.User{},
		&node.Group{},
		&node.Directory{},
		&node.EXE{},
		&node.DLL{},
		&node.Runner{},
	)
	if err != nil {
		log.Fatal(err)
	}

	gogm.SetGlobalGogm(driver)
}
