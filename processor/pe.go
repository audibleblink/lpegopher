package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
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
		"", node.Principal, "name", inode.DACL.Owner,
	).Merge(
		"", node.Principal, "name", inode.DACL.Group,
	)

	for _, ace := range inode.DACL.Aces {
		query.Merge("", node.Principal, "name", ace.Principal)
	}

	return
}

func RelatePEs(path string) (err error) {
	log := logerr.Add("pe relation")
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var inode collectors.INode

	count := 0
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 8*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		count += 1
		text := scanner.Bytes()
		err = json.Unmarshal(text, &inode)
		if err != nil {
			log.Infof("malformed json at line: %d", count)
			continue
		}

		pe, dir := "pe", "dir"

		cypherQ, err := cypher.NewQuery()
		if err != nil {
			return log.Wrap(err)
		}

		cypherQ.Merge(
			pe, inode.Type, "path", inode.Path,
		).Merge(
			dir, node.Dir, "path", inode.Parent,
		).Relate(
			dir, "CONTAINS", pe,
		)

		id := 0
		for _, ace := range inode.DACL.Aces {
			prnpl := fmt.Sprintf("p%d", id)
			cypherQ.Merge(
				prnpl, node.Principal, "name", ace.Principal,
			)
			
			for _, priv := range ace.Rights
			id++
		}

		err = cypherQ.ExecuteW()
		if err != nil {
			log.Infof("error processing line: %d %w", count, err)
			log.Debugf("failed query was: %s", cypherQ.ToString())
			continue
		}
	}
	return
}
