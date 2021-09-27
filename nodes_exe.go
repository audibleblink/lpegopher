package main

import (
	"context"
	"fmt"
)

type EXE struct {
	securableIFile
	containableIFile

	ExecutesFrom *Runner `gogm:"direction=outgoing;relationship=EXECUTED_FROM"`
	// Parent *Directory `gogm:"direction=incoming;relationship=CONTAINS"`
	// Imports      []*containableIFile   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*containableIFile   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
}

func (f *EXE) GetsRunBy(runner *Runner) error {
	f.ExecutesFrom = runner
	return f.save()
}

func (x *EXE) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "EXE"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	params := make(map[string]interface{})
	if x.Name != "" {
		params["name"] = x.Name
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, params, x)
}

func (u *EXE) SetName(name string) error {
	u.Name = name
	return u.save()
}

func (x *EXE) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := newNeoSession()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
