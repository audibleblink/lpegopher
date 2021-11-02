package processor

import (
	"encoding/json"

	"github.com/audibleblink/pegopher/cache"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
)

func CreateRunnerFromJSON(jsonLine []byte) (query *cypher.Query, err error) {
	var runner collectors.PERunner
	err = json.Unmarshal(jsonLine, &runner)
	if err != nil {
		return
	}

	nodeAlias := "d"
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	if !cache.Add(node.Runner, runner.name) {
		props := map[string]string{
			"type": runner.Type,
			"args": runner.Args,
			"exe":  runner.Exe,
		}
		query.Create(
			nodeAlias, node.Runner, "name", runner.Name,
		).Set(
			nodeAlias, props,
		)
	}

	if !cache.Add(node.Principal, runner.Context) {
		query.Merge(
			"", node.Principal, "name", runner.Context,
		)
	}

	if !cache.Add(node.Exe, runner.FullPath) {
		query.Merge(
			"", node.Exe, "path", runner.FullPath,
		)
	}

	if !cache.Add(node.Dir, runner.Parent) {
		query.Merge(
			"", node.Dir, "path", runner.Parent,
		)
	}

	return
}

func RelateRunners(jsonLine []byte) (cypherQ *cypher.Query, err error) {
	log := logerr.Add("runner relation")
	var runner collectors.PERunner

	err = json.Unmarshal(jsonLine, &runner)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	rnr, prcpl, pe, dir := "runner", "principal", "pe", "hostDir"

	cypherQ, err = cypher.NewQuery()
	if err != nil {
		err = logerr.Wrap(err)
		return
	}

	cypherQ.Match(
		rnr, node.Runner, "name", runner.Name,
	).Match(
		prcpl, node.Principal, "name", runner.Context,
	).Relate(
		rnr, "EXECUTES_AS", prcpl,
	).Match(
		pe, node.Exe, "path", runner.FullPath,
	).Relate(
		pe, "EXECUTED_BY", rnr,
	).Match(
		dir, node.Dir, "path", runner.Parent,
	).Relate(
		dir, "HOSTS_PES_FOR", rnr,
	)

	return
}
