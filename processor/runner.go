package processor

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
)

func CreateRunnerFromJSON(jsonLine []byte) (err error) {

	var runner collectors.PERunner
	err = json.Unmarshal(jsonLine, &runner)
	if err != nil {
		return
	}

	cypherQ, err := queryForRunner(&runner)
	if err != nil {
		return err
	}

	err = cypherQ.ExecuteW()
	return

}

func queryForRunner(runner *collectors.PERunner) (query *cypher.Query, err error) {
	nodeAlias := "d"
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	query.Merge(
		nodeAlias, node.Runner, "name", runner.Name,
	).Set(
		nodeAlias, "type", runner.Type,
	).Set(
		nodeAlias, "args", runner.Args,
	).Set(
		nodeAlias, "exe", runner.Exe,
	).Merge(
		"", node.Principal, "name", runner.Context,
	).Merge(
		"", node.Exe, "path", runner.FullPath,
	).Merge(
		"", node.Dir, "path", runner.Parent,
	)

	return
}

func RelateRunners(path string) (err error) {
	log := logerr.Add("runner relation")
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var runner collectors.PERunner

	count := 0
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 8*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		count += 1
		text := scanner.Bytes()
		err = json.Unmarshal(text, &runner)
		if err != nil {
			log.Infof("malformed json at line: %d", count)
			continue
		}

		rnr, prcpl, pe, dir := "runner", "principal", "pe", "hostDir"

		cypherQ, err := cypher.NewQuery()
		if err != nil {
			return log.Wrap(err)
		}

		cypherQ.Merge(
			rnr, node.Runner, "name", runner.Name,
		).Merge(
			prcpl, node.Principal, "name", runner.Context,
		).Relate(
			rnr, "EXECUTES_AS", prcpl,
		).Merge(
			pe, node.Exe, "path", runner.FullPath,
		).Relate(
			pe, "EXECUTED_BY", rnr,
		).Merge(
			dir, node.Dir, "path", runner.Parent,
		).Relate(
			dir, "HOSTS_PES_FOR", rnr,
		)

		err = cypherQ.ExecuteW()
		if err != nil {
			log.Infof("error processing line: %d %w", count, err)
			log.Debugf("failed query was: %s", cypherQ.ToString())
			continue
		}
	}
	return
}

// TODO: associations
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
