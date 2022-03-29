package collectors

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/audibleblink/getsystem"
	"github.com/audibleblink/memutils"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/capnspacehook/taskmaster"
	"golang.org/x/sys/windows/svc/mgr"
)

var regKeys = []map[registry.Key]string{
	{registry.LOCAL_MACHINE: `Software\Microsoft\Windows\CurrentVersion\Run`},
	{registry.LOCAL_MACHINE: `Software\Microsoft\Windows\CurrentVersion\RunOnce`},
	{registry.LOCAL_MACHINE: `Software\Microsoft\Windows\CurrentVersion\RunServices`},
	{registry.LOCAL_MACHINE: `Software\Microsoft\Windows\CurrentVersion\RunServicesOnce`},
	{registry.CURRENT_USER: `Software\Microsoft\Windows\CurrentVersion\Run`},
	{registry.CURRENT_USER: `Software\Microsoft\Windows\CurrentVersion\RunOnce`},
	{registry.CURRENT_USER: `Software\Microsoft\Windows\CurrentVersion\RunServices`},
	{registry.CURRENT_USER: `Software\Microsoft\Windows\CurrentVersion\RunServicesOnce`},
	{registry.CURRENT_USER: `ProgID\Software\Microsoft\Windows\CurrentVersion\Run`},
}

func Autoruns() {
	log := logerr.Add("autoruns")
	defer logerr.ClearContext()

	for _, regKey := range regKeys {
		for key, subKey := range regKey {

			key, err := registry.OpenKey(key, subKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
			if err != nil {
				log.Debugf("unable to read key: %s", err)
				continue
			}
			defer key.Close()

			info, err := key.Stat()
			if err != nil {
				log.Debugf("unable to read key info: %s", err)
				continue
			}

			valueNames, err := key.ReadValueNames(int(info.SubKeyCount))
			if err != nil {
				log.Debugf("unable to read subkeys: %s", err)
				continue
			}

			for _, valueName := range valueNames {
				val, _, err := key.GetStringValue(valueName)
				if err != nil {
					log.Debugf("unable to read value: %s", err)
					continue
				}
				path, args := util.SmoothBrainPath(util.EvaluatePath(val))

				context := &Principal{Name: "unknown"}

				exe := &INode{
					Path:   path,
					Name:   filepath.Base(path),
					Parent: filepath.Dir(path),
				}

				autorun := PERunner{
					Name:    valueName,
					Type:    "autorun",
					Args:    args,
					Exe:     exe,
					Context: context,
				}

				autorun.Exe.Write(writers[ExeFile])
				autorun.Context.Write(writers[PrincipalFile])
				autorun.Write(writers[RunnersFile])
			}
		}
	}
}

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

	var s *uint16
	h, err := windows.OpenSCManager(s, nil, windows.SC_MANAGER_CONNECT|windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		log.Error(err.Error())
		return
	}

	svcMgr := &mgr.Mgr{Handle: h}

	svcNames, err := svcMgr.ListServices()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, svcName := range svcNames {

		h, err := windows.OpenService(svcMgr.Handle, windows.StringToUTF16Ptr(svcName), windows.SERVICE_QUERY_CONFIG)
		if err != nil {
			log.Warnf("failed to open service %s: %s", svcName, err)
			continue
		}

		svc := &mgr.Service{Name: svcName, Handle: h}

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

		service.Exe.Write(writers[ExeFile])
		service.Context.Write(writers[PrincipalFile])
		service.Write(writers[RunnersFile])
	}
}

func Processes() {
	log := logerr.Add("processes")
	defer logerr.ClearContext()

	processes, err := memutils.Processes()
	if err != nil {
		log.Warnf("failed to enumerate processes: %s", err)
		return
	}

	for _, process := range processes {

		path, args := util.SmoothBrainPath(util.EvaluatePath(process.Exe))
		exe := &INode{
			Path:   path,
			Name:   filepath.Base(path),
			Parent: filepath.Dir(path),
		}

		token, err := tokenForPid(process.Pid)
		if err != nil {
			log.Warnf("failed to query token of pid %d: %s", process.Pid, err)
			continue
		}

		owner, err := getsystem.TokenOwner(token)
		if err != nil {
			log.Warnf("failed to query owner of pid %d: %s", process.Pid, err)
			continue
		}

		context := &Principal{Name: owner}

		proc := PERunner{
			Name:    process.Exe,
			Type:    "process",
			Args:    args,
			Exe:     exe,
			Context: context,
		}

		proc.Exe.Write(writers[ExeFile])
		proc.Context.Write(writers[PrincipalFile])
		proc.Write(writers[RunnersFile])
	}
}

func tokenForPid(pid int) (tokenH windows.Token, err error) {
	hProc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, true, uint32(pid))
	if err != nil {
		err = fmt.Errorf("tokenForPid | openProcess | %s", err)
		return
	}

	err = windows.OpenProcessToken(hProc, windows.TOKEN_QUERY, &tokenH)
	if err != nil {
		err = fmt.Errorf("tokenForPid | openToken | %s", err)
	}
	return
}
