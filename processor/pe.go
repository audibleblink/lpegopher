package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/audibleblink/pegopher/collectors"
	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
)

const (
	MaxBatchSize = 5000
)

var EntityCache = map[string]map[string]bool{
	node.Principal: make(map[string]bool),
	node.Exe:       make(map[string]bool),
	node.Dir:       make(map[string]bool),
	node.Dll:       make(map[string]bool),
}

var (
	BatchCount   = 1
	CurrentBatch = &strings.Builder{}
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

	if BatchCount%MaxBatchSize == 0 {
		cypherQ.Raw(CurrentBatch.String())
		err = cypherQ.ExecuteW()
		CurrentBatch.Reset()
		BatchCount = 0
	} else {
		cypherQ.Terminate()
		_, err = CurrentBatch.WriteString(cypherQ.String())
		BatchCount++
	}

	return
}

func exists(node, uniqPropValue string) bool {
	return EntityCache[node][uniqPropValue]
}

func cacheAdd(node, uniqPropValue string) {
	EntityCache[node][uniqPropValue] = true
}

// MatchOrCreate works like q.Merge, except a local cache is consulted
// instead of asking Neo4j. If cache indicates the node has already been
// created, and relate argument is false, no addition query is appened
// func MatchOrCreate(varr, label, uniqProp, value string, relate bool) *cypher.Query {
// 	fmt.Fprintf(q.b, template, varr, label, uniqProp, value)
// 	fmt.Fprintf(q.b, "\n")
// 	return q
// }

func queryForINode(inode *collectors.INode) (query *cypher.Query, err error) {
	nodeAlias := "d"
	query, err = cypher.NewQuery()
	if err != nil {
		return nil, err
	}

	props := map[string]string{
		"name":   inode.Name,
		"parent": inode.Parent,
	}

	query.Create(
		nodeAlias, inode.Type, "path", inode.Path,
	).Set(
		nodeAlias, props,
	)
	cacheAdd(inode.Type, inode.Path)

	if !exists(node.Principal, inode.DACL.Owner) {
		query.Create("", node.Principal, "name", inode.DACL.Owner)
		cacheAdd(node.Principal, inode.DACL.Owner)
	}

	if !exists(node.Principal, inode.DACL.Group) {
		query.Create("", node.Principal, "name", inode.DACL.Group)
		cacheAdd(node.Principal, inode.DACL.Group)
	}

	for _, ace := range inode.DACL.Aces {
		if !exists(node.Principal, ace.Principal) {
			query.Create("", node.Principal, "name", ace.Principal)
			cacheAdd(node.Principal, ace.Principal)
		}
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
	CurrentBatch.Reset()

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

		if count%5000 == 0 {
			err = cypherQ.ExecuteW()
			if err != nil {
				log.Infof("error processing line: %d %w", count, err)
				log.Debugf("failed query was: %s", cypherQ.String())
				continue
			}
			count++
			continue
		}
		CurrentBatch.WriteString(cypherQ.String())
	}
	return
}
