package processor

import (
	"fmt"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
)

func NewDllFromJson(jsonLine []byte) (err error) {
	line, err := gabs.ParseJSON(jsonLine)
	if err != nil {
		return
	}

	var (
		inodeName string
		path      string
		etype     string
		parent    string

		ok bool
	)

	if inodeName, ok = line.Path("Name").Data().(string); !ok {
		err = fmt.Errorf("could not create DLL with JSON property: %s", "Name")
		return
	}

	if parent, ok = line.Path("Dir").Data().(string); !ok {
		err = fmt.Errorf("could not create DLL with JSON property: %s", "Dir")
		return
	}
	parent = util.PathFix(parent)

	if path, ok = line.Path("Path").Data().(string); !ok {
		err = fmt.Errorf("could not create DLL with JSON property: %s", "Path")
		return
	}
	path = util.PathFix(path)

	if etype, ok = line.Path("Type").Data().(string); !ok {
		err = fmt.Errorf("could not create DLL with JSON property: %s", "Type")
		return
	}

	switch etype {
	case "directory":
		err = doDir(inodeName, parent, path, line)
	case "file":
		err = doDll(inodeName, parent, path, line)
	}
	return
}

func doDll(name, parent, path string, line *gabs.Container) (err error) {
	dll := &node.DLL{}
	err = dll.Merge("path", path)
	if err != nil {
		return
	}

	err = dll.UpsertName(name)
	if err != nil {
		return
	}

	parentDir := &node.Directory{}
	err = parentDir.Merge("path", parent)
	if err != nil {
		return
	}

	err = parentDir.UpsertName(filepath.Base(parent))
	if err != nil {
		return
	}

	// create the CONTAINS relationship
	err = parentDir.Add(dll)
	if err != nil {
		return
	}

	return doDACLS(dll, parentDir, line)
}
