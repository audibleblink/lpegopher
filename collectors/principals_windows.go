package collectors

import (
	winapi "github.com/gueencode/go-win64api"

	"github.com/audibleblink/logerr"
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
		principal.Type = "group"
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
	log := logerr.Add("createGroupMemberPrincipals")
	users, err := winapi.LocalGroupGetMembers(group)
	if err != nil {
		return log.Wrap(err)
	}

	for _, user := range users {
		principal := Principal{}
		principal.Name = user.DomainAndName
		principal.Group = group
		principal.Type = "user"
		principal.Write(writers[PrincipalFile])
	}

	return nil
}
