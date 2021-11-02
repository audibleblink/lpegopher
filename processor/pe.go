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

func CreatePEFromJSON(jsonLine []byte) (query *cypher.Query, err error) {
	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		return
	}

	nodeAlias := "d"
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	if !cache.Add(inode.Type, inode.Path) {
		props := map[string]string{
			"name":   inode.Name,
			"parent": inode.Parent,
		}
		query.Create(
			nodeAlias, inode.Type, "path", inode.Path,
		).Set(
			nodeAlias, props,
		)
	}

	if !cache.Add(node.Principal, inode.DACL.Owner) {
		query.Create("", node.Principal, "name", inode.DACL.Owner)
	}

	if !cache.Add(node.Principal, inode.DACL.Group) {
		query.Create("", node.Principal, "name", inode.DACL.Group)
	}

	for _, ace := range inode.DACL.Aces {
		if !cache.Add(node.Principal, ace.Principal) {
			query.Create("", node.Principal, "name", ace.Principal)
		}
	}

	return
}

func RelatePEs(jsonLine []byte) (cypherQ *cypher.Query, err error) {
	log := logerr.Add("pe relation")

	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	pe, dir := "pe", "dir"

	cypherQ, err = cypher.NewQuery()
	if err != nil {
		err = log.Wrap(err)
		return
	}

	cypherQ.Match(
		pe, inode.Type, "path", inode.Path,
	).Match(
		dir, node.Dir, "path", inode.Parent,
	).Relate(
		dir, "CONTAINS", pe,
	)

	id := 0
	for _, ace := range inode.DACL.Aces {
		prnpl := fmt.Sprintf("p%d", id)
		cypherQ.Match(
			prnpl, node.Principal, "name", ace.Principal,
		)

		for _, priv := range ace.Rights {
			if node.AbusableAces[priv] {
				cypherQ.Relate(
					prnpl, priv, pe,
				)
			}
		}
		id++
	}
	return
}
