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

package api

import (
	"github.com/kiprotect/kodex"
)

type SettingsValidator func(settings map[string]interface{}) (interface{}, error)

type APIControllerMaker func(
	config map[string]interface{},
	baseController kodex.Controller,
	definitions *Definitions) (Controller, error)

type APIControllerDefinitions map[string]APIControllerMaker

type Controller interface {
	kodex.Controller

	KodexController() kodex.Controller
	RegisterAPIPlugin(APIPlugin) error
	APIDefinitions() *Definitions

	// User provider
	UserProvider() (UserProvider, error)

	// Object roles
	CanAccess(user *User, object kodex.Model, objectRoles []string) (bool, error)
	ObjectRole(id []byte) (ObjectRole, error)
	RolesForObject(object kodex.Model) ([]ObjectRole, error)
	ObjectRolesForUser(objectType string, user *User) ([]ObjectRole, error)
	ObjectRolesForOrganizationRoles(objectType string, organizationRoles []string, organizationID []byte) ([]ObjectRole, error)
	MakeObjectRole(object kodex.Model, organization Organization) ObjectRole

	// Organizations
	MakeOrganization() Organization
	Organization(source string, sourceID []byte) (Organization, error)
	Organizations(filters map[string]interface{}) ([]Organization, error)
}
