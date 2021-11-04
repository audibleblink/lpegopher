package processor

import (
	"encoding/json"
	"fmt"

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

	nodeAlias := fmt.Sprintf("pe_%d", CurrentBatchLen)
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	if cache.Add(node.Runner, runner.Name) {
		props := map[string]string{
			node.Prop.Type:    runner.Type,
			node.Prop.Args:    runner.Args,
			node.Prop.Exe:     runner.Exe,
			node.Prop.Context: runner.Context,
			node.Prop.Parent:  runner.Parent,
		}
		query.Create(
			nodeAlias, node.Runner, node.Prop.Name, runner.Name,
		).Set(
			nodeAlias, props,
		)
	}

	if cache.Add(node.Principal, runner.Context) {
		query.Create(
			"", node.Principal, node.Prop.Name, runner.Context,
		)
	}

	if cache.Add(node.Exe, runner.FullPath) {
		query.Create(
			"", node.Exe, node.Prop.Path, runner.FullPath,
		)
	}

	if cache.Add(node.Dir, runner.Parent) {
		query.Create(
			"", node.Dir, node.Prop.Path, runner.Parent,
		)
	}

	return
}

func BulkRelateRunners() (err error) {
	log := logerr.Add("runner relation")

	cypherQ, err := cypher.NewQuery()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	// relate dirs that hosts a runner exe
	cypherQ.Raw(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(dir:Directory) WHERE r.parent = dir.path RETURN r,dir",
		"MERGE (dir)-[:HOSTS_PES_FOR]->(r)",
		{batchSize:100, parallel: true, iterateList:true})
	`)
	err = cypherQ.ExecuteW()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	// relate principals that run certain runners
	cypherQ.Raw(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(p:Principal) WHERE r.context = p.name RETURN r,p",
		"MERGE (r)-[:EXECUTES_AS]->(p)",
		{batchSize:100, parallel: true, iterateList:true})
	`)
	err = cypherQ.ExecuteW()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	// relate exes that are executed by a runner
	cypherQ.Raw(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(exe:Exe) WHERE r.path = exe.path RETURN r,exe",
		"MERGE (exe)-[:EXECUTED_FROM]->(r)",
		{batchSize:100, parallel: true, iterateList:true})
	`)
	err = cypherQ.ExecuteW()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	return
}
