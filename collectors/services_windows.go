package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/audibleblink/pegopher/logerr"
	"golang.org/x/sys/windows/svc/mgr"
)

func Services(writer io.Writer) {
	svcLog := logerr.DefaultLogger().Context("serivces")
	svcLog.Level = logerr.LogLevelWarn

	svcMgr, err := mgr.Connect()
	if err != nil {
		svcLog.Error(err.Error())
		return
	}

	svcNames, err := svcMgr.ListServices()
	if err != nil {
		svcLog.Error(err.Error())
		return
	}

	for _, svcName := range svcNames {
		svc, err := svcMgr.OpenService(svcName)
		if err != nil {
			svcLog.Warnf(svcName, err)
			svcLog.Warnf("failed to open service", svcName, err)
			continue
		}
		conf, err := svc.Config()
		if err != nil {
			svcLog.Warnf("failed to fetch service config", svcName, err)
			continue
		}

		cmdLine := conf.BinaryPathName
		splitCmd := strings.Split(cmdLine, " ")
		path := splitCmd[0]
		args := strings.Join(splitCmd[1:], " ")
		context := conf.ServiceStartName

		service := TaskResult{
			Name:     conf.DisplayName,
			Type:     "service",
			Exe:      filepath.Base(path),
			Parent:   filepath.Dir(path),
			Args:     args,
			FullPath: path,
			Context:  context,
		}

		jason, _ := json.Marshal(service)
		fmt.Fprintln(writer, string(jason))
	}
}
