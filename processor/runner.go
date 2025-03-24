package processor

import (
	"fmt"

	"github.com/audibleblink/lpegopher/cypher"
	"github.com/audibleblink/lpegopher/logerr"
	"github.com/audibleblink/lpegopher/node"
)

func execString(query string) error {
	cypherQ, err := cypher.NewQuery()
	if err != nil {
		return err
	}
	cypherQ.Raw(query)
	return cypherQ.ExecuteW()
}

// InsertAllRunners loads runner data into the graph database
func InsertAllRunners(stageURL string) (err error) {
	log := logerr.Add("runner inserts")

	template, _ := node.GetTemplateForNodeType(node.Runner)
	err = execString(fmt.Sprintf(template, dataPrefix(stageURL)))
	if err != nil {
		return log.Wrap(err)
	}
	return nil
}

// BulkRelateRunners creates relationships between runners and other nodes
func BulkRelateRunners() (err error) {
	log := logerr.Add("runner relationships")

	// relate dirs that hosts a runner exe
	log.Debugf("relating all (:Dir)-[:%s]->(:Runner)", node.HostsPesFor)
	template, _ := node.GetRelationshipTemplate(node.HostsPesFor)
	err = execString(template)
	if err != nil {
		return log.Wrap(err)
	}

	// relate principals that run certain runners
	log.Debugf("relating all (:Runner)-[:%s]->(:Principal)", node.RunsAs)
	template, _ = node.GetRelationshipTemplate(node.RunsAs)
	err = execString(template)
	if err != nil {
		return log.Wrap(err)
	}

	// relate exes that are executed by a runner
	log.Debugf("relating all (:Exe)-[:%s]->(:Runner)", node.ExecutedBy)
	template, _ = node.GetRelationshipTemplate(node.ExecutedBy)
	err = execString(template)
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

func dataPrefix(url string) (uri string) {
	if url == "" {
		return "file://"
	}
	return fmt.Sprintf("http://%s", url)
}
