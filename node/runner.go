package node

import (
	"context"
	"fmt"

	"github.com/audibleblink/pegopher/db"
	"github.com/mindstand/gogm/v2"
)

// Runners can be Tasks, Services, or AutoRuns
type Runner struct {
	// Querier
	gogm.BaseNode

	Name    string     `gogm:"name=name;unique"`
	Type    string     `gogm:"name=type"`
	Context *User      `gogm:"direction=outgoing;relationship=EXECUTES_AS"`
	Exe     *EXE       `gogm:"direction=incoming;relationship=EXECUTED_FROM"`
	ExeDir  *Directory `gogm:"direction=incoming;relationship=HOSTS_PES"`
}

func (r *Runner) RunsExeAs(user *User) error {
	r.Context = user
	return r.save()
}

// Merge will either create or retreive the node based on the key-valie pair provides
// In this case, the Runner struct designates the "name" field as unique
func (x *Runner) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "Runner"
	sess, err := db.Session()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	if x.Type != "" {
		queryTemplate = `MERGE (x:%s {%s: "%s", type: "%s"}) RETURN x`
		query = fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue, x.Type)
	}
	return sess.Query(context.Background(), query, nil, x)
}

// type Task struct{ Runner }
// type Service struct{ Runner }

func (x *Runner) SetType(name string) error {
	x.Type = name
	return x.save()
}

func (x *Runner) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := db.Session()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
