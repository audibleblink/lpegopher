package main

import (
	gogm "github.com/mindstand/gogm/v2"
)

type User struct {
	Groups []*Group `gogm:"direction=outgoing;relationship=BELONGS_TO"`
	Principal
}

type Group struct {
	Principal
}

type Principal struct {
	gogm.BaseNode
	Name         string    `gogm:"unique"`
	WriteOwner   *FSObject `gogm:"direction=outgoing;relationship=WRITE_OWNER"`
	WriteDACL    *FSObject `gogm:"direction=outgoing;relationship=WRITE_DACL"`
	WriteProp    *FSObject `gogm:"direction=outgoing;relationship=WRITE_PROP"`
	GenericAll   *FSObject `gogm:"direction=outgoing;relationship=GENERIC_ALL"`
	GenericWrite *FSObject `gogm:"direction=outgoing;relationship=GENERIC_WRITE"`
}

type FSObject struct {
	gogm.BaseNode

	Name   string `gogm:"index"`
	Path   string `gogm:"pk"`
	Parent *Directory
}

type Directory struct {
	FSObject
	Contains []*FSObject `gogm:"direction=outgoing;relationship=CONTAINS"`
}

type PE struct {
	FSObject
	Imports      []*FSObject `gogm:"direction=outgoing;relationship=IMPORTS"`
	ImportedBy   []*FSObject `gogm:"direction=outgoing;relationship=IMPORTED_BY"`
	ExecutedFrom *Runner     `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
}

type EXE struct{ PE }
type DLL struct{ PE }

type Runner struct {
	gogm.BaseNode
	Context *User  `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Name    string `gogm:"unique,pk"`
	Exe     *EXE
	ExeDir  *Directory
}
type Task struct{ Runner }
type Service struct{ Runner }
