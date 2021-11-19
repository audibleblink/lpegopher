package processor

import (
	"fmt"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
)

func InsertAllNodes() (err error) {
	log := logerr.Add("file inserts")

	log.Info("processing exes")
	query := ` LOAD CSV FROM'file:////exes.csv' AS line
		WITH line
		CREATE (:Exe:INode {
			nid: line[0], 
			name: line[1],
			type: line[2],
			path: line[3],
			exe: line[4],
			parent: line[5],
			context: line[6],
			runlevel: line[7]})`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	log.Info("processing exes")
	query = ` LOAD CSV FROM'file:////dlls.csv' AS line
		WITH line
		CREATE (:Dll:INode {
			nid: line[0], 
			name: line[1],
			type: line[2],
			path: line[3],
			exe: line[4],
			parent: line[5],
			context: line[6],
			runlevel: line[7]})`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Info("processing dirs")
	query = ` LOAD CSV FROM'file:////dirs.csv' AS line
		WITH line
		CREATE (:Directory:INode {
			nid: line[0], 
			name: line[1],
			type: line[2],
			path: line[3],
			exe: line[4],
			parent: line[5],
			context: line[6],
			runlevel: line[7]})`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Info("processing dependencies")
	query = ` 
	LOAD CSV FROM'file:////deps.csv' AS line
		WITH line CREATE (:Dep {nid: line[0], name: line[1]})`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Info("processing principals")
	query = ` 
	LOAD CSV FROM'file:////principals.csv' AS line
		WITH line CREATE (:Principal {nid: line[0], name: line[1]})`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	return
}

func BulkRelateFileTree() (err error) {
	log := logerr.Add("filetree relationships")
	for _, typ := range []string{node.Dir, node.Exe, node.Dll} {
		log.Infof("relating all (:Dir)-[:CONTAINS]-(:%s)", typ)
		err = execString(fmt.Sprintf(`
			CALL apoc.periodic.iterate(
				"MATCH (node:%s),(dir:Directory) WHERE node.parent = dir.path RETURN node,dir",
				"MERGE (dir)-[:CONTAINS]->(node)",
				{batchSize:1000})
			`, typ))
		if err != nil {
			err = log.Wrap(err)
			return
		}
	}
	return
}
