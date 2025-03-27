package processor

import (
	"fmt"

	"github.com/audibleblink/logerr"
	"github.com/audibleblink/lpegopher/node"
)

// InsertAllNodes loads node data into the graph database
func InsertAllNodes(stageURL string) (err error) {
	log := logerr.Add("file inserts")
	urlPrefix := dataPrefix(stageURL)

	// Process exes
	log.Debug("processing exes")
	template, _ := node.GetTemplateForNodeType(node.Exe)
	err = execString(fmt.Sprintf(template, urlPrefix))
	if err != nil {
		return log.Wrap(err)
	}

	// Process dlls
	log.Debug("processing dlls")
	template, _ = node.GetTemplateForNodeType(node.Dll)
	err = execString(fmt.Sprintf(template, urlPrefix))
	if err != nil {
		return log.Wrap(err)
	}

	// Process dirs
	log.Debug("processing dirs")
	template, _ = node.GetTemplateForNodeType(node.Dir)
	err = execString(fmt.Sprintf(template, urlPrefix))
	if err != nil {
		return log.Wrap(err)
	}

	// Process deps
	log.Debug("processing forwards")
	template, _ = node.GetTemplateForNodeType(node.Dep)
	err = execString(fmt.Sprintf(template, urlPrefix))
	if err != nil {
		return log.Wrap(err)
	}

	// Process principals
	log.Debug("processing principals")
	template, _ = node.GetTemplateForNodeType(node.Principal)
	err = execString(fmt.Sprintf(template, urlPrefix))
	if err != nil {
		return log.Wrap(err)
	}

	return nil
}

// BulkRelateFileTree creates relationships between files and directories
func BulkRelateFileTree() (err error) {
	log := logerr.Add("filetree relationships")
	template, _ := node.GetRelationshipTemplate(node.Contains)

	for _, typ := range []string{node.Dir, node.Exe, node.Dll} {
		log.Debugf("relating all (:Dir)-[:%s]-(:%s)", node.Contains, typ)
		err = execString(fmt.Sprintf(template, typ))
		if err != nil {
			return log.Wrap(err)
		}
	}
	return nil
}

// RelateOwnership creates ownership relationships between principals and nodes
func RelateOwnership() (err error) {
	log := logerr.Add("ownership creation")
	log.Debugf("relating all (:Principal)-[:%s]-(:INode)", node.Owns)

	template, _ := node.GetRelationshipTemplate(node.Owns)
	err = execString(template)
	if err != nil {
		return log.Wrap(err)
	}
	return nil
}

// RelateMembership creates group membership relationships between principals
func RelateMembership() (err error) {
	log := logerr.Add("membership creation")
	log.Debugf("relating all (:Principal)-[:%s]-(:Principal)", node.MemberOf)

	template, _ := node.GetRelationshipTemplate(node.MemberOf)
	err = execString(template)
	if err != nil {
		return log.Wrap(err)
	}
	return nil
}

// RelateACLs creates access control relationships between nodes
func RelateACLs(stageURL string) (err error) {
	log := logerr.Add("acl relationships")
	log.Debug("relating all (:Principal)-[$ACE]-(:INodes)")

	// ACL relationships are custom and don't use a specific relationship type
	query := `CALL apoc.periodic.iterate("
			LOAD CSV FROM '%s/relationships.csv' AS line RETURN line
		","
			MATCH (a:Principal {nid: line[0]}), (b:INode {nid: line[2]})
			CALL apoc.create.relationship(a, line[1], {}, b) YIELD rel RETURN rel
		", {batchSize: 20000});
		`
	err = execString(fmt.Sprintf(query, dataPrefix(stageURL)))
	if err != nil {
		return log.Wrap(err)
	}
	return nil
}

// RelateDependecies creates dependency relationships between nodes
func RelateDependecies(stageURL string) (err error) {
	log := logerr.Add("dependecy relationships")
	log.Debugf("relating (:INode)-[:%s]-(:Dep)", node.ImportedBy)

	template, _ := node.GetRelationshipTemplate(node.ImportedBy)
	err = execString(fmt.Sprintf(template, dataPrefix(stageURL)))
	if err != nil {
		return log.Wrap(err)
	}
	return nil
}
