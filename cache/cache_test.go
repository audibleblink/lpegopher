package cache

import (
	"testing"

	"github.com/audibleblink/lpegopher/node"
)

func TestAdd(t *testing.T) {
	// Reset the cache store for testing
	resetStore()

	t.Run("returns true for new item", func(t *testing.T) {
		// Add a new item to the cache
		result := Add(node.Exe, "c:/windows/system32/cmd.exe")

		// Should return true for a new item
		if !result {
			t.Error("Add should return true for a new item")
		}

		// Verify item is in cache
		if !store[node.Exe]["c:/windows/system32/cmd.exe"] {
			t.Error("Item should be added to cache")
		}
	})

	t.Run("returns false for existing item", func(t *testing.T) {
		// Add the same item again
		result := Add(node.Exe, "c:/windows/system32/cmd.exe")

		// Should return false for an existing item
		if result {
			t.Error("Add should return false for an existing item")
		}
	})

	t.Run("handles multiple node types", func(t *testing.T) {
		// Add items of different types
		Add(node.Dll, "c:/windows/system32/kernel32.dll")
		Add(node.Dir, "c:/windows/system32")
		Add(node.Principal, "SYSTEM")
		Add(node.Runner, "Task Scheduler")

		// Verify all items are in cache
		if !store[node.Dll]["c:/windows/system32/kernel32.dll"] {
			t.Error("DLL should be added to cache")
		}
		if !store[node.Dir]["c:/windows/system32"] {
			t.Error("Directory should be added to cache")
		}
		if !store[node.Principal]["SYSTEM"] {
			t.Error("Principal should be added to cache")
		}
		if !store[node.Runner]["Task Scheduler"] {
			t.Error("Runner should be added to cache")
		}
	})
}

// Helper function to reset the store between tests
func resetStore() {
	store = map[string]map[string]bool{
		node.Principal: make(map[string]bool),
		node.Exe:       make(map[string]bool),
		node.Dir:       make(map[string]bool),
		node.Dll:       make(map[string]bool),
		node.Runner:    make(map[string]bool),
	}
}
