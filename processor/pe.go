package processor

import (
	"encoding/json"
	"fmt"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
)

func NewPEFromJSON(jsonLine []byte) (err error) {

	var inode collectors.INode
	err = json.Unmarshal(jsonLine, &inode)
	if err != nil {
		return
	}

	cypherQ := queryForNode(&inode)
	fmt.Println(cypherQ.ToString())
	return
}

func queryForNode(inode *collectors.INode) (query *cypher.Query) {
	varr := "d"
	query = cypher.NewQuery()
	query.Merge(
		varr, inode.Type, "path", inode.Path,
	).Set(
		varr, "name", inode.Name,
	).Set(
		varr, "parent", inode.Parent,
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
