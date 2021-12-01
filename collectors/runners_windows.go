package collectors

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/capnspacehook/taskmaster"
	"github.com/minio/highwayhash"
	"golang.org/x/sys/windows/svc/mgr"
)

func Tasks() {
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

			exe := &INode{
				Path:   fullPath,
				Name:   filepath.Base(fullPath),
				Parent: filepath.Dir(fullPath),
			}

			taschzk := PERunner{
				Name:     task.Name,
				Type:     "task",
				Args:     args,
				Exe:      exe,
				Context:  &Principal{Name: task.Definition.Context},
				RunLevel: task.Definition.Principal.RunLevel.String(),
			}

			taschzk.Exe.Write(writers[ExeFile])
			taschzk.Context.Write(writers[PrincipalFile])
			taschzk.Write(writers[RunnersFile])
		}
	}
}

func Services() {
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
			log.Warnf("failed to open service %s: %s", svcName, err)
			continue
		}
		conf, err := svc.Config()
		if err != nil {
			log.Warnf("failed to fetch service config: %s: %s", svcName, err)
			continue
		}

		path, args := util.SmoothBrainPath(conf.BinaryPathName)
		context := &Principal{Name: conf.ServiceStartName}

		exe := &INode{
			Path:   path,
			Name:   filepath.Base(path),
			Parent: filepath.Dir(path),
		}

		service := PERunner{
			Name:    conf.DisplayName,
			Type:    "service",
			Args:    args,
			Exe:     exe,
			Context: context,
		}

		// if strings.HasSuffix(context.Name, "ystem") {
		// 	fmt.Print(1)
		// }

		service.Exe.Write(writers[ExeFile])
		service.Context.Write(writers[PrincipalFile])
		service.Write(writers[RunnersFile])
	}
}

func hashFor(data string) string {
	data = util.PathFix(data)
	hash, err := highwayhash.New(key)
	if err != nil {
		fmt.Printf("Failed to create HighwayHash instance: %v", err)
		os.Exit(1)
	}

	txt := strings.NewReader(data)
	if _, err = io.Copy(hash, txt); err != nil {
		fmt.Printf("hash reader creation failed: %v", err)
		os.Exit(1)
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}
