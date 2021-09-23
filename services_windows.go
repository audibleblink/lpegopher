package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

func listServices() {
	svcMgr, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	svcNames, err := svcMgr.ListServices()

	for _, svcName := range svcNames {
		svc, err := svcMgr.OpenService(svcName)
		if err != nil {
			log.Fatal(err)
		}
		conf, err := svc.Config()
		if err != nil {
			log.Fatal(err)
		}
		cmdLine := conf.BinaryPathName
		splitCmd := strings.Split(cmdLine, " ")
		path := splitCmd[0]
		args := strings.Join(splitCmd[1:], " ")
		context := conf.ServiceStartName

		service := taskResult{
			Name:     conf.DisplayName,
			Type:     "service",
			Exe:      filepath.Base(path),
			Parent:   filepath.Dir(path),
			Args:     args,
			FullPath: path,
			Context:  context,
		}

		jason, _ := json.Marshal(service)
		fmt.Println(string(jason))
	}
}
