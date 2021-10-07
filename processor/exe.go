package processor

import (
	"fmt"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
)

func NewExeFromJson(jsonLine []byte) (err error) {
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
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Name")
		return
	}

	if parent, ok = line.Path("Dir").Data().(string); !ok {
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Dir")
		return
	}
	parent = util.PathFix(parent)

	if path, ok = line.Path("Path").Data().(string); !ok {
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Path")
		return
	}
	path = util.PathFix(path)

	if etype, ok = line.Path("Type").Data().(string); !ok {
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Type")
		return
	}

	switch etype {
	case "directory":
		err = doDir(inodeName, parent, path, line)
	case "file":
		err = doExe(inodeName, parent, path, line)
	}
	return
}

func doExe(name, parent, path string, line *gabs.Container) (err error) {
	exe := &node.EXE{}
	err = exe.Merge("path", path)
	if err != nil {
		return
	}

	err = exe.UpsertName(name)
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
	err = parentDir.Add(exe)
	if err != nil {
		return
	}

	return doDACLS(exe, parentDir, line)
}

func doDir(name, parent, path string, line *gabs.Container) (err error) {

	dir := &node.Directory{}
	err = dir.Merge("path", path)
	if err != nil {
		return
	}

	err = dir.UpsertName(name)
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

	err = parentDir.Add(dir)
	if err != nil {
		return
	}

	err = doDACLS(dir, parentDir, line)
	if err != nil {
		return
	}
	return

}

func doDACLS(object interface{}, parent *node.Directory, line *gabs.Container) (err error) {
	aces, err := line.Search("DACL", "Aces").Children()
	if len(aces) == 0 {
		return nil
	}
	if err != nil {
		return
	}

	for _, ace := range aces {
		principalName := ace.Path("Principal").Data().(string)
		principal := &node.User{}
		err = principal.Merge("name", util.Lower(principalName))
		if err != nil {
			return
		}

		rights, err := ace.Path("Rights").Children()
		if err != nil {
			return err
		}

		for _, right := range rights {
			rightStr := right.Data().(string)
			if !node.RelevantRights[rightStr] {
				continue
			}

			switch rightStr {
			case node.WriteOwner:
				principal.SetPermCanWriteOwner(object)
				return principal.Save()
			case node.WriteDACL:
				principal.SetPermCanWriteDACL(object)
				return principal.Save()
			case node.WriteProp:
				principal.SetPermCanWriteProp(object)
				return principal.Save()
			case node.GenericAll:
				principal.SetPermGenericAll(object)
				return principal.Save()
			case node.GenericWrite:
				principal.SetPermGenericWrite(object)
				return principal.Save()
			default:
				continue
			}
		}
	}
	return
}
