package main

import (
	gogm "github.com/mindstand/gogm/v2"
)

type containableIFile struct {
	Parent *Directory `gogm:"direction=incoming;relationship=CONTAINS"`
}

type securableIFile struct {
	gogm.BaseNode
	Name string `gogm:"name=name"`
	Path string `gogm:"name=path;unique"`
	// Parent *Directory `gogm:"direction=incoming;relationship=CONTAINS"`

	PermsWriteOwnerG   *Group `gogm:"direction=incoming;relationship=WRITE_OWNER"`
	PermsWriteOwnerU   *User  `gogm:"direction=incoming;relationship=WRITE_OWNER"`
	PermsWriteDACLG    *Group `gogm:"direction=incoming;relationship=WRITE_DACL"`
	PermsWriteDACLU    *User  `gogm:"direction=incoming;relationship=WRITE_DACL"`
	PermsWritePropG    *Group `gogm:"direction=incoming;relationship=WRITE_PROP"`
	PermsWritePropU    *User  `gogm:"direction=incoming;relationship=WRITE_PROP"`
	PermsGenericAllG   *Group `gogm:"direction=incoming;relationship=GENERIC_ALL"`
	PermsGenericAllU   *User  `gogm:"direction=incoming;relationship=GENERIC_ALL"`
	PermsGenericWriteG *Group `gogm:"direction=incoming;relationship=GENERIC_WRITE"`
	PermsGenericWriteU *User  `gogm:"direction=incoming;relationship=GENERIC_WRITE"`
}

type Directory struct {
	securableIFile
	ContainedExes []*EXE  `gogm:"direction=outgoing;relationship=CONTAINS"`
	ContainedDlls []*DLL  `gogm:"direction=outgoing;relationship=CONTAINS"`
	HostsPEs      *Runner `gogm:"direction=outgoing;relationship=HOSTS_PES"`
}

func (d *Directory) Hosts(runner *Runner) {
	d.HostsPEs = runner
}

func (d *Directory) Add(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		d.ContainedExes = append(d.ContainedExes, f)
	case *DLL:
		d.ContainedDlls = append(d.ContainedDlls, f)
	}
}

type EXE struct {
	securableIFile
	containableIFile

	ExecutesFrom *Runner `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
	// Parent *Directory `gogm:"direction=incoming;relationship=CONTAINS"`
	// Imports      []*containableIFile   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*containableIFile   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
}

func (f *EXE) GetsRunBy(runner *Runner) {
	f.ExecutesFrom = runner
}

type DLL struct {
	containableIFile
	securableIFile
	// Imports      []*PE   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*PE   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
	// ExecutedFrom *Runner `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
}
