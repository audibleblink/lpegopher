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
