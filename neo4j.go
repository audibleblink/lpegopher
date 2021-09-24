package main

import (
	"github.com/mindstand/gogm/v2"
)

func newNeoSession() (sess gogm.SessionV2, err error) {
	config := gogm.Config{
		IndexStrategy:     gogm.ASSERT_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
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
		&Task{},
		&Service{},
	)
	if err != nil {
		return
	}

	sessConf := gogm.SessionConfig{
		AccessMode:   gogm.AccessModeWrite,
		DatabaseName: args.Process.Database,
	}

	return driver.NewSessionV2(sessConf)
}
