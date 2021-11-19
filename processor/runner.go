package processor

import (
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
)

func execString(query string) error {
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		return err
	}
	cypherQ.Raw(query)
	return cypherQ.ExecuteW()
}

func InsertAllRunners() (err error) {
	log := logerr.Add("runner inserts")
	err = execString(`
	LOAD CSV FROM'file:////runners.csv' AS line
	WITH line
	CREATE (e:Runner {
		id: line[0], 
		name: line[1],
		type: line[2],
		path: line[3],
		exe: line[4],
		parent: line[5],
		context: line[6],
		runlevel: line[7]})
	`)
	if err != nil {
		err = log.Wrap(err)
	}
	return
}

func BulkRelateRunners() (err error) {
	log := logerr.Add("runner relationships")

	// relate dirs that hosts a runner exe
	log.Infof("relating all (:Dir)-[:HOSTS_PES_FOR]->(:Runner)")
	err = execString(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(dir:Directory) WHERE r.parent = dir.path RETURN r,dir",
		"MERGE (dir)-[:HOSTS_PES_FOR]->(r)",
		{batchSize:100, parallel: true, iterateList:true})
	`)
	if err != nil {
		return log.Wrap(err)
	}

	// relate principals that run certain runners
	log.Infof("relating all (:Runner)-[:RUNS_AS]->(:Principal)")
	err = execString(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(p:Principal) WHERE r.context = p.name RETURN r,p",
		"MERGE (r)-[:RUNS_AS]->(p)",
		{batchSize:1000})
	`)
	if err != nil {
		return log.Wrap(err)
	}

	// relate exes that are executed by a runner
	log.Infof("relating all (:Exe)-[:EXECUTED_BY]->(:Runner)")
	err = execString(`
	CALL apoc.periodic.iterate(
		"MATCH (r:Runner),(exe:Exe) WHERE r.parent+'/'+r.exe = exe.path RETURN r,exe",
		"MERGE (exe)-[:EXECUTED_BY]->(r)",
		{batchSize:100})
	`)

	return
}
