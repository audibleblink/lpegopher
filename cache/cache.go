package cache

import "github.com/audibleblink/lpegopher/node"

// TODO: Generic - This could be defined as map[string]map[Key]Value[bool] for type safety
var store = map[string]map[string]bool{
	node.Principal: make(map[string]bool),
	node.Exe:       make(map[string]bool),
	node.Dir:       make(map[string]bool),
	node.Dll:       make(map[string]bool),
	node.Runner:    make(map[string]bool),
}

// Add adds a node to the cache and returns whether it was newly added
func Add(nodeType, uniqPropValue string) bool {
	if _, exists := store[nodeType]; !exists {
		return false
	}

	if store[nodeType][uniqPropValue] {
		return false
	}

	store[nodeType][uniqPropValue] = true
	return true
}
