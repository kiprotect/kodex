// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	Email        string
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

	org := api.MakeUserOrganization("test", apiOrg.Name(), apiOrg.Description(), apiOrg.SourceID())

	token := api.MakeAccessToken(u.Scopes)

	roles := api.MakeOrganizationRoles(org, u.Roles)

	user := api.MakeUser("test", u.Email, u.SuperUser, []*api.OrganizationRoles{roles}, u.Limits, token)

	return user, nil
}

func (u User) Teardown(fixture interface{}) error {
	return nil
}
