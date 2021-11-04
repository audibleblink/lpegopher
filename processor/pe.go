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

func CreatePEFromJSON(jsonLine []byte) (cypherQ *cypher.Query, err error) {
	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		return
	}

	inodeAlias := fmt.Sprintf("pe%d", CurrentBatchLen)
	cypherQ, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	if cache.Add(inode.Type, inode.Path) {
		props := map[string]string{
			node.Prop.Name:   inode.Name,
			node.Prop.Parent: inode.Parent,
		}
		cypherQ.Create(
			inodeAlias, inode.Type, node.Prop.Path, inode.Path,
		).Set(
			inodeAlias, props,
		)
	}

	// if cache.Add(node.Principal, inode.DACL.Owner) {
	// 	cypherQ.Create("", node.Principal, node.Prop.Name, inode.DACL.Owner)
	// }

	// if cache.Add(node.Principal, inode.DACL.Group) {
	// 	cypherQ.Create("", node.Principal, node.Prop.Name, inode.DACL.Group)
	// }
	count := 0
	for _, ace := range inode.DACL.Aces {
		prince := fmt.Sprintf("prpl%d", count)
		cypherQ.Merge(prince, node.Principal, node.Prop.Name, ace.Principal)

		abusables := []string{}
		for _, priv := range ace.Rights {
			if node.AbusableAces[priv] {
				abusables = append(abusables, priv)
			}
		}

		if len(abusables) > 0 {
			for idx, priv := range abusables {
				cypherQ.Relate(prince, priv, inodeAlias)
				if idx+1 == len(abusables) {
					cypherQ.EndMerge()
				}
			}
		}
		count++
	}

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

func RelatePEs(jsonLine []byte) (cypherQ *cypher.Query, err error) {
	log := logerr.Add("pe relation")

	MaxBatchSize = 10

	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	nLbl := fmt.Sprintf("pe%d", CurrentBatchLen)

	cypherQ, err = cypher.NewQuery()
	if err != nil {
		err = log.Wrap(err)
		return
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
			if id == 0 {
				// Only match the inode if we're going to merge abusables
				// since we can't 'undo' the cypherQ Match if we do it
				// earlier
				cypherQ.Match(nLbl, inode.Type, node.Prop.Path, inode.Path)
			}
			prnpl := fmt.Sprintf("p_%d", id)
			cypherQ.Match(prnpl, node.Principal, node.Prop.Name, ace.Principal)

			for idx, priv := range abusables {
				cypherQ.Relate(prnpl, priv, nLbl)
				if idx+1 == len(abusables) {
					cypherQ.EndMerge()
				}
			}
			id++
		}

	}
	return
}
