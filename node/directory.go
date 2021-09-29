package node

import (
	"context"
	"fmt"

	"github.com/audibleblink/pegopher/db"
)

type Directory struct {
	securableIFile

	ContainedExes []*EXE  `gogm:"direction=outgoing;relationship=CONTAINS"`
	ContainedDlls []*DLL  `gogm:"direction=outgoing;relationship=CONTAINS"`
	HostsPEs      *Runner `gogm:"direction=outgoing;relationship=HOSTS_PES"`
}

func (d *Directory) HostsPEsFor(runner *Runner) error {
	d.HostsPEs = runner
	return d.save()
}

func (d *Directory) Add(ifile interface{}) error {
	switch f := ifile.(type) {
	case *EXE:
		d.ContainedExes = append(d.ContainedExes, f)
	case *DLL:
		d.ContainedDlls = append(d.ContainedDlls, f)
	}

	return d.save()
}

// Merge will either create or retreive the node based on the key-valie pair provides
// In this case, the Directory struct designates the "path" field as unique
func (x *Directory) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "Directory"
	sess, err := db.Session()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	if x.Name != "" {
		queryTemplate = `MERGE (x:%s {%s: "%s", name: "%s"}) RETURN x`
		query = fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue, x.Name)
	}
	return sess.Query(context.Background(), query, nil, x)
}

func (x *Directory) SetName(name string) error {
	x.Name = name
	return x.save()
}

func (x *Directory) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := db.Session()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
