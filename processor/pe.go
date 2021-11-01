package processor

import (
	"encoding/json"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
)

func CreatePEFromJSON(jsonLine []byte) (err error) {

	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		return
	}

	cypherQ, err := queryForINode(&inode)
	if err != nil {
		return err
	}

	err = cypherQ.ExecuteW()
	return
}

func queryForINode(inode *collectors.INode) (query *cypher.Query, err error) {
	nodeAlias := "d"
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	query.Merge(
		nodeAlias, inode.Type, "path", inode.Path,
	).Set(
		nodeAlias, "name", inode.Name,
	).Set(
		nodeAlias, "parent", inode.Parent,
	).Merge(
		"", collectors.Principal, "name", inode.DACL.Owner,
	).Merge(
		"", collectors.Principal, "name", inode.DACL.Group,
	)

	for _, ace := range inode.DACL.Aces {
		query.Merge("", collectors.Principal, "name", ace.Principal)
	}

	return
}
