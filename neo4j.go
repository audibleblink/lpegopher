package main

import (
	"github.com/mindstand/gogm/v2"
)

func init() {

	config := gogm.Config{
		IndexStrategy:     gogm.IGNORE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:          50,
		Port:              args.Process.Port,
		IsCluster:         false, //tells it whether or not to use `bolt+routing`
		Host:              args.Process.Host,
		Password:          args.Process.Password,
		Username:          args.Process.Username,
		Protocol:          args.Process.Protocol,
		UseSystemCertPool: true,
		EnableLogParams:   true,
	}

	driver, err := gogm.New(
		&config,
		gogm.DefaultPrimaryKeyStrategy,
		&User{},
		&Group{},
		&Directory{},
		&EXE{},
		&DLL{},
		&Runner{},
		// &Task{},
		// &Service{},
	)
	if err != nil {
		panic(err)
	}

	gogm.SetGlobalGogm(driver)
}

func newNeoSession() (sess gogm.SessionV2, err error) {

	sessConf := gogm.SessionConfig{
		AccessMode:   gogm.AccessModeWrite,
		DatabaseName: args.Process.Database,
	}

	return gogm.G().NewSessionV2(sessConf)
}
