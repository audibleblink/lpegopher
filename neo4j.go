package main

import (
	"github.com/mindstand/gogm/v2"
)

func newNeoSession() *gogm.SessionV2 {
	config := gogm.Config{
		IndexStrategy:     gogm.IGNORE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:          50,
		Port:              7687,
		IsCluster:         false, //tells it whether or not to use `bolt+routing`
		Host:              "neo4j.i.hyrule.link",
		Password:          "password",
		Username:          "neo4j",
		Protocol:          "bolt+s",
		UseSystemCertPool: true,
		EnableLogParams:   true,
	}

	_gogm, err := gogm.New(&config, gogm.UUIDPrimaryKeyStrategy, &User{}, &Group{}, &Directory{}, &EXE{}, &DLL{}, &Task{}, &Service{})
	if err != nil {
		panic(err)
	}

	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := _gogm.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	return &sess
}
