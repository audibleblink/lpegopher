package collectors

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/audibleblink/pegopher/util"
)

const (
	WriteOwner    = "WRITE_OWNER"
	WriteDACL     = "WRITE_DACL"
	GenericAll    = "GENERIC_ALL"
	GenericWrite  = "GENERIC_WRITE"
	ControlAccess = "CONTROL_ACCESS"
	Owns          = "OWNS"

	HostsPEFor = "HOSTS_PE_FOR"
	Contains   = "CONTAINS"
	ExecutedBy = "EXECUTED_BY"
	RunsAs     = "RUNS_AS"

	Imports    = "IMPORTS"
	Forwards   = "FORWARS"
	ImportedBy = "IMPORTED_BY"
)

// INode contains the parsed import and exports of the INode
type INode struct {
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	Parent   string `json:"Dir"`
	Type     string `json:"Type"`
	Imports  []*Dep `json:"Imports"`
	Forwards []*Dep `json:"Forwards"`
	DACL     DACL   `json:"DACL"`
}

func (i INode) ID() string {
	return hashFor(i.Path)
}

func (i INode) Write(file *bufio.Writer) string {
	id := i.ID()
	csv := i.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, csv)
	if !cacheHit {
		file.WriteString(csv)
	}
	return id
}

func (i INode) ToCSV() string {
	o := "NULL"
	g := "NULL"
	if i.DACL.Group != nil {
		g = i.DACL.Group.Name
	}
	if i.DACL.Owner != nil {
		o = i.DACL.Owner.Name
	}

	fields := make([]string, 6)
	fields[0] = i.ID()
	fields[1] = util.PathFix(i.Name)
	fields[2] = util.PathFix(i.Path)
	fields[3] = util.PathFix(i.Parent)
	fields[4] = hashFor(o)
	fields[5] = hashFor(g)
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}

type DACL struct {
	Owner *Principal    `json:"Owner"`
	Group *Principal    `json:"Group"`
	Aces  []ReadableAce `json:"Aces"`
}

type ReadableAce struct {
	Principal *Principal `json:"Principal"`
	Rights    []string   `json:"Rights"`
}

// Principal represents Users or Groups
type Principal struct {
	Name string `json:"Name"`
}

func (p Principal) ID() string {
	return hashFor(p.Name)
}

func (p Principal) ToCSV() string {
	fields := make([]string, 6)
	fields[0] = p.ID()
	fields[1] = util.PathFix(p.Name)
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}
func (p Principal) Write(file *bufio.Writer) string {
	id, csv := p.ID(), p.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, csv)
	if !cacheHit {
		file.WriteString(csv)
	}
	return id
}

type Rel struct {
	Start string
	Rel   string
	End   string
}

func (r Rel) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s\n", r.Start, r.Rel, r.End)
}

func (r Rel) ID() string {
	return hashFor(r.ToCSV())
}

func (r Rel) Write(file *bufio.Writer) string {
	id, csv := r.ID(), r.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, csv)
	if !cacheHit {
		file.WriteString(csv)
	}
	return id
}

type Dep struct {
	Name string `json:"Name"`
}

func (i Dep) ID() string {
	return hashFor(i.Name)
}

func (i Dep) Write(file *bufio.Writer) string {
	id, name := i.ID(), i.Name
	_, cacheHit := cache.LoadOrStore(id, name)
	if !cacheHit {
		file.WriteString(name)
	}
	return id
}

type PERunner struct {
	Name     string     `json:"Name"`
	Type     string     `json:"Type"`
	Exe      *INode     `json:"FullPath"` // PE.Path
	Args     string     `json:"Args"`
	Context  *Principal `json:"Context"` // Principal.Name
	RunLevel string     `json:"RunLevel"`
}

func (r PERunner) ID() string {
	return hashFor(r.Name)
}

func (r PERunner) ToCSV() string {
	fields := make([]string, 8)
	fields[0] = r.ID()
	fields[1] = util.PathFix(r.Name)       // runner name
	fields[2] = r.Type                     // service or task or runkey
	fields[3] = util.PathFix(r.Exe.Path)   // full path to executed exe
	fields[4] = util.PathFix(r.Exe.Name)   // exe name
	fields[5] = util.PathFix(r.Exe.Parent) // exe parent dir
	fields[6] = r.Context.Name             // executin Principal
	fields[7] = r.RunLevel                 // runlevel
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}

func (r PERunner) Write(file *bufio.Writer) string {
	id, csv := r.ID(), r.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, csv)
	if !cacheHit {
		file.WriteString(csv)
	}
	return id
}
