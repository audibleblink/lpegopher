package main

type DLL struct {
	containableIFile
	securableIFile
	// Imports      []*PE   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*PE   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
	// ExecutedFrom *Runner `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
}
