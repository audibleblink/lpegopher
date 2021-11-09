package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/capnspacehook/taskmaster"
	"golang.org/x/sys/windows/svc/mgr"
)

type PERunner struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Exe      string `json:"Exe"`
	Parent   string `json:"Parent"`   // Directory.Path
	FullPath string `json:"FullPath"` // PE.Path
	Args     string `json:"Args"`
	Context  string `json:"Context"` // Principal.Name
	RunLevel string `json:"RunLevel"`
}

func Tasks(writer io.Writer) {
	log := logerr.Add("tasks")

	svc, err := taskmaster.Connect()
	if err != nil {
		log.Fatalf("could not connect to tasks scheduler: %s", err)
	}
	tasks, err := svc.GetRegisteredTasks()
	if err != nil {
		log.Fatalf("could not fetch registered tasks: %s", err)
	}

	for _, task := range tasks {

		if task.Enabled {

			var execAction taskmaster.ExecAction
			actionType := task.Definition.Actions[0].GetType()

			switch actionType {
			case taskmaster.TASK_ACTION_EXEC:
				execAction = task.Definition.Actions[0].(taskmaster.ExecAction)
			default:
				log.Debugf("task %s has no action, continuing", task.Name)
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
	log := logerr.Add("services")
	defer logerr.ClearContext()

	svcMgr, err := mgr.Connect()
	if err != nil {
		log.Error(err.Error())
		return
	}

	svcNames, err := svcMgr.ListServices()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, svcName := range svcNames {
		svc, err := svcMgr.OpenService(svcName)
		if err != nil {
			log.Warnf("failed to open service: %s: %s", svcName, err)
			continue
		}
		conf, err := svc.Config()
		if err != nil {
			log.Warnf("failed to fetch service config: %s: %s", svcName, err)
			continue
		}

		path, args := util.SmoothBrainPath(conf.BinaryPathName)
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
