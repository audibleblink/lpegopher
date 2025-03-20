package node

import (
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Abusable ACE privilege constants
const (
	WriteOwner    = "WRITE_OWNER"    // Permission to change ownership
	WriteDACL     = "WRITE_DACL"     // Permission to modify access control list
	GenericAll    = "GENERIC_ALL"    // Full control permission
	GenericWrite  = "GENERIC_WRITE"  // Write permission
	ControlAccess = "CONTROL_ACCESS" // Control access permission
)

// AbusableAces maps privilege names to a boolean indicating they are abusable
var AbusableAces = map[string]bool{
	WriteOwner:    true,
	WriteDACL:     true,
	GenericAll:    true,
	GenericWrite:  true,
	ControlAccess: true,
}

// Node type constants
const (
	Dll       = "Dll"       // Dynamic Link Library
	Exe       = "Exe"       // Executable
	Dir       = "Directory" // Directory
	Runner    = "Runner"    // Auto-executing program
	Principal = "Principal" // Security principal
	Dep       = "Dep"       // Dependency
	INode     = "INode"     // Base node type for files and directories
)

// Relationship type constants
const (
	Contains      = "CONTAINS"       // Directory contains a file or subdirectory
	Owns          = "OWNS"           // Principal owns a node
	MemberOf      = "MEMBER_OF"      // User is a member of a group
	HostsPesFor   = "HOSTS_PES_FOR"  // Directory hosts executables for a runner
	RunsAs        = "RUNS_AS"        // Runner runs as a principal
	ExecutedBy    = "EXECUTED_BY"    // Executable is executed by a runner
	ImportedBy    = "IMPORTED_BY"    // Dependency is imported by a node
)

// Basic property name constants for nodes
var Prop = struct {
	Name     string
	Dir      string
	Parent   string
	Path     string
	Type     string
	Args     string
	Exe      string
	Context  string
	Nid      string
	Owner    string
	Group    string
	RunLevel string
}{
	"name",
	"dir",
	"parent",
	"path",
	"type",
	"args",
	"exe",
	"context",
	"nid",
	"owner",
	"group",
	"runlevel",
}

// Node schema index and constraint definitions
var Schema = struct {
	// Unique constraints
	UniqueConstraints map[string]string
	// BTREE indices
	BTREEIndices map[string][]string
	// Constraint query template
	UniqueConstraintTemplate string
	// BTREE index query template
	BTREEIndexTemplate string
}{
	UniqueConstraints: map[string]string{
		INode:     Prop.Nid,
		Principal: Prop.Nid,
		Runner:    Prop.Nid,
		Dep:       Prop.Nid,
		Exe:       Prop.Path,
		Dll:       Prop.Path,
		Dir:       Prop.Path,
	},
	BTREEIndices: map[string][]string{
		INode: {
			Prop.Owner,
			Prop.Group,
			Prop.Name,
		},
		Exe: {
			Prop.Parent,
		},
		Dll: {
			Prop.Parent,
		},
		Dir: {
			Prop.Parent,
		},
		Runner: {
			Prop.Parent,
			Prop.Exe,
			Prop.Context,
		},
		Principal: {
			Prop.Name,
		},
	},
	UniqueConstraintTemplate: "CREATE CONSTRAINT ON (a:%s) ASSERT a.%s IS UNIQUE;",
	BTREEIndexTemplate:       "CREATE BTREE INDEX FOR (n:%s) ON (n.%s)",
}

// Node property maps for each node type
var PropMaps = struct {
	INode     []string
	Principal []string
	Runner    []string
	Dep       []string
}{
	INode: []string{
		Prop.Nid,
		Prop.Name,
		Prop.Path,
		Prop.Parent,
		Prop.Owner,
		Prop.Group,
	},
	Principal: []string{
		Prop.Nid,
		Prop.Name,
		Prop.Group,
	},
	Runner: []string{
		Prop.Nid,
		Prop.Name,
		Prop.Type,
		Prop.Path,
		Prop.Exe,
		Prop.Parent,
		Prop.Context,
		Prop.RunLevel,
	},
	Dep: []string{
		Prop.Nid,
		Prop.Name,
	},
}

// Cypher query templates for node operations
var CypherTemplates = struct {
	// Node creation templates
	CreateExe      string
	CreateDll      string
	CreateDir      string
	CreateDep      string
	CreatePrincipal string
	CreateRunner   string
	// Relationship creation templates
	RelateFileTree string
	RelateOwnership string
	RelateMembership string
	RelateRunnerDir string
	RelateRunnerPrincipal string
	RelateRunnerExe string
	RelateDependency string
}{
	CreateExe: `LOAD CSV FROM '%s/exes.csv' AS line
		WITH line
		CREATE (:Exe:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`,
			
	CreateDll: `LOAD CSV FROM '%s/dlls.csv' AS line
		WITH line
		CREATE (:Dll:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`,
			
	CreateDir: `LOAD CSV FROM '%s/dirs.csv' AS line
		WITH line
		CREATE (:Directory:INode {
			nid: line[0], 
			name: line[1],
			path: line[2],
			parent: line[3],
			owner: line[4],
			group: line[5] })`,
			
	CreateDep: `LOAD CSV FROM '%s/deps.csv' AS line
		WITH line CREATE (:Dep {nid: line[0], name: line[1]})`,
		
	CreatePrincipal: `LOAD CSV FROM '%s/principals.csv' AS line
		WITH line CREATE (:Principal {nid: line[0], name: line[1], group: line[2]})`,
		
	CreateRunner: `LOAD CSV FROM '%s/runners.csv' AS line
		WITH line
		CREATE (e:Runner {
			nid: line[0], 
			name: line[1],
			type: line[2],
			path: line[3],
			exe: line[4],
			parent: line[5],
			context: line[6],
			runlevel: line[7]})`,
			
	RelateFileTree: `
		CALL apoc.periodic.iterate(
			"MATCH (node:%s),(dir:Directory) WHERE node.parent = dir.path RETURN node,dir",
			"MERGE (dir)-[:CONTAINS]->(node)",
			{batchSize:1000})
		`,
		
	RelateOwnership: `
		CALL apoc.periodic.iterate("
			MATCH (pcpl:Principal),(inode:INode) 
			WHERE pcpl.nid = inode.owner or pcpl.nid = inode.group 
			RETURN pcpl, inode
		","
			MERGE (pcpl)-[:OWNS]->(inode)
		", {batchSize: 1000})
		`,
		
	RelateMembership: `
		CALL apoc.periodic.iterate("
			MATCH (group:Principal),(user:Principal) 
			WHERE user.group = group.name 
			RETURN user, group
		","
			MERGE (user)-[:MEMBER_OF]->(group)
		", {batchSize: 10})
		`,
		
	RelateRunnerDir: `
		CALL apoc.periodic.iterate(
			"MATCH (r:Runner),(dir:Directory) WHERE r.parent = dir.path RETURN r,dir",
			"MERGE (dir)-[:HOSTS_PES_FOR]->(r)",
			{batchSize:100, parallel: true, iterateList:true})
		`,
		
	RelateRunnerPrincipal: `
		CALL apoc.periodic.iterate(
			"MATCH (r:Runner),(p:Principal) WHERE r.context = p.name RETURN r,p",
			"MERGE (r)-[:RUNS_AS]->(p)",
			{batchSize:100, iterateList: true})
		`,
		
	RelateRunnerExe: `
		CALL apoc.periodic.iterate(
			"MATCH (r:Runner),(exe:Exe) WHERE r.parent+'/'+r.exe = exe.path RETURN r,exe",
			"MERGE (exe)-[:EXECUTED_BY]->(r)",
			{batchSize:100})
		`,
		
	RelateDependency: `CALL apoc.periodic.iterate("
			LOAD CSV FROM '%s/imports.csv' AS line RETURN line
		","
			MATCH (a:INode {nid: line[0]}), (b:Dep {nid: line[2]})
			MERGE (b)-[:IMPORTED_BY]->(a)
		", {batchSize: 20000});
		`,
}

// NodeSchema represents a Neo4j graph schema for nodes
type NodeSchema struct {
	tx neo4j.Transaction
}

// NewNodeSchema creates a new NodeSchema with the given transaction
func NewNodeSchema(tx neo4j.Transaction) *NodeSchema {
	return &NodeSchema{tx: tx}
}

// CreateUniqueConstraints creates unique constraints for all node types in the schema
func (ns *NodeSchema) CreateUniqueConstraints() error {
	for nodeType, prop := range Schema.UniqueConstraints {
		query := fmt.Sprintf(Schema.UniqueConstraintTemplate, nodeType, prop)
		if _, err := ns.tx.Run(query, nil); err != nil {
			return fmt.Errorf("failed to create unique constraint for %s.%s: %w", nodeType, prop, err)
		}
	}
	return nil
}

// CreateBTreeIndices creates BTree indices for all node types in the schema
func (ns *NodeSchema) CreateBTreeIndices() error {
	for nodeType, props := range Schema.BTREEIndices {
		for _, prop := range props {
			query := fmt.Sprintf(Schema.BTREEIndexTemplate, nodeType, prop)
			if _, err := ns.tx.Run(query, nil); err != nil {
				return fmt.Errorf("failed to create btree index for %s.%s: %w", nodeType, prop, err)
			}
		}
	}
	return nil
}

// FormatNodeQuery formats a query based on a template and parameters
func FormatNodeQuery(template string, params ...interface{}) string {
	query := fmt.Sprintf(template, params...)
	return strings.TrimSpace(query)
}

// GetTemplateForNodeType returns the appropriate creation template for a given node type
func GetTemplateForNodeType(nodeType string) (string, error) {
	switch nodeType {
	case Exe:
		return CypherTemplates.CreateExe, nil
	case Dll:
		return CypherTemplates.CreateDll, nil
	case Dir:
		return CypherTemplates.CreateDir, nil
	case Dep:
		return CypherTemplates.CreateDep, nil
	case Principal:
		return CypherTemplates.CreatePrincipal, nil
	case Runner:
		return CypherTemplates.CreateRunner, nil
	default:
		return "", fmt.Errorf("no template available for node type: %s", nodeType)
	}
}

// GetRelationshipTemplate returns the template for a relationship type
func GetRelationshipTemplate(relType string) (string, error) {
	switch relType {
	case Contains:
		return CypherTemplates.RelateFileTree, nil
	case Owns:
		return CypherTemplates.RelateOwnership, nil
	case MemberOf:
		return CypherTemplates.RelateMembership, nil
	case HostsPesFor:
		return CypherTemplates.RelateRunnerDir, nil
	case RunsAs:
		return CypherTemplates.RelateRunnerPrincipal, nil
	case ExecutedBy:
		return CypherTemplates.RelateRunnerExe, nil
	case ImportedBy:
		return CypherTemplates.RelateDependency, nil
	default:
		return "", fmt.Errorf("no template available for relationship type: %s", relType)
	}
}
