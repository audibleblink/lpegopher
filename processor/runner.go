package processor

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
)

func NewRunnerFromJson(jsonLine []byte) (err error) {

	var runner collectors.PERunner
	err = json.Unmarshal(jsonLine, &runner)
	if err != nil {
		return
	}

	query := cypher.NewQuery()

	query.Merge(
		"r", collectors.Runner, "name", runner.Name,
	).Set(
		"r", "type", runner.Type,
	).Set(
		"r", "args", runner.Args,
	).Merge(
		"pe", collectors.Exe, "path", runner.FullPath,
	).Set(
		"pe", "name", filepath.Base(runner.FullPath),
	).Set(
		"pe", "parent", runner.Parent,
	).Merge(
		"p", collectors.Principal, "name", runner.Context,
	)

	fmt.Println(query.ToString())
	return

	// Associate a directory that hosts PEs for a particular Runner
	// err = dir.HostsPEsFor(runner)
	// if err != nil {
	// 	return
	// }

	// Creats the edges for Directories CONTAINS a PE
	// err = dir.Add(exeNode)
	// if err != nil {
	// 	return
	// }

	// Associate the Context (or User) from which the Runner executes
	// err = runner.RunsExeAs(user)
	// if err != nil {
	// 	return
	// }

	// Associate the EXE as getting run by a service
	// err = exeNode.GetsRunBy(runner)
	// if err != nil {
	// 	return
	// }
}
