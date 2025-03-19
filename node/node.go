package node

// Abusable ACE privilege constants
const (
	WriteOwner    = "WRITE_OWNER"    // Permission to change ownership
	WriteDACL     = "WRITE_DACL"     // Permission to modify access control list
	GenericAll    = "GENERIC_ALL"    // Full control permission
	GenericWrite  = "GENERIC_WRITE"  // Write permission
	ControlAccess = "CONTROL_ACCESS" // Control access permission
)

// AbusableAces maps privilege names to a boolean indicating they are abusable
var AbusableAces = map[string]bool{
	WriteOwner:    true,
	WriteDACL:     true,
	GenericAll:    true,
	GenericWrite:  true,
	ControlAccess: true,
}

// Node type constants
const (
	Dll       = "Dll"       // Dynamic Link Library
	Exe       = "Exe"       // Executable
	Dir       = "Directory" // Directory
	Runner    = "Runner"    // Auto-executing program
	Principal = "Principal" // Security principal
	Dep       = "Dep"       // Dependency
)

// Prop contains property name constants for nodes
var Prop = struct {
	Name    string
	Dir     string
	Parent  string
	Path    string
	Type    string
	Args    string
	Exe     string
	Context string
}{
	"name",
	"dir",
	"parent",
	"path",
	"type",
	"args",
	"exe",
	"context",
}
