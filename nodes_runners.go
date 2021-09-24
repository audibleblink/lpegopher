package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Jeffail/gabs/v2"

	"github.com/mindstand/gogm/v2"
)

// Runners can be Tasks, Services, or AutoRuns
type Runner struct {
	gogm.BaseNode

	Context *User      `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Name    string     `gogm:"name=name"`
	Exe     *EXE       `gogm:"direction=incoming;relationship=EXECUTED_FROM"`
	Parent  *Directory `gogm:"direction=incoming;relationship=HOSTS_PES"`
	Type    string     `gogm:"name=type"`
}

func NewRunnerFromJson(jsonLine []byte) (runner *Runner, err error) {
	line, err := gabs.ParseJSON(jsonLine)
	if err != nil {
		return
	}

	var (
		userName string
		exe      string
		parent   string
		fullPath string

		ok bool
	)
	runner = &Runner{}
	if runner.Name, ok = line.Path("Name").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "Name")
		return
	}

	if userName, ok = line.Path("Context").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "userName")
		return
	}

	if exe, ok = line.Path("Exe").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "exeName")
		return
	}

	if fullPath, ok = line.Path("FullPath").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "fullPath")
		return
	}

	if parent, ok = line.Path("Parent").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "parent")
		return
	}

	user := &User{}
	user.Name = userName

	exeNode := &EXE{}
	exeNode.Name = exe
	exeNode.Path = fullPath

	dir := &Directory{}
	dir.Name = filepath.Base(parent)
	dir.Path = parent

	runner.Context = user
	runner.Exe = exeNode
	runner.Parent = dir

	return
}

func (r *Runner) Save() (err error) {
	sess, err := newNeoSession()
	if err != nil {
		return
	}

	return sess.Save(context.Background(), r)
}

type Task struct{ Runner }
type Service struct{ Runner }
