package cypher

import (
	"fmt"
	"strings"

	"github.com/audibleblink/pegopher/util"
)

const (
	MergeTmpl  = `MERGE (%s:%s { %s: "%s" })`
	RelateTmpl = `(%s)-[%s]->(%s)`
	SetTmpl    = `SET %s.%s = "%s"`
)

type Query struct {
	b *strings.Builder
}

func NewQuery() *Query {
	return &Query{
		b: &strings.Builder{},
	}
}

func (q *Query) Merge(varr, label, property, value string) *Query {
	value = util.PathFix(value)
	fmt.Fprintf(q.b, MergeTmpl, varr, label, property, value)
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

func (q *Query) ToString() string {
	return q.b.String()
}
