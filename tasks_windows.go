package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func evaluatePath(path string) (out string) {
	// https://gitlab.com/stu0292/windowspathenv
	out = path
	// return unless strings starts with %
	if !strings.HasPrefix(path, "%") {
		return
	}
	// return unless there's a second %
	trim := strings.TrimPrefix(path, "%")
	i := strings.Index(trim, "%")
	if i == -1 {
		return
	}
	// check if substr between two % is the name of an existing env var
	val, ok := os.LookupEnv(trim[:i])
	if !ok {
		return
	}
	// env var value will use os path separator
	remainder := filepath.FromSlash(trim[i+1:])
	// check the remainder starts with path separateor
	if !strings.HasPrefix(remainder, "\\") {
		return
	}
	// prepend the value to the remainder of the path
	return val + remainder
}
