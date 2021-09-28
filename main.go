package main

import (
	"github.com/alexflint/go-arg"
	"github.com/audibleblink/pegopher/args"
	"github.com/audibleblink/pegopher/node"
	"github.com/mindstand/gogm/v2"
)

var (
	argv = args.ArgType{}
	cli  = arg.MustParse(&argv)
)

func main() {
	switch {
	case argv.Collect != nil:
		err := doCollectCmd(cli, argv)
		if err != nil {
			cli.Fail(err.Error())
		}
	case argv.Process != nil:
		err := doProcessCmd(cli, argv)
		if err != nil {
			cli.Fail(err.Error())
		}
	}
}

func init() {

	config := gogm.Config{
		IndexStrategy:     gogm.IGNORE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:          50,
		Port:              argv.Process.Port,
		IsCluster:         false, //tells it whether or not to use `bolt+routing`
		Host:              argv.Process.Host,
		Password:          argv.Process.Password,
		Username:          argv.Process.Username,
		Protocol:          argv.Process.Protocol,
		UseSystemCertPool: true,
		EnableLogParams:   true,
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
		// &Task{},
		// &Service{},
	)
	if err != nil {
		panic(err)
	}

	gogm.SetGlobalGogm(driver)
}
