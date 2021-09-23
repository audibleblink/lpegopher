package main

import (
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

type Task struct{ Runner }
type Service struct{ Runner }
