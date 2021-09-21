package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/capnspacehook/taskmaster"
)

type taskResult struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Exe      string `json:"Exe"`
	Parent   string `json:"Parent"`
	FullPath string `json:"FullPath"`
	Args     string `json:"Args"`
	Context  string `json:"Context"`
	RunLevel string `json:"RunLevel"`
}

func listTasks() {
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

			fullPath := evaluatePath(execAction.Path)

			taschzk := taskResult{
				Name:     task.Name,
				Type:     "task",
				Exe:      filepath.Base(fullPath),
				Parent:   filepath.Dir(fullPath),
				FullPath: fullPath,
				Context:  task.Definition.Context,
				RunLevel: task.Definition.Principal.RunLevel.String(),
			}

			jason, _ := json.Marshal(taschzk)
			fmt.Println(string(jason))
			// fmt.Println(taschzk.Cwd)
		}
	}
}
