package collectors

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/minio/highwayhash"

	"github.com/audibleblink/lpegopher/logerr"
	"github.com/audibleblink/lpegopher/util"
)

// Writer defines the interface for types that can be written to output streams
type Writer interface {
	// ID returns a unique identifier for the item
	ID() string

	// ToCSV converts the item to a CSV formatted string
	ToCSV() string

	// Write outputs the item to the given writer and returns its ID
	Write(io.Writer) string
}

// KeyedWriter is a helper type for types that need caching by a specific key
type KeyedWriter interface {
	Writer
	// CacheKey returns the key used for caching this item
	CacheKey() string
}

// GenericWriteOp is a generic function that handles the common Write pattern
// for all collector types. T must implement the Writer interface.
func GenericWriteOp[T Writer](item T, file io.Writer, cacheKey string) string {
	id, csv := item.ID(), item.ToCSV()
	_, cacheHit := cache.LoadOrStore(id, cacheKey)
	if !cacheHit {
		_, err := io.WriteString(file, csv)
		if err != nil {
			return ""
		}
	}
	return id
}

// WriteToFile is a convenience function for types that implement KeyedWriter
func WriteToFile[T KeyedWriter](item T, file io.Writer) string {
	return GenericWriteOp(item, file, item.CacheKey())
}

// WriteItems is a generic function to write a batch of items
func WriteItems[T KeyedWriter](items []T, file io.Writer) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = WriteToFile(item, file)
	}
	return ids
}

// Constants for relationship types
const (
	WriteOwner    = "WRITE_OWNER"    // Permission to change the owner of an object
	WriteDACL     = "WRITE_DACL"     // Permission to modify the discretionary access control list
	GenericAll    = "GENERIC_ALL"    // Full control permission
	GenericWrite  = "GENERIC_WRITE"  // Write permission
	ControlAccess = "CONTROL_ACCESS" // Right to control access to an object
	Owns          = "OWNS"           // Ownership relationship

	HostsPEFor = "HOSTS_PE_FOR" // Host-PE relationship
	Contains   = "CONTAINS"     // Container relationship
	ExecutedBy = "EXECUTED_BY"  // Execution relationship
	RunsAs     = "RUNS_AS"      // Execution context relationship

	Imports    = "IMPORTS"     // Import relationship
	Forwards   = "FORWARDS"    // Forwarding relationship
	ImportedBy = "IMPORTED_BY" // Reverse import relationship

	Null = "NULL" // Null or empty value
)

// INode contains the parsed import and exports of a node
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

// ID returns the unique identifier for an INode
func (i INode) ID() string {
	if i.id != "" {
		return i.id
	}
	i.id = hashFor(i.Path)
	return i.id
}

// CacheKey returns the key to use for caching an INode
func (i INode) CacheKey() string {
	return i.Path
}

// Write outputs the INode data to the provided writer and returns its ID
func (i INode) Write(file io.Writer) string {
	return GenericWriteOp(i, file, i.CacheKey())
}

// ToCSV converts the INode to a CSV formatted string
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

// DACL represents a Discretionary Access Control List
type DACL struct {
	Owner *Principal    `json:"Owner"`
	Group *Principal    `json:"Group"`
	Aces  []ReadableAce `json:"Aces"`
}

// ReadableAce represents a readable access control entry
type ReadableAce struct {
	Principal *Principal `json:"Principal"`
	Rights    []string   `json:"Rights"`
}

// Principal represents Users or Groups in access control
type Principal struct {
	Name  string `json:"Name"`
	Group string `json:"Group"`
	Type  string `json:"Type"`
	id    string
}

// ID returns the unique identifier for a Principal
func (p Principal) ID() string {
	if p.id != "" {
		return p.id
	}
	p.id = hashFor(p.Name)
	return p.id
}

// CacheKey returns the key to use for caching a Principal
func (p Principal) CacheKey() string {
	return p.Name
}

// ToCSV converts the Principal to a CSV formatted string
func (p Principal) ToCSV() string {
	fields := make([]string, 6)
	fields[0] = p.ID()
	fields[1] = util.PathFix(p.Name)
	fields[2] = util.PathFix(p.Group)
	fields[3] = p.Type
	row := fmt.Sprintf("%s\n", strings.Join(fields, ","))
	return row
}

// Write outputs the Principal data to the provided writer and returns its ID
func (p Principal) Write(file io.Writer) string {
	return GenericWriteOp(p, file, p.CacheKey())
}

// Rel represents a relationship between two entities
type Rel struct {
	Start string
	Rel   string
	End   string

	id string
}

// ToCSV converts the relationship to a CSV formatted string
func (r Rel) ToCSV() string {
	return fmt.Sprintf("%s,%s,%s\n", r.Start, r.Rel, r.End)
}

// ID returns the unique identifier for a relationship
func (r Rel) ID() string {
	if r.id != "" {
		return r.id
	}
	r.id = hashFor(r.ToCSV())
	return r.id
}

// CacheKey returns the key to use for caching a Rel
func (r Rel) CacheKey() string {
	return r.ToCSV()
}

// Write outputs the relationship data to the provided writer and returns its ID
func (r Rel) Write(file io.Writer) string {
	return GenericWriteOp(r, file, r.CacheKey())
}

// Dep represents a dependency with a name
type Dep struct {
	Name string `json:"Name"`
	id   string
}

// ID returns the unique identifier for a dependency
func (d Dep) ID() string {
	if d.id != "" {
		return d.id
	}
	d.id = hashFor(d.Name)
	return d.id
}

// CacheKey returns the key to use for caching a Dep
func (d Dep) CacheKey() string {
	return d.Name
}

// ToCSV converts the dependency to a CSV formatted string
func (d Dep) ToCSV() string {
	return fmt.Sprintf("%s,%s\n", d.ID(), d.Name)
}

// Write outputs the dependency data to the provided writer and returns its ID
func (d Dep) Write(file io.Writer) string {
	return GenericWriteOp(d, file, d.CacheKey())
}

// PERunner represents an executable runner such as a service or scheduled task
type PERunner struct {
	Name     string     `json:"Name"`
	Type     string     `json:"Type"`
	Exe      *INode     `json:"FullPath"` // PE.Path
	Args     string     `json:"Args"`
	Context  *Principal `json:"Context"` // Principal.Name
	RunLevel string     `json:"RunLevel"`

	id string
}

// ID returns the unique identifier for a PERunner
func (r PERunner) ID() string {
	if r.id != "" {
		return r.id
	}
	r.id = hashFor(fmt.Sprintf("%s:%s", r.Type, r.Name))
	return r.id
}

// CacheKey returns the key to use for caching a PERunner
func (r PERunner) CacheKey() string {
	return r.Name
}

// ToCSV converts the PERunner to a CSV formatted string
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

// Write outputs the PERunner data to the provided writer and returns its ID
func (r PERunner) Write(file io.Writer) string {
	return GenericWriteOp(r, file, r.CacheKey())
}

// hashFor generates a hash string for the given data after normalizing with PathFix
// This is a more efficient implementation that writes directly to the hash
func hashFor(data string) string {
	// Normalize the data using PathFix
	data = util.PathFix(data)
	
	// Create a new hash instance
	hash, err := highwayhash.New(key)
	if err != nil {
		// Use logerr package instead of direct fmt.Printf and don't exit
		// This could be enhanced further to return an error, but keeping signature compatible
		log := logerr.Add("hash")
		log.Errorf("Failed to create HighwayHash instance: %v", err)
		return "" // Return empty string instead of crashing the program
	}

	// Write data directly to hash - much more efficient than using io.Copy and strings.NewReader
	_, err = hash.Write([]byte(data))
	if err != nil {
		log := logerr.Add("hash")
		log.Errorf("Hash write failed: %v", err)
		return ""
	}
	
	// Get the checksum and encode to hex string
	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}
