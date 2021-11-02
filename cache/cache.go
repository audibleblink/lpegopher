package cache

import "github.com/audibleblink/pegopher/node"

var store = map[string]map[string]bool{
	node.Principal: make(map[string]bool),
	node.Exe:       make(map[string]bool),
	node.Dir:       make(map[string]bool),
	node.Dll:       make(map[string]bool),
}

func Add(nodeType, uniqPropValue string) bool {
	if store[nodeType][uniqPropValue] {
		return false
	}

	store[nodeType][uniqPropValue] = true
	return true
}
