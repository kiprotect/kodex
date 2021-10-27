// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

type UserProfile interface {
	Source() string
	SourceID() []byte
	EMail() string
	SuperUser() bool
	DisplayName() string
	AccessToken() AccessToken
	Roles() []OrganizationRoles
	Limits() map[string]interface{}
}

type AccessToken interface {
	Scopes() []string
}

type OrganizationRoles interface {
	Roles() []string
	Organization() UserOrganization
}

type UserOrganization interface {
	Name() string
	Source() string
	Default() bool
	Description() string
	ID() []byte
	ApiOrganization(Controller) (Organization, error)
}

type UserProfileProvider interface {
	Get(string) (UserProfile, error)
	Start()
	Stop()
}
