package main

import (
	"context"
	"fmt"

	gogm "github.com/mindstand/gogm/v2"
)

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

func (p *Principal) CanWriteOwner(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.WriteOwnerExe = f
	case *DLL:
		p.WriteOwnerDll = f
	case *Directory:
		p.WriteOwnerDir = f
	}
}

func (p *Principal) CanWriteDACL(ifile interface{}) {
	switch f := ifile.(type) {
	case *EXE:
		p.WriteOwnerExe = f
	case *DLL:
		p.WriteOwnerDll = f
	case *Directory:
		p.WriteOwnerDir = f
	}
}

type User struct {
	Principal

	Groups       []*Group `gogm:"direction=outgoing;relationship=MEMBER_OF"`
	ExecutedFrom *Runner  `gogm:"direction=incoming;relationship=EXECUTES_AS"`
}

func (x *User) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "User"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, nil, x)
}

func (u *User) JoinGroup(group *Group) {
	u.Groups = append(u.Groups, group)
}

type Group struct {
	Principal

	Users []*User `gogm:"direction=incoming;relationship=MEMBER_OF"`
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
