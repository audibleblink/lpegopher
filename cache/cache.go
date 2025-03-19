package cache

import "github.com/audibleblink/lpegopher/node"

var store = map[string]map[string]bool{
	node.Principal: make(map[string]bool),
	node.Exe:       make(map[string]bool),
	node.Dir:       make(map[string]bool),
	node.Dll:       make(map[string]bool),
	node.Runner:    make(map[string]bool),
}

// Add adds a node to the cache and returns whether it was newly added
func Add(nodeType, uniqPropValue string) bool {
	if store[nodeType][uniqPropValue] {
		return false
	}

	store[nodeType][uniqPropValue] = true
	return true
}
