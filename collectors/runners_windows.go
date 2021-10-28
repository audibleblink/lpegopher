package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/capnspacehook/taskmaster"
	"golang.org/x/sys/windows/svc/mgr"
)

type PERunner struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Exe      string `json:"Exe"`
	Parent   string `json:"Parent"`
	FullPath string `json:"FullPath"`
	Args     string `json:"Args"`
	Context  string `json:"Context"`
	RunLevel string `json:"RunLevel"`
}

func Tasks(writer io.Writer) {
	logerr.Context("tasks")
	defer logerr.ClearContext()

	svc, _ := taskmaster.Connect()
	tasks, _ := svc.GetRegisteredTasks()

	for _, task := range tasks {

		if task.Enabled {

			var execAction taskmaster.ExecAction
			actionType := task.Definition.Actions[0].GetType()

			switch actionType {
			case taskmaster.TASK_ACTION_EXEC:
				execAction = task.Definition.Actions[0].(taskmaster.ExecAction)
			default:
				continue
			}

			if execAction.Path == "" {
				continue
			}

			fullPath := util.EvaluatePath(execAction.Path)
			args := execAction.Args

			taschzk := PERunner{
				Name:     task.Name,
				Type:     "task",
				Exe:      filepath.Base(fullPath),
				Args:     args,
				Parent:   filepath.Dir(fullPath),
				FullPath: fullPath,
				Context:  task.Definition.Context,
				RunLevel: task.Definition.Principal.RunLevel.String(),
			}

			jason, _ := json.Marshal(taschzk)
			fmt.Fprintln(writer, string(jason))
		}
	}
}

func Services(writer io.Writer) {
	logerr.Context("services")
	defer logerr.ClearContext()

	svcMgr, err := mgr.Connect()
	if err != nil {
		logerr.Error(err.Error())
		return
	}

	svcNames, err := svcMgr.ListServices()
	if err != nil {
		logerr.Error(err.Error())
		return
	}

	for _, svcName := range svcNames {
		svc, err := svcMgr.OpenService(svcName)
		if err != nil {
			logerr.Warnf(svcName, err)
			logerr.Warnf("failed to open service", svcName, err)
			continue
		}
		conf, err := svc.Config()
		if err != nil {
			logerr.Warnf("failed to fetch service config", svcName, err)
			continue
		}

		cmdLine := conf.BinaryPathName
		splitCmd := strings.Split(cmdLine, " ")
		path := splitCmd[0]
		args := strings.Join(splitCmd[1:], " ")
		context := conf.ServiceStartName

		service := PERunner{
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
