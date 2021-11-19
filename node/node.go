package node

const (
	WriteOwner    = "WRITE_OWNER"
	WriteDACL     = "WRITE_DACL"
	GenericAll    = "GENERIC_ALL"
	GenericWrite  = "GENERIC_WRITE"
	ControlAccess = "CONTROL_ACCESS"
)

var AbusableAces = map[string]bool{
	WriteOwner:    true,
	WriteDACL:     true,
	GenericAll:    true,
	GenericWrite:  true,
	ControlAccess: true,
}

const (
	Dll       = "Dll"
	Exe       = "Exe"
	Dir       = "Directory"
	Runner    = "Runner"
	Principal = "Principal"
	Dep       = "Dep"
)

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
