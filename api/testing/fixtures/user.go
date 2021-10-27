// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

type TestUserProfile struct {
	api.BaseUserProfile
	email       string
	displayName string
	sourceID    []byte
	superUser   bool
	limits      map[string]interface{}
	token       *TestAccessToken
	roles       []api.OrganizationRoles
}

func (t *TestUserProfile) Source() string {
	return "test"
}

func (t *TestUserProfile) Limits() map[string]interface{} {
	return t.limits
}

func (t *TestUserProfile) SourceID() []byte {
	return t.sourceID
}

func (t *TestUserProfile) SuperUser() bool {
	return t.superUser
}

func (t *TestUserProfile) EMail() string {
	return t.email
}

func (t *TestUserProfile) DisplayName() string {
	return t.displayName
}

func (t *TestUserProfile) AccessToken() api.AccessToken {
	return t.token
}

func (t *TestUserProfile) Roles() []api.OrganizationRoles {
	return t.roles
}

type TestAccessToken struct {
	api.BaseAccessToken
	scopes []string
}

func (t *TestAccessToken) Scopes() []string {
	return t.scopes
}

type TestOrganization struct {
	api.BaseUserOrganization
	name        string
	description string
	org         api.Organization
	id          []byte
}

func (t *TestOrganization) Default() bool {
	return true
}

func (t *TestOrganization) Name() string {
	return t.name
}

func (t *TestOrganization) Source() string {
	return "test"
}

func (t *TestOrganization) Description() string {
	return t.description
}

func (t *TestOrganization) ID() []byte {
	return t.id
}

type TestOrganizationRoles struct {
	api.BaseOrganizationRoles
	organization *TestOrganization
	roles        []string
}

func (t *TestOrganizationRoles) Roles() []string {
	return t.roles
}

func (t *TestOrganizationRoles) Organization() api.UserOrganization {
	return t.organization
}

type Organization struct {
	Name string
}

func (o Organization) Setup(fixtures map[string]interface{}) (interface{}, error) {
	controller, err := GetController(fixtures)
	if err != nil {
		return nil, err
	}
	org := controller.MakeOrganization()
	org.SetSourceID(kodex.RandomID())
	org.SetSource("test")
	org.SetDescription("")
	org.SetName(o.Name)
	if err := org.Save(); err != nil {
		return nil, err
	}
	return org, nil
}

func (o Organization) Teardown(fixture interface{}) error {
	return nil
}

type User struct {
	Limits       map[string]interface{}
	EMail        string
	SuperUser    bool
	Organization string
	Roles        []string
	Scopes       []string
}

func (u User) Setup(fixtures map[string]interface{}) (interface{}, error) {
	apiOrg, ok := fixtures[u.Organization].(api.Organization)
	if !ok {
		return nil, fmt.Errorf("no organization defined")
	}

	org := &TestOrganization{
		name:        apiOrg.Name(),
		description: apiOrg.Description(),
		id:          apiOrg.SourceID(),
	}

	org.Self = org

	token := &TestAccessToken{
		scopes: u.Scopes,
	}

	token.Self = token

	roles := &TestOrganizationRoles{
		organization: org,
		roles:        u.Roles,
	}

	roles.Self = roles

	profile := &TestUserProfile{
		limits:    u.Limits,
		email:     u.EMail,
		superUser: u.SuperUser,
		roles:     []api.OrganizationRoles{roles},
		token:     token,
	}

	profile.Self = profile

	return profile, nil
}

func (u User) Teardown(fixture interface{}) error {
	return nil
}
