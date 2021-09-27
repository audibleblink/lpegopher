package main

import (
	"context"
	"fmt"
)

type User struct {
	Principal

	Groups       []*Group `gogm:"direction=outgoing;relationship=MEMBER_OF"`
	ExecutedFrom *Runner  `gogm:"direction=incoming;relationship=EXECUTES_AS"`
}

func (x *User) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "User"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, lower(propValue))
	return sess.Query(context.Background(), query, nil, x)
}

func (u *User) JoinGroup(group *Group) {
	u.Groups = append(u.Groups, group)
}
