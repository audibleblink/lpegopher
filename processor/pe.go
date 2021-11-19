package processor

import (
	"fmt"

	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
)

func BulkRelateFileTree() (err error) {
	log := logerr.Add("filetree relationships")

	cypherQ, err := cypher.NewQuery()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	nodeTypes := []string{node.Exe, node.Dir, node.Dll}

	for _, typ := range nodeTypes {

		// relate file/dir heirarchy
		cypherQ.Raw(fmt.Sprintf(`
			CALL apoc.periodic.iterate(
				"MATCH (node:%s),(dir:Directory) WHERE node.parent = dir.path RETURN node,dir",
				"MERGE (dir)-[:CONTAINS]->(node)",
				{batchSize:100, parallel: true, iterateList:true})
			`, typ))
		err = cypherQ.ExecuteW()
		if err != nil {
			err = log.Wrap(err)
			return
		}

	}
	return
}
