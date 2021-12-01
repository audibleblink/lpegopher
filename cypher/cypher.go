package cypher

import (
	"fmt"
	"strings"

	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const (
	CreateTmpl = `CREATE (%s:%s { %s: '%s' })`
	MatchTmpl  = `MATCH (%s:%s { %s: '%s' })`
	MergeTmpl  = `MERGE (%s:%s { %s: '%s' })`
	RelateTmpl = `MERGE (%s)-[:%s]->(%s) `
	SetTmpl    = `SET %s.%s = '%s'`
	SetProps   = `, %s.%s = '%s'`
)

var Driver neo4j.Driver

type Query struct {
	b *strings.Builder
	d neo4j.Driver
	l *logerr.Logger
}

func InitDriver(host, user, passwd string) (err error) {
	log := logerr.G
	log.SetContext("neo4j driver init")

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
	return q.getAction(MergeTmpl, varr, label, uniqProp, value)
}

func (q *Query) Create(varr, label, uniqProp, value string) *Query {
	return q.getAction(CreateTmpl, varr, label, uniqProp, value)
}

func (q *Query) Match(varr, label, uniqProp, value string) *Query {
	return q.getAction(MatchTmpl, varr, label, uniqProp, value)
}

func (q *Query) getAction(template, varr, label, uniqProp, value string) *Query {
	value = util.PathFix(value)
	fmt.Fprintf(q.b, template, varr, label, uniqProp, value)
	fmt.Fprintf(q.b, " ")
	return q
}

func (q *Query) Append(query string) *Query {
	q.b.WriteString(query)
	q.b.WriteString("\n")
	return q
}

func (q *Query) With(label string) *Query {
	return q.Append(fmt.Sprintf("WITH %s", label))
}

func (q *Query) EndMerge() *Query {
	return q.Append("WITH count(*) as dummy")
}

func (q *Query) Return() *Query {
	return q.Append("RETURN count(*)")
}

func (q *Query) Terminate() *Query {
	return q.Append("")
}

func (q *Query) Set(varr string, props map[string]string) *Query {
	first := true
	for key, value := range props {
		value = util.PathFix(value)
		if first {
			fmt.Fprintf(q.b, SetTmpl, varr, key, value)
			first = false
			continue
		}
		fmt.Fprintf(q.b, SetProps, varr, key, value)
	}
	fmt.Fprintf(q.b, " ")
	return q
}

func (q *Query) Relate(var1, rel, var2 string) *Query {
	fmt.Fprintf(q.b, RelateTmpl, var1, rel, var2)
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
		q.l.Debugf("failed to create summary %v:", result)
	}

	if q.l.Level == logerr.LogLevelDebug {
		res := map[string]int{
			"created":   summary.Counters().NodesCreated(),
			"props set": summary.Counters().PropertiesSet(),
			"new rels":  summary.Counters().RelationshipsCreated(),
		}
		q.l.Debugf("query result %#v:", res)
	}

	return nil
}

func (q *Query) Raw(query string) *Query {
	q.b.Reset()
	fmt.Fprint(q.b, query)
	return q
}

func (q *Query) Reset() *Query {
	q.b.Reset()
	return q
}

func (q *Query) String() string {
	return q.b.String()
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

func (q *Query) Begin() (neo4j.Transaction, error) {
	sess := q.d.NewSession(neo4j.SessionConfig{})
	return sess.BeginTransaction()
}
