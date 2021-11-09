package processor

import (
	"encoding/json"
	"fmt"

	"github.com/audibleblink/pegopher/cache"
	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
)

func CreatePEFromJSON(jsonLine []byte) (cypherQ *cypher.Query, err error) {
	var inode collectors.INode
	rand := util.Rand()

	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		return
	}

	nVar := fmt.Sprintf("pe_%s", rand)
	cypherQ, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	if cache.Add(inode.Type, inode.Path) {
		props := map[string]string{
			node.Prop.Name:   inode.Name,
			node.Prop.Parent: inode.Parent,
		}
		label := fmt.Sprintf("%s:iNode", inode.Type)
		cypherQ.Create(
			nVar, label, node.Prop.Path, inode.Path,
		).Set(
			nVar, props,
		).With(nVar)
	} else {
		return
	}

	owner := fmt.Sprintf("owner_%s", rand)
	if cache.Add(node.Principal, inode.DACL.Owner) {
		cypherQ.Merge(
			owner, node.Principal, node.Prop.Name, inode.DACL.Owner,
		).Relate(owner, "OWNS", nVar).With(nVar)

	} else {
		cypherQ.Match(
			owner, node.Principal, node.Prop.Name, inode.DACL.Owner,
		).Relate(owner, "OWNS", nVar).With(nVar)

	}

	id := 0
	for _, ace := range inode.DACL.Aces {

		abusables := []string{}
		for _, priv := range ace.Rights {
			if node.AbusableAces[priv] {
				abusables = append(abusables, priv)
			}
		}

		if len(abusables) > 0 {
			pid := fmt.Sprintf("prncpl_%d", id)
			newPrcpl := cache.Add(node.Principal, ace.Principal)
			if newPrcpl {
				cypherQ.Create(
					pid, node.Principal, node.Prop.Name, ace.Principal,
				).With(nVar)
			} else {
				cypherQ.Match(pid, node.Principal, node.Prop.Name, ace.Principal)
			}

			for idx, priv := range abusables {
				cypherQ.Relate(pid, priv, nVar)
				if idx+1 == len(abusables) {
					cypherQ.With(nVar)
				}
			}
			id++
		}
	}

	// impID := fmt.Sprintf("imp%d", id*CurrentBatchLen)
	// for _, imp := range inode.Imports {
	// 	cypherQ.Merge(impID, "Import", "name", imp.Host)
	// 	// imp.Functions
	// 	for idx, fn := range imp.Functions {
	// 		q := fmt.Sprintf(
	// 			"MERGE (%s)-[:%s {fn: %s}]->(%s)",
	// 			inodeAlias,
	// 			"IMPORTS",
	// 			fn,
	// 			imp.Host,
	// 		)
	// 		cypherQ.Append(q)

	// 		if idx+1 == len(imp.Functions) {
	// 			cypherQ.EndMerge()
	// 		}
	// 	}
	// }

	return
}

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
