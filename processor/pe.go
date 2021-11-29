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
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	log.Info("processing dlls")
	query = ` LOAD CSV FROM'file:////dlls.csv' AS line
		WITH line
		CREATE (:Dll:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`
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
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`
	err = execString(query)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Info("processing forwards")
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

func RelateOwnership() (err error) {
	log := logerr.Add("ownership creation")
	log.Info("relating all (:Principal)-[:OWNS]-(:INode)")
	err = execString(`
			CALL apoc.periodic.iterate("
				MATCH (pcpl:Principal),(inode:INode) WHERE pcpl.nid = inode.owner or pcpl.nid = inode.group RETURN pcpl, inode
			","
				MERGE (pcpl)-[:OWNS]->(inode)
			", {batchSize:1000})
			`)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}

func RelateACLs() (err error) {
	log := logerr.Add("acl relationships")
	log.Info("relating all (:Principal)-[$ACE]-(:INodes)")
	err = execString(`
		CALL apoc.periodic.iterate("
			CALL apoc.load.csv('relationships.csv',{}) yield list as line return line
		","
			MATCH (a {id: line[0]}), (b {id: line[2]})
			CALL apoc.create.relationship(a, line[1], {}, b) YIELD rel RETURN *
		", {batchSize:1000, iterateList:true, parallel:true});
		`)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}
