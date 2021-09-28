package node

import (
	"context"
	"fmt"

	"github.com/audibleblink/pegopher/db"
)

type User struct {
	Principal

	Groups       []*Group `gogm:"direction=outgoing;relationship=MEMBER_OF"`
	ExecutedFrom *Runner  `gogm:"direction=incoming;relationship=EXECUTES_AS"`
}

func (x *User) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "User"
	sess, err := db.Session()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, propValue)
	return sess.Query(context.Background(), query, nil, x)
}

func (u *User) JoinGroup(group *Group) {
	u.Groups = append(u.Groups, group)
}

func (x *User) Save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := db.Session()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}
