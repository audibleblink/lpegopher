package collectors

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/audibleblink/pegopher/util"
	"github.com/minio/highwayhash"
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
	Forwards   = "FORWARDS"
	ImportedBy = "IMPORTED_BY"

	Null = "NULL"
)

// INode contains the parsed import and exports of the INode
type INode struct {
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	Parent   string `json:"Dir"`
	Type     string `json:"Type"`
	Forwards []*Dep `json:"Forwards"`
	Imports  []*Dep `json:"Imports"`
	DACL     DACL   `json:"DACL"`

	id string
}

func (i INode) ID() string {
	if i.id != "" {
		return i.id
	}
	i.id = hashFor(i.Path)
	return i.id
}

func (i INode) Write(file io.Writer) string {
	id := i.ID()
	csv := i.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, i.Path)
	if !cacheHit {
		io.WriteString(file, csv)
	}
	return id
}

func (i INode) ToCSV() string {
	o := Null
	g := Null
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
	id   string
}

func (p Principal) ID() string {
	if p.id != "" {
		return p.id
	}
	p.id = hashFor(p.Name)
	return p.id
}

func (p Principal) ToCSV() string {
	fields := make([]string, 6)
	fields[0] = p.ID()
	fields[1] = util.PathFix(p.Name)
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}

func (p Principal) Write(file io.Writer) string {
	id, csv := p.ID(), p.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, p.Name)
	if !cacheHit {
		io.WriteString(file, csv)
	}
	return id
}

type Rel struct {
	Start string
	Rel   string
	End   string

	id string
}

func (r Rel) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s\n", r.Start, r.Rel, r.End)
}

func (r Rel) ID() string {
	if r.id != "" {
		return r.id
	}
	r.id = hashFor(r.ToCSV())
	return r.id
}

func (r Rel) Write(file io.Writer) string {
	id, csv := r.ID(), r.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, csv)
	if !cacheHit {
		io.WriteString(file, csv)
	}
	return id
}

type Dep struct {
	Name string `json:"Name"`
	id   string
}

func (d Dep) ID() string {
	if d.id != "" {
		return d.id
	}
	d.id = hashFor(d.Name)
	return d.id
}

func (d Dep) ToCSV() string {
	return fmt.Sprintf("%s,%s\n", d.ID(), d.Name)
}

func (i Dep) Write(file io.Writer) string {
	id, csv := i.ID(), i.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, i.Name)
	if !cacheHit {
		io.WriteString(file, csv)
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

	id string
}

func (r PERunner) ID() string {
	if r.id != "" {
		return r.id
	}
	r.id = hashFor(fmt.Sprintf("%s:%s", r.Type, r.Name))
	return r.id
}

func (r PERunner) ToCSV() string {
	fields := make([]string, 8)
	fields[0] = r.ID()
	fields[1] = util.PathFix(r.Name)       // runner name
	fields[2] = r.Type                     // service or task or runkey
	fields[3] = util.PathFix(r.Exe.Path)   // full path to executed exe
	fields[4] = util.PathFix(r.Exe.Name)   // exe name
	fields[5] = util.PathFix(r.Exe.Parent) // exe parent dir
	fields[6] = util.Lower(r.Context.Name) // executin Principal
	fields[7] = r.RunLevel                 // runlevel
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}

func (r PERunner) Write(file io.Writer) string {
	id, csv := r.ID(), r.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, r.Name)
	if !cacheHit {
		io.WriteString(file, csv)
	}
	return id
}

func hashFor(data string) string {
	data = util.PathFix(data)
	hash, err := highwayhash.New(key)
	if err != nil {
		fmt.Printf("Failed to create HighwayHash instance: %v", err)
		os.Exit(1)
	}

	txt := strings.NewReader(data)
	if _, err = io.Copy(hash, txt); err != nil {
		fmt.Printf("hash reader creation failed: %v", err)
		os.Exit(1)
	}

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}
