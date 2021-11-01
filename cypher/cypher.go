package cypher

import (
	"fmt"
	"strings"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const (
	MergeTmpl  = `MERGE (%s:%s { %s: '%s' })`
	RelateTmpl = `MERGE (%s)-[:%s]->(%s)`
	SetTmpl    = `SET %s.%s = '%s'`
)

var Driver neo4j.Driver

type Query struct {
	b *strings.Builder
	d neo4j.Driver
	l *logerr.Logger
}

func InitDriver(host, user, passwd string) (err error) {
	log := logerr.Add("neo4j driver init")

	Driver, err = neo4j.NewDriver(
		host,
		neo4j.BasicAuth(user, passwd, ""),
	)
	if err != nil {
		err = log.Wrap(err)
		return
	}

	log.Infof("global driver initialized for %s@%s", user, host)
	return
}

func NewQuery() (*Query, error) {
	log := logerr.Add("neo4j query")

	if Driver == nil {
		return &Query{}, log.Wrap(fmt.Errorf("uninitiailized driver"))
	}

	return &Query{
		b: &strings.Builder{},
		d: Driver,
		l: &log,
	}, nil
}

func (q *Query) Merge(varr, label, uniqProp, value string) *Query {
	value = util.PathFix(value)
	fmt.Fprintf(q.b, MergeTmpl, varr, label, uniqProp, value)
	fmt.Fprintf(q.b, "\n")
	return q
}

func (q *Query) Set(varr, prop, value string) *Query {
	value = util.PathFix(value)
	fmt.Fprintf(q.b, SetTmpl, varr, prop, value)
	fmt.Fprintf(q.b, "\n")
	return q
}

func (q *Query) Relate(var1, rel, var2 string) *Query {
	fmt.Fprintf(q.b, RelateTmpl, var1, rel, var2)
	fmt.Fprintf(q.b, "\n")
	return q
}

func (q *Query) ExecuteW() error {
	sess := q.d.NewSession(neo4j.SessionConfig{})
	result, err := sess.WriteTransaction(q.txWork)
	if err != nil {
		return q.l.Wrap(err)
	}

	var (
		ok      bool
		summary neo4j.ResultSummary
	)

	if summary, ok = result.(neo4j.ResultSummary); !ok {
		q.l.Debugf("failed to cast %v:", result)
	}

	if q.l.Level <= logerr.LogLevelDebug {
		res := map[string]int{
			"created":   summary.Counters().NodesCreated(),
			"props set": summary.Counters().PropertiesSet(),
			"new rels":  summary.Counters().RelationshipsCreated(),
		}
		q.l.Debugf("query result %#v:", res)
	}

	return nil
}

func (q *Query) txWork(tx neo4j.Transaction) (interface{}, error) {
	result, err := tx.Run(q.b.String(), nil)
	if err != nil {
		return nil, err
	}
	summary, err := result.Consume()
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (q *Query) ExecuteR() error {
	return nil
}

func (q *Query) ToString() string {
	return q.b.String()
}
