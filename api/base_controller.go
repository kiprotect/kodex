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

package api

import (
	"bytes"
	"github.com/kiprotect/kodex"
)

type BaseController struct {
	Definitions_  *Definitions
	UserProvider_ UserProvider
	Self          Controller
}

func (b *BaseController) APIDefinitions() *Definitions {
	return b.Definitions_
}

func (b *BaseController) UserProvider() UserProvider {
	return b.UserProvider_
}

func (b *BaseController) ObjectRolesForUser(objectType string, user *User) ([]ObjectRole, error) {
	objectRoles := make([]ObjectRole, 0)
	for _, organizationRoles := range user.Roles() {
		// we retrieve the organization of the user
		apiOrg, err := organizationRoles.Organization().ApiOrganization(b.Self)
		if err != nil {
			return nil, err
		}
		newRoles, err := b.Self.ObjectRolesForOrganizationRoles(objectType, organizationRoles.Roles(), apiOrg.ID())
		if err != nil {
			return nil, err
		}
		objectRoles = append(objectRoles, newRoles...)
	}
	return objectRoles, nil
}

func (b *BaseController) CanAccess(user *User, object kodex.Model, objectRoles []string) (bool, error) {

	// we retrive all organization roles for this object
	roles, err := b.Self.RolesForObject(object)

	kodex.Log.Info(roles)

	if err != nil {
		return false, err
	}

	for _, organizationRoles := range user.Roles() {

		apiOrg, err := organizationRoles.Organization().ApiOrganization(b.Self)
		if err != nil {
			return false, err
		}
		organizationID := apiOrg.ID()
		userRoles := organizationRoles.Roles()

		for _, role := range roles {
			if !bytes.Equal(organizationID, role.OrganizationID()) {
				// this role is for another organization
				continue
			}
			for _, userRole := range userRoles {
				// users with the "superuser" role always have access to the object
				if userRole == "superuser" {
					return true, nil
				}
				if userRole == role.OrganizationRole() {
					// the user has a matching role
					for _, objectRole := range objectRoles {
						if objectRole == role.ObjectRole() {
							// there's a matching object role
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}

func (b *BaseController) RegisterAPIPlugin(plugin APIPlugin) error {
	b.Definitions_.Routes = append(b.Definitions_.Routes, plugin.InitializeAPI)
	if err := plugin.InitializeAdaptors(b.Definitions_.ObjectAdaptors); err != nil {
		return err
	}
	return nil
}
