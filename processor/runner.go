package processor

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
)

func NewRunnerFromJson(jsonLine []byte) (err error) {
	line, err := gabs.ParseJSON(jsonLine)
	if err != nil {
		return
	}

	var (
		userName   string
		exe        string
		parent     string
		runnerName string
		runnerType string
		fullPath   string

		ok bool
	)

	if runnerName, ok = line.Path("Name").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "Name")
		return
	}

	if userName, ok = line.Path("Context").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "userName")
		return
	}

	if exe, ok = line.Path("Exe").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "exeName")
		return
	}

	if fullPath, ok = line.Path("FullPath").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "fullPath")
		return
	}
	fullPath = util.PathFix(fullPath)

	if parent, ok = line.Path("Parent").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "parent")
		return
	}
	parent = util.PathFix(parent)

	if runnerType, ok = line.Path("Type").Data().(string); !ok {
		err = fmt.Errorf("could not create Runner with JSON property: %s", "parent")
		return
	}

	user := &node.User{}
	err = user.Merge("name", util.Lower(userName))
	if err != nil {
		return
	}

	exeNode := &node.EXE{}
	err = exeNode.Merge("path", fullPath)
	if err != nil {
		return
	}
	err = exeNode.SetName(util.Lower(exe))
	if err != nil {
		return
	}

	dir := &node.Directory{}
	dir.Path = parent
	err = dir.Merge("path", dir.Path)
	if err != nil {
		return
	}
	err = dir.SetName(dir.Path)
	if err != nil {
		return
	}

	runner := &node.Runner{}
	err = runner.Merge("name", (runnerName))
	if err != nil {
		return
	}
	err = runner.SetType(runnerType)
	if err != nil {
		return
	}

	// Associate a directory that hosts PEs for a particular Runner
	err = dir.HostsPEsFor(runner)
	if err != nil {
		return
	}

	// Creats the edges for Directories CONTAINS a PE
	err = dir.Add(exeNode)
	if err != nil {
		return
	}

	// Associate the Context (or User) from which the Runner executes
	err = runner.RunsExeAs(user)
	if err != nil {
		return
	}

	// Associate the EXE as getting run by a service
	err = exeNode.GetsRunBy(runner)
	if err != nil {
		return
	}

	return
}
