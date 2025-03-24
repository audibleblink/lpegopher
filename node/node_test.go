package node

import (
	"fmt"
	"strings"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestNodeConstants(t *testing.T) {
	// Test node type constants
	nodeTypes := map[string]string{
		"Dll":       Dll,
		"Exe":       Exe,
		"Directory": Dir,
		"Runner":    Runner,
		"Principal": Principal,
		"Dep":       Dep,
		"INode":     INode,
	}

	for expected, actual := range nodeTypes {
		if expected != actual {
			t.Errorf("Expected node type constant %s to be %s, got %s", expected, expected, actual)
		}
	}
}

func TestRelationshipConstants(t *testing.T) {
	// Test relationship type constants
	relTypes := map[string]string{
		"CONTAINS":      Contains,
		"OWNS":          Owns,
		"MEMBER_OF":     MemberOf,
		"HOSTS_PES_FOR": HostsPesFor,
		"RUNS_AS":       RunsAs,
		"EXECUTED_BY":   ExecutedBy,
		"IMPORTED_BY":   ImportedBy,
	}

	for expected, actual := range relTypes {
		if expected != actual {
			t.Errorf(
				"Expected relationship type constant %s to be %s, got %s",
				expected,
				expected,
				actual,
			)
		}
	}
}

func TestPropStructValues(t *testing.T) {
	// Test property name constants
	propTests := map[string]string{
		"name":     Prop.Name,
		"dir":      Prop.Dir,
		"parent":   Prop.Parent,
		"path":     Prop.Path,
		"type":     Prop.Type,
		"args":     Prop.Args,
		"exe":      Prop.Exe,
		"context":  Prop.Context,
		"nid":      Prop.Nid,
		"owner":    Prop.Owner,
		"group":    Prop.Group,
		"runlevel": Prop.RunLevel,
	}

	for expected, actual := range propTests {
		if expected != actual {
			t.Errorf("Expected Prop.%s to be %s, got %s", expected, expected, actual)
		}
	}
}

func TestAbusableAcesConstants(t *testing.T) {
	// Test abusable ACE privilege constants
	aceTests := map[string]string{
		"WRITE_OWNER":    WriteOwner,
		"WRITE_DACL":     WriteDACL,
		"GENERIC_ALL":    GenericAll,
		"GENERIC_WRITE":  GenericWrite,
		"CONTROL_ACCESS": ControlAccess,
	}

	for expected, actual := range aceTests {
		if expected != actual {
			t.Errorf("Expected ACE constant %s to be %s, got %s", expected, expected, actual)
		}
	}
}

func TestAbusableAcesMap(t *testing.T) {
	// Test that all abusable ACEs are marked as true in the map
	for _, ace := range []string{
		WriteOwner,
		WriteDACL,
		GenericAll,
		GenericWrite,
		ControlAccess,
	} {
		if !AbusableAces[ace] {
			t.Errorf("Expected %s to be marked as abusable in AbusableAces map", ace)
		}
	}

	// Test that a non-existent ACE returns false
	if AbusableAces["NON_EXISTENT_ACE"] {
		t.Errorf("Expected non-existent ACE to return false from AbusableAces map")
	}
}

func TestSchemaDefinition(t *testing.T) {
	// Test unique constraints
	expectedUniqueConstraints := map[string]string{
		INode:     Prop.Nid,
		Principal: Prop.Nid,
		Runner:    Prop.Nid,
		Dep:       Prop.Nid,
		Exe:       Prop.Path,
		Dll:       Prop.Path,
		Dir:       Prop.Path,
	}

	for nodeType, prop := range expectedUniqueConstraints {
		actualProp, exists := Schema.UniqueConstraints[nodeType]
		if !exists {
			t.Errorf("Expected unique constraint for node type %s, but none found", nodeType)
			continue
		}
		if actualProp != prop {
			t.Errorf(
				"Expected unique constraint for %s to be on property %s, got %s",
				nodeType,
				prop,
				actualProp,
			)
		}
	}

	// Test BTREE indices
	expectedBTreeIndices := map[string][]string{
		INode:     {Prop.Owner, Prop.Group, Prop.Name},
		Exe:       {Prop.Parent},
		Dll:       {Prop.Parent},
		Dir:       {Prop.Parent},
		Runner:    {Prop.Parent, Prop.Exe, Prop.Context},
		Principal: {Prop.Name},
	}

	for nodeType, expectedProps := range expectedBTreeIndices {
		actualProps, exists := Schema.BTREEIndices[nodeType]
		if !exists {
			t.Errorf("Expected BTREE indices for node type %s, but none found", nodeType)
			continue
		}
		if len(actualProps) != len(expectedProps) {
			t.Errorf(
				"Expected %d BTREE indices for %s, got %d",
				len(expectedProps),
				nodeType,
				len(actualProps),
			)
			continue
		}
		for i, prop := range expectedProps {
			if actualProps[i] != prop {
				t.Errorf(
					"Expected BTREE index %d for %s to be on property %s, got %s",
					i,
					nodeType,
					prop,
					actualProps[i],
				)
			}
		}
	}
}

func TestPropMaps(t *testing.T) {
	// Test property maps for each node type
	expectedProps := map[string][]string{
		"INode":     {Prop.Nid, Prop.Name, Prop.Path, Prop.Parent, Prop.Owner, Prop.Group},
		"Principal": {Prop.Nid, Prop.Name, Prop.Group},
		"Runner": {
			Prop.Nid,
			Prop.Name,
			Prop.Type,
			Prop.Path,
			Prop.Exe,
			Prop.Parent,
			Prop.Context,
			Prop.RunLevel,
		},
		"Dep": {Prop.Nid, Prop.Name},
	}

	// Test INode properties
	testPropertyList(t, "INode", PropMaps.INode, expectedProps["INode"])

	// Test Principal properties
	testPropertyList(t, "Principal", PropMaps.Principal, expectedProps["Principal"])

	// Test Runner properties
	testPropertyList(t, "Runner", PropMaps.Runner, expectedProps["Runner"])

	// Test Dep properties
	testPropertyList(t, "Dep", PropMaps.Dep, expectedProps["Dep"])
}

func testPropertyList(t *testing.T, nodeType string, actual, expected []string) {
	if len(actual) != len(expected) {
		t.Errorf("Expected %s to have %d properties, got %d", nodeType, len(expected), len(actual))
		return
	}
	for i, prop := range expected {
		if actual[i] != prop {
			t.Errorf("Expected %s property %d to be %s, got %s", nodeType, i, prop, actual[i])
		}
	}
}

func TestCypherTemplates(t *testing.T) {
	// Test node creation templates
	templates := []struct {
		name     string
		template string
		expected []string
	}{
		{
			"CreateExe",
			CypherTemplates.CreateExe,
			[]string{
				"exes.csv",
				"CREATE",
				"Exe",
				"INode",
				"nid",
				"name",
				"path",
				"parent",
				"owner",
				"group",
			},
		},
		{
			"CreateDll",
			CypherTemplates.CreateDll,
			[]string{
				"dlls.csv",
				"CREATE",
				"Dll",
				"INode",
				"nid",
				"name",
				"path",
				"parent",
				"owner",
				"group",
			},
		},
		{
			"CreateDir",
			CypherTemplates.CreateDir,
			[]string{
				"dirs.csv",
				"CREATE",
				"Directory",
				"INode",
				"nid",
				"name",
				"path",
				"parent",
				"owner",
				"group",
			},
		},
		{
			"CreateDep",
			CypherTemplates.CreateDep,
			[]string{"deps.csv", "CREATE", "Dep", "nid", "name"},
		},
		{
			"CreatePrincipal",
			CypherTemplates.CreatePrincipal,
			[]string{"principals.csv", "CREATE", "Principal", "nid", "name", "group"},
		},
		{
			"CreateRunner",
			CypherTemplates.CreateRunner,
			[]string{
				"runners.csv",
				"CREATE",
				"Runner",
				"nid",
				"name",
				"type",
				"path",
				"exe",
				"parent",
				"context",
				"runlevel",
			},
		},
		{
			"RelateFileTree",
			CypherTemplates.RelateFileTree,
			[]string{"MATCH", "node", "dir", "Directory", "parent", "path", "MERGE", "CONTAINS"},
		},
		{
			"RelateOwnership",
			CypherTemplates.RelateOwnership,
			[]string{"MATCH", "Principal", "INode", "nid", "owner", "group", "MERGE", "OWNS"},
		},
		{
			"RelateMembership",
			CypherTemplates.RelateMembership,
			[]string{"MATCH", "Principal", "group", "name", "MERGE", "MEMBER_OF"},
		},
		{
			"RelateRunnerDir",
			CypherTemplates.RelateRunnerDir,
			[]string{"MATCH", "Runner", "Directory", "parent", "path", "MERGE", "HOSTS_PES_FOR"},
		},
		{
			"RelateRunnerPrincipal",
			CypherTemplates.RelateRunnerPrincipal,
			[]string{"MATCH", "Runner", "Principal", "context", "name", "MERGE", "RUNS_AS"},
		},
		{
			"RelateRunnerExe",
			CypherTemplates.RelateRunnerExe,
			[]string{"MATCH", "Runner", "Exe", "parent", "exe", "path", "MERGE", "EXECUTED_BY"},
		},
		{
			"RelateDependency",
			CypherTemplates.RelateDependency,
			[]string{"imports.csv", "MATCH", "INode", "Dep", "nid", "MERGE", "IMPORTED_BY"},
		},
	}

	for _, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			for _, expected := range tt.expected {
				if !strings.Contains(tt.template, expected) {
					t.Errorf(
						"Expected %s template to contain %q, but it doesn't",
						tt.name,
						expected,
					)
				}
			}
		})
	}
}

func TestFormatNodeQuery(t *testing.T) {
	template := "MATCH (n:%s) WHERE n.%s = '%s' RETURN n"
	result := FormatNodeQuery(template, "Person", "name", "John Doe")
	expected := "MATCH (n:Person) WHERE n.name = 'John Doe' RETURN n"

	if result != expected {
		t.Errorf("Expected FormatNodeQuery to return %q, got %q", expected, result)
	}
}

func TestGetTemplateForNodeType(t *testing.T) {
	tests := []struct {
		nodeType string
		wantErr  bool
	}{
		{Exe, false},
		{Dll, false},
		{Dir, false},
		{Dep, false},
		{Principal, false},
		{Runner, false},
		{"UnknownType", true},
	}

	for _, tt := range tests {
		t.Run(tt.nodeType, func(t *testing.T) {
			template, err := GetTemplateForNodeType(tt.nodeType)
			if tt.wantErr {
				if err == nil {
					t.Errorf(
						"Expected GetTemplateForNodeType to return error for type %q, but it didn't",
						tt.nodeType,
					)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for node type %q: %v", tt.nodeType, err)
				}
				if template == "" {
					t.Errorf("Expected non-empty template for node type %q", tt.nodeType)
				}
			}
		})
	}
}

func TestGetRelationshipTemplate(t *testing.T) {
	tests := []struct {
		relType string
		wantErr bool
	}{
		{Contains, false},
		{Owns, false},
		{MemberOf, false},
		{HostsPesFor, false},
		{RunsAs, false},
		{ExecutedBy, false},
		{ImportedBy, false},
		{"UnknownRelationship", true},
	}

	for _, tt := range tests {
		t.Run(tt.relType, func(t *testing.T) {
			template, err := GetRelationshipTemplate(tt.relType)
			if tt.wantErr {
				if err == nil {
					t.Errorf(
						"Expected GetRelationshipTemplate to return error for type %q, but it didn't",
						tt.relType,
					)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for relationship type %q: %v", tt.relType, err)
				}
				if template == "" {
					t.Errorf("Expected non-empty template for relationship type %q", tt.relType)
				}
			}
		})
	}
}

// Mock Neo4j Transaction for testing
type mockTransaction struct{}

func (m *mockTransaction) Run(cypher string, params map[string]any) (neo4j.Result, error) {
	if strings.Contains(cypher, "ERROR") {
		return nil, fmt.Errorf("mock error")
	}
	return nil, nil
}

func (m *mockTransaction) Commit() error {
	return nil
}

func (m *mockTransaction) Rollback() error {
	return nil
}

func (m *mockTransaction) Close() error {
	return nil
}

func TestNodeSchema(t *testing.T) {
	mockTx := &mockTransaction{}
	schema := NewNodeSchema(mockTx)

	// Test CreateUniqueConstraints
	err := schema.CreateUniqueConstraints()
	if err != nil {
		t.Errorf("Unexpected error from CreateUniqueConstraints: %v", err)
	}

	// Test CreateBTreeIndices
	err = schema.CreateBTreeIndices()
	if err != nil {
		t.Errorf("Unexpected error from CreateBTreeIndices: %v", err)
	}
}

