package node

import (
	"context"
	"fmt"

	"github.com/audibleblink/pegopher/db"
)

type DLL struct {
	containableIFile
	securableIFile
	// Imports      []*PE   `gogm:"direction=outgoing;relationship=IMPORTS"`
	// ImportedBy   []*PE   `gogm:"direction=incoming;relationship=IMPORTED_BY"`
	// ExecutedFrom *Runner `gogm:"direction=outgoing;relationship=DLLCUTED_FROM"`
}

func (x *DLL) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "DLL"
	sess, err := db.Session()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, nil, x)
}

func (u *DLL) UpsertName(name string) error {
	if u.Name == name {
		return nil
	}
	u.Name = name
	return u.save()
}

func (x *DLL) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := db.Session()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
