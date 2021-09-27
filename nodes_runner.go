package main

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"

	"github.com/mindstand/gogm/v2"
)

// Runners can be Tasks, Services, or AutoRuns
type Runner struct {
	// Querier
	gogm.BaseNode

	Name    string     `gogm:"name=name;unique"`
	Context *User      `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Exe     *EXE       `gogm:"direction=incoming;relationship=EXECUTED_FROM"`
	ExeDir  *Directory `gogm:"direction=incoming;relationship=HOSTS_PES"`
	Type    string     `gogm:"name=type"`
}

func (r *Runner) RunsExeAs(user *User) error {
	r.Context = user
	return r.save()
}

func NewRunnerFromJson(jsonLine []byte) (err error) {
	line, err := gabs.ParseJSON(jsonLine)
	if err != nil {
		return
	}

	var (
		userName   string
		exe        string
		parent     string
		runnerName string
		runnerType string
		fullPath   string

		ok bool
	)
	if runnerName, ok = line.Path("Name").Data().(string); !ok {
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

	if runnerType, ok = line.Path("Type").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "parent")
		return
	}

	user := &User{}
	err = user.Merge("name", userName)
	if err != nil {
		return
	}

	exeNode := &EXE{}
	err = exeNode.Merge("path", fullPath)
	if err != nil {
		return
	}
	err = exeNode.SetName(exe)
	if err != nil {
		return
	}

	dir := &Directory{}
	dir.Path = parent
	err = dir.Merge("path", dir.Path)
	if err != nil {
		return
	}
	err = dir.SetName(winPathBase(parent))
	if err != nil {
		return
	}

	runner := &Runner{}
	err = runner.Merge("name", runnerName)
	if err != nil {
		return
	}
	err = runner.SetType(runnerType)
	if err != nil {
		return
	}

	// Associate a directory that hosts PEs for a particular Runner
	err = dir.HostsPEsFor(runner)
	if err != nil {
		return
	}

	// Creats the edges for Directories CONTAINS a PE
	err = dir.Add(exeNode)
	if err != nil {
		return
	}

	// Associate the Context (or User) from which the Runner executes
	err = runner.RunsExeAs(user)
	if err != nil {
		return
	}

	// Associate the EXE as getting run by a service
	err = exeNode.GetsRunBy(runner)
	if err != nil {
		return
	}

	return
}

// Merge will either create or retreive the node based on the key-valie pair provides
// In this case, the Runner struct designates the "name" field as unique
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

// type Task struct{ Runner }
// type Service struct{ Runner }

func (x *Runner) SetType(name string) error {
	x.Type = name
	return x.save()
}

func (x *Runner) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := newNeoSession()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
