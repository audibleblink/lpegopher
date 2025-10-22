package node

import (
	"strings"
	"testing"
)

func TestMembershipRelationshipTemplate(t *testing.T) {
	// Test that the RelateMembership template creates the correct Cypher query
	template, err := GetRelationshipTemplate(MemberOf)
	if err != nil {
		t.Fatalf("Failed to get MEMBER_OF relationship template: %v", err)
	}

	// Verify the template contains the expected components
	expectedComponents := []string{
		"MATCH",
		"(group:Principal)",
		"(user:Principal)",
		"user.group = group.name",
		"MERGE",
		"(user)-[:MEMBER_OF]->(group)",
		"apoc.periodic.iterate",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(template, component) {
			t.Errorf("Expected MEMBER_OF template to contain '%s', but it doesn't. Template: %s", component, template)
		}
	}
}

func TestMembershipRelationshipDirection(t *testing.T) {
	// Test that the relationship direction is correct: (user)-[:MEMBER_OF]->(group)
	template, _ := GetRelationshipTemplate(MemberOf)
	
	// The template should have user as source and group as target
	if !strings.Contains(template, "(user)-[:MEMBER_OF]->(group)") {
		t.Errorf("Expected MEMBER_OF relationship to go from user to group, but template doesn't show this pattern: %s", template)
	}
}

func TestMembershipLogic(t *testing.T) {
	// Test that the matching logic is correct: user.group should equal group.name
	template, _ := GetRelationshipTemplate(MemberOf)
	
	if !strings.Contains(template, "user.group = group.name") {
		t.Errorf("Expected MEMBER_OF template to match user.group with group.name, but pattern not found: %s", template)
	}
}

func TestPrincipalGroupProperty(t *testing.T) {
	// Verify that Principal nodes have a group property in their property map
	principalProps := PropMaps.Principal
	
	hasGroupProp := false
	for _, prop := range principalProps {
		if prop == Prop.Group {
			hasGroupProp = true
			break
		}
	}
	
	if !hasGroupProp {
		t.Errorf("Expected Principal property map to include group property, but it doesn't. Props: %v", principalProps)
	}
}

func TestPrincipalNodeCreationIncludesGroup(t *testing.T) {
	// Verify that the Principal creation template includes the group field
	template, err := GetTemplateForNodeType(Principal)
	if err != nil {
		t.Fatalf("Failed to get Principal node template: %v", err)
	}
	
	// Should include line[2] for group property based on PropMaps.Principal order
	expectedComponents := []string{
		"CREATE",
		":Principal",
		"nid: line[0]",
		"name: line[1]", 
		"group: line[2]",
		"principals.csv",
	}
	
	for _, component := range expectedComponents {
		if !strings.Contains(template, component) {
			t.Errorf("Expected Principal creation template to contain '%s', but it doesn't. Template: %s", component, template)
		}
	}
}