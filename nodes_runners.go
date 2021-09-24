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
	// Querier
	gogm.BaseNode

	Context *User      `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Name    string     `gogm:"name=name"`
	Exe     *EXE       `gogm:"direction=incoming;relationship=EXECUTED_FROM"`
	ExeDir  *Directory `gogm:"direction=incoming;relationship=HOSTS_PES"`
	Type    string     `gogm:"name=type"`
}

func (r *Runner) RunsExeAs(user *User) error {
	r.Context = user
	return r.save()
}

func (r *Runner) save() (err error) {
	sess, err := newNeoSession()
	return sess.Save(context.Background(), r)
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
	err = user.Merge("name", user.Name)
	if err != nil {
		return
	}

	exeNode := &EXE{}
	exeNode.Name = exe
	exeNode.Path = fullPath
	err = exeNode.Merge("path", exeNode.Path)
	if err != nil {
		return
	}

	dir := &Directory{}
	dir.Name = filepath.Base(parent)
	dir.Path = parent
	err = dir.Merge("path", dir.Path)
	if err != nil {
		return
	}

	err = runner.Merge("name", runner.Name)
	if err != nil {
		return
	}

	err = dir.Hosts(runner)
	if err != nil {
		return
	}
	err = dir.Add(exeNode)
	if err != nil {
		return
	}
	err = runner.RunsExeAs(user)
	if err != nil {
		return
	}
	return
}

func (x *Runner) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "Runner"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, nil, x)
}

type Task struct{ Runner }
type Service struct{ Runner }
