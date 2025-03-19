package processor

import (
	"fmt"

	"github.com/audibleblink/lpegopher/logerr"
	"github.com/audibleblink/lpegopher/node"
)

// InsertAllNodes loads node data into the graph database
func InsertAllNodes(stageURL string) (err error) {
	log := logerr.Add("file inserts")

	log.Debug("processing exes")
	query := `LOAD CSV FROM '%s/exes.csv' AS line
		WITH line
		CREATE (:Exe:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`

	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}
	log.Debug("processing dlls")
	query = `LOAD CSV FROM '%s/dlls.csv' AS line
		WITH line
		CREATE (:Dll:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Debug("processing dirs")
	query = `LOAD CSV FROM '%s/dirs.csv' AS line
		WITH line
		CREATE (:Directory:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Debug("processing forwards")
	query = `LOAD CSV FROM '%s/deps.csv' AS line
		WITH line CREATE (:Dep {nid: line[0], name: line[1]})`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Debug("processing principals")
	query = `LOAD CSV FROM '%s/principals.csv' AS line
		WITH line CREATE (:Principal {nid: line[0], name: line[1], group: line[2]})`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}

	return
}

// BulkRelateFileTree creates relationships between files and directories
func BulkRelateFileTree() (err error) {
	log := logerr.Add("filetree relationships")
	for _, typ := range []string{node.Dir, node.Exe, node.Dll} {
		log.Debugf("relating all (:Dir)-[:CONTAINS]-(:%s)", typ)
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

// RelateOwnership creates ownership relationships between principals and nodes
func RelateOwnership() (err error) {
	log := logerr.Add("ownership creation")
	log.Debug("relating all (:Principal)-[:OWNS]-(:INode)")
	err = execString(`
			CALL apoc.periodic.iterate("
				MATCH (pcpl:Principal),(inode:INode) 
				WHERE pcpl.nid = inode.owner or pcpl.nid = inode.group 
				RETURN pcpl, inode
			","
				MERGE (pcpl)-[:OWNS]->(inode)
			", {batchSize: 1000})
			`)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}

// RelateMembership creates group membership relationships between principals
func RelateMembership() (err error) {
	log := logerr.Add("membership creation")
	log.Debug("relating all (:Principal)-[:MEMBER_OF]-(:Principal)")
	err = execString(`
			CALL apoc.periodic.iterate("
				MATCH (group:Principal),(user:Principal) 
				WHERE user.group = group.name 
				RETURN user, group
			","
				MERGE (user)-[:MEMBER_OF]->(group)
			", {batchSize: 10})
			`)
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}

// RelateACLs creates access control relationships between nodes
func RelateACLs(stageURL string) (err error) {
	log := logerr.Add("acl relationships")
	log.Debug("relating all (:Principal)-[$ACE]-(:INodes)")
	query := `CALL apoc.periodic.iterate("
			LOAD CSV FROM '%s/relationships.csv' AS line RETURN line
		","
			MATCH (a:Principal {nid: line[0]}), (b:INode {nid: line[2]})
			CALL apoc.create.relationship(a, line[1], {}, b) YIELD rel RETURN rel
		", {batchSize: 20000});
		`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}

// RelateDependecies creates dependency relationships between nodes
func RelateDependecies(stageURL string) (err error) {
	log := logerr.Add("dependecy relationships")
	log.Debug("relating (:INode)-[:IMPORTS]-(:Dep)")
	query := `CALL apoc.periodic.iterate("
			LOAD CSV FROM '%s/imports.csv' AS line RETURN line
		","
			MATCH (a:INode {nid: line[0]}), (b:Dep {nid: line[2]})
			MERGE (b)-[:IMPORTED_BY]->(a)
		", {batchSize: 20000});
		`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		err = log.Wrap(err)
		return
	}
	return
}
