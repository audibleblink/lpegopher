package node

import (
	gogm "github.com/mindstand/gogm/v2"
)

const (
	WriteOwner   = "WRITE_OWNER"
	WriteDACL    = "WRITE_DACL"
	WriteProp    = "WRITE_PROP"
	GenericAll   = "GENERIC_ALL"
	GenericWrite = "GENERIC_WRITE"
)

var RelevantRights = map[string]bool{
	WriteOwner:   true,
	WriteDACL:    true,
	WriteProp:    true,
	GenericAll:   true,
	GenericWrite: true,
	// "CONTROL_ACCESS": true,
}

type Principal struct {
	gogm.BaseNode
	Name string `gogm:"unique;name=name"`

	WriteOwnerDir   *Directory `gogm:"direction=outgoing;relationship=WRITE_OWNER"`
	WriteDACLDir    *Directory `gogm:"direction=outgoing;relationship=WRITE_DACL"`
	WritePropDir    *Directory `gogm:"direction=outgoing;relationship=WRITE_PROP"`
	GenericAllDir   *Directory `gogm:"direction=outgoing;relationship=GENERIC_ALL"`
	GenericWriteDir *Directory `gogm:"direction=outgoing;relationship=GENERIC_WRITE"`

	WriteOwnerExe   *EXE `gogm:"direction=outgoing;relationship=WRITE_OWNER"`
	WriteDACLExe    *EXE `gogm:"direction=outgoing;relationship=WRITE_DACL"`
	WritePropExe    *EXE `gogm:"direction=outgoing;relationship=WRITE_PROP"`
	GenericAllExe   *EXE `gogm:"direction=outgoing;relationship=GENERIC_ALL"`
	GenericWriteExe *EXE `gogm:"direction=outgoing;relationship=GENERIC_WRITE"`

	WriteOwnerDll   *DLL `gogm:"direction=outgoing;relationship=WRITE_OWNER"`
	WriteDACLDll    *DLL `gogm:"direction=outgoing;relationship=WRITE_DACL"`
	WritePropDll    *DLL `gogm:"direction=outgoing;relationship=WRITE_PROP"`
	GenericAllDll   *DLL `gogm:"direction=outgoing;relationship=GENERIC_ALL"`
	GenericWriteDll *DLL `gogm:"direction=outgoing;relationship=GENERIC_WRITE"`
}

func (p *Principal) SetPermCanWriteOwner(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.WriteOwnerExe = f
	case *DLL:
		p.WriteOwnerDll = f
	case *Directory:
		p.WriteOwnerDir = f
	}
}

func (p *Principal) SetPermCanWriteDACL(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.WriteOwnerExe = f
	case *DLL:
		p.WriteOwnerDll = f
	case *Directory:
		p.WriteOwnerDir = f
	}
}

func (p *Principal) SetPermCanWriteProp(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.WritePropExe = f
	case *DLL:
		p.WritePropDll = f
	case *Directory:
		p.WritePropDir = f
	}
}

func (p *Principal) SetPermGenericAll(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.GenericAllExe = f
	case *DLL:
		p.GenericAllDll = f
	case *Directory:
		p.GenericAllDir = f
	}
}

func (p *Principal) SetPermGenericWrite(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.GenericWriteExe = f
	case *DLL:
		p.GenericWriteDll = f
	case *Directory:
		p.GenericWriteDir = f
	}
}

//////////////////////////////////
// Edges
//////////////////////////////////

// type EdgeBelongsTo struct {
// 	User  *User
// 	Group *Group
// }
//
// func (e *EdgeBelongsTo) GetStartNode() interface{} {
// 	return e.User
// }
//
// func (e *EdgeBelongsTo) GetStartNodeType() reflect.Type {
// 	return reflect.TypeOf(&User{})
// }
//
// func (e *EdgeBelongsTo) SetStartNode(v interface{}) error {
// 	val, ok := v.(*User)
// 	if !ok {
// 		return fmt.Errorf("unable to cast [%T] to *VertexA", v)
// 	}
//
// 	e.User = val
// 	return nil
// }
//
// func (e *EdgeBelongsTo) GetEndNode() interface{} {
// 	return e.Group
// }
//
// func (e *EdgeBelongsTo) GetEndNodeType() reflect.Type {
// 	return reflect.TypeOf(&Group{})
// }
//
// func (e *EdgeBelongsTo) SetEndNode(v interface{}) error {
// 	val, ok := v.(*Group)
// 	if !ok {
// 		return fmt.Errorf("unable to cast [%T] to *VertexB", v)
// 	}
//
// 	e.Group = val
// 	return nil
// }
