package node

import (
	"testing"
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
	}

	for expected, actual := range nodeTypes {
		if expected != actual {
			t.Errorf("Expected node type constant %s to be %s, got %s", expected, expected, actual)
		}
	}
}

func TestPropStructValues(t *testing.T) {
	// Test property name constants
	propTests := map[string]string{
		"name":    Prop.Name,
		"dir":     Prop.Dir,
		"parent":  Prop.Parent,
		"path":    Prop.Path,
		"type":    Prop.Type,
		"args":    Prop.Args,
		"exe":     Prop.Exe,
		"context": Prop.Context,
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