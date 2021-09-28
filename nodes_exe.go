package main

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
)

type EXE struct {
	securableIFile
	containableIFile

	ExecutesFrom *Runner `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
	// Imports      []*containableIFile   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*containableIFile   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
}

func (f *EXE) GetsRunBy(runner *Runner) error {
	f.ExecutesFrom = runner
	return f.save()
}

func (x *EXE) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "EXE"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, nil, x)
}

func (u *EXE) SetName(name string) error {
	u.Name = name
	return u.save()
}

func (x *EXE) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := newNeoSession()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}

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
		// ownerUser  string
		// ownerGroup string

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
	parent = pathFix(parent)

	if path, ok = line.Path("Path").Data().(string); !ok {
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Path")
		return
	}
	path = pathFix(path)

	if etype, ok = line.Path("Type").Data().(string); !ok {
		err = fmt.Errorf("could not create EXE with JSON property: %s", "Type")
		return
	}

	// if ownerUser, ok = line.Path("DACL.Owner").Data().(string); !ok {
	// 	err = fmt.Errorf("could not create EXE with JSON property: %s", "DACL.Owner")
	// 	return
	// }

	// if ownerGroup, ok = line.Path("DACL.Group").Data().(string); !ok {
	// 	err = fmt.Errorf("could not create EXE with JSON property: %s", "DACL.Group")
	// 	return
	// }

	switch etype {
	case "directory":
		return doDir(inodeName, parent, path, line)
	case "file":
		return doFile(inodeName, parent, path, line)
	}
	return
}

func doFile(name, parent, path string, line *gabs.Container) (err error) {

	exe := &EXE{}
	err = exe.Merge("path", path)
	if err != nil {
		return
	}
	err = exe.SetName(name)
	if err != nil {
		return
	}

	dir := &Directory{}
	err = dir.Merge("path", parent)
	if err != nil {
		return
	}
	err = dir.Add(exe)
	if err != nil {
		return
	}

	for _, ace := range line.Search("DACL", "Aces").Children() {
		principalName := ace.Path("Principal").Data().(string)
		principal := &User{}
		err = principal.Merge("name", lower(principalName))
		if err != nil {
			return
		}

		for _, right := range ace.Path("Rights").Children() {
			rightStr := right.Data().(string)
			if !RelevantRights[rightStr] {
				continue
			}

			switch rightStr {
			case WriteOwner:
				principal.SetPermCanWriteOwner(exe)
				return principal.save()
			case WriteDACL:
				principal.SetPermCanWriteDACL(exe)
				return principal.save()
			case WriteProp:
				principal.SetPermCanWriteProp(exe)
				return principal.save()
			case GenericAll:
				principal.SetPermGenericAll(exe)
				return principal.save()
			case GenericWrite:
				principal.SetPermGenericWrite(exe)
				return principal.save()
			default:
				continue
			}
		}
	}
	return
}

func doDir(name, parent, path string, line *gabs.Container) (err error) {

	dir := &Directory{}
	err = dir.Merge("path", path)
	if err != nil {
		return
	}
	err = dir.SetName(name)
	if err != nil {
		return
	}

	parentDir := &Directory{}
	err = parentDir.Merge("path", parent)
	if err != nil {
		return
	}
	err = parentDir.Add(dir)
	if err != nil {
		return
	}

	for _, ace := range line.Search("DACL", "Aces").Children() {
		principalName := ace.Path("Principal").Data().(string)
		principal := &User{}
		err = principal.Merge("name", lower(principalName))
		if err != nil {
			return
		}

		for _, right := range ace.Path("Rights").Children() {
			rightStr := right.Data().(string)
			if !RelevantRights[rightStr] {
				continue
			}

			switch rightStr {
			case WriteOwner:
				principal.SetPermCanWriteOwner(dir)
				return principal.save()
			case WriteDACL:
				principal.SetPermCanWriteDACL(dir)
				return principal.save()
			case WriteProp:
				principal.SetPermCanWriteProp(dir)
				return principal.save()
			case GenericAll:
				principal.SetPermGenericAll(dir)
				return principal.save()
			case GenericWrite:
				principal.SetPermGenericWrite(dir)
				return principal.save()
			default:
				continue
			}
		}
	}
	return
}
