package main

import (
	"os"
	"path/filepath"
	"strings"

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

func lower(str string) string {
	return strings.ToLower(str)
}

func pathFix(str string) string {
	str = strings.Trim(str, `"`)
	str = resolveEnvPath(str)
	str = strings.ReplaceAll(str, `\`, "/")
	// swap slack direction to avoid cross-platform issues
	return lower(str)
}

func resolveEnvPath(path string) (out string) {

	// return the original filepath unchanged unless we get to the end
	out = path

	// return unless strings starts with %
	if !strings.HasPrefix(path, "%") {
		return
	}

	// return unless there's a second %
	trim := strings.TrimPrefix(path, "%")
	i := strings.Index(trim, "%")
	if i == -1 {
		return
	}

	// check if substr between two % is the name of an existing env var
	val, ok := os.LookupEnv(trim[:i])
	if !ok {
		return
	}

	// env var value will use os path separator
	remainder := filepath.FromSlash(trim[i+1:])

	// check the remainder starts with path separateor
	if !strings.HasPrefix(remainder, "\\") {
		return
	}

	// prepend the value to the remainder of the path
	return val + remainder
}
