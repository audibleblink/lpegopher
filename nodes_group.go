package main

type Group struct {
	Principal

	Users []*User `gogm:"direction=incoming;relationship=MEMBER_OF"`
}
