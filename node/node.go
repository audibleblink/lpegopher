package node

const (
	WriteOwner   = "WRITE_OWNER"
	WriteDACL    = "WRITE_DACL"
	WriteProp    = "WRITE_PROP"
	GenericAll   = "GENERIC_ALL"
	GenericWrite = "GENERIC_WRITE"
)

var AbusableAces = map[string]bool{
	WriteOwner:   true,
	WriteDACL:    true,
	WriteProp:    true,
	GenericAll:   true,
	GenericWrite: true,
	// "CONTROL_ACCESS": true,
}

const (
	Dll       = "Dll"
	Exe       = "Exe"
	Dir       = "Directory"
	Runner    = "Runner"
	Principal = "Principal"
)
