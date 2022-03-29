package collectors

import (
	"fmt"

	"github.com/audibleblink/pegopher/logerr"
	winapi "github.com/gueencode/go-win64api"
)

func CreateGroupPrincipals() error {
	log := logerr.Add("createGroupPrincipals")
	groups, err := winapi.ListLocalGroups()
	if err != nil {
		return log.Add("listLocalGroups").Wrap(err)
	}

	for _, group := range groups {
		principal := Principal{}
		principal.Name = group.Name
		principal.Write(writers[PrincipalFile])

		err := CreateGroupMemberPrincipals(group.Name)
		if err != nil {
			log.Infof("%v", err)
			continue
		}
	}

	return nil
}

func CreateGroupMemberPrincipals(group string) error {
	log := logerr.Add("createGroupPrincipals")
	users, err := winapi.LocalGroupGetMembers(group)
	if err != nil {
		return log.Add("listLocalGroupMembers").Wrap(err)
	}

	for _, user := range users {
		principal := Principal{}
		principal.Name = fmt.Sprintf(`%s\%s`, user.Domain, user.Name)
		principal.Group = group
		principal.Write(writers[PrincipalFile])
	}

	return nil
}
