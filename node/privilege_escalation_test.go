package node

import (
	"strings"
	"testing"
)

func TestMembershipScenarios(t *testing.T) {
	// Test realistic privilege escalation scenarios that would use MEMBER_OF relationships
	
	t.Run("User member of group privilege escalation", func(t *testing.T) {
		// Scenario: A regular user is a member of a group that has system access
		// This should be detected by the privilege escalation queries
		
		// The MEMBER_OF relationship should allow Neo4j to traverse:
		// (user:Principal)-[:MEMBER_OF]->(group:Principal)-[:OWNS|ACE]->(system_resource)
		
		template, _ := GetRelationshipTemplate(MemberOf)
		
		// Verify the template will match this scenario
		expectedPatterns := []string{
			"MATCH (group:Principal),(user:Principal)",
			"WHERE user.group = group.name", 
			"MERGE (user)-[:MEMBER_OF]->(group)",
		}
		
		for _, pattern := range expectedPatterns {
			if !strings.Contains(template, pattern) {
				t.Errorf("MEMBER_OF template missing pattern for user-group escalation: %s", pattern)
			}
		}
	})
	
	t.Run("Group hierarchy support", func(t *testing.T) {
		// Test that the membership template can handle scenarios where:
		// - User 'alice' has group property 'Developers'  
		// - Group 'Developers' exists as a Principal with name 'Developers'
		// - This creates: (alice)-[:MEMBER_OF]->(Developers)
		
		template, _ := GetRelationshipTemplate(MemberOf)
		
		// The template should create relationships based on property matching
		if !strings.Contains(template, "user.group = group.name") {
			t.Error("Template should match user.group property with group.name property")
		}
	})
	
	t.Run("Privilege escalation query compatibility", func(t *testing.T) {
		// Test that MEMBER_OF relationships work with the privilege escalation queries from queries.md
		// The "GetSystem" query uses: shortestPath((low:Principal)-[*..5]->(hi:Principal))
		
		// MEMBER_OF relationships should be included in the path traversal
		// since they connect Principal nodes to other Principal nodes
		
		relType := MemberOf
		if relType != "MEMBER_OF" {
			t.Errorf("Expected relationship type to be MEMBER_OF for privilege escalation queries, got %s", relType)
		}
	})
}

func TestPrivilegeEscalationPathSupport(t *testing.T) {
	// Test that MEMBER_OF relationships enable the privilege escalation scenarios described in the issue
	
	t.Run("User to system via group membership", func(t *testing.T) {
		// Issue scenario: "A group might have a path to system, and the user is a member of that group"
		// This requires: (user)-[:MEMBER_OF]->(group)-[privilege_path*]->(system)
		
		template, _ := GetRelationshipTemplate(MemberOf)
		
		// Verify the relationship direction allows privilege escalation traversal
		if !strings.Contains(template, "(user)-[:MEMBER_OF]->(group)") {
			t.Error("MEMBER_OF relationship should go from user to group to enable privilege escalation")
		}
	})
	
	t.Run("Integration with existing privilege relationships", func(t *testing.T) {
		// MEMBER_OF should work alongside other privilege relationships like OWNS
		
		// Verify MEMBER_OF constant is available for use in privilege queries
		memberOfRel := MemberOf
		ownsRel := Owns
		
		if memberOfRel == ownsRel {
			t.Error("MEMBER_OF and OWNS should be different relationship types")
		}
		
		// Both should be valid relationship types
		_, err1 := GetRelationshipTemplate(memberOfRel)
		_, err2 := GetRelationshipTemplate(ownsRel)
		
		if err1 != nil || err2 != nil {
			t.Errorf("Both MEMBER_OF and OWNS should have valid templates, got errors: %v, %v", err1, err2)
		}
	})
}