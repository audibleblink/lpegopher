package main

import (
	"github.com/Jeffail/gabs/v2"
	gogm "github.com/mindstand/gogm/v2"
)

// Runners can be Tasks, Services, or AutoRuns
type Runner struct {
	gogm.BaseNode

	Context *User      `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Name    string     `gogm:"name=name"`
	Exe     *EXE       `gogm:"direction=incoming;relationship=EXECUTED_FROM"`
	ExeDir  *Directory `gogm:"direction=incoming;relationship=HOSTS_PES"`
	Type    string     `gogm:"name=type"`
}

func NewRunnerFromJson(jsonLine []byte) (runner *Runner, err error) {
	line, err := gabs.ParseJSON(jsonLine)
	if err != nil {
		return
	}

	if name, ok := line.Path("Name").Data().(string); !ok {
		return
	}

	if context, ok := line.Path("Context").Data().(string); !ok {
		return
	}

	if exe, ok := line.Path("Exe").Data().(string); !ok {
		return
	}

	if parent, ok := line.Path("Parent").Data().(string); !ok {
		return
	}

	if fullPath, ok := line.Path("FullPath").Data().(string); !ok {
		return
	}

	return
}

func (r *Runner) Save() error {

}

type Task struct{ Runner }
type Service struct{ Runner }
